package feishu

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	"github.com/go-zoox/oauth2/feishu"
	"github.com/go-zoox/random"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
	"github.com/go-zoox/zoox/middleware"
)

// Config is the basic config
type Config struct {
	Port int64

	// SecretKey uses for session and backend jwt token
	SecretKey string

	// ClientID is the feishu client id
	ClientID string
	// ClientSecret is the feishu client secret
	ClientSecret string
	// RedirectURI is the feishu redirect uri
	RedirectURI string

	// Upstream is the upstream service
	// Example: http://httpbin:8080
	Upstream string
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	if cfg.SecretKey != "" {
		app.Config.SecretKey = cfg.SecretKey
	}

	client, err := feishu.New(&feishu.FeishuConfig{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURI:  cfg.RedirectURI,
	})
	if err != nil {
		return fmt.Errorf("failed to create feishu client: %v", err)
	}

	app.Use(func(ctx *zoox.Context) {
		if ctx.Method != "GET" {
			if ctx.Session().Get("oauth2.user") == "" {
				if ctx.AcceptJSON() {
					ctx.JSON(http.StatusUnauthorized, zoox.H{
						"code":    401000,
						"message": "Unauthorized",
					})
					return
				}

				ctx.String(http.StatusUnauthorized, "Unauthorized")
				return
			}

			ctx.Next()
			return
		}

		// login
		if ctx.Path == "/login" || ctx.Path == "/login/feishu" {
			originState := random.String(8)
			originFrom := ctx.Query().Get("from").String()
			client.Authorize(originState, func(loginUrl string) {
				if originFrom != "" {
					ctx.Session().Set("from", originFrom)
				}

				ctx.Session().Set("oauth2.state", originState)
				ctx.Redirect(loginUrl)
			})
			return
		}

		// callback
		if ctx.Path == "/login/callback" || ctx.Path == "/login/feishu/callback" {
			code := ctx.Query().Get("code").String()
			state := ctx.Query().Get("state").String()

			originState := ctx.Session().Get("oauth2.state")
			if state != originState {
				logger.Errorf("invalid oauth2 state, expect %s, but got %s", originState, state)
				time.Sleep(1 * time.Second)
				ctx.Redirect(fmt.Sprintf("/login?reason=%s", "invalid_oauth2_state"))
				return
			}
			originFrom := ctx.Session().Get("from")
			if originFrom == "" {
				originFrom = "/"
			}

			client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
				userSessionKey := fmt.Sprintf("user:%s", user.ID)

				ctx.Cache().Set(userSessionKey, user, ctx.App.Config.Session.MaxAge)

				ctx.Session().Set("oauth2.user", userSessionKey)
				// ctx.Session().Set("oauth2.token", token.AccessToken)

				ctx.Redirect(originFrom)
			})
			return
		}

		// logout
		if ctx.Path == "/logout" {
			client.Logout(func(logoutUrl string) {
				ctx.Session().Del("oauth2.user")
				ctx.Redirect(logoutUrl)
			})
			return
		}

		if ctx.Path == "/api/user" {
			userSessionKey := ctx.Session().Get("oauth2.user")
			if userSessionKey == "" {
				if ctx.AcceptJSON() {
					ctx.JSON(http.StatusUnauthorized, zoox.H{
						"code":    401001,
						"message": "unauthorized",
					})
					return
				}

				ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Path), "user not login or token expired"))
				return
			}
			user := oauth2.User{}
			if err := ctx.Cache().Get(userSessionKey, &user); err != nil {
				time.Sleep(1 * time.Second)
				ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Path), "user cache not found"))
				return
			}

			ctx.Success(user)
			return
		}

		if ctx.Session().Get("oauth2.user") == "" {
			originFrom := ctx.Request.RequestURI
			time.Sleep(1 * time.Second)
			ctx.Redirect(fmt.Sprintf("/login?from=%s", url.QueryEscape(originFrom)))
			return
		}

		userSessionKey := ctx.Session().Get("oauth2.user")
		user := oauth2.User{}
		if err := ctx.Cache().Get(userSessionKey, &user); err != nil {
			time.Sleep(1 * time.Second)
			ctx.Redirect(fmt.Sprintf("/login?reason=%s", "user cache not found"))
			return
		}

		token, err := ctx.Jwt().Sign(map[string]interface{}{
			"user_id":       user.ID,
			"user_nickname": user.Nickname,
			"user_avatar":   user.Avatar,
			"user_email":    user.Email,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.Request.Header.Set("X-GZAuth-Token", token)

		ctx.Next()
	})

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
