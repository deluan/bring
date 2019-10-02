package bring

import "log"

// Logger interface used by this package. It is compatible with Logrus,
// but anything implementing this interface can be used
type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Simple console logger
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

// Logger that discards all messages
type DiscardLogger struct{}

func (d *DiscardLogger) Tracef(format string, args ...interface{}) {}

func (d *DiscardLogger) Debugf(format string, args ...interface{}) {}

func (d *DiscardLogger) Infof(format string, args ...interface{}) {}

func (d *DiscardLogger) Warnf(format string, args ...interface{}) {}

func (d *DiscardLogger) Errorf(format string, args ...interface{}) {}
