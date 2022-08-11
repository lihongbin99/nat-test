package log

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	consoleColor  bool
	globalLogType string
)

func init() {
	flag.BoolVar(&consoleColor, "color", true, "consoleColor")
	flag.StringVar(&globalLogType, "log", "nil", "debug or trace")
}

func getTime() string {
	return time.Now().Format("2006:01:02 15:04:05")
}

type log struct {
	from string

	fatalColorStart string
	errorColorStart string
	infoColorStart  string
	warnColorStart  string
	debugColorStart string
	traceColorStart string
	colorEnd        string

	debug bool
	trace bool
}

func Log(from string) *log {
	result := log{from, "", "", "", "", "", "", "", false, false}
	if consoleColor {
		result.fatalColorStart = "\033[31m"
		result.errorColorStart = "\033[31m"
		result.infoColorStart = "\033[32m"
		result.warnColorStart = "\033[33m"
		result.debugColorStart = "\033[34m"
		result.traceColorStart = "\033[34m"
		result.colorEnd = "\033[0m"
	}
	switch globalLogType {
	case "debug":
		result.debug = true
	case "trace":
		result.debug = true
		result.trace = true
	}
	return &result
}

func (l *log) Fatal(err error) {
	fmt.Printf("%sFatal %s -> %s: %v%s\n", l.fatalColorStart, getTime(), l.from, err, l.colorEnd)
	os.Exit(1)
}

func (l *log) Error(err error) {
	if err != nil {
		fmt.Printf("%sError %s -> %s: %v%s\n", l.errorColorStart, getTime(), l.from, err, l.colorEnd)
	} else {
		fmt.Printf("%sError %s -> %s%s\n", l.errorColorStart, getTime(), l.from, l.colorEnd)
	}
}

func (l *log) Info(message interface{}) {
	if message != nil {
		fmt.Printf("%sInfo %s -> %s: %v%s\n", l.infoColorStart, getTime(), l.from, message, l.colorEnd)
	} else {
		fmt.Printf("%sInfo %s -> %s%s\n", l.infoColorStart, getTime(), l.from, l.colorEnd)
	}
}

func (l *log) Warn(cid int32, sid int32, message interface{}) {
	if message != nil {
		fmt.Printf("%sWarn %s -> %s(%d -> %d): %v%s\n", l.warnColorStart, getTime(), l.from, cid, sid, message, l.colorEnd)
	} else {
		fmt.Printf("%sWarn %s -> %s(%d -> %d)%s\n", l.warnColorStart, getTime(), l.from, cid, sid, l.colorEnd)
	}
}

func (l *log) Debug(cid int32, sid int32, message interface{}) {
	if l.debug {
		if message != nil {
			fmt.Printf("%sDebug %s -> %s(%d -> %d): %v%s\n", l.debugColorStart, getTime(), l.from, cid, sid, message, l.colorEnd)
		} else {
			fmt.Printf("%sDebug %s -> %s(%d -> %d)%s\n", l.debugColorStart, getTime(), l.from, cid, sid, l.colorEnd)
		}
	}
}

func (l *log) Trace(cid int32, sid int32, message interface{}) {
	if l.trace {
		if message != nil {
			fmt.Printf("%sTrace %s -> %s(%d -> %d): %v%s\n", l.traceColorStart, getTime(), l.from, cid, sid, message, l.colorEnd)
		} else {
			fmt.Printf("%sTrace %s -> %s(%d -> %d)%s\n", l.traceColorStart, getTime(), l.from, cid, sid, l.colorEnd)
		}
	}
}
