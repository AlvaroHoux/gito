package term

import "fmt"

const (
	WARNING_COLOR string = "\x1b[33m"
	ERROR_COLOR   string = "\x1b[35m"
	SUCCESS_COLOR string = "\033[1;92m"
)

func Warning(data ...any) {
	printContent(WARNING_COLOR, data...)
}

func Error(err error) {
	printContent(ERROR_COLOR, err)
}

func Log(data ...any) {
	printContent("", data...)
}

func Success(data ...any) {
	printContent(SUCCESS_COLOR, data...)
}

func printContent(color string, data ...any) {
	fmt.Printf("%s🐙 Gito: ", color)
	for i := range data {
		fmt.Print(data[i], " ")
	}
	fmt.Print("\033[0m\n")
}
