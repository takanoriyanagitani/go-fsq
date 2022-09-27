package aq

import (
	fq "github.com/takanoriyanagitani/go-fsq"
)

type NameChecker func(unchecked string) (checked string)

// NameCheckerNoCheck is NameChecker which does not check name.
// Use this for trusted user input(filename).
var NameCheckerNoCheck NameChecker = fq.Identity[string]
