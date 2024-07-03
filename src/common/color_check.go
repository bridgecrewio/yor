package common

type ColorStruct struct {
	NoColor bool
	Reset   string
	Green   string
	Yellow  string
	Blue    string
	Purple  string
}

func NoColorCheck(noColorBool bool) *ColorStruct {
	var colors ColorStruct
	colors = ColorStruct{
		NoColor: true,
		Reset:   "",
		Green:   "",
		Yellow:  "",
		Blue:    "",
		Purple:  "",
	}
	if !noColorBool {
		colors = ColorStruct{
			NoColor: false,
			Reset:   "\033[0m",
			Green:   "\033[32m",
			Yellow:  "\033[33m",
			Blue:    "\033[34m",
			Purple:  "\033[35m",
		}
	}
	return &colors
}
