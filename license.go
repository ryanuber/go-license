package license

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// Recognized license types
	LicenseMIT      = "MIT"
	LicenseBSD      = "BSD"
	LicenseNewBSD   = "NewBSD"
	LicenseFreeBSD  = "FreeBSD"
	LicenseApache20 = "Apache-2.0"
	LicenseMPL20    = "MPL-2.0"
	LicenseGPL20    = "GPL-2.0"
	LicenseGPL30    = "GPL-3.0"
	LicenseLGPL21   = "LGPL-2.1"
	LicenseLGPL30   = "LGPL-3.0"
	LicenseCDDL10   = "CDDL-1.0"
	LicenseEPL10    = "EPL-1.0"

	// Various errors
	ErrNoLicenseFile       = "license: unable to find any license file"
	ErrUnrecognizedLicense = "license: could not guess license type"
)

// A set of reasonable license file names to use when guessing where the
// license may be during a crate pack.
var DefaultLicenseFiles = []string{
	"LICENSE", "LICENSE.txt", "LICENSE.md", "license.txt",
	"COPYING", "COPYING.txt", "COPYING.md", "copying.txt",
}

// A slice of standardized license abbreviations
var KnownLicenses = []string{
	LicenseMIT,
	LicenseBSD,
	LicenseNewBSD,
	LicenseFreeBSD,
	LicenseApache20,
	LicenseMPL20,
	LicenseGPL20,
	LicenseGPL30,
	LicenseLGPL21,
	LicenseLGPL30,
	LicenseCDDL10,
	LicenseEPL10,
}

// License describes a software license
type License struct {
	Type string // The type of license in use
	Text string // License text data
}

// New creates a new License from explicitly passed license type and data
func New(licenseType, licenseText string) *License {
	l := &License{
		Type: licenseType,
		Text: licenseText,
	}
	return l
}

// NewFromFile will attempt to load a license from a file on disk, and guess the
// type of license based on the bytes read.
func NewFromFile(path string) (*License, error) {
	licenseText, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	l := &License{Text: string(licenseText)}

	licenseType, err := l.guessType()
	if err != nil {
		return nil, err
	}
	l.Type = licenseType

	return l, nil
}

// NewFromDir will search a directory for well-known and accepted license file
// names, and if one is found, read in its content and guess the license type.
func NewFromDir(dir string) (*License, error) {
	fileName, err := guessFile(dir)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(dir, fileName)
	return NewFromFile(filePath)
}

// Recognized determines if the license is known to go-license.
func (l *License) Recognized() bool {
	for _, license := range KnownLicenses {
		if license == l.Type {
			return true
		}
	}
	return false
}

// guessFile searches a given directory (non-recursively) for files with well-
// established names that indicate license content.
func guessFile(dir string) (string, error) {
	d, err := os.Stat(dir)
	if err != nil {
		return "", err
	}

	if !d.IsDir() {
		return "", fmt.Errorf("license: cannot search %s: not a directory", dir)
	}

	for _, file := range DefaultLicenseFiles {
		filePath := filepath.Join(dir, file)
		_, err := os.Stat(filePath)
		if err == nil {
			return file, nil
		}
	}
	return "", errors.New(ErrNoLicenseFile)
}

// guessType will scan license text and attempt to guess what license type it
// describes. It will return the license type on success, or an error if it
// cannot accurately guess the license type.
func (l *License) guessType() (string, error) {
	switch {
	case scanLeft(l.Text, "The MIT License"):
		return LicenseMIT, nil

	case scanLeft(l.Text, "Apache License"):
		switch {
		case scanLeft(l.Text, "Version 2.0"):
			return LicenseApache20, nil
		}

	case scanLeft(l.Text, "GNU GENERAL PUBLIC LICENSE"):
		switch {
		case scanLeft(l.Text, "Version 2"):
			return LicenseGPL20, nil
		case scanLeft(l.Text, "Version 3"):
			return LicenseGPL30, nil
		}

	case scanLeft(l.Text, "GNU LESSER GENERAL PUBLIC LICENSE"):
		switch {
		case scanLeft(l.Text, "Version 2.1"):
			return LicenseLGPL21, nil
		case scanLeft(l.Text, "Version 3"):
			return LicenseLGPL30, nil
		}

	case scanLeft(l.Text, "Mozilla Public License Version 2.0"):
		return LicenseMPL20, nil

	case scanLeft(l.Text, "Redistribution and use in source and binary forms"):
		switch {
		case scanLeft(l.Text, "4. Neither"):
			return LicenseBSD, nil
		case scanLeft(l.Text, "* Redistribution"):
			return LicenseNewBSD, nil
		case scanRight(l.Text, "FreeBSD Project."):
			return LicenseFreeBSD, nil
		}

	case scanRight(l.Text, "(CDDL)"):
		switch {
		case scanRight(l.Text, "Version 1.0"):
			return LicenseCDDL10, nil
		}

	case scanLeft(l.Text, "Eclipse Public License - v 1.0"):
		return LicenseEPL10, nil
	}

	return "", errors.New(ErrUnrecognizedLicense)
}
