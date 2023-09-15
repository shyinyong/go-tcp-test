package main

import (
	"encoding/binary"
	"fmt"
	"github.com/shyinyong/go-tcp-test/consts"
	"github.com/shyinyong/go-tcp-test/db"
	"github.com/shyinyong/go-tcp-test/msg_packet"
	"github.com/shyinyong/go-tcp-test/pb/cs"

	"io"
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
	listener        net.Listener
	clients         map[*Client]struct{}
	clientMutex     sync.Mutex
	messageHandlers map[uint16]func([]byte) []byte // 存储消息处理函数
}

func NewGameServer() (*GameServer, error) {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}

	return &GameServer{
		listener:        listener,
		clients:         make(map[*Client]struct{}),
		messageHandlers: make(map[uint16]func([]byte) []byte),
	}, nil
}

// 运行游戏服务器
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

	// 使用一个定时器来检测客户端活动
	// 设定一个合适的时间间隔
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

func (gs *GameServer) ReadMessages(client *Client) {
	for {
		header, err := gs.ParseHeader(client.conn)
		if err != nil {
			fmt.Printf("Error reading header: %v\n", err)
			client.disconnect <- struct{}{}
			return
		}

		body, err := gs.ParseBody(client.conn, int(header.PackageLen-uint16(consts.HeaderSize)))
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
			client.disconnect <- struct{}{}
			return
		}

		go gs.HandleMessage(client, header.MsgID, body)
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

func (gs *GameServer) ParseHeader(conn net.Conn) (*msg_packet.MessageHeader, error) {
	headerBytes := make([]byte, consts.HeaderSize)
	_, err := io.ReadFull(conn, headerBytes)
	if err != nil {
		return nil, err
	}

	header := &msg_packet.MessageHeader{
		PackageLen: binary.BigEndian.Uint16(headerBytes[:2]),
		MsgID:      binary.BigEndian.Uint16(headerBytes[2:4]),
		SeqID:      binary.BigEndian.Uint32(headerBytes[4:8]),
		MagicCode:  binary.BigEndian.Uint16(headerBytes[8:10]),
		Reserved:   binary.BigEndian.Uint16(headerBytes[10:12]),
	}

	return header, nil
}

func (gs *GameServer) ParseBody(conn net.Conn, bodySize int) ([]byte, error) {
	body := make([]byte, bodySize)
	_, err := io.ReadFull(conn, body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (gs *GameServer) HandleMessage(client *Client, msgID uint16, body []byte) {
	handler := cs.GetNetMsgHandler(msgID)
	if handler != nil {
		switch handler.HandlerId {
		case consts.PLAYER_MGR_HANDLER_ID:
			//cs.DispatchMsg(handler, hdr, playerMgr)
		case consts.PLAYER_HANDLER_ID:
			//player := playerMgr.getPlayerBySocket(hdr.GetSocket())
			//if player != nil {
			//	cs.DispatchMsg(handler, hdr, player)
			//}
		}

		//
		//response := handler(body)
		//if response != nil {
		//	client.outgoing <- response
		//}
		//
	}
}

func (gs *GameServer) RegisterHandler(msgID uint16, handler func([]byte) []byte) {
	gs.messageHandlers[msgID] = handler
}

func main() {
	db.InitDB()
	gameServer, err := NewGameServer()
	if err != nil {
		fmt.Printf("Error creating GameServer: %v\n", err)
		return
	}
	gameServer.RegisterHandler(consts.CMMessageID_CMPing, CMPing)
	gameServer.RegisterHandler(consts.CMMessageID_CMLogin, CMLogin)
	gameServer.Run()
}
