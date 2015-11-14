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

// Package exec is extension of the standard os.exec package.
// Adds a handy dandy interface and assorted other features.
package exec

import (
	"io"
	"os/exec"
	"syscall"
)

var (
	// ErrNotFound is the error resulting if a path search failed to find
	// an executable file.
	ErrNotFound = exec.ErrNotFound

	// LookPath searches for an executable binary named file in the
	// directories named by the PATH environment variable. If file contains
	// a slash, it is tried directly and the PATH is not consulted. The
	// result may be an absolute path or a path relative to the current
	// directory.
	LookPath = exec.LookPath
)

// Cmd is an exec.Cmd compatible interface.
type Cmd interface {
	// Methods provided by exec.Cmd
	CombinedOutput() ([]byte, error)
	Output() ([]byte, error)
	Run() error
	Start() error
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	Wait() error

	// Simplified wrapper for Process.Kill + Wait
	Kill() error
}

// Basic Cmd implementation based on exec.Cmd
type ExecCmd struct {
	*exec.Cmd
}

// Command creates a new ExecCmd from the give command name and arguments arg.
func Command(name string, arg ...string) *ExecCmd {
	return &ExecCmd{exec.Command(name, arg...)}
}

// Kill attempts to terminate the command, and waits for it to exit. It returns
// an error if the command did not exit successfully or was not terminated by
// SIGKILL.
func (cmd *ExecCmd) Kill() error {
	cmd.Process.Kill()

	err := cmd.Wait()
	if err == nil {
		return nil
	}

	if eerr, ok := err.(*exec.ExitError); ok {
		status := eerr.Sys().(syscall.WaitStatus)
		if status.Signal() == syscall.SIGKILL {
			return nil
		}
	}

	return err
}
