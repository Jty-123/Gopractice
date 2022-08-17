package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
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
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message

		//
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.msgc <- msg
		}
		s.mapLock.Unlock()
	}
}

//broadcast
func (s *Server) Broadcast(user *User, msg string) {
	sendmsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendmsg
}

//show all user
func (s *Server) show(user *User) {
	msg := "Now online:\n"
	for _, usr := range s.OnlineMap {
		msg += "[" + usr.Addr + "]" + usr.Name + "\n"
	}
	user.msgc <- msg
}

//handler
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("connect success")

	user := NewUser(conn)

	//user online
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//broadcast online msg
	s.Broadcast(user, "online now")

	//is user alive
	isLive := make(chan bool)
	//recevive client message
	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if n == 0 {
				s.Broadcast(user, "offline now")
				s.mapLock.Lock()
				delete(s.OnlineMap, user.Name)
				s.mapLock.Unlock()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err", err)
			}
			//delete '\n'
			msg := string(buff[:n-1])
			//broadcast message
			if msg == "show" {
				//show all online user
				s.show(user)
			} else if len(msg) > 7 && msg[:7] == "rename|" {
				//rename user name
				//rename|[newName]
				name := strings.Split(msg, "|")[1]
				_, ok := s.OnlineMap[name]
				if ok {
					user.msgc <- "already use"
				} else {
					s.mapLock.Lock()
					delete(s.OnlineMap, user.Name)
					s.OnlineMap[name] = user
					s.mapLock.Unlock()
					user.Name = name
				}
			} else if len(msg) > 4 && msg[:3] == "to|" {
				//private chat
				// to|[name]|[msg]
				name := strings.Split(msg, "|")[1]
				if name == "" {
					user.msgc <- "wrong commend"
					return
				}
				remoteUser, ok := s.OnlineMap[name]
				if !ok {
					user.msgc <- "invalid user name"
					return
				}

				content := strings.Split(msg, "|")[2]
				remoteUser.msgc <- "user " + user.Name + " say to you:" + content
			} else {
				s.Broadcast(user, msg)
			}
			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:
			//now user is active
		case <-time.After(time.Second * 300):
			//5min time off
			user.msgc <- "connect close"
			time.Sleep(1)
			close(user.msgc)
			conn.Close()
		}
	}
}

// start server interface
func (s *Server) start() {
	//socket listen
	listener, err := net.Listen("tcp", s.Ip+":"+strconv.Itoa(s.Port))
	if err != nil {
		fmt.Println("net.Listen err", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//listen user online
	go s.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}
}
