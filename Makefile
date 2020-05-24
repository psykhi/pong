all:
	GOOS=js GOARCH=wasm go build -o fileserver/web/main.wasm