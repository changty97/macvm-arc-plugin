package main

import (
	"log"
)

func main() {
	log.Println("Starting MacVM ARC Plugin...")
	plugin := NewMacVMProvider()
	if err := plugin.Serve(); err != nil {
		log.Fatalf("Plugin failed: %v", err)
	}
}
