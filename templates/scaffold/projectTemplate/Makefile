.PHONY: dylib

branch="`git rev-parse --abbrev-ref HEAD`"
commitInfo="`git log HEAD -1 --format=\"%h%d, Build AT:[%ai]\"`"
version="$(commitInfo)"
commonLdFlags="-s -w -X 'github.com/oceanho/gw.Version=$(version)'"

app:
	@go build -ldflags $(commonLdFlags) -o build/app main.go
	@chmod +x build/app && echo "Build Success."
