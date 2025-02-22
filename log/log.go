package log

import (
	"log"
	"os"
)

type Logger struct {
	loggerFile string
	log        *log.Logger
}

func NewLogger(file string) *Logger {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening/creating log file: %v", err)
	}

	l := log.New(f, "Monieshop", log.LstdFlags)
	return &Logger{loggerFile: file, log: l}
}

func (l *Logger) Println(msg ...any) {
	l.log.Println(msg...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.log.Printf(format, v...)
}
func (l *Logger) Fatalln(msg ...any) {
	l.log.Fatalln(msg...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.log.Fatalf(format, v...)
}
