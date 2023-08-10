package utils

import (
	"fmt"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ip, err := GetLocalIP()
	if err != nil {
		// 错误处理
		t.Errorf("local ip %v error:%s", ip, err)
	}
	fmt.Printf("local ip:[%s]\n", ip)
}
