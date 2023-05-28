package telebot

import (
	"fmt"
	"log"
)

type Logger interface {
	errorF(format string, v ...any)
	infoF(format string, v ...any)
}

type SysLogger struct {
}

func (sl *SysLogger) errorF(format string, v ...any) {
	log.Printf("ERROR: %s", fmt.Sprintf(format, v...))
}

func (sl *SysLogger) infoF(format string, v ...any) {
	log.Printf("INFO: %s", fmt.Sprintf(format, v...))
}
