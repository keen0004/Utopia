package tests

import (
	"errors"
	"os"
	"regexp"
	"testing"
	"utopia/internal/logger"
)

var (
	logfile = "./test.log"
)

func TestLog(t *testing.T) {
	os.Remove(logfile)
	logger.SetLogPath(logfile)
	defer os.Remove(logfile)

	logger.Debug("%s", "debug")
	logger.Info("%d", 100)
	logger.Warn("%f", 32.5)
	logger.Error("%v", errors.New("error"))

	data, err := os.ReadFile(logfile)
	if err != nil {
		t.Errorf("Open log filed failed with error: %v", err)
		return
	}

	match, err := regexp.Match("\\[DEBUG\\] [0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} logger_test.go:[0-9]+ debug\\n\\[INFO\\].*", data)
	if err != nil || !match {
		t.Errorf("Not match the format")
		return
	}
}
