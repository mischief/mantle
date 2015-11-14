// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"runtime"
)

// PortageArch returns a string of the architecture for portage, based on the
// architecture that this go code was compiled for.
func PortageArch() string {
	arch := runtime.GOARCH
	switch arch {
	case "386":
		arch = "x86"

	// Go and Portage agree for these.
	case "amd64":
	case "arm":
	case "arm64":
	case "ppc64":

	// Gentoo doesn't have a little-endian PPC port.
	case "ppc64le":
		fallthrough
	default:
		panic("No portage arch defined for " + arch)
	}
	return arch
}
