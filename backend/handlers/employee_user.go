package handlers

import (
	db "backend/database"
	"backend/models"
	"context"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AddAssociation() gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			UserEmail  string `json:"user_email"`
			EmployeeId int    `json:"employee_id"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		var employee = models.Employee{Id: requestBody.EmployeeId}
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
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		// Add association to users table in database
		err = database.Model(&employee).Association("Users").Append(&user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func DeleteAssociation() gin.HandlerFunc {
	return func(c *gin.Context) {
		//From postman's Body/form-data
		var requestBody struct {
			UserId     string `json:"user_id"`
			EmployeeId int    `json:"employee_id"`
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

		// Initialize employee with given id
		var employee = models.Employee{Id: requestBody.EmployeeId}
		// Initialize user with given id
		var user = models.User{Id: requestBody.UserId}

		// Remove association from users table in database
		err = database.Model(&employee).Association("Users").Delete(&user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

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

		c.JSON(http.StatusOK, nil)
	}
}
