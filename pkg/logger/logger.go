package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hyperbolicresearch/hlog/internal/core"
)

type Level int

const (
	// Log levels
	_           = iota
	DEBUG Level = 1 + iota
	INFO
	WARNING
	ERROR
	FATAL
)

var levelCorrespondence = map[string]Level{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARNING,
	"error": ERROR,
	"fatal": FATAL,
}

const (
	// Color codes for pretty logging
	RESET string = "\033[0m"
	GRAY  string = "\033[37m"
	WHITE string = "\033[97m"
)

type LoggerI interface {
	Log(data interface{}) error
}

type Logger struct {
	sync.RWMutex
	Level   Level
	Writers []io.Writer
}

func New(defaultLevel Level, writer io.Writer) *Logger {
	return &Logger{
		Level:   defaultLevel,
		Writers: []io.Writer{writer},
	}
}

// AddWriter adds a new writer to the logger
func (l *Logger) AddWriter(w io.Writer) error {
	l.Lock()
	for _, v := range l.Writers {
		if v == w {
			l.Unlock()
			return nil
		}
	}
	l.Writers = append(l.Writers, w)
	fmt.Println("New writer")
	l.Unlock()
	return nil
}

// RemoveWriter removes a writer from the logger
func (l *Logger) RemoveWriter(w io.Writer) error { 
	l.Lock()
	for i, v := range l.Writers {
		if v == w {
			l.Writers = append(l.Writers[:i], l.Writers[i+1:]...)
		}
	}
	l.Unlock()
	return nil 
}

// Log will take a Log and write it in a readable/formatted manner
// to the passed io.Writer
func (l *Logger) Log(data interface{}) error {
	l.Lock()
	defer l.Unlock()
	switch data := data.(type) {
	case []byte:
		data = append(data, '\n')
		for _, l := range l.Writers {
			_, err := l.Write(data)
			if err != nil {
				return err
			}
		}
		return nil
	case string:
		for _, l := range l.Writers {
			_, err := l.Write([]byte(data + "\n"))
			if err != nil {
				return err
			}
		}
		return nil
	case core.Log:
		if levelCorrespondence[data.Level] < l.Level {
			return nil
		}
		unix := time.Unix(data.Timestamp, 0)
		readableTime := unix.Format("2006-01-02 15:04:05")
		jsonData, err := json.Marshal(data.Data)
		if err != nil {
			return err
		}
		// Format the log and write it
		str := "["
		str += data.Channel
		str += " | "
		str += data.SenderId
		str += "]"
		str += " "
		str += readableTime
		str += " "
		str += strings.ToUpper(data.Level)
		str += " "
		str += GRAY
		str += data.Message
		str += " "
		str += string(jsonData)
		str += RESET
		str += "\n"

		for _, l := range l.Writers {
			if _, ok := l.(*os.File); ok {
				_, err = l.Write([]byte(str))
				if err != nil {
					return err
				}
				continue
			}
			js, err := json.Marshal(data)
			if err != nil {
				return err
			}
			_, err = l.Write([]byte(js))
			if err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("unsupported log format: %T", data)
}
