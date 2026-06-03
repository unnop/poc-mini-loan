package models

type LoanRequest struct {
	MerchantID      string  `json:"merchant_id" binding:"required"`
	ShopName        string  `json:"shop_name" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required"`
}

type LoanResponse struct {
	Message    string `json:"message"`
	WorkflowID string `json:"workflow_id"`
	RunID      string `json:"run_id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

type ApprovalRequest struct {
	MerchantID string `json:"merchant_id"`
	Action     string `json:"action"`
}
