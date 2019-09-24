package gocamole

import "log"

type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type DefaultLogger struct {
	Quiet bool
}

func (l *DefaultLogger) Tracef(format string, args ...interface{}) {
	if !l.Quiet {
		log.Printf("TRAC: "+format, args...)
	}
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	if !l.Quiet {
		log.Printf("DEBU: "+format, args...)
	}
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	if !l.Quiet {
		log.Printf("INFO: "+format, args...)
	}
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	if !l.Quiet {
		log.Printf("WARN: "+format, args...)
	}
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	if !l.Quiet {
		log.Printf("ERRO: "+format, args...)
	}
}
