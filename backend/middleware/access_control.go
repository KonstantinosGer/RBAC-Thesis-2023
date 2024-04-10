package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Authorize determines if current user has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		firebaseUUID, existed := c.Get("UUID")

		if !existed {
			// Always return a "message"!
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User hasn't logged in yet"})
			return
		}

		// Load policy from Database
		err := enforcer.LoadPolicy()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to load policy from DB"})
			return
		}

		// Casbin enforces policy
		ok, err := enforcer.Enforce(fmt.Sprint(firebaseUUID), obj, act)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred when authorizing user"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized"})
			return
		}
		c.Next()
	}
}
