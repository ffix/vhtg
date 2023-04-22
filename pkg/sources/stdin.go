package sources

import (
	"bufio"
	"log"
	"os"
)

func ProcessStdin(handler EventHandler) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		handler.ProcessLine(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %s", err)
	}
}
