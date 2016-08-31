// Copyright 2016 CoreOS, Inc.
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

// +build go1.7

package azure

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cmdReplicateImage = &cobra.Command{
		Use:   "replicate-image image offer sku version [regions...]",
		Short: "Replicate an OS image in Azure",
		RunE:  runReplicateImage,
	}

	regions = []string{
		"East US",
		"West US",
		"South Central US",
		"Central US",
		"North Central US",
		"East US 2",
		"North Europe",
		"West Europe",
		"Southeast Asia",
		"East Asia",
		"Japan West",
		"Japan East",
		"Brazil South",
		"Australia Southeast",
		"Australia East",
		"Central India",
		"South India",
		"West India",
		"Canada Central",
		"Canada East",
		"UK North",
		"UK South 2",
		"West US 2",
		"West Central US",
		"UK West",
		"UK South",
		"Central US EUAP",
		"East US 2 EUAP",
	}
)

func init() {
	Azure.AddCommand(cmdReplicateImage)
}

func runReplicateImage(cmd *cobra.Command, args []string) error {
	if len(args) < 4 {
		return fmt.Errorf("expecting at least 4 arguments")
	}

	if len(args) >= 5 {
		regions = args[4:]
	}

	return api.ReplicateImage(args[0], args[1], args[2], args[3], regions...)
}
