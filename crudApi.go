package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserDb struct {
	UserID string `json:"user_id"`
	User
}
type User struct {
	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}

type UpdateUserAccount struct {
	OldAccount string `json:"old_account"`
	NewAccount string `json:"new_account"`
}

type UpdateUserPassword struct {
	UserAccount string `json:"user_account"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// 將使用者的資料存在 slice 裡
var usersDB = []UserDb{
	{UserID: uuid.NewString(), User: User{UserAccount: "Leo", UserPassword: "abc"}},
	{UserID: uuid.NewString(), User: User{UserAccount: "Allen", UserPassword: "xyz"}},
	{UserID: uuid.NewString(), User: User{UserAccount: "Ivy", UserPassword: "123456"}},
}

var users []User // 使用者輸入時的格式依據

// CreateApi 新增
func CreateApi(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil { //BindJSON: 綁定提交的 JSON 參數信息
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, c.Errors)
		//c.String(404,err.Error())
		//c.AbortWithError(http.StatusNotFound, err.Error())
		return
	}

	// 確認資料是否符合格式
	if newUser.UserPassword != "" && newUser.UserAccount != "" {
		// 判斷帳號是否已存在
		for _, a := range usersDB {
			if a.UserAccount == newUser.UserAccount {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "帳號已存在"})
				return
			}
		}
		// 若帳號不存在，在資料庫裡新增它
		singUp(newUser)
		c.IndentedJSON(http.StatusCreated, newUser) // 顯示使用者剛剛創建的帳、密
		return
	} else {
		//若 UserAccount or UserPassword 的 key 或 value 格式有誤
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "UserAccount or UserPassword is null"})
		return
	}

}

func singUp(newData User) {
	users = []User{}
	users = append(users, newData)

	// 將使用者註冊的帳號寫入到 DB 裡
	id := uuid.NewString()
	newUser := UserDb{UserID: id, User: User{UserAccount: newData.UserAccount, UserPassword: newData.UserPassword}}
	usersDB = append(usersDB, newUser)
}

// ReadApi
// 讀取 撈所有資料(用 JSON 格式列出所有使用者們)
func ReadApi(c *gin.Context) {
	userID := c.Param("UserID")
	for _, a := range usersDB {
		if userID == a.UserID {
			c.AbortWithStatusJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
	return
}

func ReadAllApi(c *gin.Context) {
	// 顯示所有user
	c.AbortWithStatusJSON(http.StatusOK, usersDB)
	return
}

// UpdateAccountApi
// 更新: Update HTTP PUT or PATCH
func UpdateAccountApi(c *gin.Context) {
	var updateAccount UpdateUserAccount
	//updateUser.UserAccount = c.Request.FormValue("UserAccount")
	//updateUser.UserPassword = c.Request.FormValue("UserPassword")
	if err := c.BindJSON(&updateAccount); err != nil {
		c.IndentedJSON(http.StatusBadRequest, c.Errors)
		return
	}

	if updateAccount.OldAccount != "" && updateAccount.NewAccount != "" {
		for i, a := range usersDB {
			if updateAccount.OldAccount == a.UserAccount {
				a.UserAccount = updateAccount.NewAccount
				usersDB[i] = a
				c.JSON(http.StatusOK, a)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "有值為空"})
		return
	}
}

func UpdatePasswordApi(c *gin.Context) {
	var updatePassword UpdateUserPassword

	if err := c.BindJSON(&updatePassword); err != nil {
		c.IndentedJSON(http.StatusBadRequest, c.Errors)
		return
	}

	if updatePassword.UserAccount != "" && updatePassword.OldPassword != "" && updatePassword.NewPassword != "" {
		for i, a := range usersDB {
			if updatePassword.UserAccount == a.UserAccount && updatePassword.OldPassword == a.UserPassword {
				a.UserPassword = updatePassword.NewPassword
				usersDB[i] = a // 修改資料庫資料
				c.JSON(http.StatusOK, a)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "有值為空"})
		return
	}
}

// DeleteApi 刪除:Delete HTTP DELETE
func DeleteApi(c *gin.Context) {
	userID := c.Param("UserID")

	for i, a := range usersDB {
		if userID == a.UserID {
			// 修改資料庫資料
			usersDB = append(usersDB[:i], usersDB[i+1:]...)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "刪除成功"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
	return

}

// Login
// 登入 POST http://localhost:8080/login Body: {"UserID": "1"}
func Login(c *gin.Context) {
	var loginUser User

	if err := c.BindJSON(&loginUser); err != nil { //BindJSON: 綁定提交的 JSON 參數信息
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, c.Errors)
		return
	}
	// 確認資料是否符合格式
	if loginUser.UserAccount != "" {
		if loginUser.UserPassword != "" {
			// 比較 Body 內的資料是否存在 slice
			for _, a := range usersDB {
				if a.UserAccount == loginUser.UserAccount {
					if a.UserPassword == loginUser.UserPassword {
						c.IndentedJSON(http.StatusOK, a)
						return
					} else {
						c.IndentedJSON(http.StatusUnauthorized, gin.H{"": "密碼錯誤"})
						return
					}
				}
			}
		} else { //若 UserPassword 的 key 或 value 格式有誤
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "UserPassword is null"})
			return
		}
	} else { //若 UserAccount 的 key 或 value 格式有誤
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "UserAccount is null"})
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
}
