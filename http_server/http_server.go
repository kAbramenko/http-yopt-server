package httpserver

import (
	"net"
	"strings"
	"time"
)

// Version ...
type Version int

// ...
const (
	V1_0 Version = iota
	V1_1
	V2_0
	V_X
)

// HTTP ...
type HTTP struct {
	method  string
	path    string
	ver     Version
	headers map[string]string
}

// Parse ...
// GET / HTTP/1.1\r\n
// Host: localhost:8080\r\n
// User-Agent: curl/7.55.1\r\n
// Accept: */*\r\n
// \r\n
// HTTP{method: "GET", path: "/", ver: V1_1, headers: {"host": "localhost:8080", "user_agent": "curl/7.55.1", "accept": "*/*"}}
func Parse(conn net.Conn) (*HTTP, error) {
	buffer := make([]byte, 1)
	method, err := readUntilDelim(conn, buffer, ' ')
	if err != nil {
		return nil, err
	}
	method = strings.ToUpper(method)
	path, err := readUntilDelim(conn, buffer, ' ')
	if err != nil {
		return nil, err
	}
	verstr, err := readUntilDelim(conn, buffer, '\n')
	if err != nil {
		return nil, err
	}
	ver, has := map[string]Version{"HTTP/1.0": V1_0, "HTTP/1.1": V1_1, "HTTP/2.0": V2_0}[verstr]
	if !has {
		ver = V_X
	}
	headers := map[string]string{}
	for {
		headerKey, err := readUntilDelim(conn, buffer, ':')
		if err != nil {
			return nil, err
		}
		if headerKey == "" {
			break
		}
		headerValue, err := readUntilDelim(conn, buffer, '\n')
		if err != nil {
			return nil, err
		}
		headers[headerKey] = headerValue
	}
	return &HTTP{method: method, path: path, ver: ver, headers: headers}, nil
}

// GetHeader ...
func (http *HTTP) GetHeader(key string) (val string, has bool) {
	val, has = http.headers[key]
	return val, has
}

// GetMethod ...
func (http *HTTP) GetMethod() string {
	return http.method
}

// GetPath ...
func (http *HTTP) GetPath() string {
	return http.path
}

// GetVersion ...
func (http *HTTP) GetVersion() Version {
	return http.ver
}

// Parse ...
// GET / HTTP/1.1\r\n
// Host: localhost:8080\r\n
// User-Agent: curl/7.55.1\r\n
// Accept: */*\r\n
// \r\n

// GET / HTTP/1.1\r\nHost: localhost:8080\r\nUser-Agent: curl/7.55.1\r\nAccept: */*\r\n\r\n
// POST / HTTP/1.1\r\nHost: localhost:8080\r\nUser-Agent: curl/7.55.1\r\nAccept: */*\r\n\r\nuser-id: 15

func readUntilDelim(conn net.Conn, buffer []byte, delim byte) (string, error) {
	var value string
	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		_, err := conn.Read(buffer)
		if err != nil {
			return "", err
		}
		if buffer[0] == delim {
			break
		}
		if buffer[0] != '\r' {
			value += string(buffer[0])
		}
		if buffer[0] == '\n' && delim == ':' {
			return "", nil
		}
	}
	return value, nil
}
