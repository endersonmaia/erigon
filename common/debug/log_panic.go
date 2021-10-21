package debug

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/ledgerwatch/log/v3"
)

var sigc chan os.Signal

func GetSigC(sig *chan os.Signal) {
	sigc = *sig
}

var panicReplacer = strings.NewReplacer("\n", " ", "\t", "", "\r", "")

// LogPanic - does log panic to logger and to <datadir>/crashreports then stops the process
func LogPanic() {
	panicResult := recover()
	if panicResult == nil {
		return
	}

	stack := string(debug.Stack())
	log.Error("catch panic", "err", panicResult, "stack", panicReplacer.Replace(stack))
	//WriteStackTraceOnPanic(stack)
	if sigc != nil {
		sigc <- syscall.SIGINT
	}
}

// ReportPanicAndRecover - does save panic to datadir/crashreports, bud doesn't log to logger and doesn't stop the process
// it returns recovered panic as error in format friendly for our logger
// common pattern of use - assign to named output param:
//  func A() (err error) {
//	    defer func() { err = debug.ReportPanicAndRecover(err) }() // avoid crash because Erigon's core does many things
//  }
func ReportPanicAndRecover(err error) error {
	panicResult := recover()
	if panicResult == nil {
		return err
	}

	stack := string(debug.Stack())
	switch typed := panicResult.(type) {
	case error:
		err = fmt.Errorf("%w, trace: %s", typed, panicReplacer.Replace(stack))
	default:
		err = fmt.Errorf("%+v, trace: %s", typed, panicReplacer.Replace(stack))
	}
	return err
}
