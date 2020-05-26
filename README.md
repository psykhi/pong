A browser version of the game of pong, entirely written in Go and compiled to webassembly.

Running live at https://pong-wasm.web.app/

# Getting started

- Compile the frontend code: `make`
- Start a fileserver to serve that code: `cd fileserver && go run server.go`
- Start the game server : `cd server/cmd && go run main.go`
- Connect to `localhost:3000`


