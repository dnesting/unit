package unit

import (
	"fmt"
	"io"
)

var indent int

var Debug io.Writer // = os.Stderr

func printIndent() {
	for i := 0; i < indent; i++ {
		fmt.Fprintf(Debug, ". ")
	}
}

func tracein(msg string, args ...interface{}) func() {
	if Debug == nil {
		return func() {}
	}
	printIndent()
	indent++
	fmt.Fprintf(Debug, msg, args...)
	fmt.Fprintln(Debug, "(")
	return func() {
		indent--
		printIndent()
		fmt.Fprintf(Debug, ")\n")
	}
}
func tracemsg(msg string, args ...interface{}) {
	if Debug == nil {
		return
	}
	printIndent()
	fmt.Fprintf(Debug, msg, args...)
	fmt.Fprintln(Debug)
}
