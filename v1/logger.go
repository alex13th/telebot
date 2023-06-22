package telebot

import (
	"fmt"
	"log"
)

type Logger interface {
	ErrorF(format string, v ...any)
	InfoF(format string, v ...any)
}

type SysLogger struct {
}

func (sl *SysLogger) ErrorF(format string, v ...any) {
	log.Printf("ERROR: %s", fmt.Sprintf(format, v...))
}

func (sl *SysLogger) InfoF(format string, v ...any) {
	log.Printf("INFO: %s", fmt.Sprintf(format, v...))
}
