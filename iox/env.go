// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package iox

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetGopaths() []string {
	gopath := os.Getenv("GOPATH")
	var paths []string
	if runtime.GOOS == "windows" {
		gopath = strings.Replace(gopath, "\\", "/", -1)
		paths = strings.Split(gopath, ";")
	} else {
		paths = strings.Split(gopath, ":")
	}
	return paths
}

func GetSrcpath(importPath string) (_appPath string, _err error) {
	paths := GetGopaths()
	for _, p := range paths {
		if IsExist(p + "/src/" + importPath + "/") {
			_appPath = p + "/src/" + importPath + "/"
			break
		}
	}
	if len(_appPath) == 0 {
		_err = errors.New("unable to locate source folder path")
		return
	}
	_appPath = filepath.Dir(_appPath) + "/"
	if runtime.GOOS == "windows" {
		// Replace all '\' to '/'.
		_appPath = strings.Replace(_appPath, "\\", "/", -1)
	}
	return
}

func GetHomeDir() (home string, _err error) {
	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if len(home) == 0 {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
	} else {
		home = os.Getenv("HOME")
	}

	if len(home) == 0 {
		_err = errors.New("cannot specify home directory because it's empty")
		return
	}
	return
}
