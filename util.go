package daemon

import (
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type config interface {
	Start()
	Stop()
}

func filter(args []string) []string {
	var tmp []string = nil
	for _, arg := range args {
		if "-d" == arg || "-daemon" == arg {
			continue
		}

		if keys := split(arg); len(keys) > 0 {
			tmp = append(tmp, keys...)
			continue
		}
		tmp = append(tmp, arg)
	}

	return tmp
}

func split(arg string) []string {
	var tmp []string = nil
	if strings.Contains(arg, "=") {
		keys := strings.Split(arg, "=")
		for _, key := range keys {
			if "" == key {
				continue
			}
			tmp = append(tmp, key)
		}
	}

	return tmp
}

func check(args []string) bool {
	for _, arg := range args {
		if "-d" == arg || "-daemon" == arg {
			return true
		}
	}

	return false
}

func (d *node) start() {
	num := d.findProcessNum()
	var ext bool = false
	var wg sync.WaitGroup

	for {
		select {
		case <-d.ctx.Done():
			ext = true
		default:
			if atomic.LoadInt32(&num) == 0 {
				time.Sleep(time.Second)
				continue
			}

			go d.watch(&num, &wg)
			atomic.AddInt32(&num, -1)
		}

		if ext {
			break
		}
	}

	wg.Wait()
}

func (d *node) findArgs() []string {
	for _, opt := range d.opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(process_args); nil != value {
			return value.([]string)
		}
	}

	return nil
}

func (d *node) findEnvs() []string {
	for _, opt := range d.opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(process_envs); nil != value {
			return value.([]string)
		}
	}

	return nil
}

func (d *node) findProcessNum() int32 {
	var num int32 = 1
	for _, opt := range d.opts {
		if nil == opt {
			continue
		}

		if value := opt.Get(process_num); nil != value {
			num = value.(int32)
			if num == 0 {
				num = 1
			}
		}
	}

	return num
}

func (d *node) watch(num *int32, wg *sync.WaitGroup) {
	wg.Add(1)
	defer atomic.AddInt32(num, 1)
	defer wg.Done()

	var done = make(chan struct{})
	tick := time.Tick(time.Second)
	for {
		select {
		case <-d.ctx.Done():
			d.close()
			return
		case <-done:
			d.close()
			return
		case <-tick:
			d.close()
			d.process(done)
			tick = time.Tick(time.Hour * 24 * 30)
		}
	}

	return
}

func (d *node) process(done chan struct{}) {
	defer close(done)

	d.handle = newProcess(d.ctx, d.name)
	if nil == d.handle {
		d.log.Error("new process handle is nil.")
		return
	}

	if err := d.handle.Start(filter(d.findArgs()), d.findEnvs()); nil != err {
		d.log.Error(err)
	}

	return
}

func (d *node) close() {
	if nil == d.handle {
		return
	}

	if err := d.handle.Stop(); nil != err {
		d.log.Error(err)
	}
}

func (d *daemon) notify() {
	signal.Notify(d.ch, syscall.SIGSTOP, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGCHLD)

	<-d.ch

	for _, node := range d.nodes {
		if nil != node {
			node.Stop()
		}
	}

	signal.Stop(d.ch)

	return
}
