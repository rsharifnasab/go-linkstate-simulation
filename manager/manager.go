package main

import (
	"sync"

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

func loadConfig(configFile string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(configFile)
	config.AddConfigPath(".")
	pnc(config.ReadInConfig())
	return config
}
func newManagerWithConfig(configFile string) *Manager {
	config := loadConfig(configFile)
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
		manager.routers[i] = &Router{Index: i}
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
