package main

import (
	kcp "github.com/xtaci/kcp-go/v5"
	log "github.com/golang/glog"
	"errors"
	"net"
	"fmt"
)

/*************************************************************************************
 * 数据通道入口                                                                        *
*************************************************************************************/
type HandleMessagePtr func(net.Conn)

type Ingress struct {
	Protocol string
	IP string
	Port int
	HandleMessage HandleMessagePtr
}

// 返回Ingress的描述信息
func (ingress Ingress) String() string {
	return fmt.Sprintf("%s://%s:%d", ingress.Protocol, ingress.IP, ingress.Port)
}

// 创建监听句柄
func (ingress *Ingress) Listen() (net.Listener, error) {
	switch {
	case ingress.Protocol == "tcp":
		ingress.HandleMessage = HandleTCPMessage
		return net.Listen("tcp", fmt.Sprintf("%s:%d", ingress.IP, ingress.Port))
	case ingress.Protocol == "kcp":
		ingress.HandleMessage = HandleKCPMessage
		return kcp.Listen(fmt.Sprintf("%s:%d", ingress.IP, ingress.Port))
	}
	return nil, errors.New(fmt.Sprintf("invalid Protocol: %s", ingress.Protocol))
}

// 根据不同协议，定义不同的消息处理协程函数
func HandleKCPMessage(conn net.Conn) {
	log.Infof("kcp client `%s` connect...", conn.LocalAddr())
}

func HandleTCPMessage(conn net.Conn) {
	log.Infof("tcp client `%s` connect...", conn.LocalAddr())
}


/*************************************************************************************
 * 数据通道出口                                                                        *
*************************************************************************************/
type Egress struct {
	Protocol string
	IP string
	Port int
}

func (egress Egress) String() string {
	if egress.Protocol == "direct" {
		return "direct://"
	}
	return fmt.Sprintf("%s://%s:%d", egress.Protocol, egress.IP, egress.Port)
}


/*************************************************************************************
 * 数据通道                                                                           *
*************************************************************************************/
type DataChannel struct {
	ingress Ingress
	egress Egress
	name string
}

func NewDataChannel(ingress Ingress, egress Egress, name string) *DataChannel {
	return &DataChannel{
		ingress: ingress,
		egress: egress,
		name: name,
	}
}

// 启动数据通道
func (sc DataChannel) Start() {
	log.Infof("create data channel `%s`: [%s] <-> [%s]", sc.name, sc.ingress, sc.egress)
	listener, err := sc.ingress.Listen()
	if err != nil {
		log.Errorf("create listener failed, err: %s", err)
		return 
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Accept failed, err: %s", err)
			return
		}
		go sc.ingress.HandleMessage(conn)	
	}
}


