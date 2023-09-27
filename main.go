package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	server := gin.Default()

	//POST http://localhost:8080/Login Body -> {"user_account": "Leo", "user_password": "abc"}
	server.POST("/Login", Login) // 登入(讀取)
	//POST http://localhost:8080/creatUser Body -> {"user_account": "Tina", "user_password": "0123"}
	server.POST("/creatUser", CreateApi) // 註冊(新增)
	//GET http://localhost:8080/users
	server.GET("/users", ReadAllApi) // 讀取Users
	//GET http://localhost:8080/users/:ID Params -> 9d8207cb-5fa1-4a86-8085-8add7f1e6cf8
	server.GET("/users/:UserID", ReadApi) // 讀取User[i]
	//POST http://localhost:8080/updateAccount Body -> {"old_account": "Leo", "new_account": "Leo王"}
	server.POST("/updateAccount", UpdateAccountApi) // 更新帳號
	//POST http://localhost:8080/updatePassword Body -> {"user_account": "Allen", "old_password": "xyz", "new_password": "Dog"}
	server.POST("/updatePassword", UpdatePasswordApi) // 更新密碼
	//DELETE http://localhost:8080/users/:ID Params -> 9d8207cb-5fa1-4a86-8085-8add7f1e6cf8
	server.DELETE("/users/:ID", DeleteApi)

	if err := server.Run(":8080"); err != nil {
		log.Fatalln(err.Error())
		return
	}
}
