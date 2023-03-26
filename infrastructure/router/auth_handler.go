package router

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/auth0-community/go-auth0"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nicholasanthonys/go-recipe/adapter/api/action"
	"github.com/nicholasanthonys/go-recipe/adapter/presenter"
	"github.com/nicholasanthonys/go-recipe/adapter/repository"
	"github.com/nicholasanthonys/go-recipe/usecase"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"
)

type AuthHandler struct{}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

// with X-api-key
func (g ginEngine) CheckAPIKey(c *gin.Context) {
	if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "API key invalid",
		})
		c.Abort()
		return
	}
	c.Next()

}

func (g ginEngine) SignInHandler(c *gin.Context) {

	uc := usecase.NewSignInInteractor(
		repository.NewUserNoSQL(g.db),
		presenter.NewSignInPresenter(),
		g.ctxTimeout,
	)

	ac := action.NewSignInAction(uc, g.log, c)
	ac.Execute(c.Writer, c.Request)
}

func (g ginEngine) RegisterHandler(c *gin.Context) {
	uc := usecase.NewRegisterInteractor(
		repository.NewUserNoSQL(g.db),
		presenter.NewRegisterPresenter(),
		g.ctxTimeout,
	)

	ac := action.NewRegisterAction(uc, g.log)
	ac.Execute(c.Writer, c.Request)
}

func (g ginEngine) SignOutHandler(c *gin.Context) {
	fmt.Println("sign out handler triggered")
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Signed out..."})
}

// middleware
func (g ginEngine) AuthMiddleware(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}

	fmt.Println("JWT SECRET : ", os.Getenv("JWT_SECRET"))
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if tkn == nil || !tkn.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()

}

func (g ginEngine) AuthMiddlewareWithSession(c *gin.Context) {

	session := sessions.Default(c)
	sessionToken := session.Get("token")
	if sessionToken == nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Not Logged in",
		})
		c.Abort()
	}
	c.Next()
}

func (g ginEngine) RefreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if tkn == nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	if time.Unix(int64(claims.ExpiresAt.Second()), 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is not expired yet",
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(time.Unix(expirationTime.Unix(), 0))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtOutput := JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)

}

// Auth 0

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

// New instantiates the *Authenticator.
func NewAuthenticator() (*Authenticator, error) {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

// Handler for our login.
func Auth0HandlerLogin(auth *Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

// Handler for our callback.
func Auth0HandlerCallback(auth *Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}

// Handler for our logout.
func Auth0HandlerLogout(ctx *gin.Context) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

// Handler for our logged-in user page.
func Auth0UserProfileHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	ctx.HTML(http.StatusOK, "user.html", profile)
}

// Auth0 IsAuthenticated Middleware
// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
func Auth0IsAuthenticated(ctx *gin.Context) {
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		ctx.Next()
	}
}

func Auth0Middleware(ctx *gin.Context) {
	var auth0Domain = "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
	client := auth0.NewJWKClient(auth0.JWKClientOptions{
		URI: auth0Domain + ".well-known/jwks.json",
	}, nil)

	fmt.Println("Idnetifier : ", os.Getenv("AUTH0_API_IDENTIFIER"))
	configuration := auth0.NewConfiguration(
		client,
		[]string{os.Getenv("AUTH0_API_IDENTIFIER")},
		auth0Domain,
		jose.RS256,
	)

	validator := auth0.NewValidator(configuration, nil)
	_, err := validator.ValidateRequest(ctx.Request)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token", "error" : err.Error()})
		ctx.Abort()
		return
	}

	ctx.Next()

}
