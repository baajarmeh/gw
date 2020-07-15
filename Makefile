.PHONY: dylib

dylib:
	@go build -ldflags "-s -w" -o build/dylib/stor.so -buildmode=plugin dylib/stor/stor.go
	@go build -ldflags "-s -w" -o build/dylib/uap.so -buildmode=plugin dylib/uap/uap.go
	@go build -ldflags "-s -w" -o build/dylib/confsvr.so -buildmode=plugin dylib/confsvr/confsvr.go
