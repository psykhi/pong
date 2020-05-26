// +build !prod

package main

var Conf = Config{
	Port:       ":8080",
	Address:    "localhost",
	HTTPScheme: "http://",
	WsScheme:   "ws",
}
