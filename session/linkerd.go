package main

import (
	simplejson "github.com/bitly/go-simplejson"
	log "github.com/golang/glog"
	"os/signal"
	"io/ioutil"
	"runtime"
	"syscall"
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
	log.Info("kcp server run...")
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
	// 发送`quit`信号到信令通道和数据通道
	sc_quit_chan := make(chan interface{}, 1)
	// 读取信令通道配置
	protocol := js.Get("signal.channel").Get("protocol").MustString()
	ip := js.Get("signal.channel").Get("ip").MustString()
	port := js.Get("signal.channel").Get("port").MustInt()
	// 创建信令通道实例
	sc := NewSignalChannel(protocol, ip, port)
	go sc.Start(sc_quit_chan)
	// 读取数据通道配置并创建Ingress/Egress(流量入口出口)对象
	data_channel_name := js.Get("data.channel").Get("name").MustString()
	ingress := &Ingress{
		Protocol: js.Get("data.channel").Get("ingress").Get("protocol").MustString(),
		IP: js.Get("data.channel").Get("ingress").Get("ip").MustString(),
		Port: js.Get("data.channel").Get("ingress").Get("port").MustInt(),
	}
	egress := &Egress{
		Protocol: js.Get("data.channel").Get("egress").Get("protocol").MustString(),
		IP: js.Get("data.channel").Get("egress").Get("ip").MustString(),
		Port: js.Get("data.channel").Get("egress").Get("port").MustInt(),
	}
	dc := NewDataChannel(*ingress, *egress, data_channel_name)
	go dc.Start()
	// 等待用户终止信号
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
		// 关闭信令通道
		log.Info("shutdown signal channel...")
		sc_quit_chan<-0
		sc.Close()
		<-sc_quit_chan
	}
}





