.PHONY: setup, run, fmt, build

setup:
	go get github.com/oliamb/cutter
	go get github.com/Mushus/twtr

fmt:
	gofmt -w ./

run:
	go run ./reicon.go

build:
	GOARTCH=amd64 GOOS=windows go build -o ./build/windows-amd64/reicon.exe reicon.go
	GOARTCH=amd64 GOOS=darwin go build -o ./build/darwin-amd64/reicon reicon.go
	GOARTCH=amd64 GOOS=linux go build -o ./build/linux-amd64/reicon reicon.go
	GOARTCH=arm GOOS=freebsd GOARM=6 go build -o ./build/linux-arm6/reicon reicon.go
	GOARTCH=arm GOOS=freebsd GOARM=7 go build -o ./build/linux-arm7/reicon reicon.go
