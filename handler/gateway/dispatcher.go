package gateway

import "fmt"

type Dispatcher struct {
	handlers map[uint16]map[uint32]func([]byte) error
}

type HandlerFunc func(msgID uint32, body []byte)

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[uint16]map[uint32]func([]byte) error),
	}
}

func (d *Dispatcher) RegisterHandler(msgType uint16, msgID uint32, handler func([]byte) error) {
	if _, exists := d.handlers[msgType]; !exists {
		d.handlers[msgType] = make(map[uint32]func([]byte) error)
	}
	d.handlers[msgType][msgID] = handler
}

func (d *Dispatcher) Dispatch(msgType uint16, msgID uint32, body []byte) error {
	if handlers, exists := d.handlers[msgType]; exists {
		if handler, exists := handlers[msgID]; exists {
			return handler(body)
		}
	}
	return fmt.Errorf("handler not found for msgType=%d, msgID=%d", msgType, msgID)
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
