package license

import "testing"

func TestScan(t *testing.T) {
	text := "She sells sea shells by the sea shore"

	// scan matches properly
	scanPass := []string{
		"^she sells sea shells",
		".*by the sea shore$",
		"she sells sea shells by the sea shore",
		"she sells SEA SHELLS by the sea.*$",
	}
	for _, s := range scanPass {
		if !scan(text, s) {
			t.Fatalf("%s did not match during scan", s)
		}
	}

	// scan rejects messages that shouldn't match
	scanFail := []string{
		"^by the sea shore",
		"she sells sea shells$",
	}
	for _, s := range scanFail {
		if scan(text, s) {
			t.Fatalf("%s matched during scan", s)
		}
	}
}
