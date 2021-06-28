package main

import (
	log "github.com/golang/glog"
	"runtime"
	"flag"
	"fmt"
	"net"
)

/* 存储命令行参数 */
var (
	port int
	ip string
	multiple int
	echo int
)

/* 定义命令行参数 */
func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "server ip")	
	flag.IntVar(&port, "port", 4444, "server port")	
	flag.IntVar(&multiple, "multiple", 1, "DOP = multiple * num_of_cpu")	
	flag.IntVar(&echo, "echo", 1, "echo mode")	
}

/* 主函数 */
func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("tcp server run...")

	// 调整并发度 
	envFlag := runtime.GOMAXPROCS(runtime.NumCPU())
	if envFlag > -1 {
		log.Info("GOMAXPROCS = ", runtime.NumCPU())
	} else {
		log.Info("GOMAXPROCS is default!")
	}

	// 启动监听
	log.Infof("listen on: %v:%v", ip, port)
	listener,err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Error("listen failed, abort!")
		panic(err)
	}

	cid := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Infof("Accept failed, err: %v", err)
			panic(err)
		}
		log.Infof("get connection: %s -> %s...", conn.RemoteAddr(), conn.LocalAddr())
		go routine_recvmsg(cid, conn)
		cid+=1
	}
}

func routine_recvmsg(cid int, conn net.Conn) {
	// net库默认是开启nodelay=true的
	log.Info("goroutine run ...")
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		nbytes, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				log.Infof("EOF, routine exist!")
				return
			}
			log.Infof("error reading: %v", err)
			return //终止程序
		}
		if echo == 0 {
			log.Infof("#%v: %s", cid, buf)
		} else {
			log.Infof("#%v: %s", cid, buf)
			conn.Write(buf[:nbytes])
		}
	}
}
