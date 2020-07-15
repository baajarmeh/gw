.PHONY: dylib

dylib:
	@go build -o build/dylib/stor.so -buildmode=plugin dylib/stor/stor.go
	@go build -o build/dylib/uap.so -buildmode=plugin dylib/uap/uap.go
	@go build -o build/dylib/confsvr.so -buildmode=plugin dylib/confsvr/confsvr.go
