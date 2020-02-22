package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"strconv"
	"strings"
	"time"
)

//type a map.s.int32

const version = "0.0.1"

type param struct {
	Req  int
	Type int
	Name string
	Par  []string
}

type opus_param []param

var data opus_param = opus_param{
	{1, 2, "DataVersionCounter", []string{"DataVersionCounter"}},
	{2, 1, "OperatingMode", []string{
		"2_0_FloatCharAct",
		"2_1_BatTestAct",
		"2_2_BoostCharAct",
		"2_3_TempCompenAct"}},
	{3, 1, "BatteryTestState", []string{
		"3_0_PerioBatTest",
		"3_1_ManualBatTest",
		"3_2_NaturBatTest",
		"3_3_RemBatTest",
	}},
	{4, 1, "BoostChargeState", []string{
		"4_0_AutoBoostChar",
		"4_1_PerioBoostChar",
		"4_2_ManBoostChar",
		"4_3_RemBoostChar",
	}},
	{9, 2, "SystemVoltage", []string{"bcmSystemVoltage"}},
	{10, 2, "LoadCurrent", []string{"bcmLoadCurrent"}},
	{11, 2, "BatteryCurrent", []string{"bcmBatteryCurrent"}},
	{12, 2, "TotalRectifierCurrent", []string{"bcmTotalRectifierCurrent"}},
	{13, 2, "TotalInverterCurrent", []string{"bcmTotalInverterCurrent"}},
	{14, 2, "MaximumBatteryTemperature", []string{"bcmMaxBatteryTemperature"}},
	{15, 2, "MaximumSystemTemperature", []string{"bcmMaxSystemTemperature"}},
	{29, 1, "SystemVoltageAlarms", []string{
		"30_0_MainsFault",
		"30_1_PhaseFault",
		"30_2_LowSysVolt",
		"30_3_HighSysVolt",
		"30_4_FloatCharDevi",
		"30_5_InvertSysMainsFault",
	}},
	{30, 1, "SystemFaultAlarms", []string{
		"31_0_EarthFault",
		"31_1_LoadFuseFault",
		"31_2_BatFuseFault",
		"31_3_RectifOverlo",
		"31_4_InvertOverlo",
		"31_5_BusFault",
		"31_6_DCPBusFault",
		"31_7_ShuntFault",
		"31_8_SysOverTemp",
		"31_9_NoSysTempSens",
	}},
	{31, 1, "MiscellaneousSystemAlarms", []string{
		"32_0_BoostChargeAct",
		"32_1_ConfigConfl",
		"32_2_InventFull",
	}},
	{32, 1, "RectifierAlarms", []string{
		"33_0_CommError",
		"33_1_NvramFault",
		"33_2_ConfigFault",
		"33_3_ModuleFault",
		"33_4_BadFirmware",
		"33_5_RectifFault",
		"33_6_RectifOverVol",
		"33_7_RectifOverTemp",
		"33_8_RectifMainsFault",
		"33_9_RectifWrongVolVers",
	}},
	{33, 1, "InverterSystemAlarms", []string{
		"34_0_CommError",
		"34_1_NvramFault",
		"34_2_ConfigFault",
		"34_3_ModuleFault",
		"34_4_BadFirmware",
		"34_5_InverSysFault",
		"34_6_InverterFault",
		"34_7_BypassFault",
	}},
	{34, 1, "OtherModulesAlarms", []string{
		"35_0_CommError",
		"35_1_NvramFault",
		"35_2_ConfigFault",
		"35_3_ModuleFault",
		"35_4_BadFirmware",
	}},
	{35, 1, "BatteryAlarms", []string{
		"36_0_BatBloLowVolt",
		"36_1_BatBloHigVolt",
		"36_2_BatStriAsymmet",
		"36_3_AutoBoosCharFault",
		"36_4_BatTestFault",
		"36_5_BatOverTemp",
		"36_6_NoBatTempSens",
		"36_7_BatTempSensFault",
	}},
	{36, 1, "LowVoltageDisconnectionAlarms", []string{
		"37_0_LoLVDDisWarn",
		"37_1_LoLVDDisImmi",
		"37_2_BatLVDDisWarn",
		"37_3_BatLVDDisImmi",
		"37_4_ContFault",
	}},
	{37, 1, "ExternalAlarms", []string{
		"38_0_ExtAlmGr1",
		"38_1_ExtAlmGr2",
		"38_2_ExtAlmGr3",
		"38_3_ExtAlmGr4",
	}},
}

func main() {
	addressIP := flag.String("ip", "localhost", "a string")
	tcpPort := flag.String("port", "502", "a string")
	slaveID := flag.Int("id", 1, "an int")

	flag.Parse()
	serverParam := fmt.Sprint(*addressIP, ":", *tcpPort)
	s := byte(*slaveID)

	handler := modbus.NewTCPClientHandler(serverParam)
	handler.SlaveId = s
	handler.Timeout = 2 * time.Second
	// Connect manually so that multiple requests are handled in one session
	err := handler.Connect()
	defer handler.Close()
	client := modbus.NewClient(handler)

	results, err := client.ReadHoldingRegisters(0, 39)
	if err != nil {
		fmt.Printf("{\"status\":\"error\", \"error\":\"%s\", \"version\": \"%s\"}", err, version)
	}

	//fmt.Println(results)

	var tempStringArr []string //Массив сформированых стринговых данных для формирования json ответа
	for i := 0; i < len(data); i++ {
		reg := data[i].Req * 2
		regData := binary.BigEndian.Uint16(results[reg : reg+2])
		if i == 0 { // формирование первого параметра DataVersionCounter
			var par = strconv.FormatFloat(float64((float32(regData))), 'f', 2, 32)
			tmpstr := []string{"\"", data[i].Par[0], "\": ", par}
			var str = strings.Join(tmpstr, "")
			tempStringArr = append(tempStringArr, str)
		} else if i > 3 && i < 11 { // формирование флотовых параметров
			var par = strconv.FormatFloat(float64((float32(regData) / 10)), 'f', 2, 32)
			tmpstr := []string{"\"", data[i].Par[0], "\": ", par}
			var str = strings.Join(tmpstr, "")
			tempStringArr = append(tempStringArr, str)
		} else { //формирование битовых параметров
			var tStrAr string
			tStrAr = strings.Join([]string{"\"", data[i].Name, "\":{"}, "")
			for l := 0; l < len(data[i].Par); l++ {
				var t = strings.Join([]string{"\"", data[i].Par[l], "\": "}, "")
				if regData&(1<<l) == 0 {
					if l == len(data[i].Par)-1 {
						t = strings.Join([]string{t, "0", "}"}, "")
					} else {
						t = strings.Join([]string{t, "0", ","}, "")
					}
				} else {
					if l == len(data[i].Par)-1 {
						t = strings.Join([]string{t, "1", "}"}, "")
					} else {
						t = strings.Join([]string{t, "1", ","}, "")
					}
				}
				tStrAr = strings.Join([]string{tStrAr, t}, "")
			}
			tempStringArr = append(tempStringArr, tStrAr)
		}
	}
	for m := 0; m < len(tempStringArr); m++ { //Вывод данных/формирование json ответа
		if m == 0 { //Печать первого параметра
			fmt.Print("{")
			fmt.Print(tempStringArr[m])
			fmt.Print(",")
		} else if m == len(tempStringArr)-1 { //Печать последнего параметра
			fmt.Print(tempStringArr[m])
			fmt.Printf(", \"version\":\"%s\"} \n", version)
		} else {
			fmt.Print(tempStringArr[m])
			fmt.Print(",")
		}
	}
}
