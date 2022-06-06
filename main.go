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

	isStalled := make(chan bool)

	go ntrip_receiver.Receive(casterAddr, mountPoint, user, pwd, isStalled)
	parser.ReadNmea(isStalled)
}
