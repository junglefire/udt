all: perf

# tcp/udp/kcp/quic/sctp performance test
perf: udp_server udp_client kcp_server kcp_client tcp_server tcp_client quic_server quic_client

udp_server: src/udp_server.go 
	go build -o bin/udp_server src/udp_server.go

udp_client: src/udp_client.go 
	go build -o bin/udp_client src/udp_client.go

kcp_server: src/kcp_server.go 
	go build -o bin/kcp_server src/kcp_server.go

kcp_client: src/kcp_client.go 
	go build -o bin/kcp_client src/kcp_client.go

tcp_server: src/tcp_server.go 
	go build -o bin/tcp_server src/tcp_server.go

tcp_client: src/tcp_client.go 
	go build -o bin/tcp_client src/tcp_client.go

quic_server: src/quic_server.go 
	go build -o bin/quic_server src/quic_server.go

quic_client: src/quic_client.go 
	go build -o bin/quic_client src/quic_client.go

# test suite
test: test/test_linker.go test/test_evio.go
	go build -o bin/test_linkerd test/test_linker.go
	go build -o bin/test_evio test/test_evio.go

run:
ifdef app
	${app} --logtostderr
else
	@echo 'no app around'
endif

clean:
	rm bin/*
	rm data/*
