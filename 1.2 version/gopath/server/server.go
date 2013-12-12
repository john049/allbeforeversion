package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

//批量建立变量
var (
	host   = "127.0.0.1"       //默认IP地址
	port   = "8080"            //默认端口号
	remote = host + ":" + port //远程地址
	fName  = "config.txt"      //IP地址配置文件，格式172.0.0.1：8080
)

//建立结构体
type ClientInfo struct {
	NiceName string
	Conn     net.Conn
}

//建立map
var clientDB map[net.Conn]ClientInfo

//定义全局变量
var clientnum = 0 //初始化用户数量

func main() {

	readFile(fName)

	fmt.Println("Server address : ", remote)
	ColorEdit(12)
	fmt.Println("Initiating server... (Ctrl-C to stop)")

	clientDB = make(map[net.Conn]ClientInfo) //初始化map

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

//接收信息函数
func receiveMsg(con net.Conn) {
	//获取所接收客户端的信息
	client, ok := clientDB[con]

	if !ok {
		fmt.Println("Did not find client.")
	}

	var (
		data = make([]byte, 1024) //接收到的数据
		res  string               //接收到的转换后信息
	)
	ColorEdit(11)
	fmt.Println("New connection: ", con.RemoteAddr(), "(", client.NiceName, ")")
	fmt.Printf("%d clients online!\n", clientnum)
	ColorEdit(10)
	for {
		length, err := con.Read(data) //读取客户端发送的信息，返回信息长度和错误信息

		//如果报错，则用户已经掉线或退出，向所有用户发送此通告
		if err != nil {
			ColorEdit(11)
			fmt.Printf("Client %v quit.\n", con.RemoteAddr())
			clientnum--
			fmt.Printf("%d clients online!\n", clientnum)
			delete(clientDB, con)
			sendMsg(clientDB, fmt.Sprintf("Client %s(%v) quit.\nnow %d clients online", client.NiceName, con.RemoteAddr(), clientnum), con)
			con.Close()
			ColorEdit(10)
			return
		}
		res = string(data[0:length]) //类型转换
		fmt.Printf("%s said: %s\n", client.NiceName, res)

		//向所有用户发送此用户所写信息
		sendMsg(clientDB, fmt.Sprintf("%s said: %s", client.NiceName, res), con)
	}
}

//群发信息函数
func sendMsg(cdb map[net.Conn]ClientInfo, str string, con net.Conn) {
	for _, v := range cdb { //遍历map
		if v.Conn == con { //如果遍历所得用户端为当前发送信息用户端则不对自己发送信息
			continue
		}
		if v.Conn != nil { //如果不为空
			v.Conn.Write([]byte(str)) //发送信息
		}
	}
}

//控制台彩色输出，1深蓝，2深绿，3深青，10绿色，11青色，12红色
func ColorEdit(i int) {
	kernel32, _ := syscall.LoadDLL("kernel32.dll")
	defer kernel32.Release()
	proc, _ := kernel32.FindProc("SetConsoleTextAttribute")
	proc.Call(uintptr(syscall.Stdout), uintptr(i)) //12 Red light
}

//读取配置文件
func readFile(filename string) {
	fin, err := os.Open(filename)
	defer fin.Close()
	if err != nil {
		ColorEdit(12)
		fmt.Println(filename, err)
		ColorEdit(11)
		return
	}
	buf := make([]byte, 1024)
	n, _ := fin.Read(buf)
	if 0 == n {
		ColorEdit(12)
		fmt.Println("The IPCONFIG file is null.")
		ColorEdit(11)
		return
	}
	remote = string(buf[:n])
}
