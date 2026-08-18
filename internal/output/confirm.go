package output

import (
	"fmt"
	"io"
	"strings"
)

func PlainConfirm(w io.Writer, r io.Reader, confirmationText string) (bool, error) {
	if _, err := fmt.Fprintf(w, "%s [y/N]", confirmationText); err != nil {
		return false, err
	}

	var answer string
	if _, err := fmt.Fscanln(r, &answer); err != nil {
		return false, err
	}

	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes", nil
}
