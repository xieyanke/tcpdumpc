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

package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPIDByContainerID(containerID string) (int, error) {
	stateFile, err := searchContainerStateFile(containerID)
	if err != nil {
		return 0, err
	}

	stateData, err := os.ReadFile(stateFile)
	if err != nil {
		return 0, err
	}

	m := make(map[string]any)
	err = json.Unmarshal(stateData, &m)
	if err != nil {
		return 0, err
	}

	if m["init_process_pid"] != nil {
		var pid int

		pidF64, ok := m["init_process_pid"].(float64)
		if ok {
			pid = int(pidF64)
		} else {
			return 0, fmt.Errorf("type assert 'init_process_pid: %v' failed", m["init_process_pid"])
		}

		return pid, nil
	} else {
		return 0, fmt.Errorf("field 'init_process_pid' does not exist in '%s' file", stateFile)
	}
}

func searchContainerStateFile(containerID string) (string, error) {
	mobyDir := "/var/run/docker/runtime-runc/moby"
	files, err := os.ReadDir("/var/run/docker/runtime-runc/moby")
	if err != nil {
		return "", err
	}

	var stateFiles []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), containerID) {
			stateFiles = append(stateFiles, filepath.Join(mobyDir, file.Name(), "state.json"))
		}
	}

	if len(stateFiles) > 1 {
		return "", fmt.Errorf("container id '%s' is not unique", containerID)
	} else if len(stateFiles) == 0 {
		return "", fmt.Errorf("can't find container by container id '%s'", containerID)
	} else {
		return stateFiles[0], nil
	}
}
