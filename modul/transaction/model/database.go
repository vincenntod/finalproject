package model

import "time"

type Account struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type Transaction struct {
	Id               int       `json:"id"`
	OdaNumber        string    `json:"oda_number"`
	BankAccountNo    string    `json:"bank_account_no"`
	BillingCycleDate time.Time `json:"billing_cycle_date"`
	PaymentDueDate   time.Time `json:"payment_due_date"`
	OverflowAmount   float32   `json:"overflow_amount"`
	BillAmount       float32   `json:"bill_amount"`
	PrincipalAmount  float32   `json:"principal_amount"`
	InterestAmount   float32   `json:"interest_amount"`
	TotalFeeAmount   float32   `json:"total_fee_amount"`
	Status           string    `json:"status"`
	PaymentMethod    string    `json:"payment_method"`
	AutoDebetCounter int       `json:"auto_debet_counter"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IsHold           bool      `json:"is_hold"`
	IsFstlPending    bool      `json:"is_fstl_pending"`
	IsHstlPending    bool      `json:"is_hstl_pending"`
	IsLaaPositif     bool      `json:"is_laa_positif"`
	PaymentAmount    float32   `json:"payment_amount"`
}
