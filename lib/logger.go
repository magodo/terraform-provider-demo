package lib

import (
	"io"
	"log"
	"os"
)

const logOutput = "/tmp/logoutput.txt"

type Logger struct {
	io.WriteCloser
}

func (l *Logger) Log(s string) {
	l.Write([]byte(s + "\n"))

}

func (l *Logger) Close() {
	l.Close()
}

var logger *Logger

func init() {
	f, err := os.Create(logOutput)
	if err != nil {
		log.Fatal(err)
	}
	logger = &Logger{f}
}
