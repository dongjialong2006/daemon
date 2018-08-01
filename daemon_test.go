package daemon

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDaemon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Daemon Suite")
}

var _ = Describe("Daemon", func() {
	Specify("daemon start test", func() {
		// daemon := New(NewNode(context.Background(), "./bin/sslvpn-agent", WithArgs([]string{"-config", "./bin/204.ini", "-d"}), WithProcessNum(1)), NewNode(context.Background(), "./bin/sslvpn-agent", WithArgs([]string{"-config", "./bin/205.ini", "-d"}), WithProcessNum(1)))
		daemon := New(NewNode(context.Background(), "./bin/sslvpn-agent", WithArgs([]string{"-config", "./bin/204.ini", "-d"}), WithProcessNum(1)))
		go func() {
			time.Sleep(time.Second * 20)
			daemon.Stop()
		}()

		daemon.Start()
	})

	/*
		Specify("daemon start test", func() {
			daemon := NewNode(context.Background(), "./bin/sslvpn-agent", WithArgs([]string{"-config", "./bin/205.ini", "-d"}), WithProcessNum(3))
			go func() {
				time.Sleep(time.Second * 5)
				daemon.Stop()
			}()

			daemon.Start()
		})
	*/
})
