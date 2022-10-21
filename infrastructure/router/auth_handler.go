package router

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nicholasanthonys/go-recipe/adapter/api/action"
	"github.com/nicholasanthonys/go-recipe/adapter/presenter"
	"github.com/nicholasanthonys/go-recipe/adapter/repository"
	"github.com/nicholasanthonys/go-recipe/usecase"
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

type User struct {
	Password string `json:"password"`
	Username string `json:"username"`
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

// with jwt
func (g ginEngine) SignInHandler(c *gin.Context) {

	uc := usecase.NewSignInInteractor(
		repository.NewUserNoSQL(g.db),
		presenter.NewSignInPresenter(),
		g.ctxTimeout,
	)

	ac := action.NewSignInAction(uc, g.log, c)
	ac.Execute(c.Writer, c.Request)

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.Username != "admin" || user.Password != "password" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := Claims{
		user.Username,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(
		os.Getenv("JWT_SECRET"),
	))

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
