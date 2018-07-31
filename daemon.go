package daemon

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dongjialong2006/log"
)

type daemon struct {
	ctx    context.Context
	opts   []option
	name   string
	num    int32
	log    *log.Entry
	handle *process
	cancel context.CancelFunc
}

func New(ctx context.Context, name string, opts ...option) *daemon {
	ctx, cancel := context.WithCancel(ctx)
	return &daemon{
		ctx:    ctx,
		log:    log.New("daemon"),
		name:   name,
		opts:   opts,
		cancel: cancel,
	}
}

func (d *daemon) Start() {
	if !check(d.findArgs()) {
		return
	}

	d.num = d.findProcessNum()

	var ext bool = false
	var wg sync.WaitGroup
	for {
		select {
		case <-d.ctx.Done():
			ext = true
		default:
			if atomic.LoadInt32(&d.num) == 0 {
				time.Sleep(time.Second)
				continue
			}

			go d.watch(&wg)
			atomic.AddInt32(&d.num, -1)
		}

		if ext {
			break
		}
	}

	wg.Wait()
}

func (d *daemon) Stop() {
	if nil != d.cancel {
		d.cancel()
	}
}

func (d *daemon) Notify() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGSTOP, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	<-sigChan
	d.Stop()

	signal.Stop(sigChan)

	return
}
