package main

import (
	"fmt"
	"net"
	"os"
)

//批量建立变量
var (
	host   = "127.0.0.1"       //IP地址
	port   = "8080"            //端口
	remote = host + ":" + port //远程地址
)

var conns [3]net.Conn //定义一个1000位net.Conn类型的数组
var num = 0           //初始化数组位置0
var clientnum = 0     //初始化用户数量

func main() {

	fmt.Println("Initiating server... (Ctrl-C to stop)")

	lis, err := net.Listen("tcp", remote) //监听远程地址的TCP协议
	defer lis.Close()                     //延迟处理监听关闭

	//如果错误不为空，则输出错误信息并退出
	if err != nil {
		fmt.Println("Error when listen: ", remote)
		os.Exit(-1)
	}

	//循环等待接收客户端
	for {
		conn, err := lis.Accept() //等待接受连接，并返回这个连接conn

		//如果错误不为空，则输出错误信息并退出
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
			os.Exit(0)
		}

		conns[num] = conn //将此客户连接放入数组中
		num++             //数组位下移
		clientnum++       //用户数加1

		welmsg := fmt.Sprintf("welcome %s join.\nnow %d clients online", conn.RemoteAddr(), clientnum) //定义欢迎信息
		sendMsg(conns, welmsg, nil)                                                                    //向所有用户发送欢迎信息

		go receiveMsg(conn) //建立针对此用户的接收goroutine

	}
}

//接收信息函数
func receiveMsg(con net.Conn) {
	var (
		data = make([]byte, 1024) //接收到的数据
		res  string               //接收到的转换后信息
	)
	fmt.Println("New connection: ", con.RemoteAddr())
	fmt.Printf("%d clients online!\n ", clientnum)
	for {
		length, err := con.Read(data) //读取客户端信息，返回信息长度和错误信息

		//如果报错，则用户已经掉线或退出，向所有用户发送信息
		if err != nil {
			fmt.Printf("Client %v quit.\n", con.RemoteAddr())
			clientnum--
			fmt.Printf("%d clients online!\n ", clientnum)
			sendMsg(conns, fmt.Sprintf("Client %v quit.\nnow %d clients online", con.RemoteAddr(), clientnum), con)
			con.Close()
			return
		}
		res = string(data[0:length]) //类型转换
		fmt.Printf("%s said: %s\n", con.RemoteAddr(), res)

		//向所有用户发送此用户所写信息
		sendMsg(conns, fmt.Sprintf("%s said: %s", con.RemoteAddr(), res), con)
	}
}

//群发信息函数
func sendMsg(cons [3]net.Conn, str string, con net.Conn) {
	for i, v := range cons { //遍历数组
		if cons[i] == con {
			continue
		}
		if cons[i] != nil { //如果不为空
			v.Write([]byte(str)) //发送信息
		}
	}

}
