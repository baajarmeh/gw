package logger

import "testing"

func TestLogAll(t *testing.T) {
	Info("Info tester")
	Debug("Debug tester")
	Warn("Warn tester")
	Error("Error tester\n")

	SetLogLevel(DEBUG)
	Info("Info tester") // should be not output.
	Debug("Debug tester")
	Warn("Warn tester")
	Error("Error tester\n")

	SetLogLevel(WARN)
	Info("Info tester")   // should be not output.
	Debug("Debug tester") // should be not output.
	Warn("Warn tester")
	Error("Error tester\n")

	SetLogLevel(ERROR)
	Info("Info tester")   // should be not output.
	Debug("Debug tester") // should be not output.
	Warn("Warn tester")   // should be not output.
	Error("Error tester\n")
}
