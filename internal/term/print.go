package term

import "fmt"

const (
	WARN_COLOR string = "\x1b[33m"
	ERROR_COLOR string = "\x1b[35m"
)

func Warn(data ...any) {
	printContent(WARN_COLOR, data...)
}

func Error(err error) {
	printContent(ERROR_COLOR, err)
}

func Log(data ...any) {
	printContent("", data...)
}


func printContent(color string, data ...any) {
	fmt.Printf("%s🐙 Gito: ", color)
	for i := range data {
		fmt.Print(data[i], " ")
	}
	fmt.Print("\x1b[39m\n")
}