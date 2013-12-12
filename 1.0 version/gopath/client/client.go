package main

import (
	"fmt"
	"net"
	"os"
)

var str string
var msg = make([]byte, 1024)

//批量建立变量
var (
	host   = "127.0.0.1"       //IP地址
	port   = "8080"            //端口号
	remote = host + ":" + port //远程地址
)

func main() {
	//输入昵称
	//fmt.Printf("Enter a nicename:")
	//fmt.Scanf("%s\n", &str) //屏幕输入

	con, err := net.Dial("tcp", remote) //与远程服务器监理连接，返回连接con
	defer con.Close()                   //延迟关闭连接

	if err != nil {
		fmt.Println("Server not found.")
		os.Exit(-1)
	}
	fmt.Println("Connection OK.")

	//in, err := con.Write([]byte(str)) //将昵称发送给服务器端
	//if err != nil {
	//	fmt.Printf("Error when send to server: %d\n", in)
	//	os.Exit(0)
	//}

	//接收欢迎信息
	length, err := con.Read(msg)
	if err != nil {
		fmt.Printf("Error when read from server.\n")
		os.Exit(0)
	}
	str = string(msg[0:length]) //格式转换
	fmt.Println(str)            //将信息显示到屏幕

	go receiveMsg(con) //接收信息goroutine
	sendMsg(con)       //发送信息，不用goroutine的原因是为了防止main结束

}

//接收信息函数
func receiveMsg(con net.Conn) {
	for {
		length, err := con.Read(msg)
		if err != nil {
			fmt.Printf("Error when read from server.\n")
			os.Exit(0)
		}
		str = string(msg[0:length]) //格式转换
		fmt.Println(str)            //将信息显示到屏幕
	}

}

//发送信息函数
func sendMsg(con net.Conn) {
	for {
		//fmt.Printf("Enter a sentence:")
		fmt.Scanf("%s\n", &str) //屏幕输入
		//如果输入quit则结束连接
		if str == "quit" {
			fmt.Println("Communication terminated.")
			os.Exit(1)
		}

		in, err := con.Write([]byte(str)) //将输入的内容发送给服务器端
		if err != nil {
			fmt.Printf("Error when send to server: %d\n", in)
			os.Exit(0)
		}
	}
}
