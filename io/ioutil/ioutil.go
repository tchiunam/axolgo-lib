package ioutil

import (
	"gopkg.in/ini.v1"
)

// Read the file given by filepath as an INI file
func ReadIniFile(filepath string) (*ini.File, error) {
	return ini.Load(filepath)
}
