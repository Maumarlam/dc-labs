package api

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	////////IMPORTS DE CADA FILE REST REQUEST
	"github.com/Maumarlam/dc-labs/challenges/final/controller"

	"github.com/gin-gonic/gin"
)

type user struct {
	username  string
	password  string
	tokenAuth string
}

var userList []user //To store all the users created

// var data = gin.H{
// 	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433", "lastToken": ""},
//   "mau":    gin.H{"email": "xxx@bar.com", "phone": "321456", "lastToken": ""}
// }

//Checar parametros que necesita para empezar
func Start() {
	log.Printf("Entro apistart")
	r := gin.Default()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"username": "password",
		"mau":      "mau123",
	}))

	authorized.GET("/login", login)

	r.GET("/login", login)
	r.GET("/logout", logout)
	r.GET("/status", status)
	r.POST("/upload", uploadImage)
	r.GET("/status/:worker", workerStatus)
	r.Run()
}

func tokenGen() string {
	x := make([]byte, 8)
	rand.Read(x)
	return fmt.Sprintf("%x", x)
}

func logout(c *gin.Context) {
	isThere := authentication(c)
	var index int
	var userID string
	inList := false
	for i, j := range userList {
		if j.tokenAuth == isThere {
			index = i
			userID = j.username
			inList = true
		}
	}

	if inList == true {
		userList[index] = userList[len(userList)-1] //como  un delete del  array
		userList = userList[:len(userList)-1]
		c.JSON(http.StatusOK, gin.H{"Message": "Bye " + userID + ", your token has been revoked"})
	}
}

func uploadImage(c *gin.Context) {
	isThere := authentication(c)
	inList := false
	var userToken string
	for _, i := range userList {
		if isThere == i.tokenAuth {
			//Have to get the token
			userToken = i.tokenAuth
			inList = true
		}
	}
	if inList == true {
		file, err := c.FormFile("data") //funcion del tutorial
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"Message": "Error uploading image", "error": err.Error})
			return
		}
		timeATM := time.Now()
		timeUpdated := timeATM.Format("2006-01-02 15:04:05")
		upImage := filepath.Base(file.Filename) //Get the name of the file location

		newFile := userToken + "_" + timeUpdated + "_" + upImage
		if err := c.SaveUploadedFile(file, newFile); err != nil {
			c.JSON(http.StatusOK, gin.H{"Message": "Error while uploading image", "error": err.Error})
			return
		}
		//If the image saves   successfully
		c.JSON(http.StatusOK, gin.H{"Message": "Image successfully uploaded", "filename": upImage, "size": strconv.FormatInt(file.Size, 10) + "bytes"})
	}
}

func workerStatus(c *gin.Context) {
	worker := c.Param("worker") //How does this work
	tags, status, usage, url, port := controller.WorkerStatus(worker)
	//Get the status of  the  worker withall its attributes
	c.JSON(http.StatusOK, gin.H{"Worker": worker, "Tags": tags, "Status": status, "Usage": usage, "URL": url, "Port": port})

}

func status(c *gin.Context) {
	isThere := authentication(c)
	for _, i := range userList {
		if isThere == i.tokenAuth {
			timeATM := time.Now()
			c.JSON(http.StatusOK, gin.H{"Message": "Hi " + i.username + ", system is running", "time": timeATM.Format("2006-01-02 15:04:05")})
		}
	}
}

func authentication(c *gin.Context) string {
	authToken := c.Request.Header["Authorization"] ///DUDAaaa
	fullToken := authToken[0]
	splitToken := strings.Split(fullToken, " ")
	return string(splitToken[1])
}

func login(c *gin.Context) {
	var newUser user
	inList := false
	userID, password, _ := c.Request.BasicAuth()
	userToken := tokenGen()

	for i, j := range userList {
		if userID == j.username {
			userList[i].tokenAuth = userToken
			c.JSON(http.StatusOK, gin.H{"Message": "New token generated", "token": userToken})
			inList = true
		}
	}
	if inList != true { //IF the user wasn't on the list we need to register it
		newUser.username = userID
		newUser.password = password
		newUser.tokenAuth = userToken
		userList = append(userList, newUser) //APPEND FUNCIONA???
		c.JSON(http.StatusOK, gin.H{"Message": "Hi " + userID + ", welcome", "token": userToken})
	}
}
