package main

import (
	log "github.com/golang/glog"
	"flag"
	"time"
	"net"
)

/* 存储命令行参数 */
var (
	port int
	ip string
	dop int
	number int
)

/* 定义命令行参数 */
func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "Server ip")	
	flag.IntVar(&port, "port", 4444, "Server port")	
	flag.IntVar(&dop, "dop", 100, "Degree of Parallelism")	
	flag.IntVar(&number, "number", 1000, "Number of messages per coroutine")	
}

/* 主函数 */
func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("udp client run...")

	addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip) }

	// 启动多个coroutine接收数据
	ch_alarm := make(chan struct{})
	for i := 0; i < dop; i++ {
		go sendmsg(i, &addr, ch_alarm)
	}

	for i := 0; i < dop; i++ {
		msg := <-ch_alarm
		log.Infof("get message#%v from `ch_alarm`: %v", i, msg)
	}
}

func sendmsg(cid int, destAddr *net.UDPAddr, ch_alarm chan struct{}) {
	log.Infof("R#%v send message...", cid)
	srcAddr := &net.UDPAddr{IP:net.IPv4zero, Port:0}
	// 连接服务器
	conn, err := net.DialUDP("udp", srcAddr, destAddr)
	if err != nil {
		log.Infof("DialUDP() failed, err: %v", err)
		ch_alarm <- struct{}{}
		return 
	}
	defer conn.Close()
	// 发送信息
	for i := 0; i < number; i++ {
		t1 := time.Now().UnixNano() / 1e6
		conn.Write([]byte("hello"))
		data := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(1*time.Second))
		n, err := conn.Read(data)
		if err != nil {
			log.Infof("read message failed, err: %v", err)
			log.Infof("[interval]: %v(ms)", -1)
			continue
		}
		t2 := time.Now().UnixNano() / 1e6
		if i % 100 == 0 {
			log.Infof("read %v message `%s` from <%s>", i, data[:n], conn.RemoteAddr())
		}
		log.Infof("[interval]: %v(ms)", t2-t1)
	}
	log.Infof("R#%v exist...", cid)
	ch_alarm <- struct{}{}
}
