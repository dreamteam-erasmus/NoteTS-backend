package tools

import "webtools"

var Logger *webtools.ConsoleLogger

func INITLogger() {
	Logger = webtools.NewConsoleLogger("MAIN", 0)
}
