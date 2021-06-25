package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
			append(manager.NetConns[edge.Node1], Edge{Dest: edge.Node2, Cost: edge.Cost})
		manager.NetConns[edge.Node2] =
			append(manager.NetConns[edge.Node2], Edge{Dest: edge.Node1, Cost: edge.Cost})
	}
	return manager
}

type Router struct {
	reader *bufio.Reader
	writer *bufio.Writer
	index  int
	port   int
}

func newRouterConnection(routerIndex int, conn net.Conn) *Router {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	router := &Router{
		reader: reader,
		writer: writer,
		index:  routerIndex,
	}
	return router
}

func (router *Router) readString() string {
	str, err := router.reader.ReadString('\n')
	pnc(err)
	return strings.TrimSpace(str)
}

func (router *Router) readInt() int {
	num, err := strconv.Atoi(router.readString())
	pnc(err)
	return num
}

func (router *Router) writeAsString(obj interface{}) {
	router.writer.WriteString(fmt.Sprintf("%v\n", obj))
	router.writer.Flush()
}

func (router *Router) writeAsBytes(obj interface{}) {
	marshalled, err := json.Marshal(obj)
	pnc(err)
	_, err = router.writer.Write(marshalled)
	pnc(err)
	_, err = router.writer.Write([]byte("\n"))
	pnc(err)
	router.writer.Flush()
}

func (manager *Manager) handleRouter(routerIndex int, conn net.Conn) {
	router := newRouterConnection(routerIndex, conn)
	router.port = router.readInt()
	log.Printf("router #%v connected, udp port: %v\n", router.index, router.port)

	// send connectivity table
	router.writeAsString(manager.RoutersCount)
	router.writeAsBytes(manager.NetConns[router.index])
}
