package main

import (
	"golang/modul/account/auth"
	"golang/modul/account/controller"
	"golang/modul/account/model"

	"github.com/gin-gonic/gin"
)

// func GetAllTransactions(c *gin.Context) {
// 	var transactions []Transaction
// 	var count int64
// 	id := c.Param("id")
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", id))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

// 	DB.Model(&transactions).Count(&count)
// 	DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions)
// 	c.JSON(200, gin.H{"message": &transactions})
// }

// func GetTransactionByStatus(c *gin.Context) {
// 	var transactions []Transaction
// 	status := c.Param("status")
// 	if err := DB.Where("status = ?", status).Find(&transactions).Error; err != nil {
// 		c.JSON(500, gin.H{"message": "Error"})
// 		return

// 	}
// 	c.JSON(200, gin.H{"message": &transactions})
// }

func main() {
	model.ConnectDatabase()
	r := gin.Default()
	Admin := r.Group("/api", auth.MiddlewareAdmin)
	{
		Admin.GET("/data-user", controller.GetDataUser)
		Admin.GET("/data-user/:id", controller.GetDataUserById)
		Admin.PUT("/data-user/:id", controller.EditDataUser)
		Admin.DELETE("/data-user/:id", controller.DeleteDataUser)

		// Admin.GET("/get-transactions/:id", controller.GetAllTransactions)
		// Admin.GET("/get-transaction/:status", controller.GetTransactionByStatus)
		Admin.GET("/logout", controller.Logout)
	}
	r.POST("/create-user", controller.CreateAccount)
	r.POST("/login", controller.Login)

	r.Run(":8081")
}
