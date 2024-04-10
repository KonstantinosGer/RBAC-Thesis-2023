package handlers

import (
	db "backend/database"
	"backend/models"
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "strconv"
	"strings"
)

func ToggleCustomerUserAccess(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			CustomerId   int    `json:"customer_id"`
			UserId       string `json:"user_id"`
			AccessObject string `json:"access_object"`
			HasAccess    bool   `json:"has_access"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		var customer = models.Customer{Id: requestBody.CustomerId}

		type userModel struct {
			Id           string `json:"id"`
			AccessObject string `json:"access_object"`
			HasAccess    bool   `json:"has_access"`
		}
		var user = userModel{Id: requestBody.UserId, AccessObject: requestBody.AccessObject, HasAccess: requestBody.HasAccess}

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

		// Make sure user exists in database
		var count int64
		err = database.Model(&models.User{}).Where(&models.User{Id: requestBody.UserId}).Count(&count).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err.Error())
			return
		}

		if count == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Please save an email before choosing permission!", "type": "warning"})
			return
		}

		// Permissions for specific customer
		permissionName := fmt.Sprintf("portal::data::%d::%s", customer.Id, user.AccessObject)
		if user.HasAccess {
			// Add permission, if not exists (enforcer checks if exists)
			_, err = enforcer.AddPermissionForUser(user.Id, permissionName, "read")
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			// Remove permission, if exists (enforcer checks if exists)
			_, err = enforcer.DeletePermissionForUser(user.Id, permissionName, "read")
			if err != nil {
				log.Println(err.Error())
			}
		}

		// Permissions for general financial or performance access
		permissionName = fmt.Sprintf("portal::data::customer::%s", user.AccessObject)
		if user.HasAccess {
			// add permission, if not exists (enforcer checks if exists)
			_, err = enforcer.AddPermissionForUser(user.Id, permissionName, "read")
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			// remove permission, if exists (enforcer checks if exists)
			_, err = enforcer.DeletePermissionForUser(user.Id, permissionName, "read")
			if err != nil {
				log.Println(err.Error())
			}
		}

		c.JSON(http.StatusOK, nil)
	}
}

func AddCustomerUserAssociation(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			CustomerId int    `json:"id"`
			UserEmail  string `json:"email"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		var customer = models.Customer{Id: requestBody.CustomerId}
		var user = models.User{Email: requestBody.UserEmail}

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

		// Make sure given user email is not empty
		if strings.TrimSpace(user.Email) == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Please choose an email first!", "type": "warning"})
			return
		}

		//Get user associated with passed user's email
		err = database.Debug().Where(&models.User{Email: user.Email}).First(&user).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		// Check if customerId - userId association already exists in customer_user table
		// If true return message, if false add new association to customer_user
		var count int64
		err = database.Table("customer_user").Where("customer_id = ? AND user_id = ?", customer.Id, user.Id).Count(&count).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}
		// Association already exists
		if count > 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "This association already exists!\nPlease select another email."})
			log.Println(err.Error())
			return
		}
		// Association does not exist, so add new association to customer_user table in database
		err = database.Model(&customer).Association("Users").Append(&user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		// Register user as a customer
		// Returns false if the user already has the permission
		_, err = enforcer.AddPermissionForUser(user.Id, "portal::data::customer", "read")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func DeleteCustomerUserAssociation(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			UserId     string `json:"user_id"`
			CustomerId int    `json:"customer_id"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
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

		// Make sure user exists in database
		var count int64
		err = database.Model(&models.User{}).Where(&models.User{Id: requestBody.UserId}).Count(&count).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err.Error())
			return
		}

		if count == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Please select an existing user to delete.", "type": "warning"})
			return
		}

		// Initialize customer with given id
		var customer = models.Customer{Id: requestBody.CustomerId}
		// Initialize user with given id
		var user = models.User{Id: requestBody.UserId}

		// First check if user associates with multiple customers
		// If true, do NOT delete him
		// If false, delete him completely
		var associatedCustomers []models.Customer
		err = database.Joins("JOIN customer_user ON customers.id = customer_user.customer_id").
			Where("customer_user.user_id = ?", user.Id).
			Find(&associatedCustomers).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err.Error())
			return
		}

		//Delete user association with this customer from "customer_user" table in DB
		err = database.Debug().Model(&user).Association("Customers").Delete(&customer)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//
		// Delete user's permissions, specifically for this customer
		//
		// Remove financial policy from user, if exists
		if enforcer.HasPolicy(user.Id, fmt.Sprintf("portal::data::%d::finance", customer.Id), "read") {
			_, err = enforcer.RemovePolicy(user.Id, fmt.Sprintf("portal::data::%d::finance", customer.Id), "read")
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				log.Println(err)
				return
			}
		}
		// Remove performance policy from user, if exists
		if enforcer.HasPolicy(user.Id, fmt.Sprintf("portal::data::%d::performance", customer.Id), "read") {
			_, err = enforcer.RemovePolicy(user.Id, fmt.Sprintf("portal::data::%d::performance", customer.Id), "read")
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				log.Println(err)
				return
			}
		}

		// User is associated only with this customer
		// So delete user completely
		if !(len(associatedCustomers) > 1) {
			// User is not associated with any customer anymore
			_, err = enforcer.DeletePermissionForUser(user.Id, "portal::data::customer", "read")
			if err != nil {
				log.Println(err.Error())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			// Remove general financial assess from user
			enforcer.DeletePermissionForUser(user.Id, "portal::data::customer::finance", "read")
			// Remove general performance assess from user
			enforcer.DeletePermissionForUser(user.Id, "portal::data::customer::performance", "read")

			// Delete user from "users" table
			err = database.Debug().Delete(&user).Error
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				log.Println(err)
				return
			}

			// Initialize firebaseAuth
			firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
			// Delete user from firebase
			err = firebaseAuth.DeleteUser(context.Background(), user.Id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				log.Println(err)
				return
			}

		} else {
			// User is also associated with other customers
			var hasFinancialAccessToAnotherCustomer = false
			var hasPerformanceAccessToAnotherCustomer = false
			for _, associatedCustomer := range associatedCustomers {
				// Ignore current customer
				if associatedCustomer.Id != customer.Id {
					if enforcer.HasPolicy(user.Id, fmt.Sprintf("portal::data::%d::finance", associatedCustomer.Id), "read") {
						hasFinancialAccessToAnotherCustomer = true
					}
					if enforcer.HasPolicy(user.Id, fmt.Sprintf("portal::data::%d::performance", associatedCustomer.Id), "read") {
						hasPerformanceAccessToAnotherCustomer = true
					}
				}
			}

			if !hasFinancialAccessToAnotherCustomer {
				enforcer.DeletePermissionForUser(user.Id, "portal::data::customer::finance", "read")
			}
			if !hasPerformanceAccessToAnotherCustomer {
				enforcer.DeletePermissionForUser(user.Id, "portal::data::customer::performance", "read")
			}
		}

		c.JSON(http.StatusOK, nil)
	}
}
