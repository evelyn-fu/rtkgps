package ntrip_receiver

import (
	"fmt"
	"log"
	"io"
	"bufio"

	"github.com/de-bkg/gognss/pkg/ntrip"
	"github.com/go-gnss/rtcm/rtcm3"
	"github.com/jacobsa/go-serial/serial"
)

// attempts to connect to ntrip client until successful connection or timeout
func Connect(casterAddr string, user string, pwd string, maxAttempts int) (*ntrip.Client, error) {
	success := false
	attempts := 0

	var c *ntrip.Client
	var err error

	fmt.Print("Connecting")
	for !success && attempts < maxAttempts {
		fmt.Print("...")
		c, err = ntrip.NewClient(casterAddr, ntrip.Options{Username: user, Password: pwd})
		if err == nil {
			success = true
		}
		attempts++
	}
	fmt.Print("\n")

	if err != nil {
		log.Fatal(err)
	}

	return c, err
}

// attempts to connect to ntrip streak until successful connection or timeout
func GetStream(c *ntrip.Client, mountPoint string, maxAttempts int) (io.ReadCloser, error) {
	success := false
	attempts := 0

	var rc io.ReadCloser
	var err error

	fmt.Print("Getting Stream")

	for !success && attempts < maxAttempts {
		fmt.Print(("..."))
		rc, err = c.GetStream(mountPoint)
		if err == nil {
			success = true
		}
		attempts++
	}
	fmt.Print("\n")

	if err != nil {
		log.Fatal(err)
	}

	return rc, err
}

// Connects to ntrip client and reads stream for specific mountpoint
func Receive(casterAddr string, mountPoint string, user string, pwd string, isStalled chan<- bool) {
	isStalled <- true
    c, err := Connect(casterAddr, user, pwd, 10)
    defer c.CloseIdleConnections()

    if !c.IsCasterAlive() {
        log.Printf("caster %s seems to be down", casterAddr)
    }

	options := serial.OpenOptions{
		PortName: "/dev/serial/by-id/usb-u-blox_AG_-_www.u-blox.com_u-blox_GNSS_receiver-if00",
		BaudRate: 38400,
		DataBits: 8,
		StopBits: 1,
		MinimumReadSize: 1,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}
	defer port.Close()

	w := bufio.NewWriter(port)

	rc, err := GetStream(c, mountPoint, 10)
	defer rc.Close()

	r := io.TeeReader(rc, w)
	scanner := rtcm3.NewScanner(r)
	isStalled <- false

	fmt.Print("Stream: ")
	for err == nil {
		msg, err := scanner.NextMessage()
		if err != nil {
			if msg == nil {
				isStalled <- true
				fmt.Println("No message... reconnecting to stream...")
				rc, err = GetStream(c, mountPoint, 10)
				defer rc.Close()

				r = io.TeeReader(rc, w)
				scanner = rtcm3.NewScanner(r)
				isStalled <- false
				continue
			}
			log.Fatal(err, msg)
		}
	
		// fmt.Printf("%T\t", msg)
		// fmt.Println(msg)
	}

}