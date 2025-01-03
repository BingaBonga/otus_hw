package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envValue := range env {
		err := os.Unsetenv(envName)
		if err != nil {
			log.Println(err)
			return 1
		}

		if !envValue.NeedRemove {
			err := os.Setenv(envName, envValue.Value)
			if err != nil {
				log.Println(err)
				return 1
			}
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = os.Environ()
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		log.Println(err)
	}

	return command.ProcessState.ExitCode()
}
