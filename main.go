package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type (
	Evaluator map[string]float64
)

func NewEvaluator() Evaluator {
	return make(Evaluator)
}

func (e Evaluator) Eval(input string) string {
	l := NewLexer(input, 10).Run()

	root, err := NewParser(l.out, e).Parse()
	if err != nil {
		return err.Error()
	}

	val, err := root.Eval()
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("= %f", val)
}

func (e Evaluator) String() string {
	b := strings.Builder{}
	b.WriteString("Values:")
	for k, v := range e {
		b.WriteString(fmt.Sprintf(" %s=%f", k, v))
	}
	return b.String()
}

func (e Evaluator) PrettyString() string {
	b := strings.Builder{}
	b.WriteString("Values:\n")
	for k, v := range e {
		b.WriteString(fmt.Sprintf(" %s = %f\n", k, v))
	}

	return b.String()
}

func main() {
	e := NewEvaluator()

	// check for one shot mode
	if len(os.Args) > 1 {
		maxLen := 0
		for _, arg := range os.Args[1:] {
			if argL := len(arg); argL > maxLen {
				maxLen = argL
			}
		}
		for _, arg := range os.Args[1:] {
			fmt.Printf("%-*s %s\n", maxLen, arg, e.Eval(arg))
		}
		return
	}

	// do interactive mode
	fmt.Println("Type 'exit' followed by pressing return/enter to exit the program.")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		switch {
		case input == "exit":
			fmt.Println("exiting...")
			return
		case input == "clear":
			e = NewEvaluator()
		case input == "list":
			fmt.Println(e)
		case input == "list-pretty":
			fmt.Println(e.PrettyString())
		// case: TODO input starts with precision, followed by number, update precision for ouputs
		default:
			fmt.Println(e.Eval(input))
		}
	}
}
