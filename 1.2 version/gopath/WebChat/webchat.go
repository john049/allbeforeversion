// WebChat project main.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

var (
	host   = "127.0.0.1"       //IP地址
	port   = "8080"            //端口
	remote = host + ":" + port //远程地址
)

var (
	str string               //消息输出变量
	msg = make([]byte, 1024) //C和S传递消息数据
)

//建立结构体
type ClientInfo struct {
	NiceName string   //昵称
	Conn     net.Conn //连接
}

//建立用户map
var clientDB map[net.Conn]ClientInfo

//定义全局变量
var clientnum = 0 //初始化连接用户数量

func main() {

}

//服务端
func server() {
	fmt.Println("Initiating server... (Ctrl-C to stop)")

	clientDB = make(map[net.Conn]ClientInfo) //初始化MAP

	lis, err := net.Listen("tcp", remote) //监听远程地址的TCP协议
	defer lis.Close()                     //延迟处理监听关闭

	//如果错误不为空，则输出错误信息并退出
	if err != nil {
		ColorEdit(12)
		fmt.Println("Error when listen: ", remote)
		os.Exit(-1)
	}

	//循环等待接收客户端
	for {
		conn, err := lis.Accept() //等待接受连接，并返回这个连接conn

		//如果错误不为空，则输出错误信息并退出
		if err != nil {
			ColorEdit(12)
			fmt.Println("Error accepting client: ", err.Error())
			os.Exit(0)
		}

		//接收用户输入的昵称
		var (
			ndata    = make([]byte, 1024) //接收到的数据
			nicename string               //接收到的转换后信息
		)
		length, err := conn.Read(ndata)    //读取客户端信息，返回信息长度和错误信息
		nicename = string(ndata[0:length]) //类型转换

		//将用户放入MAP中
		clientDB[conn] = ClientInfo{nicename, conn}

		clientnum++ //用户数加1

		welmsg := fmt.Sprintf("welcome %s(%s) join.\nnow %d clients online", nicename, conn.RemoteAddr(), clientnum) //定义欢迎信息
		sendMsg(clientDB, welmsg, nil)                                                                               //向所有用户发送欢迎信息

		go receiveMsg(conn) //建立针对此用户的接收goroutine

	}
}

//服务端接收信息函数
func SerReceiveMsg(con net.Conn) {
	//获取所接收客户端的信息
	client, ok := clientDB[con]
	if !ok {
		fmt.Println("Did not find client.")
	}

	var (
		data = make([]byte, 1024) //接收到的数据
		res  string               //接收到的转换后信息
	)
	fmt.Println("New connection: ", con.RemoteAddr(), "(", client.NiceName, ")")
	fmt.Printf("%d clients online!\n", clientnum)
	for {
		length, err := con.Read(data) //读取客户端信息，返回信息长度和错误信息

		//如果报错，则用户已经掉线或退出，向所有用户发送信息
		if err != nil {
			fmt.Printf("Client %v quit.\n", con.RemoteAddr())
			clientnum--
			fmt.Printf("%d clients online!\n", clientnum)
			delete(clientDB, con)
			sendMsg(clientDB, fmt.Sprintf("Client %s(%v) quit.\nnow %d clients online", client.NiceName, con.RemoteAddr(), clientnum), con)
			con.Close()
			return
		}
		res = string(data[0:length]) //类型转换
		fmt.Printf("%s said: %s\n", client.NiceName, res)

		//向所有用户发送此用户所写信息
		sendMsg(clientDB, fmt.Sprintf("%s said: %s", client.NiceName, res), con)
	}
}

//服务端群发信息给客户端函数
func SerSendMsg(cdb map[net.Conn]ClientInfo, str string, con net.Conn) {
	for _, v := range cdb { //遍历数组
		if v.Conn == con {
			continue
		}
		if v.Conn != nil { //如果不为空
			v.Conn.Write([]byte(str)) //发送信息
		}
	}
}

//客户端
func client() {
	//输入昵称
	for {
		fmt.Printf("Enter a nicename:")
		fmt.Scanf("%s\n", &str) //屏幕输入
		if str == "" {
			fmt.Printf("Please input a right nicename!\n")
			continue
		}
		break
	}

	con, err := net.Dial("tcp", remote) //与远程服务器监理连接，返回连接con
	defer con.Close()                   //延迟关闭连接

	if err != nil {
		ColorEdit(12)
		fmt.Println("Server not found.")
		os.Exit(-1)
	}

	in, err := con.Write([]byte(str)) //将昵称发送给服务器端
	if err != nil {
		ColorEdit(12)
		fmt.Printf("Error when send to server: %d\n", in)
		os.Exit(0)
	}

	ColorEdit(11)
	fmt.Println("Connection OK.")

	//接收欢迎信息
	length, err := con.Read(msg)
	if err != nil {
		ColorEdit(12)
		fmt.Printf("Error when read from server.\n")
		os.Exit(0)
	}
	str = string(msg[0:length]) //格式转换
	fmt.Println(str)            //将信息显示到屏幕

	ColorEdit(10)
	go receiveMsg(con) //接收信息goroutine
	sendMsg(con)       //发送信息，不用goroutine的原因是为了防止main结束
}

//客户端接收服务器信息函数
func CliReceiveMsg(con net.Conn) {
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

//客户端发送信息给服务器函数
func CliSendMsg(con net.Conn) {
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
