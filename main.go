package main

import (
	"fmt"
	"github.com/shyinyong/go-tcp-test/consts"
	"github.com/shyinyong/go-tcp-test/db"
	"github.com/shyinyong/go-tcp-test/msg_packet"
	"github.com/shyinyong/go-tcp-test/pb/cs"

	"net"
	"sync"
	"time"
)

type Client struct {
	conn       net.Conn
	incoming   chan []byte
	outgoing   chan []byte
	disconnect chan struct{}
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:       conn,
		incoming:   make(chan []byte),
		outgoing:   make(chan []byte),
		disconnect: make(chan struct{}),
	}
}

type GameServer struct {
	listener    net.Listener
	clients     map[*Client]struct{}
	clientMutex sync.Mutex
}

func NewGameServer() (*GameServer, error) {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}

	return &GameServer{
		listener: listener,
		clients:  make(map[*Client]struct{}),
	}, nil
}

func (gs *GameServer) Run() {
	for {
		conn, err := gs.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		client := NewClient(conn)
		gs.AddClient(client)
		go gs.HandleClient(client)
	}
}

func (gs *GameServer) AddClient(client *Client) {
	gs.clientMutex.Lock()
	defer gs.clientMutex.Unlock()
	gs.clients[client] = struct{}{}
}

func (gs *GameServer) RemoveClient(client *Client) {
	gs.clientMutex.Lock()
	defer gs.clientMutex.Unlock()
	delete(gs.clients, client)
	close(client.incoming)
	close(client.outgoing)
	client.conn.Close()
}

func (gs *GameServer) HandleClient(client *Client) {
	defer gs.RemoveClient(client)
	go gs.ReadMessages(client)
	go gs.WriteMessages(client)

	clientActivityTimer := time.NewTimer(time.Minute)
	defer clientActivityTimer.Stop()
	for {
		select {
		case <-client.disconnect:
			return
		case <-clientActivityTimer.C:
			fmt.Println("Client is idle. Sleeping...")
			select {
			case <-client.disconnect:
				return
			case <-client.incoming:
				// 客户端发送消息，重新激活定时器
				clientActivityTimer.Reset(time.Minute)
			}
		}
	}
}

func (gs *GameServer) WriteMessages(client *Client) {
	for {
		select {
		case message, ok := <-client.outgoing:
			if !ok {
				return
			}

			_, err := client.conn.Write(message)
			if err != nil {
				fmt.Printf("Error writing message: %v\n", err)
				client.disconnect <- struct{}{}
				return
			}
		}
	}
}

func (gs *GameServer) ReadMessages(client *Client) {
	for {
		header, err := msg_packet.ParseHeader(client.conn)
		if err != nil {
			fmt.Printf("Error reading header: %v\n", err)
			client.disconnect <- struct{}{}
			return
		}

		body, err := msg_packet.ParseBody(client.conn, int(header.PackageLen-uint16(consts.HeaderSize)))
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
			client.disconnect <- struct{}{}
			return
		}

		go gs.HandleMessage(client, header.MsgID, body)
	}
}

func (gs *GameServer) HandleMessage(client *Client, msgID uint16, body []byte) {
	handler := cs.GetNetMsgHandler(msgID)
	if handler == nil {
		fmt.Printf("No handler found for message ID %d\n", msgID)
		return
	}

	// 构建消息头信息
	hdr := &msg_packet.MsgHdr{
		MsgId: msgID,
		Conn:  client.conn,
		Data:  body,
	}

	// 解析消息并调用处理函数
	msgHandler := handler.ParseCb(body)
	msgHandlerTyped, ok := msgHandler.(cs.MsgHandler)
	if !ok {
		fmt.Printf("Invalid message handler type for message ID %d\n", msgID)
		return
	}
	// func DispatchMsg(handler *CsNetMsgHandler, hdr *msg_packet.MsgHdr, msgHandler MsgHandler) {
	cs.DispatchMsg(handler, hdr, msgHandlerTyped)

	// client.outgoing <- response

	//switch handler.HandlerId {
	//	msg := cs.ParsePb(msgID, body)
	//	response := handler(body)
	//	if response != nil {
	//		client.outgoing <- response
	//	}
	//}
}

func (gs *GameServer) RegisterHandler(msgID uint16, modelId int) {
	cs.RegHandlerId(msgID, modelId)
}

func main() {
	db.InitDB()
	gameServer, err := NewGameServer()
	if err != nil {
		fmt.Printf("Error creating GameServer: %v\n", err)
		return
	}
	gameServer.RegisterHandler(consts.CMMessageID_CMPing, 1)
	gameServer.RegisterHandler(consts.CMMessageID_CMLogin, 1)
	gameServer.Run()
}
