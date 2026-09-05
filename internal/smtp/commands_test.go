package smtp

import "testing"

func TestParseCommand(t *testing.T) {
	cmd := parseCommand("MAIL FROM:<a@b.com>")
	if cmd.verb != "MAIL" || cmd.args != "FROM:<a@b.com>" {
		t.Fatalf("got %+v", cmd)
	}
}

func TestParseAddress(t *testing.T) {
	cases := []struct {
		prefix, args, want string
	}{
		{"FROM:", "FROM:<a@b.com>", "a@b.com"},
		{"TO:", "TO:<a@b.com> SIZE=100", "a@b.com"},
		{"FROM:", "FROM:a@b.com", "a@b.com"},
	}
	for _, c := range cases {
		got, err := parseAddress(c.prefix, c.args)
		if err != nil {
			t.Fatalf("parseAddress(%q, %q): %v", c.prefix, c.args, err)
		}
		if got != c.want {
			t.Errorf("parseAddress(%q, %q) = %q, want %q", c.prefix, c.args, got, c.want)
		}
	}
}

func TestParseAddressRejectsMissingPrefix(t *testing.T) {
	if _, err := parseAddress("FROM:", "TO:<a@b.com>"); err == nil {
		t.Fatal("expected error for mismatched prefix")
	}
}
