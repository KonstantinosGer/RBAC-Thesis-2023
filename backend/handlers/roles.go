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

func UpdateRole(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {

		//From postman's Body/form-data
		var requestBody struct {
			UserId  string `json:"id"`
			NewRole string `json:"role"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//From db
		oldRole, err := enforcer.GetRolesForUser(requestBody.UserId)
		if err != nil {
			panic(fmt.Sprintf("failed to get roles for user %s: %v", requestBody.UserId, err))
		}

		if len(oldRole) == 0 {
			enforcer.AddGroupingPolicy(requestBody.UserId, requestBody.NewRole)
		} else {
			enforcer.UpdateGroupingPolicy([]string{requestBody.UserId, oldRole[0]}, []string{requestBody.UserId, requestBody.NewRole})
		}
	}
}

func GetAllRoles() gin.HandlerFunc {
	return func(c *gin.Context) {

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

		var roles []models.Role

		// Select all roles
		err = database.Debug().Find(&roles).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, roles)
	}
}

func AddRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Initialize role with given parameters from postman's Body/form-data
		var role models.Role
		if err := c.ShouldBindJSON(&role); err != nil {
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

		// Add new role or update customer's name
		err = database.Debug().Save(&role).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func DeleteRole(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			Role string `json:"role"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		var role = models.Role{Role: requestBody.Role}

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

		//Get number of users with this role from "cabin_rule" table in DB
		var countUsersWithRole int64

		err = database.Model(&models.CasbinRule{}).Where(&models.CasbinRule{V1: role.Role}).Count(&countUsersWithRole).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		if countUsersWithRole != 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Please remove this role from all users before deleting it"})
			return
		}

		//Delete from "role" table in DB
		err = database.Delete(&role).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err)
			return
		}

		//Delete from "casbin_rule" table in DB
		_, err = enforcer.RemoveFilteredPolicy(0, role.Role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}
