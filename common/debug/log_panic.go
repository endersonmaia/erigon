package debug

import (
	"os"
	"runtime/debug"
	"syscall"

	"github.com/ledgerwatch/erigon-lib/common/dbg"
	"github.com/ledgerwatch/log/v3"
)

var sigc chan os.Signal

func GetSigC(sig *chan os.Signal) {
	sigc = *sig
}

// LogPanic - does log panic to logger and to <datadir>/crashreports then stops the process
func LogPanic() {
	panicResult := recover()
	if panicResult == nil {
		return
	}

	stack := string(debug.Stack())
	log.Error("catch panic", "err", panicResult, "stack", dbg.PanicReplacer.Replace(stack))
	//WriteStackTraceOnPanic(stack)
	if sigc != nil {
		sigc <- syscall.SIGINT
	}
}
