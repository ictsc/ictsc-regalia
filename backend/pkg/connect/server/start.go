package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/defval/di"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// Start サーバーを起動する
//
//	"github.com/defval/di" でInvokeするための関数
func Start(ctx context.Context, srv []*http.Server, container *di.Container) error {
	defer container.Cleanup()

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer cancel()

	errChan := make(chan error)

	for _, s := range srv {
		go func(s *http.Server) {
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}(s)
	}

	log.Println("Server started")

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		return errors.Wrap(errors.ErrUnknown, err)
	}
}
