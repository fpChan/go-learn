execName=generateAddr
go mod tidy
go build -race  -o ${execName} cmd/wallet/main.go
./generateAddr