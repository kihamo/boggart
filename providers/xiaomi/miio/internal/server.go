package internal

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
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
	size     int
	md5sum   string
}

func NewServer(file io.ReadSeeker, hostname string) (*Server, error) {
	size := 0

	// md5
	h := md5.New()

	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}

		if n > 0 {
			size += n
			h.Write(buf[:n])
		}
	}

	md5sum := hex.EncodeToString(h.Sum(nil))

	file.Seek(0, 0)

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
		md5sum:   md5sum,
	}

	go func() {
		http.Serve(listener, server)
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
			w.Header().Set("Content-Length", strconv.Itoa(s.size))
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

func (s *Server) MD5() string {
	return s.md5sum
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
