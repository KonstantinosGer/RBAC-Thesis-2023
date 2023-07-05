package handlers

import (
	db "backend/database"
	"backend/models"
	"context"
	"firebase.google.com/go/auth"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetAllFirebaseUsers() gin.HandlerFunc {
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

		var users []models.User

		// Select all customers
		err = database.Debug().Select("id", "email").Where("id LIKE ? OR email LIKE ?", "%"+filters.Keyword+"%", "%"+filters.Keyword+"%").Find(&users).Error
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func AddFirebaseUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize user to add with given email and password
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
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

		// Initialize firebaseAuth
		firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)

		// Create new user in firebase with the following parameters
		params := (&auth.UserToCreate{}).
			Email(requestBody.Email).
			Password(requestBody.Password).
			Disabled(false)
		// If email already exists in firebase, return error
		firebaseUser, err := firebaseAuth.CreateUser(context.Background(), params)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "User already exists in Firebase! Please use another email."})
			log.Println(err)
			return
		}

		// Add user to users table in database
		var user = models.User{Id: firebaseUser.UserInfo.UID, Email: requestBody.Email}

		err = database.Debug().Create(&user).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func DeleteFirebaseUser(enforcer *casbin.SyncedEnforcer) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Bind from url, passing url's passed id directly to a user object with id = url's passed id
		var user models.User
		err := c.ShouldBindUri(&user)
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

		// Remove all user's permissions from casbin rule
		_, err = enforcer.DeleteUser(user.Id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			log.Println(err.Error())
			return
		}

		//Delete all user's associations, if exists, from "customer_user" table in DB
		err = database.Debug().Model(&user).Association("Customers").Clear()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		// Delete user from "users" table
		// In this way employee-user association is deleted too, if exists
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
