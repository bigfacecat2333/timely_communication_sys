package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听message广播channel的进程，一有消息就发送给全部user，通知message可用
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		// 发送给所有在线user (真正的广播)
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	// 存入message的channel中，还未发送，交给监听自动完成
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// 业务
	// fmt.Println("连接成功")

	user := NewUser(conn, this)

	user.Online()

	// 监听是否活跃的channel
	isLive := make(chan bool)

	// 基本业务接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) // 活跃处理
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}

			// 去除"\n"
			msg := string(buf[:n-1])

			// 用户对消息处理(广播)
			user.DoMessage(msg)

			// 用户的任意消息代表活跃
			isLive <- true
		}
	}()

	// 当前handler阻塞，不阻塞user就没了（指针）
	for {
		select {
		case <-isLive:
			// 当前用户活跃，重置定时器, 激活了select，可以自动执行after
		case <-time.After(500 * time.Second): // 自动重置
			// 已经超时，当前user强制关闭
			user.SendMsg("用户超时未使用，被销毁")

			close(user.C)
			conn.Close()

			return // 或 runtime.Goexit()
		}
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	// listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}

	// close
	defer listener.Close()

	// 启动监听message的goroutine
	go this.ListenMessage()

	for {
		// accept, 代表用户已上线
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}
