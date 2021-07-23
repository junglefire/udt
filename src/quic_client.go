package main

import (
	quic "github.com/lucas-clemente/quic-go"
	log "github.com/golang/glog"
	"crypto/tls"
	// "strings"
	"context"
	"flag"
	// "time"
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

const message = "foobar"

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
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	session, err := quic.DialAddr(fmt.Sprintf("%s:%d", ip, port), tlsConf, nil)
	if err != nil {
		return
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return 
	}

	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		return 
	}

	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return 
	}
	fmt.Printf("Client: Got '%s'\n", buf)
	log.Infof("R#%v exist...", cid)
	ch_alarm <- struct{}{}
}
