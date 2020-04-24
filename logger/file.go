package loggers

import (
	"fmt"
	"os"

	"github.com/keenfury/axenda/config"
)

type File struct{}

func (f *File) SetMessage(msg string) {
	file, errOpen := os.OpenFile(config.LogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errOpen != nil {
		panic(fmt.Sprintf("Unable to open file for writing: %s", errOpen))
	}
	defer file.Close()
	if _, errWrite := file.WriteString(fmt.Sprintf("%s\n", msg)); errWrite != nil {
		panic(fmt.Sprintf("Unable to write to file: %s", errWrite))
	}
}
