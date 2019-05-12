package miio

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	fileName = "firmware.pkg"
)

type ServerFile struct {
	io.ReadSeeker
}

type Server struct {
	listener net.Listener
	file     io.Reader
	size     int64
}

func NewServer(file io.Reader, size int64, hostname string) (*Server, error) {
	if hostname == "" {
		ip, err := localIP()
		if err != nil {
			return nil, err
		}

		hostname = ip.String()
	}

	listener, err := net.Listen("tcp", hostname+":")
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		file:     file,
		size:     size,
	}

	go func() {
		http.Serve(listener, server)

		fmt.Println("Close SERVER")
	}()

	return server, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.URL.Path, "/"+fileName) {
		w.Header().Set("Content-Type", "application/octet-stream")
		// если не указать Content-Length то progress методы будут отдавать всегда 0
		// что логично, так как устройство не знает конечного размера файла
		// хуже того, если SoundInstall установит и без этого, то OTA вернет ошибку.
		// Скорее всего OTA использует это как дополнительную проверку
		if s.size > 0 {
			w.Header().Set("Content-Length", strconv.FormatInt(s.size, 10))
		}

		io.Copy(w, s.file)
		return
	}

	http.NotFound(w, r)
}

func (s *Server) URL() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   s.listener.Addr().String(),
		Path:   "/" + fileName,
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func localIP() (*net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return &ip, nil
		}
	}

	return nil, errors.New("IP address not found")
}
