package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"math"
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
		"4_0_AutoBoostChar ",
		"4_1_PerioBoostChar",
		"4_2_ManBoostChar  ",
		"4_3_RemBoostChar  ",
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
	{33, 1, "Inverter system alarms", []string{
		"34_0_Comm_error",
		"34_1_Nvram_fault",
		"34_2_Config_fault",
		"34_3_Module_fault",
		"34_4_Bad_firmware",
		"34_5_Inver_sys_fault",
		"34_6_Inverter_fault",
		"34_7_Bypass_fault",
	}},
	{34, 1, "Other modules alarms", []string{
		"35_0_Comm_error",
		"35_1_Nvram_fault",
		"35_2_Config_fault",
		"35_3_Module_fault",
		"35_4_Bad_firmware",
	}},
	{35, 1, "Battery alarms", []string{
		"36_0_Bat_blo_low_volt",
		"36_1_Bat_blo_hig_volt",
		"36_2_Bat_Stri_Asymmet",
		"36_3_Auto_boos_char_fault",
		"36_4_Bat_Test_Fault",
		"36_5_Bat_Over_Temp",
		"36_6_No_Bat_Temp_Sens",
		"36_7_Bat_Temp_Sens_Fault",
	}},
	{36, 1, "Low voltage disconnection alarms", []string{
		"37_0_Lo_LVD_Dis_Warn",
		"37_1_Lo_LVD_Dis_Immi",
		"37_2_Bat_LVD_Dis_Warn",
		"37_3_Bat_LVD_Dis_Immi",
		"37_4_Cont_Fault",
	}},
	{37, 1, "External alarms", []string{
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

	fmt.Println(results)

	var tempStringArr []string
	for i := 0; i < len(data); i++ {
		reg := data[i].Req * 2
		regData := binary.BigEndian.Uint16(results[reg : reg+2])
		if i == 0 {
			var par = strconv.FormatFloat(float64((float32(regData))), 'f', 2, 32)
			tmpstr := []string{"\"", data[i].Par[0], "\": ", par}
			var str = strings.Join(tmpstr, "")
			tempStringArr = append(tempStringArr, str)
		} else if i > 3 && i < 11 {
			var par = strconv.FormatFloat(float64((float32(regData) / 10)), 'f', 2, 32)
			tmpstr := []string{"\"", data[i].Par[0], "\": ", par}
			var str = strings.Join(tmpstr, "")
			tempStringArr = append(tempStringArr, str)
		} else {
			//var tStrAr []string
			fmt.Print(data[i].Req)
			fmt.Print(" ")
			fmt.Println(data[i].Name)
			//fmt.Println(regData)
			for l := 0; l < 16; l++ {
				if regData&(1<<l) != 0 {
					fmt.Println(data[i].Par[l])
					//mesage = append(mesage, listOfAllarm1[i+1])
				}
				//fmt.Println(data[i].Par[l])
			}
		}
	}
	//fmt.Println(len(tempStringArr))
	fmt.Println(tempStringArr)
}

/*
for i := uint(0); i < 16; i++ {
if a&(1<<i) != 0 {
//fmt.Println(i)
mesage = append(mesage, listOfAllarm1[i+1])
}
}
*/

func MekeFlotPar(reg int, data []byte) string {
	//var retunData string

	return "string"
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
