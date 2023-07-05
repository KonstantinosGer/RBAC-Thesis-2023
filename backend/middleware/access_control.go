package middleware

import (
	"fmt"
	"net/http"

	//"github.com/casbin/casbin"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Authorize determines if current user has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		firebaseUUID, existed := c.Get("UUID")

		if !existed {
			//Προσοχή! Πάντα θα επιστρέφουμε "message"! (σαν σύμβαση)
			//c.AbortWithStatusJSON(401, gin.H{"message": "User hasn't logged in yet"})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User hasn't logged in yet"})
			return
		}

		// Load policy from Database
		err := enforcer.LoadPolicy()
		if err != nil {
			//c.AbortWithStatusJSON(500, gin.H{"message": "Failed to load policy from DB"})
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to load policy from DB"})
			return
		}

		// Casbin enforces policy
		ok, err := enforcer.Enforce(fmt.Sprint(firebaseUUID), obj, act)

		if err != nil {
			//c.AbortWithStatusJSON(500, gin.H{"message": "Error occurred when authorizing user"})
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Error occurred when authorizing user"})
			return
		}

		if !ok {
			//c.AbortWithStatusJSON(401, gin.H{"message": "You are not authorized"})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized"})
			return
		}
		c.Next()
	}
}
