package middleware

import (
	"golang_email_verification/initializers"
	"golang_email_verification/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeserializeUser(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		cookie, err := ctx.Cookie("token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(token, config.TokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		subObj, ok := sub.(map[string]interface{})
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Token sub is not a valid object"})
			return
		}

		// Access the properties within "sub"
		userID, userIDOk := subObj["user_id"].(string)
		userRole, userRoleOk := subObj["user_role"].(string)

		if !userIDOk || !userRoleOk {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Token sub is missing required properties"})
			return
		}

		// Check user role and authorization
		if !checkAuthorization(userRole, requiredRole) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not authorized"})
			return
		}

		// Set user information in the context
		ctx.Set("currentUser", userID)

		// Continue with the request
		ctx.Next()
	}
}

func checkAuthorization(userRole, requiredRole string) bool {
	switch requiredRole {
	case "admin":
		return userRole == "admin"
	case "staff":
		return userRole == "admin" || userRole == "staff"
	case "employee":
		return userRole == "admin" || userRole == "staff" || userRole == "employee"
	default:
		return false
	}
}
