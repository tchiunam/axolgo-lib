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
// A Viper object will be returned if optClass is not provided.
// If optClass is provided, Unmarshalling will be performed and
// the class will be updated. No object will be returned.
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
		return nil, v.Unmarshal(&optClass[0])
	}
}
