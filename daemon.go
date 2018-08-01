package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dongjialong2006/log"
)

type daemon struct {
	ch    chan os.Signal
	nodes []config
}

type node struct {
	ctx    context.Context
	opts   []option
	name   string
	log    *log.Entry
	handle *process
	cancel context.CancelFunc
}

func New(nodes ...config) *daemon {
	return &daemon{
		ch:    make(chan os.Signal),
		nodes: nodes,
	}
}

func (d *daemon) Start() {
	for _, node := range d.nodes {
		if nil != node {
			go node.Start()
		}
	}

	d.notify()
}

func (d *daemon) Stop() {
	for _, node := range d.nodes {
		if nil != node {
			node.Stop()
		}
	}

	d.ch <- syscall.SIGINT
}

func NewNode(ctx context.Context, name string, opts ...option) config {
	ctx, cancel := context.WithCancel(ctx)
	return &node{
		ctx:    ctx,
		log:    log.New("node"),
		name:   name,
		opts:   opts,
		cancel: cancel,
	}
}

func (d *node) Start() {
	if !check(d.findArgs()) {
		return
	}

	d.start()
}

func (d *node) Stop() {
	if nil != d.cancel {
		d.cancel()
	}
}

func (d *node) Notify() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGSTOP, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	<-sigChan
	d.Stop()

	signal.Stop(sigChan)

	return
}
