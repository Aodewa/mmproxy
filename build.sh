export CGO_ENABLED=0
export GOOS=linux 
export GOARCH=amd64 
go build -o mmproxy main.go tcp.go utils.go udp.go buffers.go proxyprotocol.go config.go
tar zcvf mmproxy.tar.gz mmproxy
