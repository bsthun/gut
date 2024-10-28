package gut

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
)

var timeColor = color.New(color.FgGreen)
var callerColor = color.C256(140)
var messageColor = color.S256(110).SetOpts(color.Opts{color.OpBold})
var argKeyColor = color.S256(246).SetOpts(color.Opts{color.OpBold})
var argValueColor = color.C256(246)
var argErrorColor = color.New(color.FgRed)
var separatorColor = color.New(color.FgGray)

func PrintTime() string {
	return timeColor.Render(time.Now().Format("03:04:05 PM"))
}

func PrintCaller(skip int) string {
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		panic("No caller information")
	}
	name := runtime.FuncForPC(pc).Name()
	name = name[strings.LastIndex(name, "/")+1:]
	return callerColor.Sprint(name + ":" + strconv.Itoa(line))
}

func PrintSeparator() string {
	return separatorColor.Render(" • ")
}

func Debug(args ...any) {
	var output strings.Builder

	// Print time
	output.WriteString(PrintTime())

	// Print caller
	output.WriteString(PrintSeparator())
	output.WriteString(PrintCaller(2))

	// Print message
	if v, ok := args[0].(string); ok {
		output.WriteString(PrintSeparator())
		output.WriteString(messageColor.Sprint(v))
		args = args[1:]
	}

	// Print args
	key := true
	for len(args) > 0 {
		if v, ok := args[0].(string); ok && key {
			output.WriteString("\n")
			output.WriteString(argKeyColor.Sprint(" - " + v))
			key = false
		} else if v, ok := args[0].(error); ok {
			output.WriteString("\n")
			output.WriteString(argErrorColor.Render(" - " + v.Error()))
			key = false
		} else {
			output.WriteString(argValueColor.Sprint(" " + fmt.Sprintf("%v", args[0])))
			key = true
		}
		args = args[1:]
	}
	output.WriteString("\n")
	fmt.Print(output.String())
}

func Error(message string, err error, fatal ...bool) {
	var output strings.Builder

	// Print time
	output.WriteString(PrintTime())

	// Print caller
	output.WriteString(PrintSeparator())
	if len(fatal) > 0 {
		output.WriteString(PrintCaller(3))
	} else {
		output.WriteString(PrintCaller(2))
	}

	// Print message
	output.WriteString(PrintSeparator())
	output.WriteString(messageColor.Sprint(message))

	// Print error
	if err != nil {
		if strings.Index(err.Error(), "\n") != -1 {
			output.WriteString("\n" + argErrorColor.Render(err.Error()))
		} else {
			output.WriteString(PrintSeparator() + argErrorColor.Render(err.Error()))
		}
	}

	output.WriteString("\n")
	fmt.Print(output.String())
}

func Fatal(message string, err error) {
	Error(message, err, true)
	os.Exit(1)
}
