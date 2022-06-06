package parser

import (
	"fmt"
	"log"

	"github.com/jacobsa/go-serial/serial"
	"github.com/adrianmo/go-nmea"
)

func ReadNmea(isStalled <-chan bool) {
	options := serial.OpenOptions{
		PortName: "/dev/serial/by-id/usb-u-blox_AG_-_www.u-blox.com_u-blox_GNSS_receiver-if00",
		BaudRate: 38400,
		DataBits: 8,
		StopBits: 1,
		MinimumReadSize: 4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()

	for err == nil {
		if <-isStalled {
			fmt.Println("No correction data:")
		}

		buf := make([]byte, 256)
		n, err := port.Read(buf)
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		fmt.Println(n, buf[:n])

		sentence := buf[:n]
		s, err := nmea.Parse(string(sentence))
		if err != nil {
			log.Fatal(err)
		}

		if s.DataType() == nmea.TypeGLL {
			m := s.(nmea.GGA)
			fmt.Printf("Longitude GPS: %s\n", nmea.FormatGPS(m.Longitude))
			fmt.Printf("Longitude DMS: %s\n", nmea.FormatDMS(m.Longitude))
			fmt.Printf("Latitude GPS: %s\n", nmea.FormatGPS(m.Latitude))
			fmt.Printf("Latitude DMS: %s\n", nmea.FormatDMS(m.Latitude))
			fmt.Printf("Altitude: %f\n\n", m.Altitude)
		}
	}
}