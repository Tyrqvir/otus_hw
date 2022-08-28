package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	proxyCmd := exec.Command(cmd[0], cmd[1:]...) //#nosec G204

	for key, value := range env {
		if err := os.Unsetenv(key); err != nil {
			return 1
		}

		if value.NeedRemove {
			continue
		}
		os.Setenv(key, value.Value)
	}

	proxyCmd.Stdout = os.Stdout
	proxyCmd.Stderr = os.Stderr
	proxyCmd.Stdin = os.Stdin

	if err := proxyCmd.Run(); err != nil {
		return 1
	}

	return 0
}
