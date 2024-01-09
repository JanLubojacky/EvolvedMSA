package pkg

import (
	"fmt"
	"os"
  "math"
)

type Logger struct {
  LogFile *os.File
}

func (l *Logger) Init(filePath string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	l.LogFile = file
}

func (l *Logger) log(message string) {
	if l.LogFile != nil {	
		_, err := l.LogFile.WriteString(message)
		if err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	}
}

func (l Logger) WriteLog(fitness float64) {

  // feasible solution not found yet
  if fitness == math.Inf(1) {
    l.log("nan,")
  } else {
    l.log(fmt.Sprintf("%f,", fitness))
  }
}

// end the row in the log file by writing a new line
func (l Logger) EndRow() {
  l.LogFile.WriteString("\n")
}

func (l *Logger) Close() {
	if l.LogFile != nil {
		l.LogFile.Close()
	}
}
