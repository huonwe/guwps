package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

var msgChan chan *DeviceMsg

func UDPListen(conn *net.UDPConn) {
	var data [128]byte
	for {
		n, addr, err := conn.ReadFromUDP(data[:])
		if err != nil {
			log.Println(err.Error())
			continue
		}

		content := PRead(data[:])
		if content == nil {
			log.Println("udp data read failed")
			continue
		}

		devicePool.Update(&Device{Addr: addr.String(), ID: content.ID, MM: content.MM, LastAlive: time.Now(), Lifetime: LIFETIME})

		_, err = conn.WriteToUDP(data[:n], addr) // 原样返回
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 发布Msg
		// 如果阻塞了就跳过
		go func() {
			select {
			case msgChan <- &DeviceMsg{Addr: addr.String(), Msg: data[:n]}:
				fmt.Println("UDP: ", string(data[:]))
			default:
				fmt.Println("skipped")
			}
			// conn.Write(data[:])
		}()

	}
}

func MsgDistributor(manager *wsClientManager, msgChan chan *DeviceMsg) {
	for {
		msg := <-msgChan
		// fmt.Println("Distributor: ", string(msg.Msg))
		s := [][]byte{[]byte(msg.Addr), msg.Msg}
		// fmt.Println("all clients: ", manager.Clients)
		for _, client := range manager.Clients {
			// fmt.Println("now distributing ", client.UUID)
			client.Send <- bytes.Join(s, []byte(" -> "))
		}
	}
}
