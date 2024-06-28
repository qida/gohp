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

//go:build windows
// +build windows

package osx

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetCPUID 获取cpuid
func GetCPUID() string {
	var cpuid string
	cmd := exec.Command("wmic", "cpu", "get", "processorid")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "unkonw"
	}
	cpuid = string(b)
	cpuid = cpuid[12 : len(cpuid)-2]
	cpuid = strings.ReplaceAll(cpuid, "\n", "")
	return cpuid
}

// GetBaseBoardID 获取主板的id
func GetBaseBoardID() string {
	var cpuid string
	cmd := exec.Command("wmic", "baseboard", "get", "serialnumber")
	b, e := cmd.CombinedOutput()

	if e == nil {
		cpuid = string(b)
		cpuid = cpuid[12 : len(cpuid)-2]
		cpuid = strings.ReplaceAll(cpuid, "\n", "")
	} else {
		fmt.Printf("%+v", e)
	}

	return cpuid
}
