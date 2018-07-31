# daemon
守护进程

## 引用

- import "github.com/dongjialong2006/daemon"

## 示例
{
	package main
	
	import "context"
	import "github.com/dongjialong2006/daemon"
	
	func main() {
		dae := daemon.New(context.Background(), "./bin/sslvpn-agent", daemon.WithArgs([]string{"-config", "./bin/config.ini", "-d"}), daemon.WithProcessNum(3))
		
		go dae.Start()
		dae.Notify()
	}
}
