package config

import (
	"fmt"
	"log"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%v \n", cfg)
}
