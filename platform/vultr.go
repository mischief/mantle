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

package platform

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/coreos/mantle/network"
	"github.com/coreos/mantle/system/exec"

	vultr "github.com/JamesClonk/vultr/lib"
	"github.com/coreos/mantle/Godeps/_workspace/src/golang.org/x/crypto/ssh"
)

type vultrMachine struct {
	cluster *vultrCluster
	info    vultr.Server
}

func (vm *vultrMachine) ID() string {
	return vm.info.ID
}

func (vm *vultrMachine) IP() string {
	return vm.info.MainIP
}

func (vm *vultrMachine) PrivateIP() string {
	return vm.info.InternalIP
}

func (vm *vultrMachine) SSHClient() (*ssh.Client, error) {
	sshClient, err := vm.cluster.agent.NewClient(vm.IP())
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}

func (vm *vultrMachine) SSH(cmd string) ([]byte, error) {
	client, err := vm.SSHClient()
	if err != nil {
		return nil, err
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	defer session.Close()

	session.Stderr = os.Stderr
	out, err := session.Output(cmd)
	out = bytes.TrimSpace(out)
	return out, err
}

func (vm *vultrMachine) Destroy() error {
	err := vm.cluster.api.DeleteServer(vm.ID())
	if err != nil {
		return err
	}

	vm.cluster.delMach(vm)
	return nil
}

// VultrOptions contains Vultr-specific instance options.
type VultrOptions struct {
	APIKey string
}

type vultrCluster struct {
	mu    sync.Mutex
	api   *vultr.Client
	agent *network.SSHAgent
	machs map[string]*vultrMachine
}

func NewVultrCluster(conf VultrOptions) (Cluster, error) {
	api := vultr.NewClient(conf.APIKey, &vultr.Options{RateLimitation: 1 * time.Second})

	agent, err := network.NewSSHAgent(&net.Dialer{})
	if err != nil {
		return nil, err
	}

	vc := &vultrCluster{
		api:   api,
		agent: agent,
		machs: make(map[string]*vultrMachine),
	}

	return vc, nil
}

func (vc *vultrCluster) addMach(m *vultrMachine) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.machs[m.ID()] = m
}

func (vc *vultrCluster) delMach(m *vultrMachine) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	delete(vc.machs, m.ID())
}

func (vc *vultrCluster) NewCommand(name string, arg ...string) exec.Cmd {
	return exec.Command(name, arg...)
}

func (vc *vultrCluster) NewMachine(userdata string) (Machine, error) {
	conf, err := NewConf(userdata)
	if err != nil {
		return nil, err
	}

	keys, err := vc.agent.List()
	if err != nil {
		return nil, err
	}

	conf.CopyKeys(keys)

	// TODO(mischief): upload iPXE/iso and actually spawn a machine
	//options := &vultr.ServerOptions{}

	mach := &vultrMachine{
		cluster: vc,
	}

	if err := commonMachineChecks(mach); err != nil {
		return nil, fmt.Errorf("machine %q failed basic checks: %v", mach.ID(), err)
	}

	vc.addMach(mach)

	return mach, nil
}

func (vc *vultrCluster) Machines() []Machine {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	machs := make([]Machine, 0, len(vc.machs))
	for _, m := range vc.machs {
		machs = append(machs, m)
	}
	return machs
}

func (vc *vultrCluster) EtcdEndpoint() string {
	return ""
}

func (vc *vultrCluster) GetDiscoveryURL(size int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://discovery.etcd.io/new?size=%d", size))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (vc *vultrCluster) Destroy() error {
	machs := vc.Machines()
	for _, vm := range machs {
		vm.Destroy()
	}
	vc.agent.Close()
	return nil
}
