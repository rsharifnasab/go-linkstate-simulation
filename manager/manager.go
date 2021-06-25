package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Manager struct {
	routersCount int
	netConns     [][]Edge
	readyWG      sync.WaitGroup
	readyChannel chan struct{}
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
		routersCount: routersCount,
		netConns:     make([][]Edge, routersCount),
		readyWG:      sync.WaitGroup{},
		readyChannel: make(chan struct{}),
	}
	for i := 0; i < manager.routersCount; i++ {
		manager.netConns[i] = make([]Edge, 0)
	}

	var edges []ConfigEdge
	pnc(config.UnmarshalKey("links", &edges))

	for _, edge := range edges {
		manager.netConns[edge.Node1] =
			append(manager.netConns[edge.Node1], Edge{Dest: edge.Node2, Cost: edge.Cost})
		manager.netConns[edge.Node2] =
			append(manager.netConns[edge.Node2], Edge{Dest: edge.Node1, Cost: edge.Cost})
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
	router.writeAsString(manager.routersCount)
	router.writeAsBytes(manager.netConns[router.index])
	readiness := router.readString()
	if readiness == "READY" {
		manager.readyWG.Done()
	} else {
		panic("Router couldn't get ready.")
	}
	<-manager.readyChannel
}
