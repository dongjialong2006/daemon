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
		daemon := New(context.Background(), "./bin/sslvpn-agent", WithArgs([]string{"-config", "./bin/config.ini", "-d"}), WithProcessNum(3))
		go func() {
			time.Sleep(time.Second * 5)
			daemon.Stop()
		}()

		daemon.Start()
	})
})
