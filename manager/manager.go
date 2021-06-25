package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Manager struct {
	numberOfRouters    int
	networkConnections [][]Edge
}

type Edge struct {
	dest int
	cost int
}

func (manager *Manager) loadConfig() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
	manager.numberOfRouters = config.GetInt("number_of_routers")
	manager.networkConnections = make([][]Edge, manager.numberOfRouters)
	for i := 0; i < manager.numberOfRouters; i++ {
		manager.networkConnections[i] = make([]Edge, 0)
	}
	var edges []ConfigEdge
	_ = config.UnmarshalKey("links", &edges)
	for _, edge := range edges {
		manager.networkConnections[edge.Node1] =
			append(manager.networkConnections[edge.Node1], Edge{dest: edge.Node2, cost: edge.Cost})
		manager.networkConnections[edge.Node2] =
			append(manager.networkConnections[edge.Node2], Edge{dest: edge.Node1, cost: edge.Cost})
	}
}

type ConfigEdge struct {
	Node1 int `mapstructure:"node1"`
	Node2 int `mapstructure:"node2"`
	Cost  int `mapstructure:"cost"`
}

func (manager *Manager) handleRouter(routerIndex int, conn net.Conn) {
	defer conn.Close()
	log.Printf("Handling connection for router #%v\n", routerIndex)
	reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(conn)
	portStr, err := reader.ReadString('\n')
	pnc(err)
	port, _ := strconv.Atoi(strings.TrimSpace(portStr))
	log.Printf("client sent port: %v\n", port)
}
