VERSION = 0.1.0
GO_BUILD_FLAG = -trimpath -ldflags "-w -s"

tcpdumpc:
	GO111MODULE=on CGO_ENABLED=0 go build $(GO_BUILD_FLAG) -o ./_output/tcpdumpc main.go

docker-tcpdumpc:
	docker build -t xieyanke/tcpdumpc:$(VERSION)-alpine .

docker-latest: docker-tcpdumpc
	docker tag xieyanke/tcpdumpc:$(VERSION)-alpine xieyanke/tcpdumpc:latest

clean:
	rm -rf ./_output

.PHONY: clean tcpdumpc docker-tcpdumpc docker-latest