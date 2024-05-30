package commandline

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	clear  map[string]func()
	reader *bufio.Reader
)

func Setup() {
	fmt.Println(WELCOME_MESSAGE)
	reader = bufio.NewReader(os.Stdin)

	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ReadInput() (text string) {
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return
}

func PressEnterToContinue() {
	fmt.Println(PRESS_ENTER_TO_CONTINUE)
	ReadInput()
}

func CleanTerminal() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic(UNSUPPORTED_OS)
	}
}

func PrintText(text string) {
	fmt.Println(LINE)
	fmt.Println(text)
	fmt.Println(LINE)
}

func GetLoader(i int) (text string) {
	text = STILL_PROCESSING_LOADER
	if i > 3 && i <= 8 {
		text += PROCESS_TAKING_TOO_LONG
	} else if i > 8 {
		text += PROCESS_TRY_AGAIN_LATER
	} else {
		text += PROCESS_WAIT_DEFAULT
	}
	return
}

func EndProgram() {
	CleanTerminal()
	fmt.Println(END_PROGRAM_MESSAGE)
}
