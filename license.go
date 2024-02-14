package main

import (
	"fmt"
	"strings"
)

func findLicense(name string) (string, error) {
	name = strings.ToLower(name)

	if _, ok := licenses[name]; ok {
		return name, nil
	}

	for n, matches := range licenses {
		for _, m := range matches {
			if strings.EqualFold(m, name) {
				return n, nil
			}
		}
	}

	return "", fmt.Errorf("no license matching: %s", name)
}

var licenses = map[string][]string{
	"":        {"none"},
	"agpl":    {"agpl-3.0", "affero gpl", "gnu agpl"},
	"apache":  {"apache-2.0", "apache20", "apache 2.0", "apache2.0"},
	"bsd":     {"bsd-3-Clause", "newbsd", "3 clause bsd", "3-clause bsd"},
	"freebsd": {"bsd-2-Clause", "simpbsd", "simple bsd", "2-clause bsd", "2 clause bsd", "simplified bsd license"},
	"gpl2":    {"gpl-2.0", "gnu gpl2", "gplv2"},
	"gpl3":    {"gpl-3.0", "gplv3", "gpl", "gnu gpl3", "gnu gpl"},
	"lgpl":    {"lgpl-3.0", "lesser gpl", "gnu lgpl"},
	"mit":     {"mit"},
}
