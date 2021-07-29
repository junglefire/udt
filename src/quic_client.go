package main

import (
	quic "github.com/lucas-clemente/quic-go"
	log "github.com/golang/glog"
	"crypto/tls"
	// "strings"
	"context"
	"flag"
	"time"
	"fmt"
	// "net"
	"io"
)

/* 存储命令行参数 */
var (
	port int
	ip string
	loop int
	dop int
	number int
)

const message = "PING"

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
	log.Info("quic client run...")

	// 启动多个coroutine接收数据
	ch_alarm := make(chan struct{})

	for l := 0; l < loop; l++ {
		log.Infof("LOOP#%v...", l)
		for i := 0; i < dop; i++ {
			go sendmsg(i, ch_alarm)
		}
	}

	for i := 0; i < dop; i++ {
		msg := <-ch_alarm
		log.Infof("get message#%v from `ch_alarm`: %v", i, msg)
	}
}

func sendmsg(cid int, ch_alarm chan struct{}) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	session, err := quic.DialAddr(fmt.Sprintf("%s:%d", ip, port), tlsConf, nil)
	if err != nil {
		log.Errorf("connect quic server failed")
		return
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Errorf("open sync stream failed")
		return 
	}

	buffer := make([]byte, 1024)
	for i := 0; i < number; i++ {
		t1 := time.Now().UnixNano() / 1e6

		_, err = stream.Write([]byte(message))
		if err != nil {
			log.Errorf("send message failed, abort!")
			break 
		}

		log.Infof("send message to server ok")
		
		_, err = io.Reader.Read(stream, buffer)
		if err != nil {
			return 
		}

		t2 := time.Now().UnixNano() / 1e6
		log.Infof("[interval]: %v(ms)", t2-t1)
	}

	log.Infof("R#%v exist...", cid)
	ch_alarm <- struct{}{}
}
