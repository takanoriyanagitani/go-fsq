package aq

type TempnameBuilder func(filename string) (tempname string)

func TempnameBuilderSuffixNew(suffix string) TempnameBuilder {
	return func(filename string) (tmpname string) {
		return filename + suffix
	}
}

var TempnameBuilderSimple TempnameBuilder = TempnameBuilderSuffixNew(".tmp")
