package tests

import (
	"fmt"
	"testing"

	"email-wizard/backend/logger"
)

func panic_func() {
	x := 0
	y := 1 / x
	fmt.Println(y)
}

func TestLogger(t *testing.T) {
	logger.InitLogger("log", "backend", 1, 7, "INFO")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("capturing uncaught panic")
		}
	}()
	defer logger.LogErrorStackTrace()
	logger.Info("this is a info message")
	logger.Error("this is a error message")
	logger.Warn("this is a warning message")
	panic_func()
}