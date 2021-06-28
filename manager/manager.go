package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Manager struct {
	routersCount        int
	routers             []*Router
	netConns            [][]*Edge
	readyWG             sync.WaitGroup
	networkReadyWG      sync.WaitGroup
	readyChannel        chan struct{}
	networkReadyChannel chan struct{}
}

type Edge struct {
	Dest int
	Cost int
}

type ConfigEdge struct {
	Node1 int `mapstructure:"node1"`
	Node2 int `mapstructure:"node2"`
	Cost  int `mapstructure:"cost"`
}

type Packet struct {
	From int    `mapstructure:"from"`
	To   int    `mapstructure:"to"`
	Data string `mapstructure:"data"`
}

func loadYaml(fileName string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	pnc(config.ReadInConfig())
	return config
}
func newManagerWithConfig(configFile string) *Manager {
	config := loadYaml(configFile)
	routersCount := config.GetInt("number_of_routers")
	manager := &Manager{
		routersCount:        routersCount,
		netConns:            make([][]*Edge, routersCount),
		readyWG:             sync.WaitGroup{},
		networkReadyWG:      sync.WaitGroup{},
		readyChannel:        make(chan struct{}),
		networkReadyChannel: make(chan struct{}),
		routers:             make([]*Router, routersCount),
	}
	for i := 0; i < manager.routersCount; i++ {
		manager.netConns[i] = make([]*Edge, 0)
		manager.routers[i] = &Router{Index: i, packetChannel: make(chan string)}
	}

	var configEdges []ConfigEdge
	pnc(config.UnmarshalKey("links", &configEdges))

	for _, configEdge := range configEdges {
		manager.netConns[configEdge.Node1] =
			append(manager.netConns[configEdge.Node1], &Edge{Dest: configEdge.Node2, Cost: configEdge.Cost})
		manager.netConns[configEdge.Node2] =
			append(manager.netConns[configEdge.Node2], &Edge{Dest: configEdge.Node1, Cost: configEdge.Cost})
	}
	return manager
}

func (manager *Manager) sendTestPackets() {
	tests := loadYaml("tests")
	packets := make([]Packet, 0)
	pnc(tests.UnmarshalKey("tests", &packets))
	for _, packet := range packets {
		manager.routers[packet.From].packetChannel <- serializePacket(&packet)
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	for i := range manager.routers {
		manager.routers[i].packetChannel <- "QUIT"
	}
}

func serializePacket(packet *Packet) string {
	return fmt.Sprintf("%d %d %s", packet.From, packet.To, packet.Data)
}
