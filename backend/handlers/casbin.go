package handlers

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFrontendPermission(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get current user/subject
		firebaseUUID, existed := c.Get("UUID")

		if !existed {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "No subject found"})
			return
		}

		permissions, err := enforcer.GetImplicitPermissionsForUser(firebaseUUID.(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not find permissions for user"})
			return
		}

		c.JSON(http.StatusOK, permissions)
	}
}
