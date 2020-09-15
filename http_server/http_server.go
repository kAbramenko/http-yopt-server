package httpserver

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// Version ...
type Version int
type state int

// ...
const (
	V1_0 Version = iota
	V1_1
	V2_0
	V_X
)
const (
	initialize state = iota
	header
	body
)

// HTTPRequest ...
type HTTPRequest struct {
	method  string
	path    string
	ver     Version
	headers map[string]string
	socket  net.Conn
}

// HTTPResponse ...
type HTTPResponse struct {
	code    int
	ver     Version
	socket  net.Conn
	headers map[string]string
	state   state
}

// Parse ...
// GET / HTTP/1.1\r\n
// Host: localhost:8080\r\n
// User-Agent: curl/7.55.1\r\n
// Accept: */*\r\n
// \r\n
// HTTP{method: "GET", path: "/", ver: V1_1, headers: {"host": "localhost:8080", "user_agent": "curl/7.55.1", "accept": "*/*"}}
func Parse(conn net.Conn) (*HTTPRequest, *HTTPResponse, error) {
	buffer := make([]byte, 1)
	method, err := readUntilDelim(conn, buffer, ' ')
	if err != nil {
		return nil, nil, err
	}
	method = strings.ToUpper(method)
	path, err := readUntilDelim(conn, buffer, ' ')
	if err != nil {
		return nil, nil, err
	}
	verstr, err := readUntilDelim(conn, buffer, '\n')
	if err != nil {
		return nil, nil, err
	}
	ver, has := map[string]Version{"HTTP/1.0": V1_0, "HTTP/1.1": V1_1, "HTTP/2.0": V2_0}[verstr]
	if !has {
		ver = V_X
	}
	headers := map[string]string{}
	for {
		headerKey, err := readUntilDelim(conn, buffer, ':')
		if err != nil {
			return nil, nil, err
		}
		if headerKey == "" {
			break
		}
		headerValue, err := readUntilDelim(conn, buffer, '\n')
		if err != nil {
			return nil, nil, err
		}
		headers[strings.ToLower(strings.ReplaceAll(headerKey, "-", "_"))] = strings.TrimLeft(headerValue, " ")
	}
	return &HTTPRequest{method: method, path: path, ver: ver, headers: headers, socket: conn}, &HTTPResponse{ver: ver, headers: map[string]string{}, socket: conn, code: 200}, nil
}

// GetHeader ...
func (http *HTTPRequest) GetHeader(key string) (val string, has bool) {
	val, has = http.headers[key]
	return val, has
}

// GetMethod ...
func (http *HTTPRequest) GetMethod() string {
	return http.method
}

// GetPath ...
func (http *HTTPRequest) GetPath() string {
	return http.path
}

// GetVersion ...
func (http *HTTPRequest) GetVersion() Version {
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

// SetReadDeadline ...
func (http *HTTPRequest) SetReadDeadline(t time.Time) error {
	return http.socket.SetReadDeadline(t)
}

// Read ...
func (http *HTTPRequest) Read(b []byte) (n int, err error) {
	n, err = http.socket.Read(b)
	return
}

// WriteCode ...
func (http *HTTPResponse) WriteCode(code int) (int, error) {
	if http.state != initialize {
		return 0, errors.New("code already send")
	}
	counter := 0
	n, err := http.Write([]byte("HTTP/"))
	counter += n
	if err != nil {
		return 0, err
	}
	switch http.ver {
	case V1_0:
		n, err = http.Write([]byte("1.0"))
		counter += n
		if err != nil {
			return 0, err
		}
	case V1_1:
		n, err = http.Write([]byte("1.1"))
		counter += n
		if err != nil {
			return 0, err
		}
	case V2_0:
		n, err = http.Write([]byte("2.0"))
		counter += n
		if err != nil {
			return 0, err
		}
	}
	n, err = http.Write([]byte(" "))
	counter += n
	if err != nil {
		return 0, err
	}
	n, err = http.Write([]byte(fmt.Sprint(code)))
	counter += n
	if err != nil {
		return 0, err
	}
	n, err = http.Write([]byte("\r\n"))
	counter += n
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// WriteCodeDescription ...
func (http *HTTPResponse) WriteCodeDescription(code int, description string) (int, error) {
	if http.state != initialize {
		return 0, errors.New("code already send")
	}
	counter := 0
	n, err := http.Write([]byte("HTTP/"))
	counter += n
	if err != nil {
		return 0, err
	}
	switch http.ver {
	case V1_0:
		n, err = http.Write([]byte("1.0"))
		counter += n
		if err != nil {
			return 0, err
		}
	case V1_1:
		n, err = http.Write([]byte("1.1"))
		counter += n
		if err != nil {
			return 0, err
		}
	case V2_0:
		n, err = http.Write([]byte("2.0"))
		counter += n
		if err != nil {
			return 0, err
		}
	}
	n, err = http.Write([]byte(" "))
	counter += n
	if err != nil {
		return 0, err
	}
	n, err = http.Write([]byte(fmt.Sprint(code)))
	counter += n
	if err != nil {
		return 0, err
	}
	n, err = http.Write([]byte(description))
	counter += n
	if err != nil {
		return 0, err
	}
	n, err = http.Write([]byte("\r\n"))
	counter += n
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// AddHeader ...
func (http *HTTPResponse) AddHeader(key string, value string) {
	http.headers[key] = value
}

// Write ...
func (http *HTTPResponse) Write(b []byte) (int, error) {
	return http.socket.Write(b)
}

// WriteHeaders ...
func (http *HTTPResponse) WriteHeaders() (int, error) {
	counter := 0
	for key, value := range http.headers {
		n, err := http.Write([]byte(key))
		counter += n
		if err != nil {
			return 0, err
		}
		n, err = http.Write([]byte(": "))
		counter += n
		if err != nil {
			return 0, err
		}
		n, err = http.Write([]byte(value))
		counter += n
		if err != nil {
			return 0, err
		}
	}
	return counter, nil
}
