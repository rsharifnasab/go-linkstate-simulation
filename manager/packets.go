package main

import (
	"fmt"
	"time"
)

type Packet struct {
	From int    `mapstructure:"from"`
	To   int    `mapstructure:"to"`
	Data string `mapstructure:"data"`
}

func (manager *Manager) sendTestPackets() {
	tests := loadYaml("tests")
	packets := make([]Packet, 0)
	pnc(tests.UnmarshalKey("tests", &packets))
	for _, packet := range packets {
		// put packet in sending queue
		manager.routers[packet.From].packetChannel <- serializePacket(&packet)
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	for i := range manager.routers {
		manager.routers[i].packetChannel <- "QUIT"
	}
}

func serializePacket(packet *Packet) string {
	return fmt.Sprintf("%d %d %s", packet.From, packet.To, packet.Data)
}
