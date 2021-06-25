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
	RoutersCount int
	NetConns     [][]Edge
}

type Edge struct {
	dest int
	cost int
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
func createManagerFromConfig(configFile string) *Manager {

	config := loadConfig(configFile)
	routersCount := config.GetInt("number_of_routers")
	manager := &Manager{
		RoutersCount: routersCount,
		NetConns:     make([][]Edge, routersCount),
	}
	for i := 0; i < manager.RoutersCount; i++ {
		manager.NetConns[i] = make([]Edge, 0)
	}

	var edges []ConfigEdge
	pnc(config.UnmarshalKey("links", &edges))

	for _, edge := range edges {
		manager.NetConns[edge.Node1] =
			append(manager.NetConns[edge.Node1], Edge{dest: edge.Node2, cost: edge.Cost})
		manager.NetConns[edge.Node2] =
			append(manager.NetConns[edge.Node2], Edge{dest: edge.Node1, cost: edge.Cost})
	}
	return manager
}

func (manager *Manager) handleRouter(routerIndex int, conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	portStr, err := reader.ReadString('\n')
	pnc(err)
	port, err := strconv.Atoi(strings.TrimSpace(portStr))
	pnc(err)
	//log.Printf("client sent port: %v\n", port)

	log.Printf("router #%v connected, udp port: %v\n", routerIndex, port)
	_ = writer
}
