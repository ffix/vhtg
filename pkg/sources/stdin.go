package sources

import (
	"bufio"
	"io"
	"log"
	"os"
)

type StdinProcessor struct {
	input   io.Reader
	handler EventHandler
}

func NewStdinProcessor(handler EventHandler) *StdinProcessor {
	return NewIOProcessor(os.Stdin, handler)
}

func NewIOProcessor(input io.Reader, handler EventHandler) *StdinProcessor {
	return &StdinProcessor{
		input:   input,
		handler: handler,
	}
}
func (p *StdinProcessor) Process() {
	scanner := bufio.NewScanner(p.input)
	for scanner.Scan() {
		line := scanner.Text()
		p.handler.ProcessLine(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %s", err)
	}
}
