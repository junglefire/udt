package main

import (
	"flag"
	"fmt"
	log "github.com/golang/glog"
	"net"
	"strings"
	"time"
)

/* 存储命令行参数 */
var (
	port   int
	ip     string
	loop   int
	dop    int
	number int
)

/* 定义命令行参数 */
func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "Server ip")
	flag.IntVar(&port, "port", 4444, "Server port")
	flag.IntVar(&loop, "loop", 10, "Test Loop")
	flag.IntVar(&dop, "dop", 10, "Degree of Parallelism")
	flag.IntVar(&number, "number", 10, "Number of messages per coroutine")
}

/* 主函数 */
func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("tcp client run...")

	for l := 0; l < loop; l++ {
		// 启动多个coroutine接收数据
		ch_alarm := make(chan struct{})
		for i := 0; i < dop; i++ {
			go routine_sendmsg(i, ch_alarm)
		}

		for i := 0; i < dop; i++ {
			msg := <-ch_alarm
			log.Infof("get message#%v from `ch_alarm`: %v", i, msg)
		}
	}
}

func routine_sendmsg(cid int, ch_alarm chan struct{}) {
	X := strings.Repeat("@@@", 100)
	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Infof("Dial() failed, err: %v", err)
		ch_alarm <- struct{}{}
		return
	}
	defer conn.Close()
	// 发送信息
	for i := 0; i < number; i++ {
		t1 := time.Now().UnixNano() / 1e6
		conn.Write([]byte(X))
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		if err != nil {
			log.Infof("read message failed, err: %v", err)
			log.Infof("[interval]: %v(ms)", -1)
			continue
		}
		t2 := time.Now().UnixNano() / 1e6
		log.Infof("ID: %v, [interval]: %v(ms)", cid, t2-t1)
		// time.Sleep(time.Duration(2) * time.Second)
	}
	log.Infof("R#%v exist...", cid)
	ch_alarm <- struct{}{}
}
