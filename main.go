package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Account struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

var JWT_KEY = []byte("dbceria")

type Transaction struct {
	Id 	  int       `json:"id"`
	OdaNumber string `json:"oda_number"`
	BankAccountNo string `json:"bank_account_no"`
	BillingCycleDate time.Time `json:"billing_cycle_date"`
	PaymentDueDate time.Time `json:"payment_due_date"`
	OverflowAmount float32 `json:"overflow_amount"`
	BillAmount float32 `json:"bill_amount"`
	PrincipalAmount float32 `json:"principal_amount"`
	InterestAmount float32 `json:"interest_amount"`
	TotalFeeAmount float32 `json:"total_fee_amount"`
	Status string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	AutoDebetCounter int `json:"auto_debet_counter"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsHold bool `json:"is_hold"`
	IsFstlPending bool `json:"is_fstl_pending"`
	IsHstlPending bool `json:"is_hstl_pending"`
	IsLaaPositif bool `json:"is_laa_positif"`
	PaumentAmount float32 `json:"paument_amount"`

}
type JWT struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=12345 dbname=dbceria port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	}})
	if err != nil {
		fmt.Printf("Error")
		return
	}
	DB = db
}

func GetDataUser(c *gin.Context) {
	var account []Account
	err := DB.Find(&account).Error
	if err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}
	c.JSON(200, gin.H{"message": &account})
}
func GetDataUserById(c *gin.Context) {
	var account Account
	id := c.Param("id")
	err := DB.Find(&account, id).Error
	if err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}
	c.JSON(200, gin.H{"message": &account})
}
func EditDataUser(c *gin.Context) {
	var account Account
	id := c.Param("id")
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Internal Server Error",
			"error":   err,
			"data":    &account})
		return
	}
	if err := DB.Where("id = ?", id).Updates(&account).Error; err != nil {
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
	var account *Account
	id := c.Param("id")
	if err := DB.Where("id = ?", id).Delete(&account).Error; err != nil {
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
	var account Account
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
	if err := DB.Create(&account).Error; err != nil {
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
func GetAllTransactions(c *gin.Context) {
	var transactions []Transaction
	var count int64
	id := c.Param("id")
	end := c.Param("end")
	page, _ := strconv.Atoi(c.DefaultQuery("page", id))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", end))
	if err := DB.Where("id = ?", id).Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}

	DB.Model(&transactions).Count(&count)
	DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions)
	if err := DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return
	}
	c.JSON(200, gin.H{"message": &transactions, "Jumlah All Page": count})
}
func GetTransactionByStatus(c *gin.Context) {
	var transactions []Transaction
	status := c.Param("status")
	if err := DB.Where("status = ?", status).Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"message": "Error"})
		return

	}
	c.JSON(200, gin.H{"message": &transactions})
}

func Login(c *gin.Context) {
	var account Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}
	var getAcount Account
	if err := DB.Where("name = ?", account.Name).First(&getAcount).Error; err != nil {
		c.JSON(401, gin.H{"message": "Username salah"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(getAcount.Password), []byte(account.Password)); err != nil {
		c.JSON(401, gin.H{"message": "Password salah"})
		return
	}
	expiredTime := time.Now().Add(time.Minute * 30)
	claims := &JWT{
		Id:   getAcount.Id,
		Name: getAcount.Name,
		Role: getAcount.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "golang",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenGenerate, err := token.SignedString(JWT_KEY)
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
func MiddlewareAdmin(c *gin.Context) {
	tokenString, err := c.Cookie("gin_cookie")
	if err != nil {
		c.AbortWithStatus(401)
		c.JSON(401, gin.H{"message": "Silahkan Login"})
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}

		return JWT_KEY, err
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(401)
			c.JSON(401, gin.H{"message": "Silahkan Login Kembali"})
		}
		var account Account
		DB.First(&account, claims["Id"])
		if account.Id == 0 {
			c.AbortWithStatus(401)
		}
		c.Set("account", account)
		c.Next()
	} else {
		c.AbortWithStatus(401)
	}

}
func main() {
	ConnectDatabase()
	r := gin.Default()
	Admin := r.Group("/api", MiddlewareAdmin)
	{
		Admin.GET("/data-user", GetDataUser)
		Admin.GET("/data-user/:id", GetDataUserById)
		Admin.POST("/create-user", CreateAccount)
		Admin.PUT("/data-user/:id", EditDataUser)
		Admin.DELETE("/data-user/:id", DeleteDataUser)

		Admin.GET("/get-transaction/:id", GetAllTransactions)
		Admin.GET("/get-transaction-page/:id/:end", GetAllTransactions)
		Admin.GET("/get-transactions/:status", GetTransactionByStatus)
		Admin.GET("/logout", Logout)
	}
	r.POST("/login", Login)

	r.Run()
}
