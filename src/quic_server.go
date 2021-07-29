package main

import (
	quic "github.com/lucas-clemente/quic-go"
	log "github.com/golang/glog"
	"encoding/pem"
	"crypto/x509"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"math/big"
	"runtime"
	"context"
	"flag"
	"fmt"
	"io"
)

/* 存储命令行参数 */
var (
	port int
	ip string
	multiple int
)

/* 定义命令行参数 */
func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "server ip")	
	flag.IntVar(&port, "port", 4444, "server port")	
	flag.IntVar(&multiple, "multiple", 1, "DOP = multiple * num_of_cpu")	
}

const message = "PING"

/* 主函数 */
func main() {
	flag.Parse()
	defer log.Flush()
	log.Info("quic server run...")

	// 调整并发度 
	envFlag := runtime.GOMAXPROCS(runtime.NumCPU())
	if envFlag > -1 {
		log.Info("GOMAXPROCS = ", runtime.NumCPU())
	} else {
		log.Info("GOMAXPROCS is default!")
	}

	// 启动监听
	listener, err := quic.ListenAddr(fmt.Sprintf("%s:%d", ip, port), generateTLSConfig(), nil)
	if err != nil {
		return 
	}

	log.Infof("quic server listen on `%s:%d`", ip, port)

	// 接收请求，每个客户端启动一个协程处理
	cid := 0
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			return 
		}
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}
		go lstn(cid, stream)
		cid+=1
	}

	return 
}

func lstn(cid int, stream quic.Stream) {
	buf := make([]byte, 1024)

	for {
		_, err := io.Reader.Read(stream, buf)
		if err != nil {
			return 
		}
		fmt.Printf("Server: Got '%s'\n", buf)

		_, err = stream.Write([]byte(message))
		if err != nil {
			return 
		}
	}
	log.Infof("goroutine R#%v exit!", cid)
}


// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
