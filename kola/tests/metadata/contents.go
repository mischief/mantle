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

package ignition

import (
	"fmt"
	"strings"

	"github.com/coreos/mantle/kola/register"
	"github.com/coreos/mantle/platform"
)

func init() {
	register.Register(&register.Test{
		Name:        "coreos.metadata.aws",
		Run:         verifyAWS,
		ClusterSize: 1,
		Platforms:   []string{"aws"},
		UserData: `{
		               "ignitionVersion": 1,
		               "systemd": {
		                   "units": [
		                       {
		                           "name": "coreos-metadata.service",
		                           "enable": true
		                       },
		                       {
		                           "name": "metadata.target",
		                           "enable": true,
		                           "contents": "[Install]\nWantedBy=multi-user.target"
		                       }
		                   ]
		               }
		           }`,
	})

	register.Register(&register.Test{
		Name:        "coreos.metadata.azure",
		Run:         verifyAzure,
		ClusterSize: 1,
		Platforms:   []string{"azure"},
		UserData: `{
		               "ignitionVersion": 1,
		               "systemd": {
		                   "units": [
		                       {
		                           "name": "coreos-metadata.service",
		                           "enable": true
		                       },
		                       {
		                           "name": "metadata.target",
		                           "enable": true,
		                           "contents": "[Install]\nWantedBy=multi-user.target"
		                       }
		                   ]
		               }
		           }`,
	})
}

func verifyAWS(c platform.TestCluster) error {
	return verify(c, "COREOS_IPV4_LOCAL", "COREOS_IPV4_PUBLIC", "COREOS_HOSTNAME")
}

func verifyAzure(c platform.TestCluster) error {
	return verify(c, "COREOS_IPV4_LOCAL", "COREOS_IPV4_PUBLIC")
}

func verify(c platform.TestCluster, keys ...string) error {
	m := c.Machines()[0]

	out, err := m.SSH("cat /run/metadata/coreos")
	if err != nil {
		return fmt.Errorf("failed to cat /run/metadata/coreos: %s: %v", out, err)
	}

	for _, key := range keys {
		if !strings.Contains(string(out), key) {
			return fmt.Errorf("%q wasn't found in %q", key, string(out))
		}
	}

	return nil
}
