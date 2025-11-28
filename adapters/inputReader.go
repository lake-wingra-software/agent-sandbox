package adapters

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Reads all data from stdin. Returns empty string if there's no piped input.
func ReadStdin() (string, error) {
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
