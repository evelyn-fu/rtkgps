package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/jacobsa/go-serial/serial"
	"github.com/adrianmo/go-nmea"
)

func ReadNmea(isConnected <-chan bool) {
	<-isConnected
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

	incompleteSentence := ""
	
	for err == nil {
		<-isConnected
		buf := make([]byte, 256)
		n, err := port.Read(buf)
		if err != nil {
			log.Fatalf("%s\n", err)
		}

		sentences := strings.Split(string(buf[:n]), "\n")

		for i, sentence := range sentences {
			if len(sentence) == 0 {
				continue
			}

			var s nmea.Sentence

			// check if last incomplete sentence fits with first sentence
			if i == 0 && incompleteSentence != "" {
				s, err = nmea.Parse(incompleteSentence + sentence)
				// discard
				if err != nil {
					s, err = nmea.Parse(sentence)
				}
				incompleteSentence = ""
			} else {
				s, err = nmea.Parse(sentence)
			}

			if err != nil {
				incompleteSentence = sentence
				continue
			}
	
			if s.DataType() == nmea.TypeGGA {
				m := s.(nmea.GGA)
				fmt.Printf("Time: %s\n", m.Time)
				fmt.Printf("Fix Quality: %s\n", m.FixQuality)
				fmt.Printf("Longitude GPS: %s\n", nmea.FormatGPS(m.Longitude))
				fmt.Printf("Longitude DMS: %s\n", nmea.FormatDMS(m.Longitude))
				fmt.Printf("Latitude GPS: %s\n", nmea.FormatGPS(m.Latitude))
				fmt.Printf("Latitude DMS: %s\n", nmea.FormatDMS(m.Latitude))
				fmt.Printf("Altitude: %f\n\n", m.Altitude)
			}
		}
	}
}