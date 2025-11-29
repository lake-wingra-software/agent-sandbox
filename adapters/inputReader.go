package adapters

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ParseQuery() string {
	query, err := readStdin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read stdin:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(query) == "" {
		fmt.Fprintln(os.Stderr, "no input provided on stdin; pipe a prompt, e.g.: echo 'Hello' | go run main.go")
		os.Exit(2)
	}
	return query
}

// Reads all data from stdin. Returns empty string if there's no piped input.
func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	// If input is from a pipe or file, ModeCharDevice will be false
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}
	var b strings.Builder
	r := bufio.NewReader(os.Stdin)
	for {
		chunk, err := r.ReadString('\n')
		b.WriteString(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(b.String()), nil
}
