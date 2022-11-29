package util

import (
	"errors"
	"net"
)

func SendWol(mac string, ip string) error {

	// parse mac address
	macAddr, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}

	if len(macAddr) != 6 {
		// return error with message "MAC address invalid length"
		return errors.New("MAC address invalid length")

	}
	// create 102 bytes to store the magic packet
	magicPacket := make([]byte, 102)
	// set the first 6 bytes to 0xFF
	for i := 0; i < 6; i++ {
		magicPacket[i] = 0xFF
	}
	// copy the mac address 16 times into the magic packet
	for i := 1; i <= 16; i++ {
		copy(magicPacket[i*6:(i+1)*6], macAddr)
	}

	//convert the last byte of ip to a broadcast address
	broadcast := net.ParseIP(ip)
	broadcast[15] = 255
	//convert broadcast ip to a string
	broadcastString := broadcast.String()

	conn, err := net.Dial("udp", broadcastString+":9")
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(magicPacket)
	return err
}
