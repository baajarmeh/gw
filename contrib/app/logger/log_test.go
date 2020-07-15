package logger

import "testing"

func TestLogAll(t *testing.T) {
	Info("Info tester")
	Debug("Debug tester")
	Warn("Warn tester")
	Error("Error tester\n")
}
