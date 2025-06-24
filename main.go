package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
)

var addr = flag.String("addr", ":8080", "address for server")
var disableoptshandler = flag.Bool("disableoptshandler", false, "whether to desable the general options handler")

func main() {
	flag.Parse()
	mux := http.NewServeMux()
	svr := http.Server{
		Addr:    *addr,
		Handler: mux,
		// NOTE: https://www.rfc-editor.org/rfc/rfc7231#section-4.3.7
		DisableGeneralOptionsHandler: *disableoptshandler,
		// TODO: Is TLSConfig needed for deploying to AWS, GCP, etc.
		TLSConfig: nil,
		// TODO: does TLSNextProto ever need to be set?
		TLSNextProto:      nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		// NOTE: be default, will use DefaultMaxHeaderBytes;
		// Go adds an additional 4096 bytes to this maximum as a buffer
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ConnState: func(net.Conn, http.ConnState) {
			panic("TODO")
		},
		ErrorLog: &log.Logger{},
		BaseContext: func(net.Listener) context.Context {
			panic("TODO")
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			panic("TODO")
		},
		// NOTE: HTTP2 field currently has no effect.
		// see: https://github.com/golang/go/issues/67813
		HTTP2: nil,
		// NOTE: to set unencrypted HTTP/2 for gRPC,
		// do the following and assign to Protocols field
		// var protocols http.Protocols
		// protocols.SetUnencryptedHTTP2(true)
		Protocols: nil,
	}
	svr.ListenAndServe()
}
