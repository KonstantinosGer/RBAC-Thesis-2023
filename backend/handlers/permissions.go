package handlers

import (
	db "backend/database"
	"backend/models"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetPermissionsForRole(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		// ** SOS **
		// BIND TIPS!
		// How to properly Bind?
		//
		// With normal Bind, you need "form" field to all model's columns you need to read
		// (Bind is usually used in GET methods)
		//
		// With ShouldBindJSON, you need "json" field to all model's columns you need to read
		// (ShouldBindJSON is usually used in POST, PUT, DELETE methods)
		//
		// With ShouldBindUri, you need "uri" field to all model's columns you need to read
		// (ShouldBindUri is used when parameters are passed in route's url)
		//
		var role models.Role
		if err := c.Bind(&role); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//
		// Connect to RBAC Database (for gorm queries)
		//
		database, err := db.ConnectToRBACGorm()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}
		defer db.CloseDBConnectionGorm(database)

		permissionsForUser := enforcer.GetPermissionsForUser(role.Role)
		fmt.Println(permissionsForUser)

		var permissions []models.Permission

		// Select all permissions (order by category)
		err = database.Debug().Order("category").Find(&permissions).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//fmt.Println(permissions)
		var userPermissionsObject models.FrontendPolicy
		userPermissionsObject = make(map[string][]models.PermissionInfo)
		for _, globalPerm := range permissions {

			var userHasPermission bool
			for _, userPerm := range permissionsForUser {
				if globalPerm.Action == userPerm[2] && globalPerm.Resource == userPerm[1] {
					userHasPermission = true
					break
				} else {
					userHasPermission = false
				}
			}

			userPermissionsObject[globalPerm.Category] = append(userPermissionsObject[globalPerm.Category],
				models.PermissionInfo{
					PermissionId:          globalPerm.Id,
					PermissionDescription: globalPerm.Description,
					PermissionAction:      globalPerm.Action,
					PermissionResource:    globalPerm.Resource,
					HasPermission:         userHasPermission,
				},
			)
		}

		c.JSON(http.StatusOK, userPermissionsObject)
	}
}

func AddPermission(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			NewRole      string `json:"newRole"`
			NewData      string `json:"newData"`
			NewPrivilege string `json:"newPrivilege"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//Add policy
		_, err := enforcer.AddPolicy(requestBody.NewRole, requestBody.NewData, requestBody.NewPrivilege)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create casbin enforcer"})
			return
		}

	}
}

func DeletePermission(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			Role      string `json:"role"`
			Data      string `json:"data"`
			Privilege string `json:"privilege"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//Remove policy
		enforcer.RemovePolicy(requestBody.Role, requestBody.Data, requestBody.Privilege)

	}
}
