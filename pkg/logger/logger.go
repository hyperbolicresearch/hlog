package logger

import (
	"encoding/json"
	"fmt"
	"io"
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
	Level Level
	io.Writer
}

func New(defaultLevel Level, writer io.Writer) *Logger {
	return &Logger{
		Level:  defaultLevel,
		Writer: writer,
	}
}

// Log will take a Log and write it in a readable/formatted manner
// to the passed io.Writer
func (l *Logger) Log(data interface{}) error {
	switch data := data.(type) {
	case []byte:
		data = append(data, '\n')
		_, err := l.Write(data)
		if err != nil {
			return err
		}
		return nil
	case string:
		_, err := l.Write([]byte(data + "\n"))
		if err != nil {
			return err
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

		_, err = l.Write([]byte(str))
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unsupported log format: %T", data)
}
