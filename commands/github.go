package commands

import (
	"fmt"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzauth/github"
)

func RegistryGitHub(app *cli.MultipleProgram) {
	app.Register("github", &cli.Command{
		Name:  "github",
		Usage: "auth with oauth2 github",
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
				Name:    "secret-key",
				Usage:   "secret-key used for session and jwt token",
				EnvVars: []string{"SECRET_KEY"},
			},
			&cli.StringFlag{
				Name:    "client-id",
				Usage:   "GitHub Client ID",
				EnvVars: []string{"CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:    "client-secret",
				Usage:   "GitHub Client Secret",
				EnvVars: []string{"CLIENT_SECRET"},
			},
			&cli.StringFlag{
				Name:    "redirect-uri",
				Usage:   "GitHub Redirect URI",
				EnvVars: []string{"REDIRECT_URI"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			ClientID := ctx.String("client-id")
			ClientSecret := ctx.String("client-secret")
			RedirectURI := ctx.String("redirect-uri")

			if ClientID == "" || ClientSecret == "" || RedirectURI == "" {
				return fmt.Errorf("client id, secret, redirect_uri are required")
			}

			return github.Serve(&github.Config{
				Port:         ctx.Int64("port"),
				Upstream:     ctx.String("upstream"),
				SecretKey:    ctx.String("secret-key"),
				ClientID:     ClientID,
				ClientSecret: ClientSecret,
				RedirectURI:  RedirectURI,
			})
		},
	})
}
