package license

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLicense(t *testing.T) {
	l := New("MyLicense", "Some license text.")
	if l.Type != "MyLicense" {
		t.Fatalf("bad license type: %s", l.Type)
	}
	if l.Text != "Some license text." {
		t.Fatalf("bad license text: %s", l.Text)
	}
}

func TestNewFromFile(t *testing.T) {
	f, err := ioutil.TempFile("", "go-license")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Remove(f.Name())

	licenseText := "The MIT License (MIT)"

	if _, err := f.WriteString(licenseText); err != nil {
		t.Fatalf("err: %s", err)
	}

	l, err := NewFromFile(f.Name())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if l.Type != "MIT" {
		t.Fatalf("unexpected license type: %s", l.Type)
	}

	if l.Text != licenseText {
		t.Fatalf("unexpected license text: %s", l.Text)
	}
}

func TestNewFromDir(t *testing.T) {
	d, err := ioutil.TempDir("", "go-license")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(d)

	fPath := filepath.Join(d, "LICENSE")
	f, err := os.Create(fPath)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	licenseText := "The MIT License (MIT)"

	if _, err := f.WriteString(licenseText); err != nil {
		t.Fatalf("err: %s", err)
	}

	l, err := NewFromDir(d)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if l.Type != "MIT" {
		t.Fatalf("unexpected license type: %s", l.Type)
	}

	if l.Text != licenseText {
		t.Fatalf("unexpected license text: %s", l.Text)
	}
}

func TestLicenseRecognized(t *testing.T) {
	// Known licenses are recognized
	l := New("MIT", "The MIT License (MIT)")
	if !l.Recognized() {
		t.Fatalf("license was not recognized")
	}

	// Unknown licenses are not recognized
	l = New("None", "No license text")
	if l.Recognized() {
		t.Fatalf("fake license was recognized")
	}
}
