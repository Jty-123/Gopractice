package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //mode
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       9,
	}
	conn, err := net.Dial("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Println("connect error:", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

//public chat
func (client *Client) PublicChat() {
	msg := ""
	fmt.Println("Please input chat content(input exit to quit):")
	fmt.Scanln(&msg)
	for msg != "exit" {
		if len(msg) != 0 {
			send := msg + "\n"
			_, err := client.conn.Write([]byte(send))
			if err != nil {
				fmt.Println("conn write err", err)
				break
			}
		}
		msg = ""
		fmt.Println("Please input chat content(input exit to quit):")
		fmt.Scanln(&msg)
	}
}

//p c
func (client *Client) PrivateChat() {
	msg := "show\n"
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write err", err)
		return
	}
	name := ""
	fmt.Println("Please input private chat name (exit to quit):")
	fmt.Scanln(&name)
	for name != "exit" {
		fmt.Println("input message (exit to quit):")
		chatmsg := ""
		fmt.Scanln(&chatmsg)
		for chatmsg != "exit" {
			if len(chatmsg) != 0 {
				send := "to|" + name + "|" + chatmsg + "\n"
				_, err := client.conn.Write([]byte(send))
				if err != nil {
					fmt.Println("conn write err", err)
					break
				}
			}
			chatmsg = ""
			fmt.Println("Please input chat content(input exit to quit):")
			fmt.Scanln(&chatmsg)
		}

	}
}

//raname
func (client *Client) rename() bool {
	fmt.Println("input your name")
	fmt.Scanln(&client.Name)
	msg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("client conn write err")
		return false
	}
	fmt.Println("successs change name to " + client.Name + "\n")
	return true
}
func (client *Client) Run() {
	for client.flag != 4 {
		for client.meun() != true {
		}
		switch client.flag {
		case 1:
			client.PublicChat()
			//fmt.Println("公聊")
			break
		case 2:
			client.PrivateChat()
			//fmt.Println("私聊")
			break
		case 3:
			//fmt.Println("重命名")
			client.rename()
			break
		}
	}
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "connect server ip default 127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "connect server port default 8888")
}

//deal server message
func (client *Client) deal() {
	io.Copy(os.Stdout, client.conn) // for { buf :=make() client.conn.read(buf) print(buf)}

}
func (client *Client) meun() bool {
	var flag int
	fmt.Println("1.public chat")
	fmt.Println("2.private chat")
	fmt.Println("3.rename")
	fmt.Println("4.quit")

	fmt.Scanln(&flag)
	if flag >= 1 && flag <= 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println("invalid input")
		return false
	}
}

func main() {
	flag.Parse()
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("connect server fail")
		return
	}

	go client.deal()

	fmt.Println("connect server success")

	client.Run()
}
