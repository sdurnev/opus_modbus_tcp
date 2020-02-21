package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"math"
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
	{1, 2, "Data_version_counter", []string{"Data_version_counter"}},
	{2, 1, "Operating_mode", []string{
		"2_0_Float_char_act",
		"2_1_Bat_test_act",
		"2_2_Boost_char_act",
		"2_3_Temp_compen_act"}},
	{3, 1, "Battery_test_state", []string{
		"3_0_Perio_bat_test ",
		"3_1_Manual bat_test",
		"3_2_Natur_bat_test ",
		"3_3_Rem_bat_test   ",
	}},
	{4, 1, "Boost_charge_state", []string{
		"4_0_Auto_boost_char ",
		"4_1_Perio_boost_char",
		"4_2_Man_boost_char  ",
		"4_3_Rem_boost_char  ",
	}},
	{9, 2, "System_Voltage", []string{"bcmSystemVoltage"}},
	{10, 2, "Load_current", []string{"bcmLoadCurrent"}},
	{11, 2, "Battery_current", []string{"bcmBatteryCurrent"}},
	{12, 2, "Total_rectifier_current", []string{"bcmTotalRectifierCurrent"}},
	{13, 2, "Total_inverter_current", []string{"bcmTotalInverterCurrent"}},
	{14, 2, "Maximum_battery_temperature", []string{"bcmMaxBatteryTemperature"}},
	{15, 2, "Maximum_system_temperature", []string{"bcmMaxSystemTemperature"}},
	{29, 1, "System voltage alarms", []string{
		"30_0_Mains_fault",
		"30_1_Phase_fault",
		"30_2_Low_sys_volt",
		"30_3_High_sys_volt",
		"30_4_Float_char_devi",
		"30_5_Invert_sys_mains_fault",
	}},
	{30, 1, "System fault alarms", []string{
		"31_0_Earth_fault",
		"31_1_Load_fuse_fault",
		"31_2_Bat_fuse_fault",
		"31_3_Rectif_overlo",
		"31_4_Invert_overlo",
		"31_5_Bus_fault",
		"31_6_DCP_Bus_Fault",
		"31_7_Shunt_fault",
		"31_8_Sys_over_temp",
		"31_9_No_sys_temp_sens",
	}},
	{31, 1, "Miscellaneous system alarms", []string{
		"32_0_Boost_charge_act",
		"32_1_Config_confl",
		"32_2_Invent_full",
	}},
	{32, 1, "Rectifier alarms", []string{
		"33_0_Comm_error",
		"33_1_Nvram_fault",
		"33_2_Config_fault",
		"33_3_Module_fault",
		"33_4_Bad_firmware",
		"33_5_Rectif_fault",
		"33_6_Rectif_over_vol",
		"33_7_Rectif_over_temp",
		"33_8_Rectif_mains_fault",
		"33_9_Rectif_wrong_vol_vers",
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
	regQuantity := flag.Uint("q", 39, "an uint")
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

	results, err := client.ReadHoldingRegisters(0, uint16(*regQuantity))
	if err != nil {
		fmt.Printf("{\"status\":\"error\", \"error\":\"%s\", \"version\": \"%s\"}", err, version)
	}

	fmt.Println(len(results))
	fmt.Println(len(data))
	fmt.Println(results)
	fmt.Println(hex.EncodeToString(results))

	for i := 0; i < len(data); i++ {
		a := data[i].Req * 2
		d := results[a : a+2]

		f := binary.BigEndian.Uint16(d)
		fmt.Println(hex.EncodeToString(d))
		fmt.Printf("%b \n", d)
		fmt.Printf("%d \n", f)
		fmt.Println("====")
	}
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
