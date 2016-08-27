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

package azure

import (
	"fmt"
	"strings"

	"github.com/Microsoft/azure-vhd-utils-for-go/vhdcore/validator"
	"github.com/spf13/cobra"
)

var (
	cmdUploadBlob = &cobra.Command{
		Use:   "upload-blob [file]",
		Short: "Upload a blob to Azure storage",
		RunE:  runUploadBlob,
	}

	// upload blob options
	ubo struct {
		vhd       string
		container string
		blob      string
		overwrite bool
		validate  bool
	}
)

func init() {
	sv := cmdUploadBlob.Flags().StringVar
	bv := cmdUploadBlob.Flags().BoolVar

	sv(&ubo.container, "container", "", "container name")
	sv(&ubo.blob, "blob", "", "blob name")
	bv(&ubo.overwrite, "overwrite", false, "overwrite blob")
	bv(&ubo.validate, "validate", true, "validate blob as VHD file")

	Azure.AddCommand(cmdUploadBlob)
}

func runUploadBlob(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expecting one file as argument, got %d", len(args))
	}

	ubo.vhd = args[0]

	if ubo.validate {
		plog.Printf("Validating VHD %q", ubo.vhd)
		if !strings.HasSuffix(strings.ToLower(ubo.blob), ".vhd") {
			return fmt.Errorf("blob name should end with .vhd")
		}

		if err := validator.ValidateVhd(ubo.vhd); err != nil {
			return err
		}

		if err := validator.ValidateVhdSize(ubo.vhd); err != nil {
			return err
		}
	}

	err := api.UploadBlob(ubo.vhd, ubo.container, ubo.blob, ubo.overwrite)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", opts.StorageAccountName, ubo.container, ubo.blob)

	plog.Printf("blob uploaded to %q", uri)

	return nil
}
