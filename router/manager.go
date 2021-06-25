package main

import (
	"bufio"
	"fmt"
	"net"
)

type Manager struct {
	conn   *net.TCPConn
	writer *bufio.Writer
	reader *bufio.Reader
}

func NewManager(add string) *Manager {
	managerConnection, err := net.Dial("tcp", "localhost:8585")
	pnc(err)
	managerReader := bufio.NewReader(managerConnection)
	managerWriter := bufio.NewWriter(managerConnection)

	return &Manager{
		conn:   managerConnection.(*net.TCPConn),
		reader: managerReader,
		writer: managerWriter,
	}
}

func (manager *Manager) write(data interface{}) {
	manager.writer.WriteString(fmt.Sprintf("%v\n", data))
	pnc(manager.writer.Flush())
}
