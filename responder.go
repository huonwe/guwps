package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func switcherStatus(ctx *gin.Context) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	defer ws.Close()
	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println("ws data: ", string(msg))
		pData := PRead(msg)
		if pData == nil {
			log.Println("udp data read failed")
			break
		}
		log.Println("switcher ID: ", pData.ID)

		// fmt.Println("existed ID: ")
		// fmt.Println(devicePool.DevicesMap)
		// for k, _ := range devicePool.DevicesMap {
		// 	fmt.Println(k)
		// }

		device := devicePool.Find(pData.ID)
		if device == nil {
			fmt.Println("device not found")
			ws.WriteMessage(mt, []byte("device not found"))
			continue
		}
		addr, err := net.ResolveUDPAddr("udp", device.Addr)
		if err != nil {
			fmt.Println("addr convert failed")
			break
		}
		socket, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			fmt.Println("udp to client connect failed")
			break
		}
		fmt.Println("cmd writing to devices")
		_, err = socket.Write(msg)
		if err != nil {
			fmt.Println("udp send to client failed")
			break
		}
		data := make([]byte, 128)
		n, _, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("udp read from client failed")
			break
		}
		err = ws.WriteMessage(mt, data[:n])
		if err != nil {
			log.Println(err.Error())
			break
		}

	}
}

func msgRecord(ctx *gin.Context) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	fmt.Println("new ws eltablished")
	if err != nil {
		return
	}

	// defer ws.Close()

	s := time.Now().String() + strconv.Itoa(rand.Int())

	uuid := UUIDGenerate(s)
	client := new(wsClient)
	client.Init()
	client.UUID = uuid
	client.Socket = ws
	// fmt.Println("Registering ", client.UUID)
	manager.Register <- client

	client.Publish()

	// websocket end
	fmt.Println("websocket end")
	manager.Unregister <- client
	// for {
	// 	err =
	// 	if err != nil {
	// 		break
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }

}

func devicesStatus(ctx *gin.Context) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	defer ws.Close()
	for {

		// log.Println(string(msg))
		// list := listCache
		time.Sleep(1 * time.Second)
		// json, err := json.Marshal(devicePool.DevicesList)
		// if err != nil {
		// 	log.Println(err.Error())
		// 	break
		// }
		err = ws.WriteJSON(devicePool.DevicesList)
		if err != nil {
			// log.Println(err.Error())
			break
		}

	}
}
