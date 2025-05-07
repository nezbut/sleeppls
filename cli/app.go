package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nezbut/sleeppls"
	"github.com/nezbut/sleeppls/config"
)

type stopF func() bool

type App struct {
	conf       *config.Config
	notifiers  []sleeppls.Notifier
	shutdowner sleeppls.ShutDowner
}

func New(conf *config.Config, shutdowner sleeppls.ShutDowner, notifiers ...sleeppls.Notifier) *App {
	return &App{
		conf:       conf,
		notifiers:  notifiers,
		shutdowner: shutdowner,
	}
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()
	return a.StartWithContext(ctx)
}

func (a *App) StartWithContext(ctx context.Context) error {
	errCh := make(chan error)
	shutdownStop, notifyStop := a.scheduleFuncs(errCh)
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	defer close(interruptCh)
	defer shutdownStop()
	defer notifyStop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutdown...")
			return nil
		case <-interruptCh:
			var answer string
			fmt.Print("\nAre you sure you want to get out?(y/N): ")
			n, err := fmt.Scan(&answer)
			if (answer == "" || n == 0) || err != nil {
				continue
			} else if strings.ToLower(answer) == "y" {
				shutdownStop()
				notifyStop()
				return nil
			}
		case err := <-errCh:
			if err != nil {
				return err
			}
		}
	}
}

func (a *App) scheduleFuncs(errCh chan<- error) (stopF, stopF) {
	NotifyCtx, NotifyCtxStop := context.WithDeadline(context.Background(), a.conf.TimeToShutDown.Add(-a.conf.NotifyDuration))
	ShutdownCtx, ShutdownCtxStop := context.WithDeadline(context.Background(), a.conf.TimeToShutDown)
	shutdownStop := context.AfterFunc(ShutdownCtx, func() {
		defer ShutdownCtxStop()
		err := a.shutdowner.ShutDown()
		if err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	})
	notifyStop := context.AfterFunc(NotifyCtx, func() {
		defer NotifyCtxStop()
		remaining := time.Until(a.conf.TimeToShutDown)
		msg := fmt.Sprintf("Shutdown via:\n%s", formatShutdownCountdown(remaining))
		for _, notifier := range a.notifiers {
			err := notifier.Notify(msg)
			if err != nil {
				errCh <- fmt.Errorf("error during notify: %w", err)
				return
			}
		}
		errCh <- nil
	})
	return shutdownStop, notifyStop
}
