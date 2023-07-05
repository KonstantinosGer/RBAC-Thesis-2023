package middleware

import (
	db "backend/database"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
)

func UpdateUsersDB() gin.HandlerFunc {
	return func(c *gin.Context) {

		db, err := db.Connect()
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer db.Close()

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
				//log.Fatalf("error listing users: %s\n", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users from Firebase"})
				return
			}
			//log.Printf("read user user: %v\n", user.DisplayName)

			_, err = db.Exec("REPLACE INTO users (id, email, creation_timestamp, last_login_timestamp) VALUES (?, ?, ?, ?)", user.UID, user.Email, user.UserMetadata.CreationTimestamp, user.UserMetadata.LastLogInTimestamp)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch users from Firebase"})
				return
			}
			//log.Println(result)
		}

	}
}
