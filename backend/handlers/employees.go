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

func GetAllEmployees() gin.HandlerFunc {
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

		var employees []models.Employee

		// Select all employees
		err = database.Debug().Select("id", "full_name").Where("id LIKE ? OR full_name LIKE ?", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%").Find(&employees).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, employees)
	}
}

func AddEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee models.Employee
		if err := c.ShouldBindJSON(&employee); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		// Make sure given employee name is not empty
		if strings.TrimSpace(employee.FullName) == "" {
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

		//var employee = models.Employee{FullName: requestBody.FullName}
		// Add a new employee to employees table in database
		err = database.Debug().Create(&employee).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func DeleteEmployee(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee models.Employee
		err := c.ShouldBindUri(&employee)
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

		var users []models.User
		// Fetch all users associated with this employee id in "users" table
		err = database.Debug().Model(&employee).Association("Users").Find(&users)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		// If employee was associated with a user
		if len(users) > 0 {
			// Clear employee's associations from "users" table (sets employee_id as null)
			err = database.Debug().Model(&employee).Association("Users").Clear()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				log.Println(err)
				return
			}
			// Delete users from "users" table
			err = database.Debug().Delete(&users).Error
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				log.Println(err)
				return
			}

			//
			// Delete user's (or users') permissions from "casbin_rule" table in DB, using casbin
			// Also delete user's (or users') from firebase
			//
			// Initialize firebaseAuth
			firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)

			// Foreach user of deleted customer
			for _, user := range users {
				// Delete user from "casbin_rule"
				_, err = enforcer.DeleteUser(user.Id)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					log.Println(err)
					return
				}

				// Delete user from firebase
				err = firebaseAuth.DeleteUser(context.Background(), user.Id)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					log.Println(err)
					return
				}
			}

		}

		// Delete employee with id = employeeId, from "employees" table
		// First way
		//err = database.Debug().Delete(&models.Employee{Id: employeeId}).Error
		err = database.Debug().Delete(&employee).Error
		// Second way
		//err = database.Debug().Delete(&models.Employee{}, employeeId).Error
		// Third way
		//err = database.Debug().Where("id = ?", employee.Id).Delete(&models.Employee{}).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func UpdateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee models.Employee
		if err := c.ShouldBindJSON(&employee); err != nil {
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
		err = database.Debug().Save(&employee).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func GetEmployeeUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var employee models.Employee
		err := c.ShouldBindUri(&employee)
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

		// Initialize return data
		var associatedUsers []struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		}

		// Get all users associated with this employee
		err = database.Debug().Table("users").
			Joins("JOIN employees ON users.employee_id = employees.id").
			Select("users.id, users.email").
			Where("employees.id = ?", employee.Id).
			Find(&associatedUsers).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, associatedUsers)
	}
}
