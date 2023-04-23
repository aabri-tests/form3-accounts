package logging

import (
	"fmt"
	"io"
	"os"
)

const (
	LevelNull  Level = 0
	LevelError Level = 1
	LevelWarn  Level = 2
	LevelInfo  Level = 3
	LevelDebug Level = 4
)

type Level uint32

type LeveledLogger interface {
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

type Leveled struct {
	Level          Level
	StderrOverride io.Writer
	StdoutOverride io.Writer
}

func (l *Leveled) Debugf(format string, v ...interface{}) {
	if l.Level >= LevelDebug {
		fmt.Fprintf(l.stdout(), "[DEBUG] "+format+"\n", v...)
	}
}

func (l *Leveled) Errorf(format string, v ...interface{}) {
	if l.Level >= LevelError {
		fmt.Fprintf(l.stderr(), "[ERROR] "+format+"\n", v...)
	}
}

func (l *Leveled) Infof(format string, v ...interface{}) {
	if l.Level >= LevelInfo {
		fmt.Fprintf(l.stdout(), "[INFO] "+format+"\n", v...)
	}
}

func (l *Leveled) Warnf(format string, v ...interface{}) {
	if l.Level >= LevelWarn {
		fmt.Fprintf(l.stderr(), "[WARN] "+format+"\n", v...)
	}
}

func (l *Leveled) stderr() io.Writer {
	if l.StderrOverride != nil {
		return l.StderrOverride
	}

	return os.Stderr
}

func (l *Leveled) stdout() io.Writer {
	if l.StdoutOverride != nil {
		return l.StdoutOverride
	}

	return os.Stdout
}
