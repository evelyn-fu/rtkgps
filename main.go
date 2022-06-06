package main

import (
	"rtkgps/rtkgps/nmea_parser"
	"rtkgps/rtkgps/ntrip_receiver"
)

func main() {
	casterAddr := "http://rtn.dot.ny.gov:8082"
	mountPoint := "NYGC"
	user := "evelyn"
	pwd := "checkmate"

	isConnected := make(chan bool)

	go ntrip_receiver.Receive(casterAddr, mountPoint, user, pwd, isConnected)
	parser.ReadNmea(isConnected)
}
