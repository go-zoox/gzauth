package commands

import (
	"fmt"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzauth/feishu"
)

func RegistryFeishu(app *cli.MultipleProgram) {
	app.Register("feishu", &cli.Command{
		Name:  "feishu",
		Usage: "auth with oauth2 feishu",
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
				Name:     "client-id",
				Usage:    "Feishu Client ID",
				EnvVars:  []string{"CLIENT_ID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "client-secret",
				Usage:    "Feishu Client Secret",
				EnvVars:  []string{"CLIENT_SECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "redirect-uri",
				Usage:    "Feishu Redirect URI",
				EnvVars:  []string{"REDIRECT_URI"},
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			ClientID := ctx.String("client-id")
			ClientSecret := ctx.String("client-secret")
			RedirectURI := ctx.String("redirect-uri")

			if ClientID == "" || ClientSecret == "" || RedirectURI == "" {
				return fmt.Errorf("client id, secret, redirect_uri are required")
			}

			return feishu.Serve(&feishu.Config{
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
