package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/go-httpproxy/httpproxy"
	"github.com/iamlongalong/go-socks5/socks5"
)

type params struct {
	User       string `env:"PROXY_USER" envDefault:""`
	Password   string `env:"PROXY_PASSWORD" envDefault:""`
	Socks5Port string `env:"SOCKS5_PROXY_PORT" envDefault:"10801"`
	HTTPPort   string `env:"HTTP_PROXY_PORT" envDefault:"10802"`
}

func main() {
	// Working with app params
	cfg := params{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	//Initialize socks5 config
	socsk5conf := &socks5.Config{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	if cfg.User+cfg.Password != "" {
		creds := socks5.StaticCredentials{
			os.Getenv("PROXY_USER"): os.Getenv("PROXY_PASSWORD"),
		}
		cator := socks5.UserPassAuthenticator{Credentials: creds}
		socsk5conf.AuthMethods = []socks5.Authenticator{cator}
	}

	server, err := socks5.New(socsk5conf)
	if err != nil {
		log.Fatal(err)
	}

	// socks5 server
	go func() {
		log.Printf("Start listening sock5 proxy service on port %s\n", cfg.Socks5Port)
		if err := server.ListenAndServe("tcp", ":"+cfg.Socks5Port); err != nil {
			log.Fatal(err)
		}
	}()

	// http server
	prx, _ := httpproxy.NewProxy()

	hook := &HttpProxyHook{}

	if cfg.User != "" || cfg.Password != "" {
		hook.auth = &auth{
			user: cfg.User,
			pass: cfg.Password,
		}
		prx.OnAuth = hook.OnAuth
	}

	// Set handlers.
	prx.OnError = hook.OnError
	prx.OnAccept = hook.OnAccept
	prx.OnRequest = hook.OnRequest
	prx.OnResponse = hook.OnResponse

	// Listen...
	log.Printf("Start listening http proxy service on port %s\n", cfg.Socks5Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), prx))
}

type auth struct {
	user string
	pass string
}

type HttpProxyHook struct {
	auth *auth
}

func (h *HttpProxyHook) OnError(ctx *httpproxy.Context, where string,
	err *httpproxy.Error, opErr error) {
	// Log errors.
	log.Printf("ERR: %s: %s [%s]", where, err, opErr)
}

func (h *HttpProxyHook) OnAccept(ctx *httpproxy.Context, w http.ResponseWriter,
	r *http.Request) bool {
	// Handle local request has path "/info"
	if r.Method == "GET" && !r.URL.IsAbs() && r.URL.Path == "/info" {
		w.Write([]byte("This is go-httpproxy."))
		return true
	}
	return false
}

func (h *HttpProxyHook) OnAuth(ctx *httpproxy.Context, authType string, user string, pass string) bool {
	// Auth test user.
	if h.auth != nil {
		return user == h.auth.user && pass == h.auth.pass
	}

	return true
}

func (h *HttpProxyHook) OnConnect(ctx *httpproxy.Context, host string) (
	ConnectAction httpproxy.ConnectAction, newHost string) {
	// Apply "Man in the Middle" to all ssl connections. Never change host.
	return httpproxy.ConnectMitm, host
}

func (h *HttpProxyHook) OnRequest(ctx *httpproxy.Context, req *http.Request) (
	resp *http.Response) {
	// Log proxying requests.
	log.Printf("INFO: Proxy: %s %s", req.Method, req.URL.String())
	return
}

func (h *HttpProxyHook) OnResponse(ctx *httpproxy.Context, req *http.Request,
	resp *http.Response) {
	// Add header "Via: go-httpproxy".
	resp.Header.Add("Via", "go-httpproxy")
}
