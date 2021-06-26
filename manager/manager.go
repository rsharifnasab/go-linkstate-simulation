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
	routersCount        int
	routers             []*Router
	netConns            [][]Edge
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
		netConns:            make([][]Edge, routersCount),
		readyWG:             sync.WaitGroup{},
		networkReadyWG:      sync.WaitGroup{},
		readyChannel:        make(chan struct{}),
		networkReadyChannel: make(chan struct{}),
		routers:             make([]*Router, routersCount),
	}
	for i := 0; i < manager.routersCount; i++ {
		manager.netConns[i] = make([]Edge, 0)
		manager.routers[i] = &Router{Index: i}
	}

	var configEdges []ConfigEdge
	pnc(config.UnmarshalKey("links", &configEdges))

	for _, configEdge := range configEdges {
		manager.netConns[configEdge.Node1] =
			append(manager.netConns[configEdge.Node1], Edge{Dest: configEdge.Node2, Cost: configEdge.Cost})
		manager.netConns[configEdge.Node2] =
			append(manager.netConns[configEdge.Node2], Edge{Dest: configEdge.Node1, Cost: configEdge.Cost})
	}
	return manager
}

type Router struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	Index  int
	Port   int
}

// change name to set connection
func (router *Router) setConnection(conn net.Conn) {
	router.conn = conn
	router.reader = bufio.NewReader(conn)
	router.writer = bufio.NewWriter(conn)
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
	buf := make([]byte, 0)
	buf = append(buf, marshalled...)
	buf = append(buf, '\n')
	_, err = router.writer.Write(buf)
	pnc(err)
	router.writer.Flush()
}

func (manager *Manager) handleRouter(routerIndex int, conn net.Conn) {
	router := manager.routers[routerIndex]
	router.Port = router.readInt()
	log.Printf("router #%v connected, udp port: %v\n", router.Index, router.Port)
	router.writeAsString(router.Index)
	// send connectivity table
	router.writeAsString(manager.routersCount)
	router.writeAsBytes(manager.netConns[router.Index])
	manager.getReadySignalFromRouter(router)
	<-manager.readyChannel
	router.writeAsString("safe")
	manager.sendPortMap(router)
	manager.getAcksReceivedFromRouter(router)
	<-manager.networkReadyChannel
	router.writeAsString("NETWORK_READY")
}

func (manager *Manager) getReadySignalFromRouter(router *Router) {
	readiness := router.readString()
	if readiness == "READY" {
		manager.readyWG.Done()
	} else {
		panic("Router couldn't get ready.")
	}
}

func (manager *Manager) sendPortMap(router *Router) {
	portMap := make(map[int]int)
	for _, edge := range manager.netConns[router.Index] {
		portMap[edge.Dest] = manager.routers[edge.Dest].Port
	}
	// marshalledPortMap, err := json.Marshal(portMap)
	// pnc(err)
	// fmt.Printf("portMap is %v\n. it was encoded into: %v\n", portMap, string(marshalledPortMap))
	router.writeAsBytes(portMap)
}

func (manager *Manager) getAcksReceivedFromRouter(router *Router) {
	str := router.readString()
	if str != "ACKS_RECEIVED" {
		panic(fmt.Sprintf("router #%v didn't receive acks: %v", router.Index, str))
	}
	manager.networkReadyWG.Done()
}
