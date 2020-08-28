package unit

import (
	"fmt"
	"io"
)

var indent int

var writer io.Writer // = os.Stderr

func printIndent() {
	for i := 0; i < indent; i++ {
		fmt.Fprintf(writer, ". ")
	}
}

func tracein(msg string, args ...interface{}) func() {
	if writer == nil {
		return func() {}
	}
	printIndent()
	indent++
	fmt.Fprintf(writer, msg, args...)
	fmt.Fprintln(writer, "(")
	return func() {
		indent--
		printIndent()
		fmt.Fprintf(writer, ")\n")
	}
}
func tracemsg(msg string, args ...interface{}) {
	if writer == nil {
		return
	}
	printIndent()
	fmt.Fprintf(writer, msg, args...)
	fmt.Fprintln(writer)
}
