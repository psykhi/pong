all:
	GOOS=js GOARCH=wasm go build -o server/web/main.wasm