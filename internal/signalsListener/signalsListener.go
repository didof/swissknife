/*
<https://www.abilityrush.com/how-to-catch-system-signals-in-golang-for-a-graceful-exit/>
*/

package signalslistener

import (
	"context"
	"log"
	"os"
	"os/signal"
)

type signalsListener struct {
	signals []os.Signal
}

func NewSignalsListener(signals ...os.Signal) *signalsListener {
	return &signalsListener{
		signals: signals,
	}
}

func (sl *signalsListener) AddSignals(signals ...os.Signal) {
	sl.signals = append(sl.signals, signals...)
}

func (sl *signalsListener) Listen(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, sl.signals...)

	go func() {
		<-doneCh
		log.Println("exiting...")
		cancel()
	}()

	return ctx
}
