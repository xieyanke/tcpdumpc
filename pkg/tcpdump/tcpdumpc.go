/*
Copyright © 2023 XieYanke

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tcpdump

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/vishvananda/netns"
	"github.com/xieyanke/tcpdumpc/pkg/runtime/docker"
)

type TcpdumpC struct {
	containerID      string
	containerRuntime string
	args             []string
	stdOut           io.Writer
	stdErr           io.Writer
}

func (c *TcpdumpC) Run() (err error) {
	var pid int
	switch c.containerRuntime {
	case "docker":
		pid, err = docker.GetPIDByContainerID(c.containerID)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("not yet implement runtime '%s'", c.containerRuntime)
	}

	nsHandle, err := netns.GetFromPid(pid)
	if err != nil {
		return err
	}

	err = netns.Set(nsHandle)
	if err != nil {
		return err
	}

	cmd := exec.Command("tcpdump", c.args...)
	cmd.Stdout = c.stdOut
	cmd.Stderr = c.stdErr

	return cmd.Run()
}

func NewTcpdumpC(id, runtime string, stdOut, stdErr io.Writer, args []string) *TcpdumpC {
	if runtime == "" {
		runtime = "docker"
	}

	return &TcpdumpC{
		containerID:      id,
		containerRuntime: runtime,
		args:             args,
		stdOut:           stdOut,
		stdErr:           stdErr,
	}
}

func CheckTcpdumpExist() bool {
	err := exec.Command("tcpdump", "--version").Run()
	return err == nil
}
