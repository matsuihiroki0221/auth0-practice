package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/jwk"
)

var (
	auth0Domain = os.Getenv("VITE_AUTH0_DOMAIN")  // Auth0のドメイン
	audience    = os.Getenv("AUTH0_API_AUDIENCE") // Auth0で設定したAPIのオーディエンス
)

// JWTのバリデータを初期化
var jwtValidator *validator.Validator

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth0Domain = os.Getenv("VITE_AUTH0_DOMAIN") // Auth0のドメイン
	audience = os.Getenv("AUTH0_API_AUDIENCE")   // Auth0で設定したAPIのオーディエンス
	keyFunc := func(ctx context.Context) (interface{}, error) {
		resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("failed to fetch JWKS")
		}

		var jwks struct {
			Keys []jwk.Key `json:"keys"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
			return nil, err
		}

		if len(jwks.Keys) == 0 {
			return nil, errors.New("no keys found in JWKS")
		}

		var key *rsa.PublicKey
		if err := jwks.Keys[0].Raw(&key); err != nil {
			return nil, err
		}

		return key, nil
	}

	jwtValidator, err = validator.New(
		keyFunc,
		validator.RS256,
		auth0Domain,
		[]string{audience},
	)
	if err != nil {
		panic("Failed to create JWT validator: " + err.Error())
	}
}

// Auth0認証ミドルウェア
func Auth0Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// トークンの検証
		token, err := jwtValidator.ValidateToken(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 検証に成功した場合、クレームをコンテキストに保存
		c.Set("user", token)
		c.Next()
	}
}
