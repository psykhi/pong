An online, multiplayer version of the game of pong, entirely written in Go and compiled to webassembly.

Running live at https://pong-wasm.web.app/. You can open another tab if no one wants to play with you ;(


# Getting started

- Compile the frontend code: `make`
- Start a fileserver to serve that code: `cd fileserver && go run server.go`
- Start the game server : `cd server/cmd && go run main.go`
- Connect to `localhost:3000`


# Architecture

Client and server run the game engine, running at 128 ticks. Client sends its keyboard/touch input at a fixed frequency to the server,
and both predict the next state of the game. When the client receives a server update, it reconciles both predictions by "replaying" the events that have happened since the last server packet was received.