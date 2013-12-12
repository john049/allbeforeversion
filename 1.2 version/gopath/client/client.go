package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

var str string
var msg = make([]byte, 1024)

//批量建立变量
var (
	host   = "127.0.0.1"       //默认IP地址
	port   = "8080"            //默认端口号
	remote = host + ":" + port //远程地址
	fName  = "config.txt"      //IP地址配置文件，格式172.0.0.1：8080
)

func main() {

	readFile(fName)

	fmt.Println("You will connect the server : ", remote)

	//输入昵称
	ColorEdit(12)
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

//接收信息函数
func receiveMsg(con net.Conn) {
	for {
		length, err := con.Read(msg)
		if err != nil {
			ColorEdit(12)
			fmt.Printf("Error when read from server.\n")
			os.Exit(0)
		}
		str = string(msg[0:length]) //格式转换
		ColorEdit(11)
		fmt.Println(str) //将信息显示到屏幕
		ColorEdit(10)
	}

}

//发送信息函数
func sendMsg(con net.Conn) {
	for {
		//fmt.Printf("Enter a sentence:")
		fmt.Scanf("%s\n", &str) //屏幕输入
		//如果输入quit则结束连接
		if str == "quit" {
			ColorEdit(12)
			fmt.Println("Communication terminated.")
			os.Exit(1)
		}

		in, err := con.Write([]byte(str)) //将输入的内容发送给服务器端
		if err != nil {
			ColorEdit(12)
			fmt.Printf("Error when send to server: %d\n", in)
			os.Exit(0)
		}
	}
}

//控制台彩色输出，1深蓝，2深绿，3深青，10绿色，11青色，12红色
func ColorEdit(i int) {
	kernel32, _ := syscall.LoadDLL("kernel32.dll")          //调用dll
	defer kernel32.Release()                                //延迟释放
	proc, _ := kernel32.FindProc("SetConsoleTextAttribute") //找到dll中的方法，并返回此方法
	proc.Call(uintptr(syscall.Stdout), uintptr(i))          //执行此方法
}

//读取配置文件
func readFile(filename string) {
	fin, err := os.Open(filename) //打开文件
	defer fin.Close()             //延迟关闭
	if err != nil {               //是否报错
		ColorEdit(12)
		fmt.Println(filename, err) //输出错误信息
		return
		ColorEdit(11)
	}
	buf := make([]byte, 1024) //声明1024位byte型数组切片
	n, _ := fin.Read(buf)     //读取文件信息，并存入buf，返回信息长度
	if 0 == n {               //如果信息长度为0
		ColorEdit(12)
		fmt.Println("The IPCONFIG file is null.")
		ColorEdit(11)
		return
	}
	remote = string(buf[:n]) //将读取的信息赋值给远程地址变量
}
