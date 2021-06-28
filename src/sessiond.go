package main

import (
	simplejson "github.com/bitly/go-simplejson"
	evio "github.com/tidwall/evio"
	log "github.com/golang/glog"
	// "os/signal"
	"io/ioutil"
	"runtime"
	// "syscall"
	"flag"
	"os"
)

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
	log.Info("session deamon run...")
	// 调整并发度 
	envFlag := runtime.GOMAXPROCS(runtime.NumCPU())
	if envFlag > -1 {
		log.Info("GOMAXPROCS = ", runtime.NumCPU())
	} else {
		log.Info("GOMAXPROCS is default!")
	}
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
	// 启动会话监听服务
	var events evio.Events
	ip := js.Get("signal.channel").Get("ip").MustString()
	port := js.Get("signal.channel").Get("port").MustInt()
	parallel := js.Get("signal.channel").Get("parallel").MustInt()
	sc := NewSignalChannel(events, js.Get("session.database").Get("dsn").MustString())
	sc.Start(ip, port, parallel)
}





