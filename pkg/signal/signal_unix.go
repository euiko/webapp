//go:build linux || darwin || !windows
// +build linux darwin !windows

package signal

import (
	"os"
	"syscall"
)

var signals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGKILL,
	syscall.SIGTERM,
	syscall.SIGSTOP,
	syscall.SIGHUP,
	syscall.SIGQUIT,
}
