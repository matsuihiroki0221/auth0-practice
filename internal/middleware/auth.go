package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

var (
	auth0Domain = os.Getenv("VITE_AUTH0_DOMAIN")       // Auth0のドメイン
	audience    = os.Getenv("VITE_AUTH0_API_AUDIENCE") // Auth0で設定したAPIのオーディエンス
	// jwtValidator    *validator.Validator
	auth0Middleware *jwtmiddleware.JWTMiddleware
)

type customClaims struct {
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions"`
}

func (c customClaims) Validate(ctx context.Context) error {
	return nil
}

// JWT バリデータの作成
func init() {
	var err error
	issuerURL, err := url.Parse("https://" + auth0Domain + "/")
	if err != nil {
		log.Fatalf("failed to parse issuer URL: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &customClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("failed to initialize JWT validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
	}

	auth0Middleware = jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)
}

// 認証ミドルウェア
func Auth0Middleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// http.Handler を Gin コンテキストに適用
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
		})

		httpHandler := auth0Middleware.CheckJWT(next)
		httpHandler.ServeHTTP(c.Writer, c.Request)

		if c.Writer.Status() == http.StatusUnauthorized {
			c.Abort()
			return
		}

		validatedClaims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		customClaims, ok := validatedClaims.CustomClaims.(*customClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid custom claims"})
			c.Abort()
			return
		}

		// Gin Contextに保存
		c.Set("permissions", customClaims.Permissions)
		c.Set("registeredClaims", validatedClaims.RegisteredClaims)

		// 次のミドルウェアまたはハンドラーを呼び出す
		c.Next()
	}
}
