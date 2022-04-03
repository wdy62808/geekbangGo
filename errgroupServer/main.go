package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	mux := http.NewServeMux()
	mux.HandleFunc("ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	})
	serverQuit := make(chan struct{})
	mux.HandleFunc("quit", func(writer http.ResponseWriter, request *http.Request) {
		serverQuit <- struct{}{}
	})
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errorgroup quit")
		case <-serverQuit:
			log.Println("Server exit")
		}
		rctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		log.Println("shutdown server")
		return server.Shutdown(rctx)
	})
	g.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return errors.Errorf("get os signal:%v", sig)
		}
	})
	fmt.Printf("errgroup exinting:%+v\n", g.Wait())
}

type Group struct {
	// context 的 cancel 方法
	cancel func()

	// 复用 WaitGroup
	wg sync.WaitGroup

	// 用来保证只会接受一次错误
	errOnce sync.Once
	// 保存第一个返回的错误
	err error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}
