package daemon

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/dongjialong2006/log"
)

type process struct {
	sync.RWMutex
	cmd  *exec.Cmd
	ctx  context.Context
	log  *log.Entry
	pid  int
	name string
}

func newProcess(ctx context.Context, name string) *process {
	return &process{
		ctx:  ctx,
		name: name,
		log:  log.New("process"),
	}
}

func (p *process) pipeLine(r io.Reader) {
	br := bufio.NewReader(r)
	for {
		l, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				p.log.WithField("pid", p.pid).Warn("pipe line is closed.")
				return
			}
			p.log.Error(err)
			return
		}

		p.output(string(l))
	}
}

func (p *process) output(value string) {
	value = strings.Replace(value, "  ", " ", -1)
	value = strings.Replace(value, "\t", " ", -1)
	value = strings.Trim(value, " ")
	if "" == value {
		return
	}
	temp := strings.SplitN(value, " ", 3)
	if len(temp) < 3 {
		p.log.Info(value)
		return
	}

	switch temp[1] {
	case "INFO":
		p.log.WithField("pid", p.pid).Info(temp[2])
	case "WARN", "WARNING":
		p.log.WithField("pid", p.pid).Warn(temp[2])
	case "DEBUG":
		p.log.WithField("pid", p.pid).Debug(temp[2])
	case "ERROR", "PANIC", "FATAL":
		p.log.WithField("pid", p.pid).Error(temp[2])
	}

	return
}

func (p *process) Start(args []string, envs []string) error {
	cmd := exec.CommandContext(p.ctx, p.name, args...)

	if 0 != len(envs) {
		cmd.Env = envs
		p.log.Infof("start command args:%s, envs:%s.", strings.Join(cmd.Args, " "), strings.Join(envs, " "))
	} else {
		p.log.Infof("start command args:%s.", strings.Join(cmd.Args, " "))
	}

	if nil == cmd {
		return fmt.Errorf("exec command handle error.")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go p.pipeLine(stdout)
	go p.pipeLine(stderr)

	if err = cmd.Start(); nil != err {
		return err
	}
	p.Lock()
	p.cmd = cmd
	p.pid = cmd.Process.Pid
	p.Unlock()

	return cmd.Wait()
}

func (p *process) Stop() error {
	p.RLock()
	defer p.RUnlock()
	process, err := os.FindProcess(p.pid)
	if nil != err {
		return nil
	}

	return process.Signal(syscall.SIGKILL)
}
