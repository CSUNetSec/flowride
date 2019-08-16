flcap: cmd/flcap/flcap.go *.go
	go build -tags mlx -o flcap ./cmd/flcap/flcap.go

all: flcap
