package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	files, err := os.ReadDir("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			t.Run(file.Name(), func(t *testing.T) {
				var buf bytes.Buffer

				inputFile := "./testdata/" + file.Name()
				err := run(inputFile, &buf)
				if err != nil {
					t.Fatal(err)
				}

				outputFile := "./testdata/" + strings.TrimSuffix(file.Name(), ".txt") + ".out"
				expected, err := os.ReadFile(outputFile)
				if err != nil {
					t.Fatal(err)
				}

				if !bytes.Equal(buf.Bytes(), expected) {
					t.Errorf("Output for %s does not match expected output", inputFile)
				}
			})
		}
	}
}
