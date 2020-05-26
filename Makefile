dev:
	GOOS=js GOARCH=wasm go build -o fileserver/web/main.wasm ./client

prod:
	GOOS=js GOARCH=wasm go build -tags prod -o fileserver/web/main.wasm ./client