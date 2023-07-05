package handlers

import (
	db "backend/database"
	"backend/models"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
)

func GetUnassignedUsers() gin.HandlerFunc {
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

		// Initialize return data
		var unassignedUserEmails []string

		err = database.
			Debug().
			Table("users").
			Where("employee_id IS NULL AND users.id NOT IN ( SELECT user_id FROM customer_user )").
			Select("email").
			Scan(&unassignedUserEmails).Error
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, unassignedUserEmails)
	}
}

func GetUsersEmails() gin.HandlerFunc {
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

		// Initialize return data
		var userEmails []string

		err = database.
			Debug().
			Table("users").
			Where("users.employee_id IS NULL AND users.email NOT LIKE '%@digitalminds.com%'").
			Select("email").
			Scan(&userEmails).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, userEmails)
	}
}

func GetAllUsers() gin.HandlerFunc {
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

		var users []struct {
			ID       string `json:"id"`
			FullName string `json:"full_name"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}

		// Get all employees and their roles (excluding those who haven't been assigned a firebase user yet)
		err = database.Debug().Table("employees").Joins("JOIN users ON users.employee_id = employees.id").
			Joins("LEFT JOIN casbin_rule ON casbin_rule.v0 = users.id").
			Select("users.id, employees.full_name, users.email, casbin_rule.v1 AS role").
			Where("users.id LIKE ? OR employees.full_name LIKE ? OR users.email LIKE ? OR casbin_rule.v1 LIKE ?", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%").
			Find(&users).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func SyncUsersWithFirebase() gin.HandlerFunc {
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

		// iterate firebase users
		// Note, behind the scenes, the Roles() iterator will retrieve 1000 Roles at a time through the API
		firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)

		iter := firebaseAuth.Users(c, "")
		for {
			user, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users from Firebase"})
				return
			}

			// Update users
			err = database.Debug().Omit("EmployeeID").Save(&models.User{Id: user.UID, Email: user.Email, CreationTimestamp: int(user.UserMetadata.CreationTimestamp), LastLoginTimestamp: int(user.UserMetadata.LastLogInTimestamp)}).Error
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users from Firebase"})
				return
			}
		}

		c.JSON(http.StatusOK, nil)

	}
}
