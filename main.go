package main

import (
	"golang/modul/transaction/auth"
	"golang/modul/transaction/controller"
	"golang/modul/transaction/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetDataUser(c *gin.Context) {
	var account []model.Account
	err := model.DB.Find(&account).Error
	if err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}
	c.JSON(200, gin.H{"message": &account})
}
func GetDataUserById(c *gin.Context) {
	var account model.Account
	id := c.Param("id")
	err := model.DB.Find(&account, id).Error
	if err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}
	c.JSON(200, gin.H{"message": &account})
}
func EditDataUser(c *gin.Context) {
	var account model.Account
	id := c.Param("id")
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
			"error":   err,
			"data":    &account})
		return
	}
	if err := model.DB.Where("id = ?", id).Updates(&account).Error; err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
			"error":   err,
			"data":    &account})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success Updated",
		"error":   "",
		"data":    &account})

}
func DeleteDataUser(c *gin.Context) {
	var account *model.Account
	id := c.Param("id")
	if err := model.DB.Where("id = ?", id).Delete(&account).Error; err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Bad Request",
			"error":   err,
			"data":    ""})
		return
	}
	c.JSON(200, gin.H{"message": "Delete Success"})
}
func CreateAccount(c *gin.Context) {
	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
			"error":   err,
			"data":    &account})
		return
	}
	HashPassword, err := bcrypt.GenerateFromPassword([]byte((account.Password)), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Bad Request",
			"error":   err,
			"data":    &account})
		return
	}
	account.Password = string(HashPassword)
	if err := model.DB.Create(&account).Error; err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "Bad Request",
			"error":   err,
			"data":    ""})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success Created",
		"error":   err,
		"data":    &account})
}

func Login(c *gin.Context) {
	var account model.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}
	var getAcount model.Account
	if err := model.DB.Where("name = ?", account.Name).First(&getAcount).Error; err != nil {
		c.JSON(401, gin.H{"message": "Username salah"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(getAcount.Password), []byte(account.Password)); err != nil {
		c.JSON(401, gin.H{"message": "Password salah"})
		return
	}
	expiredTime := time.Now().Add(time.Minute * 30)
	claims := &auth.JWT{
		Id:   getAcount.Id,
		Name: getAcount.Name,
		Role: getAcount.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "golang",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenGenerate, err := token.SignedString(auth.JWT_KEY)
	if err != nil {
		c.JSON(401, gin.H{"message": "Login Error"})
	}
	c.SetCookie("gin_cookie", tokenGenerate, 3600, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Login Berhasil"})
}
func Logout(c *gin.Context) {
	c.SetCookie("gin_cookie", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Berhasil Logout"})
}

func main() {
	model.ConnectDatabase()
	r := gin.Default()
	Admin := r.Group("/api", auth.MiddlewareAdmin)
	{
		Admin.GET("/data-user", GetDataUser)
		Admin.GET("/data-user/:id", GetDataUserById)
		Admin.POST("/create-user", CreateAccount)
		Admin.PUT("/data-user/:id", EditDataUser)
		Admin.DELETE("/data-user/:id", DeleteDataUser)

		Admin.GET("/get-transactions", controller.GetAllTransactions)
		Admin.GET("/get-transaction/:id", controller.GetAllTransactionsByParam)
		Admin.GET("/get-transaction-page/:id/:end", controller.GetAllTransactionsByParam)
		Admin.GET("/get-transactions/:status", controller.GetTransactionByStatus)
		Admin.GET("/logout", Logout)
	}
	r.POST("/login", Login)

	r.Run()
}
