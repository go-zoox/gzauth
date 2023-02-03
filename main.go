package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzauth/commands"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:    "gzauth",
		Usage:   "gzauth is a portable auth cli",
		Version: Version,
	})

	// basic
	commands.RegistryBasic(app)
	// bear
	commands.RegistryBear(app)
	// oauth2
	//	=> doreamon
	commands.RegistryDoreamon(app)
	//	=> github
	commands.RegistryGitHub(app)

	app.Run()
}
