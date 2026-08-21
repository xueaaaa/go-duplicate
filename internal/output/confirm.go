package output

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// PlainConfirm writes confirmationText to w followed by " [y/N] ", then
// reads a single line from r and reports whether the answer means yes.
// Only "y" and "yes" (case-insensitive, surrounding whitespace trimmed)
// count as yes; any other input, including an empty line (just pressing
// Enter) or no input at all (r is already at EOF), is treated as no and
// does not produce an error.
//
// Any error writing to w, or any read error from r other than a clean EOF,
// is returned as-is; in that case the returned bool has no meaningful
// value. PlainConfirm reads via [bufio.Scanner], so a single "line" longer
// than [bufio.MaxScanTokenSize] will fail with [bufio.ErrTooLong] — not a
// concern for typical interactive input, but worth knowing if r is backed
// by something other than a terminal.
func PlainConfirm(w io.Writer, r io.Reader, confirmationText string) (bool, error) {
	if _, err := fmt.Fprintf(w, "%s [y/N] ", confirmationText); err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}

		return false, nil
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}
