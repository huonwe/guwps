package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
)

func PRead(data []byte) *UDPData {
	var err error
	s := string(data)
	r := new(UDPData)

	r.Length, err = strconv.ParseUint(s[0:4], 10, 32)
	if err != nil || r.Length < 41 {
		fmt.Println("parse uint failed")
		return nil
	}
	// fmt.Println("length: ", r.Length)

	r.Version = s[4:6]
	r.DeviceType = s[6:7]
	r.DataType = s[7:9]
	r.ID = s[9:19]
	r.MM = s[19:35]
	r.Tag = s[35:39]
	r.STA = s[39 : r.Length-2]
	r.Identifier = s[r.Length-2:]
	return r
}

func UUIDGenerate(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	hashInBytes := hash.Sum(nil)
	// fmt.Printf("%x\n", h.Sum(nil))
	// hs := fmt.Sprintf("%x", h.Sum(nil))
	hashValue := hex.EncodeToString(hashInBytes)
	return hashValue
}

func LocalHost() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		// fmt.Println(err)
		return "", err
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// fmt.Print(ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("failed")
}
