package main

import (
	"crypto/tls"
	"flag"
	"goproxy-server/auth"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func main() {
	verbose := flag.Bool("v", false, "verbose log")
	port := flag.String("port", "18881", "listening port")
	certFile := flag.String("cert", "example.crt", "certificate PEM file")
	keyFile := flag.String("key", "example.key", "key PEM file")
	user := flag.String("user", "", "user name")
	pass := flag.String("pass", "", "password")
	flag.Parse()
	addr := ":" + *port
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	})

	if *user != "" && *pass != "" {
		log.Println("Proxy with Basic Auth enabled.")
		auth.ProxyBasic(proxy, "my_realm", func(_user string, _pass string) bool {
			return *user == _user && *pass == _pass
		})
	}
	log.Printf("Proxy is listening on %s\n", addr)
	srv := &http.Server{
		Handler:      proxy,
		Addr:         addr,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)), //disable http2
	}
	log.Fatal(srv.ListenAndServeTLS(*certFile, *keyFile))
}
