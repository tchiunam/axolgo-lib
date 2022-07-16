/*
Copyright © 2022 tchiunam

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package ioutil

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
)

// Read the file given by filepath as an INI file
func ReadIniFile(filepath string) (*ini.File, error) {
	return ini.Load(filepath)
}

// Supported options for reading config files
type ConfigFileOptions struct {
	Type  string
	Class interface{}
}

// ConfigFileOptionsFunc is a type alias for ConfigFileOptions functional option
type ConfigFileOptionsFunc func(f *ConfigFileOptions) error

// WithCFOType is a helper function to construct functional options
// that sets config file type on config's ConfigFileOptions.
func WithCFOType(fileType string) ConfigFileOptionsFunc {
	return func(f *ConfigFileOptions) error {
		f.Type = fileType
		return nil
	}
}

// WithCFOClass is a helper function to construct functional options
// that sets class for unmarshalling on config's ConfigFileOptions.
func WithCFOClass(class interface{}) ConfigFileOptionsFunc {
	return func(f *ConfigFileOptions) error {
		f.Class = class
		return nil
	}
}

// Read the file given by filepath. Viper instance is returned
// if optClass is not provided. Otherwise Unmarshalling will be
// performed.
// If optClass is provided, Unmarshalling will be performed and
// the class will be updated. No object will be returned.
// optClass is expected to be a pointer to a struct.
//
// For example:
// 	type Config struct {
// 		Name string
// 	}
// 	var config Config
// 	err := ReadConfigFile("config.yaml", WithClass(&config))
// 	if err != nil {
//      doSomethingWithError(err)
// 	}
func ReadConfigFile(filepath string, optFns ...func(*ConfigFileOptions) error) (interface{}, error) {
	options := ConfigFileOptions{}

	for _, optFn := range optFns {
		if err := optFn(&options); err != nil {
			return nil, fmt.Errorf("Fail to configure ReadConfigFile options: %v", err)
		}
	}

	// Read the file using Viper
	v := viper.New()
	v.SetConfigFile(filepath)
	if options.Type != "" {
		v.SetConfigType(options.Type)
	}
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if options.Class == nil {
		return v, nil
	} else {
		// Unmarshal the config into the class
		return nil, v.Unmarshal(&options.Class)
	}
}

// Read the file given by filepath as a YAML file
func ReadYamlFile(filepath string, optFns ...func(*ConfigFileOptions) error) (interface{}, error) {
	// Add the last option to read the file as YAML
	// which is the default behavior of ReadYAMLFile
	optFns = append(optFns, WithCFOType("yaml"))

	return ReadConfigFile(filepath, optFns...)
}
