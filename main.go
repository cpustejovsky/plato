package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cpustejovsky/plato/router"
)

const ServerAddres string = "address"

var addr = flag.String("addr", ":8080", "address for server")
var disableoptshandler = flag.Bool("d", false, "whether to desable the general options handler")
var idleTimeout = flag.Int64("i", 0, "Server IdleTimeout property (in seconds)")
var readHeaderTimeout = flag.Int64("r", 0, "Server ReadHeaderTimeout property (in seconds)")

func main() {
	flag.Parse()
	router, err := router.New()
	if err != nil {
		slog.Error("Error creating router", "error message", err.Error())
		os.Exit(1)
	}
	// TODO: determine if writer for logger should be a variable determined by a flag
	errorLog := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelError)
	svrCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svr := http.Server{
		Addr:              *addr,
		Handler:           router,
		IdleTimeout:       time.Duration(*idleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(*readHeaderTimeout) * time.Second,
		// NOTE: be default, will use DefaultMaxHeaderBytes;
		// Go adds an additional 4096 bytes to this maximum as a buffer
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ErrorLog:       errorLog,
		// TODO: Set up ConnState for monitoring
		// ConnState: func(net.Conn, http.ConnState) {
		// 	panic("TODO")
		// },
		// TODO: make sure this is the correct way to use BaseContext
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(svrCtx, ServerAddres, l.Addr().String())
		},
		// TODO: determine how to use this
		// ConnContext: func(ctx context.Context, c net.Conn) context.Context {
		// 	panic("TODO")
		// },
		// NOTE: https://www.rfc-editor.org/rfc/rfc7231#section-4.3.7
		DisableGeneralOptionsHandler: *disableoptshandler,
		// NOTE: TLSConfig not needed for deploying to AWS, GCP, etc.
		TLSConfig: nil,
		// NOTE: only set if you won't use HTTP/2
		TLSNextProto: nil,
		// NOTE: handlers should handle their own read and write timeouts
		// WriteTimeout: 0,
		// ReadTimeout:  0,
		// NOTE: HTTP2 field currently has no effect.
		// see: https://github.com/golang/go/issues/67813
		HTTP2: nil,
		// NOTE: to set unencrypted HTTP/2 for gRPC,
		// do the following and assign to Protocols field
		// var protocols http.Protocols
		// protocols.SetUnencryptedHTTP2(true)
		Protocols: nil,
	}
	l, err := net.Listen("tcp", svr.Addr)
	if err != nil {
		slog.Error("failed to listen on network",
			"network", "tcp", "address", *addr, "error message", err.Error())
		os.Exit(1)
	}
	// run server in a goroutine so we can multiplex between signal and error
	// handling below.
	errCh := make(chan error, 1)
	go func() {
		slog.Info("Server Started", "network", "tcp", "address", *addr)
		errCh <- svr.Serve(l)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	case <-ctx.Done():
		slog.Error("server shutting down", "error", ctx.Err())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := svr.Shutdown(ctx)
		if err != nil {
			slog.Error("failed to shutdown server, exiting anyway", "error", err)
			os.Exit(1)

		}
		slog.Info("Server shut down successfully")

	}
}
