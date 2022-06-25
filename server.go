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
	HTTP_WS_PORT := "8080"
	UDP_PORT := 7101

	LIFETIME = 30
	msgChan = make(chan *DeviceMsg)
	devicePool.Init()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: UDP_PORT,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	go UDPListen(conn)
	log.Println("UDP Server Started Successfully")

	manager = new(wsClientManager)
	manager.Init()
	go MsgDistributor(manager, msgChan) // msg分发给各个websocket前端
	log.Println("MsgDistributor Started Successfully")
	go manager.Start()
	log.Println("WebSocket Client Manager Started Successfully")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.LoadHTMLGlob("html/*")
	router.Static("static", "./static")

	router.GET("/switcher", switcherStatus) // web控制开关
	router.GET("/devices", devicesStatus)
	router.GET("/msgRecord", msgRecord) // 显示终端设备发来的消息

	router.GET("/dashboard", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "front.html", nil)
	})
	router.GET("/switch", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "switch.html", nil)
	})
	localhost, err := LocalHost()
	if err != nil {
		localhost = "0.0.0.0"
	}
	log.Printf("UDP Listening Port %d\n", UDP_PORT)
	log.Printf("WebSocket & HTTP Listening Port %s\n", HTTP_WS_PORT)
	log.Printf("Dashboard Running on %s", "http://"+localhost+":"+HTTP_WS_PORT+"/dashboard\n")
	log.Printf("Switch Running on %s", "http://"+localhost+":"+HTTP_WS_PORT+"/switch\n")

	log.Println("Server Started Successfully")
	router.Run("0.0.0.0:" + HTTP_WS_PORT)

}
