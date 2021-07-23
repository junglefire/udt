package main

import (
	msgpack "github.com/vmihailenco/msgpack/v5"
	redis "github.com/go-redis/redis/v8"
	evio "github.com/tidwall/evio"
	log "github.com/golang/glog"
	"context"
	"fmt"
	// "net"
)

// 信令通道
type SignalChannel struct {
	ev evio.Events
	sdb *redis.Client
	ctx context.Context
	ip string
	port int
	url string
	dsn string
}

// 构造函数
func NewSignalChannel(events evio.Events, dsn string) *SignalChannel {
	return &SignalChannel{
		ev: events,
		dsn: dsn,
	}
}

func (sc *SignalChannel) Start(ip string, port int, parallel int) error {
	sc.ip = ip
	sc.port = port
	sc.ev.NumLoops = parallel
	// 拼装URL
	sc.url = fmt.Sprintf("tcp://%s:%d?reuseport=true", ip, port)
	// 注册事件回调函数
	sc.ev.Serving = sc.serving
	sc.ev.Data = sc.data
	// 连接Session Database
	sc.sdb = redis.NewClient(&redis.Options{Addr:sc.dsn, Password:"", DB:0,})
	sc.ctx = context.Background()
	_, err := sc.sdb.Ping(sc.ctx).Result()
	log.Infof("connect session db, dsn: %s", sc.dsn)
	if err != nil {
		log.Infof("connect session db failed, err: %v", err)
		return err
	} 
	evio.Serve(sc.ev, sc.url)
	return nil
}

// 事件处理回调函数：服务启动
func (sc SignalChannel) serving(srv evio.Server) (action evio.Action) {
	log.Infof("sessiond `(session deamon)` started on url: %s ...", sc.url)
	return
}

// 事件处理回调函数：IO事件
func (sc SignalChannel) data(c evio.Conn, in []byte) (out []byte, action evio.Action) {
	// 解析命令
	var header Header
    err := msgpack.Unmarshal(in, &header)
    if err != nil {
        panic(err)
    }
	log.Infof("command: %s", header.Command)
	// 处理命令，这里先简化处理
	switch header.Command {
	case "login":
		ret := sc.login(header.User, header.AccessKey)
		out = []byte(ret)
	case "logout":
		ret := sc.logout(header.User)
		out = []byte(ret)
	case "connect":
		ret := sc.connect(header.User)
		out = []byte(ret)
	}
	return
}

// 信令处理函数：`login`
func (sc SignalChannel) login(user string, access_token string) string {
	log.Infof("user `%s` login...", user)
	return ""
}

// 信令处理函数：`logout`
func (sc SignalChannel) logout(user string) string {
	log.Infof("user `%s` logout...", user)
	return fmt.Sprintf("{'errno': 0, 'msg':'logout ok'}")
}

// 信令处理函数：`connect`
func (sc SignalChannel) connect(user string) string {
	log.Infof("user `%s` connect...", user)
	return fmt.Sprintf("{'errno': 0, 'msg':'connect ok'}")
}









