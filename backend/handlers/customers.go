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
	"strings"
)

func GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get filters from frontend
		var filters models.Filter
		if err := c.Bind(&filters); err != nil {
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

		var customers []models.Customer

		// Select all customers
		err = database.Debug().Select("id", "full_name").Where("id LIKE ? OR full_name LIKE ?", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%").Find(&customers).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, customers)
	}
}

func AddCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		// Make sure given customer name is not empty and given customer id is a positive integer
		if strings.TrimSpace(customer.FullName) == "" || customer.Id <= 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Make sure all fields are filled in correctly!", "type": "warning"})
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

		// Add a new customer to customers table in database
		err = database.Debug().Create(&customer).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func DeleteCustomer(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind from url, passing url's passed id directly to a customer object with id = url's passed id
		var customer models.Customer
		err := c.ShouldBindUri(&customer)
		if err != nil {
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

		// Fetch number of associated customers with this user id
		//= database.Debug().Model(&user).Association("Customers").Find(&customers)

		var users []models.User
		// Fetch all users associated with this customer podio id in "users" table
		err = database.Debug().Model(&customer).Association("Users").Find(&users)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		// First check if user associates with multiple customers
		// If true, do NOT delete him
		// If false, delete him completely

		// Delete firebase users that are only associated with this customer
		// Keep those users that have one association, to delete them after
		var associatedCustomers int
		var usersToDelete []models.User
		for _, user := range users {
			associatedCustomers = int(database.Debug().Model(&user).Association("Customers").Count())

			if associatedCustomers == 1 {
				usersToDelete = append(usersToDelete, user)
			}
		}

		// Clear all customer's associations from "customer_user" table (deletes all associations with customer_id = this customer's id)
		err = database.Debug().Model(&customer).Association("Users").Clear()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//
		// Delete permissions from "casbin_rule" table in DB, using casbin
		// For example, for customer 1996, delete "portal::data::1996::finance" and "portal::data::1996::performance" permissions, if exist
		//

		for _, user := range users {
			// Remove financial policy for user, if exists
			if enforcer.HasPolicy(user.Id, fmt.Sprintf("portal::data::%d::finance", customer.Id), "read") {
				_, err = enforcer.RemovePolicy(user.Id, fmt.Sprintf("portal::data::%d::finance", customer.Id), "read")
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					log.Println(err)
					return
				}
			}
			// Remove performance policy for user, if exists
			if enforcer.HasPolicy(user.Id, fmt.Sprintf("portal::data::%d::performance", customer.Id), "read") {
				_, err = enforcer.RemovePolicy(user.Id, fmt.Sprintf("portal::data::%d::performance", customer.Id), "read")
				if err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					log.Println(err)
					return
				}
			}
		}

		// Delete customer from "customers" table
		err = database.Debug().Delete(&customer).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		//
		// If there are users that need to be completely deleted
		//
		if usersToDelete != nil {

			// Delete those users from "users" table
			err = database.Debug().Delete(&usersToDelete).Error
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				log.Println(err)
				return
			}

			//
			// Completely delete user (or users) from "casbin_rule" table in DB, using casbin
			// Also delete user (or users) from firebase
			//
			// Initialize firebaseAuth
			firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)

			// Foreach user of deleted customer
			for _, user_to_delete := range usersToDelete {
				// Delete user from "casbin_rule"
				_, err = enforcer.DeleteUser(user_to_delete.Id)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					log.Println(err)
					return
				}

				// Delete user from firebase
				err = firebaseAuth.DeleteUser(context.Background(), user_to_delete.Id)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					log.Println(err)
					return
				}
			}

		}

		c.JSON(http.StatusOK, nil)
	}
}

func UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": fmt.Sprintf("Invalid request body: %s", err.Error()),
			})
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

		// Update customer's name
		err = database.Debug().Save(&customer).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func GetCustomerUsers(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer models.Customer
		err := c.ShouldBindUri(&customer)
		if err != nil {
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

		var users []struct {
			Id                   string `json:"id" db:"id"`
			Email                string `json:"email" db:"email"`
			HasPerformanceAccess bool   `json:"has_performance_access"`
			HasFinancialAccess   bool   `json:"has_financial_access"`
		}

		err = database.Debug().
			Table("users").
			Joins("JOIN customer_user ON customer_user.user_id = users.id").
			Joins("JOIN customers ON customers.id = customer_user.customer_id").
			Where("customers.id = ?", customer.Id).
			Select("users.id, users.email, 0 AS has_performance_access, 0 AS has_financial_access").
			Scan(&users).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		for i, user := range users {

			permissions, err := enforcer.GetImplicitPermissionsForUser(user.Id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not find permissions for user"})
				log.Println(err.Error())
				return
			}

			if permissions != nil {
				for _, permission := range permissions {

					if permission[1] == fmt.Sprintf("portal::data::%d::performance", customer.Id) {
						users[i].HasPerformanceAccess = true
					}
					if permission[1] == fmt.Sprintf("portal::data::%d::finance", customer.Id) {
						users[i].HasFinancialAccess = true
					}
				}
			}

		}

		c.JSON(http.StatusOK, users)
	}
}
