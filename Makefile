all: sessiond

linkerd: src/linkerd.go 
	go build -o bin/linkerd src/linkerd.go src/datachan.go src/sigchan.go

sessiond: src/sessiond.go 
	go build -o bin/sessiond src/sessiond.go src/sigchan.go src/command.go

# tcp/udp/kcp/quic/sctp performance test
perf: udp_server udp_client kcp_server kcp_client tcp_server tcp_client

udp_server: perf/udp_server.go 
	go build -o bin/udp_server perf/udp_server.go

udp_client: perf/udp_client.go 
	go build -o bin/udp_client perf/udp_client.go

kcp_server: perf/kcp_server.go 
	go build -o bin/kcp_server perf/kcp_server.go

kcp_client: perf/kcp_client.go 
	go build -o bin/kcp_client perf/kcp_client.go

tcp_server: perf/tcp_server.go 
	go build -o bin/tcp_server perf/tcp_server.go

tcp_client: perf/tcp_client.go 
	go build -o bin/tcp_client perf/tcp_client.go

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
