# Go supported OSs and Archs: https://github.com/golang/go/blob/master/src/go/build/syslist.go

build: clear linux windows

linux: setup
	export GOOS=linux

	GOARCH=386 go build -o build/lifecycledoc_linux_386 cmd/lifecycledoc/main.go
	GOARCH=amd64 go build -o build/lifecycledoc_linux_amd64 cmd/lifecycledoc/main.go

windows: setup
	export GOOS=windows

	GOARCH=386 go build -o build/lifecycledoc_windows_386 cmd/lifecycledoc/main.go
	GOARCH=amd64 go build -o build/lifecycledoc_windows_amd64 cmd/lifecycledoc/main.go

setup:
	export CGO_ENABLED=0

clear:
	rm -rf build