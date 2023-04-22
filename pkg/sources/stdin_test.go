package sources_test

import (
	"strings"
	"testing"

	"github.com/ffix/vhtg/pkg/sources"
)

type MockEventHandler struct {
	ProcessedLines []string
}

func (m *MockEventHandler) ProcessLine(line string) {
	m.ProcessedLines = append(m.ProcessedLines, line)
}

func TestProcessInput(t *testing.T) {
	// Create a buffer with test data simulating os.Stdin
	input := "line1\nline2\nline3\n"
	testStdin := strings.NewReader(input)

	// Create a mock event handler to track processed lines
	mockEventHandler := &MockEventHandler{
		ProcessedLines: make([]string, 0),
	}

	// Create a new StdinProcessor and process the input
	processor := sources.NewIOProcessor(testStdin, mockEventHandler)
	processor.Process()

	// Check if the processed lines match the input data
	expectedLines := strings.Split(strings.Trim(input, "\n"), "\n")
	if len(mockEventHandler.ProcessedLines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(mockEventHandler.ProcessedLines))
	}

	for i, line := range mockEventHandler.ProcessedLines {
		if line != expectedLines[i] {
			t.Errorf("Expected line %d to be '%s', got '%s'", i, expectedLines[i], line)
		}
	}
}
