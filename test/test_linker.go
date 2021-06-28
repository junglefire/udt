package main

import (
	simplejson "github.com/bitly/go-simplejson"
	msgpack "github.com/vmihailenco/msgpack/v5"
	log "github.com/golang/glog"
	"io/ioutil"
	"flag"
	"fmt"
	"net"
	"os"
)

type Header struct {
	Command 	string
	User		string
	AccessKey	string
}

// 命令行参数存储
var (
	config_file string
)

// 命令行参数定义
func init() {
	flag.StringVar(&config_file, "config", "etc/default.json", "linkerd config file name")	
}

func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("link run...")
	// 读取配置文件
	buf, err := ioutil.ReadFile(config_file)
	if err != nil {
		log.Errorf("read config file failed, err: %v\n", err)
		os.Exit(-1)
	}
	// 解析JSON文件
	js, err := simplejson.NewJson(buf)
	if err != nil {
		log.Errorf("parse json file failed, err: %v\n", err)
		os.Exit(-1)
	}
	// 测试
	test(js)
}

// 信令通道交互
func test(js *simplejson.Json) string {
	// 连接服务器
	protocol := js.Get("signal.channel").Get("protocol").MustString()
	ip := js.Get("signal.channel").Get("ip").MustString()
	port := js.Get("signal.channel").Get("port").MustInt()
	conn, err := net.Dial(protocol, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Infof("Dial() failed, err: %v", err)
		return ""
	}
	defer conn.Close()
	// 缓存区
	header, err := msgpack.Marshal(&Header{Command: "login", User: "alex", AccessKey: "SESS-X9AB-CDEX-HDEF"})
    if err != nil {
        panic(err)
    }
    log.Infof("%v", header)
    conn.Write(header)
	/*
	data := make([]byte, 1024)
	name, ak := "alex", "SESS-X9AB-CDEX-HDEF"
	// 1. 发送登陆请求
	conn.Write([]byte(fmt.Sprintf(LOGIN, name, ak)))
	len, err := conn.Read(data)
	if err != nil {
		log.Infof("read `login` response message failed, err: %v", err)
		return ""
	}
	log.Infof("recv %d bytes, msg: %s", len, data[:len])
	// 2. 发送建连请求
	repjs, err := simplejson.NewJson(data[:len])
	if err != nil {
		log.Errorf("parse json file failed, err: %v\n", err)
		os.Exit(-1)
	}
	conn.Write([]byte(fmt.Sprintf(CONNECT, name, ak, repjs.Get("sid").MustString())))
	len, err = conn.Read(data)
	if err != nil {
		log.Infof("read `login` response message failed, err: %v", err)
		return ""
	}
	log.Infof("recv %d bytes, msg: %s", len, data[:len])
	// 3. 发送注销请求
	conn.Write([]byte(fmt.Sprintf(LOGOUT, name)))
	len, err = conn.Read(data)
	if err != nil {
		log.Infof("read `logout` response message failed, err: %v", err)
		return ""
	}
	log.Infof("recv %d bytes, msg: %s", len, data[:len])
	*/
	return ""
}



