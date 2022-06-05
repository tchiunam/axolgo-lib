package ioutil

import (
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
)

// Read the file given by filepath as an INI file
func ReadIniFile(filepath string) (*ini.File, error) {
	return ini.Load(filepath)
}

// Read the file given by filepath as a YAML file. Viper instance
// is returned if optClass is not provided. Otherwise Unmarshalling
// will be performed and the object will be returned.
func ReadYamlFile(filepath string, optClass ...interface{}) (interface{}, error) {
	v := viper.New()
	v.SetConfigFile(filepath)
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if optClass == nil {
		return v, nil
	} else {
		err = v.Unmarshal(&optClass)
		return optClass, err
	}
}
