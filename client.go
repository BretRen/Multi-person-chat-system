package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("欢迎来到聊天室!")
	fmt.Println("请输入您的名称:")

	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name) // 去掉名称末尾的换行符

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("连接服务器时发生错误:", err)
		return
	}
	defer conn.Close()

	fmt.Println("已连接到服务器")

	// 发送名称给服务器
	conn.Write([]byte(name + "\n"))

	// 启动协程读取来自服务器的消息
	go func() {
		for {
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("从服务器接收消息时出错:", err)
				return
			}
			fmt.Print(msg)
		}
	}()

	// 从标准输入读取消息并发送给服务器
	for {
		msg, _ := reader.ReadString('\n')
		conn.Write([]byte(name + ": " + msg))
	}
}
