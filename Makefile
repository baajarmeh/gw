.PHONY: dylib

branch="`git rev-parse --abbrev-ref HEAD`"
commitInfo="`git log HEAD -1 --format=\"%h%d, Build AT:[%ai]\"`"
version="$(commitInfo)"
commonLdFlags="-s -w -X 'github.com/oceanho/gw.Version=$(version)'"

dylib:
	@go build -ldflags $(commonLdFlags) -o build/dylib/stor.so -buildmode=plugin dylib/stor/stor.go
	@go build -ldflags $(commonLdFlags) -o build/dylib/uap.so -buildmode=plugin dylib/uap/uap.go
	@go build -ldflags $(commonLdFlags) -o build/dylib/confsvr.so -buildmode=plugin dylib/confsvr/confsvr.go

dylibsvr:
	@go build -ldflags $(commonLdFlags) -o build/cmd/dylibsvr cmd/dylibsvr/main.go

gw-cli:
	@go build -ldflags $(commonLdFlags) -o build/cmd/gw-cli cmd/gwcli/main.go
	@chmod +x build/cmd/gw-cli
