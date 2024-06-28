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

//go:build linux
// +build linux

package osx

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

// GetCPUID 获取cpuid
func GetCPUID() string {
	cmd := exec.Command("/bin/sh", "-c", `sudo dmidecode -t 4 | grep ID `)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("StdoutPipe: " + err.Error())
		return ""
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe: ", err.Error())
		return ""
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("Start: ", err.Error())
		return ""
	}
	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		fmt.Println("ReadAll stderr: ", err.Error())
		return ""
	}
	if len(bytesErr) != 0 {
		fmt.Printf("stderr is not nil: %s", bytesErr)
		return ""
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll stdout: ", err.Error())
		return ""
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait: ", err.Error())
		return ""
	}
	return string(bytes)
}

// GetBaseBoardID 获取主板的id
func GetBaseBoardID() string {
	var cpuid string

	return cpuid
}
