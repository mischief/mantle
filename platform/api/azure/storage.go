package azure

import (
	"github.com/Azure/azure-sdk-for-go/management/osimage"
)

// https://msdn.microsoft.com/en-us/library/azure/jj157192.aspx
func (a *API) AddOSImage(md *osimage.OSImage) error {
	c := osimage.NewClient(a.client)
	op, err := c.AddOSImage(md)
	if err != nil {
		return err
	}

	return a.client.WaitForOperation(op, nil)
}
