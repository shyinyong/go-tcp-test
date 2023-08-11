package gateway

import "fmt"

type Dispatcher struct {
	handlers map[uint16]HandlerFunc
}

type HandlerFunc func(msgID uint32, body []byte)

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[uint16]HandlerFunc),
	}
}

func (d *Dispatcher) RegisterHandler(msgType uint16, handler HandlerFunc) {
	d.handlers[msgType] = handler
}

func (d *Dispatcher) Dispatch(msgType uint16, msgID uint32, body []byte) {
	handler, found := d.handlers[msgType]
	if !found {
		fmt.Printf("No handler found for msgType %d\n", msgType)
		return
	}

	handler(msgID, body)
}

// 示例：登录服务器处理函数
func LoginHandler(msgID uint32, body []byte) {
	// 根据 msgID 和 body 处理登录功能
	fmt.Println("Handling login message:", msgID, string(body))
}

// 示例：退出游戏处理函数
func ExitGameHandler(msgID uint32, body []byte) {
	// 根据 msgID 和 body 处理退出游戏功能
	fmt.Println("Handling exit game message:", msgID, string(body))
}
