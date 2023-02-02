package commands

import (
	"fmt"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzauth/beartoken"
)

func RegistryBear(app *cli.MultipleProgram) {
	app.Register("bear", &cli.Command{
		Name:  "bear",
		Usage: "auth with bear token auth",
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
				Name:    "token",
				Usage:   "Bear Token",
				EnvVars: []string{"TOKEN"},
			},
			&cli.StringFlag{
				Name:    "auth-service",
				Usage:   "auth service",
				EnvVars: []string{"AUTH_SERVICE"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			token := ctx.String("token")
			authService := ctx.String("auth-service")

			if token == "" && authService == "" {
				return fmt.Errorf("token or auth-service is required")
			}

			return beartoken.Serve(&beartoken.Config{
				Port:        ctx.Int64("port"),
				Token:       token,
				AuthService: authService,
				Upstream:    ctx.String("upstream"),
			})
		},
	})
}
