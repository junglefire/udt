package main

import (
	kcp "github.com/xtaci/kcp-go/v5"
	log "github.com/golang/glog"
	// "reflect"
	"runtime"
	"strconv"
	"flag"
	// "time"
	"net"
)

/* 存储命令行参数 */
var (
	port int
	ip string
	multiple int
)

/* 定义命令行参数 */
func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "server ip")	
	flag.IntVar(&port, "port", 4444, "server port")	
}

/* 主函数 */
func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("kcp server run...")

	// 调整并发度 
	envFlag := runtime.GOMAXPROCS(runtime.NumCPU())
	if envFlag > -1 {
		log.Info("GOMAXPROCS = ", runtime.NumCPU())
	} else {
		log.Info("GOMAXPROCS is default!")
	}

	// 启动监听
	log.Infof("listen on: %v:%v", ip, port)
	listener, err := kcp.Listen(ip+":"+strconv.Itoa(port))
	if err!=nil {
		log.Infof("kcp server listen failed, err: %v", err)
		panic(err)
	}

	cid := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Infof("Accept failed, err: %v", err)
			panic(err)
		}
		go lstn(cid, conn)
		cid+=1
	}
}

func lstn(cid int, conn net.Conn) {
	buffer := make([]byte, 1024)
	err := error(nil)
	count := 0
	// conn.SetDeadline(time.Now().Add(30*time.Second))
	for err == nil {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Infof("read data failed, err: %v", err)
			break
		}
		count = count+1
		if count % 10 == 0 {
			log.Infof("R#%v get %v messge from %v: %v", cid, count, conn.RemoteAddr(), buffer[:n])
		}
		conn.Write([]byte("PONG"))
	}
	log.Infof("goroutine R#%v exit!", cid)
}



