package smtp

import (
	"fmt"
	"strings"
)

// command is one parsed SMTP protocol line: a verb plus its raw arguments.
type command struct {
	verb string
	args string
}

func parseCommand(line string) command {
	line = strings.TrimSpace(line)
	parts := strings.SplitN(line, " ", 2)

	cmd := command{verb: strings.ToUpper(parts[0])}
	if len(parts) == 2 {
		cmd.args = parts[1]
	}
	return cmd
}

// parseAddress extracts the email address from a MAIL FROM:<addr> or
// RCPT TO:<addr> argument, ignoring any trailing ESMTP parameters.
func parseAddress(prefix, args string) (string, error) {
	args = strings.TrimSpace(args)
	upper := strings.ToUpper(args)
	if !strings.HasPrefix(upper, prefix) {
		return "", fmt.Errorf("expected %s", prefix)
	}
	rest := strings.TrimSpace(args[len(prefix):])

	start := strings.Index(rest, "<")
	end := strings.Index(rest, ">")
	if start == -1 || end == -1 || end < start {
		// Some clients omit the angle brackets; take the first token.
		fields := strings.Fields(rest)
		if len(fields) == 0 {
			return "", fmt.Errorf("missing address")
		}
		return fields[0], nil
	}

	addr := strings.TrimSpace(rest[start+1 : end])
	if addr == "" {
		return "", fmt.Errorf("empty address")
	}
	return addr, nil
}
