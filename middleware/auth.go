package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Parse token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing token adalah HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Ambil klaim dari token dan set ke context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Pastikan setiap klaim ada dan set ke context
			if authID, ok := claims["auth_id"].(float64); ok {
				c.Set("auth_id", int64(authID))
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "auth_id claim missing"})
				c.Abort()
				return
			}

			if accountID, ok := claims["account_id"].(float64); ok {
				c.Set("account_id", int64(accountID))
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "account_id claim missing"})
				c.Abort()
				return
			}

			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "username claim missing"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next() // Authorized, lanjut ke handler berikutnya
	}
}

// func AuthMiddleware(secretKey string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")

// 		// Parse the token
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, http.ErrAbortHandler
// 			}
// 			return []byte(secretKey), nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			c.Abort() // Stop further processing if unauthorized
// 			return
// 		}

// 		// Set the token claims to the context
// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			if authID, ok := claims["auth_id"].(float64); ok {
// 				c.Set("auth_id", int64(authID))
// 			}
// 			if accountID, ok := claims["account_id"].(float64); ok {
// 				c.Set("account_id", int64(accountID))
// 			}
// 			if username, ok := claims["username"].(string); ok {
// 				c.Set("username", username)
// 			}
// 		} else {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 			c.Abort()
// 			return
// 		}

// 		c.Next() // Authorized, Proceed to the next handler
// 	}
// }
