package main

import (
	"fmt"
	httpserver "http-server/http_server"
	"log"
	"net"
)

func main() {
	network, err := net.Listen("tcp4", ":8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer network.Close()
	for {
		connect, err := network.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handler(connect)
	}
}
func handler(conn net.Conn) {
	defer conn.Close()
	http, err := httpserver.Parse(conn)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(http.GetMethod(), http.GetVersion() == httpserver.V1_0, http.GetPath())
}
