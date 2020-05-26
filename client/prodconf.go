// +build prod

package main

var Conf = Config{
	Port:       "",
	Address:    "pong-wasm.herokuapp.com",
	HTTPScheme: "https://",
	WsScheme:   "wss",
}
