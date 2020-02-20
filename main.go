package main

import (
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"time"
)

//type a map.s.int32

/*
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!! VERSION !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
*/
const version = "0.0.1"

/*
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!! VERSION !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
*/

func main() {
	addressIP := flag.String("ip", "localhost", "a string")
	tcpPort := flag.String("port", "502", "a string")
	slaveID := flag.Int("id", 1, "an int")
	regQuantity := flag.Uint("q", 61, "an uint")
	flag.Parse()
	serverParam := fmt.Sprint(*addressIP, ":", *tcpPort)
	s := byte(*slaveID)

	//	fmt.Println(serverParam)

	handler := modbus.NewTCPClientHandler(serverParam)
	handler.SlaveId = s
	handler.Timeout = 2 * time.Second
	// Connect manually so that multiple requests are handled in one session
	err := handler.Connect()
	defer handler.Close()
	client := modbus.NewClient(handler)

	results, err := client.ReadHoldingRegisters(19000, uint16(*regQuantity)*2)
	if err != nil {
		fmt.Printf("{\"status\":\"error\", \"error\":\"%s\", \"version\": \"%s\"}", err, version)
		//fmt.Printf("%s\n", err)
	}

	fmt.Println(len(results))
	//fmt.Println(results)
}
