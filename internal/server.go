package internal

import (
	"crypto/tls"
	"io/fs"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	CLIENT_TIMEOUT   = 5 * time.Second
	CLIENT_KEEPALIVE = 5 * time.Minute
	// Default pool size for reverse proxy
	POOLSIZE = 8 * 1024
)

type Server struct {
	logger     Logger
	fileServer http.Handler
	apiServer  http.Handler
}

type pool struct {
	*sync.Pool
}

func newPool() pool {
	return pool{
		Pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, POOLSIZE)
			},
		},
	}
}

func (p pool) Get() []byte {
	return p.Pool.Get().([]byte)
}

func (p pool) Put(b []byte) {
	p.Pool.Put(b)
}

func NewServer(logger Logger, root fs.FS, target *url.URL, insecureSkipVerify bool) *Server {
	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	reverseProxy.BufferPool = newPool()
	reverseProxy.ErrorLog = logger.Logger()
	dialer := &net.Dialer{
		Timeout:   CLIENT_TIMEOUT,
		KeepAlive: CLIENT_KEEPALIVE,
	}
	reverseProxy.Transport = &http.Transport{
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       2 * CLIENT_KEEPALIVE,
		TLSHandshakeTimeout:   CLIENT_TIMEOUT,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
	}
	return &Server{
		logger:     logger,
		fileServer: http.FileServer(http.FS(root)),
		apiServer:  reverseProxy,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("[%s] %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	for _, prefix := range []string{"/siddhi-apps", "/statistics"} {
		if strings.HasPrefix(r.URL.Path, prefix) {
			s.apiServer.ServeHTTP(w, r)
			return
		}
	}
	s.fileServer.ServeHTTP(w, r)
}

/*
Sample Siddhi API calls:
curl -X GET "http://localhost:9090/siddhi-apps" -H "accept: application/json" -u admin:admin -k
curl -X GET "https://localhost:9443/siddhi-apps?isActive=true" -H "accept: application/json" -u admin:admin -k
curl -X GET "http://localhost:9090/siddhi-apps/{app-file-name}/status" -H "accept: application/json" -u admin:admin -k
curl -X DELETE "http://localhost:9090/siddhi-apps/{app-name}" -H "accept: application/json" -u admin:admin -k
curl -X PUT "https://localhost:9443/siddhi-apps/TestSiddhiApp/statistics" -H "accept: application/json" -H "Content-Type: application/json" -d "{“statsEnable”:”true”}" -u admin:admin -k
*/
