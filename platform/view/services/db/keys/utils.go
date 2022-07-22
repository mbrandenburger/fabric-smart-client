/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package keys

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/pkg/errors"
)

var (
	nsRegexp = regexp.MustCompile("^[a-zA-Z0-9._-]{1,128}$")
)

const NamespaceSeparator = "\u0000"

func ValidateKey(key string) error {
	// TODO: should we enforce a length limit?
	if !utf8.ValidString(key) {
		return fmt.Errorf("not a valid utf8 string: [%x]", key)
	}

	return nil
}

func ValidateNs(ns string) error {
	if !nsRegexp.MatchString(ns) {
		return errors.Errorf("namespace '%s' is invalid", ns)
	}

	return nil
}
