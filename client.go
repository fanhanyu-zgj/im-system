package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前client模式
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{ServerIp: serverIp, ServerPort: serverPort, flag: 999}
	//连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dial error:", err)
		return nil
	}

	client.conn = conn

	//返回对象
	return client
}

func (cli *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("4.查询在线用户")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		cli.flag = flag
		return true
	} else {
		fmt.Println("》》》》》请输入合法范围内的数字")
		return false
	}

}

func (client *Client) PublicChat() {
	//提示用户输入消息
	var chatMsg string
	fmt.Println("请输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		//发送给服务器

		//消息不为空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入聊天内容1，exit退出")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println("》》》》》请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return false
	}
	return true
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
	}
}

func (client *Client) PrivateChat() {

	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>>请输入聊天对象[用户名],exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>请输入消息内容,exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>请输入消息内容,exit退出:")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>>>请输入聊天对象[用户名],exit退出:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) DealResponse() {
	//一旦有消息就copy 到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}
func (cli *Client) Run() {
	for cli.flag != 0 {
		for cli.menu() != true {

		}
		//根据不同模式处理不同业务
		switch cli.flag {
		case 1:
			cli.PublicChat()
			break
		case 2:
			cli.PrivateChat()
			break
		case 3:
			cli.UpdateName()
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认8888）")
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("》》》》》》连接服务器失败")
		return
	}
	//单独开启一个goroutine去处理server的回执信息
	go client.DealResponse()
	//启动客户端的业务
	client.Run()

}
