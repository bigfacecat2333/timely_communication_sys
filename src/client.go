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
	flag       int // 当前客户端的模式
}

func NewClient(ServerIp string, ServerPort int) *Client {
	// 创建对象

	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		flag:       99,
	}

	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
	if err != nil {
		fmt.Println("net.Dial err", err)
		return nil
	}
	client.conn = conn
	return client
}

// 处理server的response消息，这里直接输出到stdout即可
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
	// 相当于 永久阻塞等待，类似于dup2
	// for {
	// 	buf := make(...)
	// 	client.conn.Read(buf)
	//  fmt.Println(buf)
	// }
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>不在菜单的合法范围>>>>>>")
		return false
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名],exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>输入消息内容,exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>输入聊天内容,exit退出")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象[用户名],exit退出")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) PublicChat() {
	// 提示用户输入信息
	var chatMsg string

	fmt.Println(">>>>>输入聊天内容,exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 发给服务器

		// 消息不为空则发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)
	}

}

func (client *Client) UpdateName() bool {

	fmt.Println(">>>>>请输入用户名:>>>>>")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 { // 退出
		for client.menu() != true {
		}
		// 根据flag处理业务
		switch client.flag {
		case 1:
			// 公聊
			fmt.Println("公聊模式选择。。。")
			client.PublicChat()
			break
		case 2:
			// 私聊 1 查询哪些用户在线，提示用户选择一个用户进入私聊
			client.PrivateChat()
			fmt.Println("私聊模式选择。。。")
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名选择。。。")
			client.UpdateName()
			break
		}

	}
}

var ServerIp string
var ServerPort int

// init函数会在main函数之前执行
func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "设置服务器ip地址")
	flag.IntVar(&ServerPort, "port", 8888, "设置服务器端口号")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(ServerIp, ServerPort)
	if client == nil {
		fmt.Println(">>>>>>连接服务器失败>>>>>>")
		return
	}
	// 单独开辟协程去处理server的回复消息
	go client.DealResponse()

	fmt.Println(">>>>>>连接服务器成功>>>>>>")

	// 启动客户端业务 相当于select阻塞
	client.Run()
}
