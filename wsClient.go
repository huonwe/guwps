package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	UUID   string
	Socket *websocket.Conn
	Send   chan []byte
}

type wsClientManager struct {
	sync.RWMutex
	Clients    map[string]*wsClient
	Register   chan *wsClient
	Unregister chan *wsClient
}

// func (manager *wsClientManager) Add(uuid string, conn *websocket.Conn) {
// 	manager.Clients[uuid] = client
// }

func (client *wsClient) Init() {
	client.Send = make(chan []byte)
}

func (client *wsClient) Publish() {
	// fmt.Println(client.UUID + "'s publish loop start")
	missed := 0
	for {
		message, ok := <-client.Send
		if !ok {
			return
		}
		// fmt.Println("Publish: ", string(message))
		err := client.Socket.WriteMessage(websocket.TextMessage, message)
		// 错过消息3次则注销该client
		if err != nil {
			missed += 1
			if missed > 3 {
				return
			}
		}
		// default:
		// 	fmt.Println("client no msg")
		// 	client.Send <- []byte("test msg")
		// 	time.Sleep(1 * time.Second)
	}
}

func (manager *wsClientManager) Init() {
	manager.Clients = make(map[string]*wsClient)
	manager.Register = make(chan *wsClient)
	manager.Unregister = make(chan *wsClient)
}

func (manager *wsClientManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			// fmt.Println("Register: ", client.UUID)
			manager.Lock()
			manager.Clients[client.UUID] = client
			manager.Unlock()
		case client := <-manager.Unregister:
			client.Socket.Close()
			close(client.Send)
			delete(manager.Clients, client.UUID)
		}
	}
}
