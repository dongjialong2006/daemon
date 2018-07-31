package daemon

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/dongjialong2006/log"
)

type process struct {
	cmd  *exec.Cmd
	ctx  context.Context
	log  *log.Entry
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
				p.log.Warn("pipe line is closed.")
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
		p.log.Info(temp[2])
	case "WARN", "WARNING":
		p.log.Warn(temp[2])
	case "DEBUG":
		p.log.Debug(temp[2])
	case "ERROR", "PANIC", "FATAL":
		p.log.Error(temp[2])
	}

	return
}

func (p *process) Start(args []string, envs []string) error {
	p.cmd = exec.CommandContext(p.ctx, p.name, args...)

	if 0 != len(envs) {
		p.cmd.Env = envs
		p.log.Infof("start command args:%s, envs:%s", strings.Join(p.cmd.Args, " "), strings.Join(envs, " "))
	} else {
		p.log.Infof("start command args:%s", strings.Join(p.cmd.Args, " "))
	}

	if nil == p.cmd {
		return fmt.Errorf("exec command handle error.")
	}

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}
	go p.pipeLine(stdout)
	go p.pipeLine(stderr)

	return p.cmd.Start()
}

func (p *process) Stop() error {
	if nil != p.cmd && nil != p.cmd.Process {
		err := p.cmd.Process.Signal(syscall.SIGKILL)
		if nil != err {
			return err
		}
		p.cmd.Wait()
	}

	return nil
}
