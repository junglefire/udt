package main

import (
	log "github.com/golang/glog"
	"runtime"
	"flag"
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
	flag.IntVar(&multiple, "multiple", 1, "DOP = multiple * num_of_cpu")	
}

/* 主函数 */
func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("udp server run...")

	// 调整并发度 
	envFlag := runtime.GOMAXPROCS(runtime.NumCPU())
	if envFlag > -1 {
		log.Info("GOMAXPROCS = ", runtime.NumCPU())
	} else {
		log.Info("GOMAXPROCS is default!")
	}

	// 启动监听
	log.Infof("listen on: %v:%v", ip, port)
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip) }
	connection, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Error("listen failed, abort!")
		panic(err)
	}

	// 启动多个coroutine接收数据
	ch_alarm := make(chan struct{})
	log.Infof("CPU Threads: %v, DOP: %v", runtime.NumCPU(), runtime.NumCPU()*multiple)
	for i := 0; i < runtime.NumCPU()*multiple; i++ {
		go lstn(i, connection, ch_alarm)
	}

	msg := <-ch_alarm
	log.Infof("get message from `ch_alarm`: %v", msg)
}

func lstn(cid int, connection *net.UDPConn, ch_alarm chan struct{}) {
	buffer := make([]byte, 1024)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
	count := 0
	for err == nil {
		n, remoteAddr, err = connection.ReadFromUDP(buffer)
		n, err = connection.WriteToUDP([]byte("OK"), remoteAddr)
		if err != nil {
			log.Infof("write to udp endpoint failed: %v", err.Error())
		} 
		count = count+1
		if count % 100 == 0 {
			log.Infof("R#%v get %v messge from %v: %v", cid, count, remoteAddr, buffer[:n])
		}
	}
	log.Errorf("Listener failed, err: %v", err)
	ch_alarm <- struct{}{}
}
