package aq

import (
	fq "github.com/takanoriyanagitani/go-fsq"
)

type NameChecker func(unchecked string) (checked string)

var NameCheckerNoCheck NameChecker = fq.Identity[string]
