package log

import (
	"log"
	"os"
)

var (
// LoadingLogger = &logger{}
)

type logger struct {
	fatal   *log.Logger
	error   *log.Logger
	warning *log.Logger
	info    *log.Logger
}

func (lg *logger) fatalLogger() *log.Logger {
	if lg.fatal == nil {
		lg.fatal = log.New(os.Stdout, "FATAL: ", log.Lshortfile|log.LstdFlags)
	}

	return lg.fatal
}

func (lg *logger) errorLogger() *log.Logger {
	if lg.error == nil {
		lg.error = log.New(os.Stdout, "ERROR: ", log.Lshortfile|log.LstdFlags)
	}

	return lg.error
}

func (lg *logger) warningLogger() *log.Logger {
	if lg.warning == nil {
		lg.warning = log.New(os.Stdout, "WARNING: ", log.Lshortfile|log.LstdFlags)
	}

	return lg.warning
}

func (lg *logger) infoLogger() *log.Logger {
	if lg.info == nil {
		lg.info = log.New(os.Stdout, "INFO: ", log.Lshortfile|log.LstdFlags)
	}

	return lg.info
}

func (lg *logger) Fatal(v ...interface{}) {
	lg.fatalLogger().Println(v...)
}

func (lg *logger) Fatalf(format string, v ...interface{}) {
	lg.fatalLogger().Printf(format, v...)
}

func (lg *logger) Error(v ...interface{}) {
	lg.errorLogger().Println(v...)
}

func (lg *logger) Errorf(format string, v ...interface{}) {
	lg.errorLogger().Printf(format, v...)
}

func (lg *logger) Warning(v ...interface{}) {
	lg.warningLogger().Println(v...)
}

func (lg *logger) Warningf(format string, v ...interface{}) {
	lg.warningLogger().Printf(format, v...)
}

func (lg *logger) Info(v ...interface{}) {
	lg.infoLogger().Println(v...)
}

func (lg *logger) Infof(format string, v ...interface{}) {
	lg.infoLogger().Printf(format, v...)
}
