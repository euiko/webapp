package signal

import (
	"context"
	"os"
	"os/signal"

	"github.com/euiko/webapp/pkg/log"
)

type (
	SignalHandler func(ctx context.Context, sig os.Signal) bool

	SignalNotifier struct {
		handlers []SignalHandler
	}
)

func NewSignalNotifier() *SignalNotifier {
	return &SignalNotifier{
		handlers: []SignalHandler{},
	}
}

func (sn *SignalNotifier) OnSignal(handler SignalHandler) *SignalNotifier {
	sn.handlers = append(sn.handlers, handler)
	return sn
}

func (sn SignalNotifier) Wait(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	// wait for signal
	signal.Notify(sigChan, signals...)

	// reset the watched signals
	defer signal.Ignore(signals...)

	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigChan:
			log.Trace("Signal received, calling signal handlers...")
			// exit wait if any handler returns true
			if sn.callHandlers(ctx, sig) {
				return
			}

		}
	}
}

func (sn SignalNotifier) callHandlers(ctx context.Context, sig os.Signal) bool {
	exited := false
	for _, h := range sn.handlers {
		if h(ctx, sig) {
			exited = true
		}
	}

	return exited
}
