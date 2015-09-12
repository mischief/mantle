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

package misc

import (
	"bytes"
	"fmt"
	"path"
	"time"

	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/util"

	"github.com/coreos/mantle/Godeps/_workspace/src/github.com/coreos/coreos-cloudinit/config"
	"github.com/coreos/mantle/Godeps/_workspace/src/github.com/coreos/pkg/capnslog"
)

var plog = capnslog.NewPackageLogger("github.com/coreos/mantle", "kola/tests/misc")

/*
core@nfs1 /usr/lib/systemd/system $ ls rpc*
rpc-statd-notify.service  rpc-statd.service  rpcbind.service  rpcbind.target
core@nfs1 /usr/lib/systemd/system $ ls nfs*
nfs-blkmap.service  nfs-blkmap.target  nfs-client.target  nfs-idmapd.service  nfs-mountd.service  nfs-server.service  nfs-utils.service
*/

// Test that the kernel NFS server and client work within CoreOS.
func NFS(c platform.TestCluster) error {
	/* server machine */
	c1 := config.CloudConfig{
		CoreOS: config.CoreOS{
			Units: []config.Unit{
				config.Unit{
					Name:    "rpcbind.service",
					Command: "start",
				},
				config.Unit{
					Name:    "rpc-statd.service",
					Command: "start",
				},
				config.Unit{
					Name:    "nfs-mountd.service",
					Command: "start",
				},
				config.Unit{
					Name:    "nfs-server.service",
					Command: "start",
				},
			},
		},
		WriteFiles: []config.File{
			config.File{
				Content: "/tmp	*(ro,insecure,all_squash,no_subtree_check,fsid=0)",
				Path: "/etc/exports",
			},
		},
		Hostname: "nfs1",
	}

	m1, err := c.NewMachine(c1.String())
	if err != nil {
		return fmt.Errorf("Cluster.NewMachine: %s", err)
	}

	defer m1.Destroy()

	plog.Info("NFS server booted.")

	/* poke a file in /tmp */
	tmp, err := m1.SSH("mktemp")
	if err != nil {
		return fmt.Errorf("Machine.SSH: %s", err)
	}

	plog.Infof("Test file %q created on server.", tmp)

	/* client machine */

	nfstmpl := `[Unit]
Description=NFS Client
After=network-online.target
Requires=network-online.target
After=rpc-statd.service
Requires=rpc-statd.service

[Mount]
What=%s:/tmp
Where=/mnt
Type=nfs
Options=defaults,noexec
`

	c2 := config.CloudConfig{
		CoreOS: config.CoreOS{
			Units: []config.Unit{
				config.Unit{
					Name:    "rpcbind.service",
					Command: "start",
				},
				config.Unit{
					Name:    "rpc-statd.service",
					Command: "start",
				},
				config.Unit{
					Name:    "mnt.mount",
					Command: "start",
					Content: fmt.Sprintf(nfstmpl, m1.PrivateIP()),
				},
			},
		},
		Hostname: "nfs2",
	}

	m2, err := c.NewMachine(c2.String())
	if err != nil {
		return fmt.Errorf("Cluster.NewMachine: %s", err)
	}

	defer m2.Destroy()

	plog.Info("NFS client booted.")

	var lsmnt []byte

	plog.Info("Waiting for NFS mount on client...")

	/* there's probably a better wait to check the mount */
	checker := func() error {
		lsmnt, _ = m2.SSH("ls /mnt")
		if len(lsmnt) == 0 {
			return fmt.Errorf("client /mnt is empty")
		}

		plog.Info("Got NFS mount.")
		return nil
	}

	if err = util.Retry(5, 1*time.Second, checker); err != nil {
		return err
	}

	if len(lsmnt) == 0 {
		return fmt.Errorf("Client /mnt is empty.")
	}

	if bytes.Contains(lsmnt, []byte(path.Base(string(tmp)))) != true {
		return fmt.Errorf("Client /mnt did not contain file %q from server /tmp -- /mnt: %s", tmp, lsmnt)
	}

	return nil
}
