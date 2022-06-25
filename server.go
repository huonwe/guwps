package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

var devicePool Pool

var manager *wsClientManager

var LIFETIME uint32 // seconds

func main() {
	LIFETIME = 30
	msgChan = make(chan *DeviceMsg)
	devicePool.Init()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 7101,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	go UDPListen(conn)
	log.Println("UDP Server Started")

	manager = new(wsClientManager)
	manager.Init()
	go MsgDistributor(manager, msgChan) // msg分发给各个websocket前端
	go manager.Start()

	router := gin.Default()
	router.LoadHTMLGlob("html/*")
	router.Static("static", "./static")
	router.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "switch.html", nil)
	})

	router.GET("/switcher", switcherStatus) // web控制开关
	router.GET("/devices", devicesStatus)
	router.GET("/msgRecord", msgRecord) // 显示终端设备发来的消息

	router.GET("/front", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "front.html", nil)
	})

	router.Run("0.0.0.0:8080")

}
