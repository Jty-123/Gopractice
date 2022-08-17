package main

import "net"

type User struct {
	Name string
	Addr string
	msgc chan string
	conn net.Conn
}

//Create a User
func NewUser(conn net.Conn) *User {
	userAdd := conn.RemoteAddr().String()
	user := &User{
		Name: userAdd,
		Addr: userAdd,
		msgc: make(chan string),
		conn: conn,
	}
	//start listen
	go user.ListenMessage()

	return user
}

//listen user channel
func (u *User) ListenMessage() {
	for {
		msg := <-u.msgc
		u.conn.Write([]byte(msg + "\n"))
	}
}
