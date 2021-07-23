package main

import (
	kcp "github.com/xtaci/kcp-go/v5"
	log "github.com/golang/glog"
	// "reflect"
	"strings"
	"strconv"
	"flag"
	"time"
)

/* 存储命令行参数 */
var (
	port int
	ip string
	loop int
	dop int
	number int
)

/* 定义命令行参数 */
func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "Server ip")	
	flag.IntVar(&port, "port", 4444, "Server port")	
	flag.IntVar(&loop, "loop", 10, "Test Loop")	
	flag.IntVar(&dop, "dop", 10, "Degree of Parallelism")	
	flag.IntVar(&number, "number", 100000, "Number of messages per coroutine")	
}

/* 主函数 */
func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("kcp client run...")

	for l := 0; l < loop; l++ {
		log.Infof("LOOP#%v...", l)
		// 启动多个coroutine接收数据
		ch_alarm := make(chan struct{})
		for i := 0; i < dop; i++ {
			go sendmsg(i, ch_alarm)
		}

		for i := 0; i < dop; i++ {
			msg := <-ch_alarm
			log.Infof("get message#%v from `ch_alarm`: %v", i, msg)
		}
	}
}

func sendmsg(cid int, ch_alarm chan struct{}) {
	X := strings.Repeat("@$", 100)
	log.Infof("R#%v send message: %s", cid, X)
	kcpconn, err := kcp.DialWithOptions(ip+":"+strconv.Itoa(port), nil, 0, 0)
	if err!=nil {
		log.Infof("connect `%v:%v` failed, err: %v", ip, port, err)
		panic(err)
	}
	defer kcpconn.Close()
	
	log.Infof("connect_id: %v", kcpconn.GetConv())

	kcpconn.SetNoDelay(1, 10, 2, 1)

	buffer := make([]byte, 1024)
	for i := 0; i < number; i++ {
		t1 := time.Now().UnixNano() / 1e6
		kcpconn.Write([]byte("PING"))
		_, err := kcpconn.Read(buffer)
		// log.Infof("R#%v receive message#%v from server: %v", cid, i, buffer[:n])
		if err != nil {
			log.Infof("Read failed, err: %v", err)
			panic(err)
		}
		t2 := time.Now().UnixNano() / 1e6
		log.Infof("[interval]: %v(ms)", t2-t1)
	}

	log.Infof("R#%v exist...", cid)
	ch_alarm <- struct{}{}
}



