package controller

import (
	"github.com/gin-gonic/gin"
	"golang/modul/transaction/model"
	"strconv"
)

func GetAllTransactions(c *gin.Context) {
	var transactions []model.Transaction
	if err := model.DB.Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success",
		"error":   "Tidak Ada Error",
		"data":    &transactions,
	})
}

func GetAllTransactionsByParam(c *gin.Context) {
	var transactions []model.Transaction
	var count int64
	id := c.Param("id")
	end := c.Param("end")
	page, _ := strconv.Atoi(c.DefaultQuery("page", id))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", end))
	if err := model.DB.Where("id = ?", id).Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}

	model.DB.Model(&transactions).Count(&count)
	model.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions)
	if err := model.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}
	c.JSON(200, gin.H{
		"code":            200,
		"message":         "Success",
		"error":           "Tidak Ada Error",
		"data":            &transactions,
		"Jumlah All Page": count,
	})
}
func GetTransactionByStatus(c *gin.Context) {
	var transactions []model.Transaction
	status := c.Param("status")
	if err := model.DB.Where("status = ?", status).Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return

	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success",
		"error":   "Tidak Ada Error",
		"data":    &transactions,
	})
}
