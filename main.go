package main

import (
	"fmt"
	httpserver "http-server/http_server"
	"log"
	"net"
	"os"
	"strconv"
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
	request, response, err := httpserver.Parse(conn)
	if err != nil {
		log.Println(err)
		return
	}
	response.Write([]byte("HTTP/1.1 100 Continue\r\n\r\n"))
	log.Println(request.GetMethod(), request.GetVersion() == httpserver.V1_0, request.GetPath())
	if value, has := request.GetHeader("content_length"); has == true {
		if len, err := strconv.Atoi(value); err == nil {
			var d []byte = make([]byte, len)
			request.Read(d)
			fmt.Println(d)
			f, _ := os.Create("tmp")
			f.Write(d)
			defer f.Close()
		} else {
			log.Println(err)
		}
	}
	response.WriteCodeDescription(200, "Xyec")
	response.AddHeader("host", "localhost")
	response.WriteHeaders()
}
