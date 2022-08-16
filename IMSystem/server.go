package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

//create a server interface
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.msgc <- msg
		}
		this.mapLock.Unlock()
	}
}

//broadcast
func (this *Server) Broadcast(user *User, msg string) {
	sendmsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendmsg
}

//handler
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("connect success")

	user := NewUser(conn)

	//user online
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//broadcast online msg
	this.Broadcast(user, "online now")

	//
	select {}

}

// start server interface
func (this *Server) start() {
	//socket listen
	listener, err := net.Listen("tcp", this.Ip+":"+strconv.Itoa(this.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	//close listen socket
	defer listener.Close()

	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}
}
