package commands

import (
	"fmt"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzauth/basic"
)

func RegistryBasic(app *cli.MultipleProgram) {
	app.Register("basic", &cli.Command{
		Name:  "basic",
		Usage: "auth with basic auth",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
			&cli.StringFlag{
				Name:     "upstream",
				Usage:    "upstream service",
				EnvVars:  []string{"UPSTREAM"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "username",
				Usage:   "basic username",
				EnvVars: []string{"USERNAME"},
			},
			&cli.StringFlag{
				Name:    "password",
				Usage:   "basic password",
				EnvVars: []string{"PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "auth-service",
				Usage:   "auth service",
				EnvVars: []string{"AUTH_SERVICE"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			username := ctx.String("username")
			password := ctx.String("password")
			authService := ctx.String("auth-service")

			if username == "" && authService == "" {
				return fmt.Errorf("username/password or auth-service is required")
			}

			if username != "" || password != "" {
				if username == "" || password == "" {
					return fmt.Errorf("username/password are required")
				}
			}

			return basic.Serve(&basic.Config{
				Port:        ctx.Int64("port"),
				Username:    username,
				Password:    password,
				AuthService: authService,
				Upstream:    ctx.String("upstream"),
			})
		},
	})
}
