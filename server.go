package main

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
	name string
}

var clients []Client

func main() {
	fmt.Println("聊天服务器正在启动...")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("启动服务器时发生错误:", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("接受连接时发生错误:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("新客户端已连接:", conn.RemoteAddr())

	client := Client{conn: conn}

	// 询问客户端名称
	fmt.Fprintf(conn, "请输入您的名称: ")
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("读取客户端名称时发生错误:", err)
		return
	}
	client.name = string(buf[:n])
	fmt.Printf("%s 加入了聊天\n", client.name)

	clients = append(clients, client)

	// 广播消息
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("%s 离开了聊天\n", client.name)
			break
		}
		msg := string(buf[:n])
		fmt.Printf("%s: %s\n", client.name, msg)
		broadcast(client, msg)
	}

	// 从切片中删除客户端
	for i, c := range clients {
		if c.conn == client.conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func broadcast(sender Client, msg string) {
	for _, client := range clients {
		if client.conn != sender.conn {
			fmt.Fprintf(client.conn, "%s\n", msg) // 添加换行符
		}
	}
}
