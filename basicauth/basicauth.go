package basicauth

import (
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
	"github.com/go-zoox/zoox/middleware"
)

// Config is the basic config
type Config struct {
	Port int64
	// mode: static username/password

	// Username is the basic username
	Username string
	// Password is the basic password
	Password string

	// mode: dynamic service with username and password

	// AuthService is auth service url
	// Example:
	//   POST https://example.com/api/login
	//	      Header => Content-Type: application/json
	//				Body => { "username": "username", "password": "password" }
	AuthService string

	// Upstream is the upstream service
	// Example: http://httpbin:8080
	Upstream string
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	if cfg.AuthService != "" {
		app.Use(func(ctx *zoox.Context) {
			user, pass, ok := ctx.Request.BasicAuth()
			if !ok {
				ctx.Set("WWW-Authenticate", `Basic realm="go-zoox"`)
				ctx.Status(401)
				return
			}

			response, err := fetch.Post(cfg.AuthService, &fetch.Config{
				Headers: fetch.Headers{
					"Content-Type": "application/json",
				},
				Body: map[string]string{
					"from":     "go-zoox/gzauth.basic",
					"username": user,
					"password": pass,
				},
			})
			if err != nil {
				logger.Errorf("basic auth with auth-service error: %s", err)
				fmt.PrintJSON(map[string]any{
					"request":  response.Request,
					"response": response.String(),
				})

				ctx.String(500, "internal server error")
				return
			}

			if response.Status != 200 {
				ctx.String(400, "invalid username and password: %s", response.String())
				return
			}

			ctx.Next()
		})
	} else {
		app.Use(func(ctx *zoox.Context) {
			user, pass, ok := ctx.Request.BasicAuth()
			if !ok {
				ctx.Set("WWW-Authenticate", `Basic realm="go-zoox"`)
				ctx.Status(401)
				return
			}

			if !(user == cfg.Username && pass == cfg.Password) {
				ctx.Status(401)
				return
			}

			ctx.Next()
		})
	}

	app.Use(middleware.Proxy(&middleware.ProxyConfig{
		Rewrites: middleware.ProxyGroupRewrites{
			{
				RegExp: "/(.*)",
				Rewrite: middleware.ProxyRewrite{
					Target: cfg.Upstream,
				},
			},
		},
	}))

	return app.Run()
}
