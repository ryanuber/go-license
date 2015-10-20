package license

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// Recognized license types
	LicenseMIT       = "MIT"
	LicenseISC       = "ISC"
	LicenseNewBSD    = "NewBSD"
	LicenseFreeBSD   = "FreeBSD"
	LicenseApache20  = "Apache-2.0"
	LicenseMPL20     = "MPL-2.0"
	LicenseGPL20     = "GPL-2.0"
	LicenseGPL30     = "GPL-3.0"
	LicenseLGPL21    = "LGPL-2.1"
	LicenseLGPL30    = "LGPL-3.0"
	LicenseAGPL30    = "AGPL-3.0"
	LicenseCDDL10    = "CDDL-1.0"
	LicenseEPL10     = "EPL-1.0"
	LicenseUnlicense = "Unlicense"

	// Various error messages.
	ErrNoLicenseFile       = "license: unable to find any license file"
	ErrUnrecognizedLicense = "license: could not guess license type"
	ErrMultipleLicenses    = "license: multiple license files found"
)

var (
	// Base names of guessable license files.
	fileNames = []string{
		"copying",
		"copyleft",
		"copyright",
		"license",
		"unlicense",
	}

	// License file extensions. Combined with the licenseFiles slice
	// to create a set of files we can reasonably assume contain
	// licensing information.
	fileExtensions = []string{
		"",
		".md",
		".rst",
		".txt",
	}

	// Lookup tables used for license file names and license types. We
	// use a poor man's set here to get O(1) lookups.
	fileTable    map[string]struct{}
	licenseTable map[string]struct{}

	// Regular expressions used for normalizing license text.
	newlineRegexp = regexp.MustCompile("(\r\n|\n)")
	spaceRegexp   = regexp.MustCompile("\\s{2,}")
)

// init allocates substructures
func init() {
	// Generate the list of known file names.
	size := len(fileNames) * len(fileExtensions)
	fileTable = make(map[string]struct{}, size)
	for _, file := range fileNames {
		for _, ext := range fileExtensions {
			fileTable[file+ext] = struct{}{}
		}
	}

	// Initialize the license types.
	licenseTable = make(map[string]struct{})
	for _, l := range []string{
		LicenseMIT,
		LicenseNewBSD,
		LicenseFreeBSD,
		LicenseApache20,
		LicenseMPL20,
		LicenseGPL20,
		LicenseGPL30,
		LicenseLGPL21,
		LicenseLGPL30,
		LicenseAGPL30,
		LicenseCDDL10,
		LicenseEPL10,
		LicenseUnlicense,
	} {
		licenseTable[l] = struct{}{}
	}
}

// LicenseFilesInDir will scan the given directory for files which match our
// list of known license file names.
func LicenseFilesInDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, fi := range files {
		name := fi.Name()
		lower := strings.ToLower(name)
		if _, ok := fileTable[lower]; ok {
			out = append(out, name)
		}
	}
	return out, nil
}

// License describes a software license
type License struct {
	Type string // The type of license in use
	Text string // License text data
	File string // The path to the source file, if any
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

	l := &License{
		Text: string(licenseText),
		File: path,
	}

	if err := l.GuessType(); err != nil {
		return nil, err
	}

	return l, nil
}

// NewFromDir will search a directory for well-known and accepted license file
// names, and if one is found, read in its content and guess the license type.
func NewFromDir(dir string) (*License, error) {
	file, err := GuessFile(dir)
	if err != nil {
		return nil, err
	}

	return NewFromFile(file)
}

// Recognized determines if the license is known to go-license.
func (l *License) Recognized() bool {
	_, ok := licenseTable[l.Type]
	return ok
}

// GuessFile searches a given directory (non-recursively) for files with well-
// established names that indicate license content.
func GuessFile(dir string) (string, error) {
	files, err := LicenseFilesInDir(dir)
	if err != nil {
		return "", err
	}

	switch len(files) {
	case 0:
		return "", errors.New(ErrNoLicenseFile)
	case 1:
		return filepath.Join(dir, files[0]), nil
	default:
		return "", errors.New(ErrMultipleLicenses)
	}
}

// GuessType will scan license text and attempt to guess what license type it
// describes. It will return the license type on success, or an error if it
// cannot accurately guess the license type.
//
// This method is a hack. It might be more accurate to also scan the entire body
// of license text and compare it using an algorithm like Jaro-Winkler or
// Levenshtein against a generic version. The problem is that some of the
// common licenses, such as GPL-family licenses, are quite large, and running
// these algorithms against them is considerably more expensive and is still not
// completely deterministic on which license is in play. For now, we will just
// scan until we find differentiating strings and call that good-enuf.gov.
func (l *License) GuessType() error {
	// Lower case everything to make comparison more adaptable.
	comp := strings.ToLower(l.Text)

	// Kill the newlines, since it is not clear if the provided license will
	// contain them or not, and either way it does not change the terms of the
	// license, so one is not "more correct" than the other. This just replaces
	// them with spaces.
	comp = newlineRegexp.ReplaceAllLiteralString(comp, " ")

	// Collapse all instances of multiple spaces into a single space. This
	// makes it simple to express license grammar without worrying about
	// the exact space matching.
	comp = spaceRegexp.ReplaceAllLiteralString(comp, " ")

	switch {
	case scan(comp, "permission is hereby granted, free of charge, to any "+
		"person obtaining a copy of this software"):
		l.Type = LicenseMIT

	case scan(comp, "permission to use, copy, modify, and/or distribute this "+
		"software for any"):
		l.Type = LicenseISC

	case scan(comp, "apache license version 2.0, january 2004") ||
		scan(comp, "http://www.apache.org/licenses/license-2.0"):
		l.Type = LicenseApache20

	case scan(comp, "gnu general public license version 2, june 1991"):
		l.Type = LicenseGPL20

	case scan(comp, "gnu general public license version 3, 29 june 2007"):
		l.Type = LicenseGPL30

	case scan(comp, "gnu lesser general public license version 2.1, "+
		"february 1999"):
		l.Type = LicenseLGPL21

	case scan(comp, "gnu lesser general public license version 3, "+
		"29 june 2007"):
		l.Type = LicenseLGPL30

	case scan(comp, "gnu affero general public license "+
		"version 3, 19 november 2007"):
		l.Type = LicenseAGPL30

	case scan(comp, "mozilla public license") && scan(comp, "version 2.0"):
		l.Type = LicenseMPL20

	case scan(comp, "redistribution and use in source and binary forms"):
		switch {
		case scan(comp, "neither the name of"):
			l.Type = LicenseNewBSD
		default:
			l.Type = LicenseFreeBSD
		}

	case scan(comp, "common development and distribution license (cddl) "+
		"version 1.0"):
		l.Type = LicenseCDDL10

	case scan(comp, "eclipse public license - v 1.0"):
		l.Type = LicenseEPL10

	case scan(comp, "this is free and unencumbered software released into "+
		"the public domain"):
		l.Type = LicenseUnlicense

	default:
		return errors.New(ErrUnrecognizedLicense)
	}

	return nil
}

// scan is used to find substrings. It type-casts to byte slices because
// bytes is an order of magnitude faster than its strings counterpart.
func scan(text, pattern string) bool {
	return bytes.Contains([]byte(text), []byte(pattern))
}
