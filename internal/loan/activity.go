package loan

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"poc-mini-loan/models"
)

type LoanActivities struct {
	DB *sql.DB
}

func (a *LoanActivities) SaveToDatabase(ctx context.Context, req models.LoanRequest, status string) error {
	log.Printf("🐘 [Activity: SaveToDatabase] -> Saving transaction status [%s] from %s...", status, req.ShopName)

	query := `
		INSERT INTO loans (merchant_id, shop_name, principal_amount, current_balance, status, updated_at)
		VALUES ($1, $2, $3, $3, $4, NOW())
		ON CONFLICT (merchant_id) 
		DO UPDATE SET status = $4, updated_at = NOW();
	`
	_, err := a.DB.ExecContext(ctx, query, req.MerchantID, req.ShopName, req.PrincipalAmount, status)
	if err != nil {
		log.Printf("❌ [Activity: SaveToDatabase] Query faield : %v", err)
		return err
	}

	log.Printf("[Activity: SaveToDatabase] Save success")
	return nil
}

func (a *LoanActivities) PreScreening(ctx context.Context, req models.LoanRequest) (bool, error) {
	log.Printf("🕵️ [Activity: Pre-Screening] -> Pre-Screen loans: %s", req.ShopName)

	if req.PrincipalAmount > 1000000 {
		log.Printf("❌ [Pre-Screening] The shop %s does not eligible for loan (Loan amount more than 1,000,000)", req.ShopName)
		return false, nil
	}

	time.Sleep(30 * time.Second)
	log.Printf("[Pre-Screening] The shop %s pre-screen passed", req.ShopName)
	return true, nil
}

func (a *LoanActivities) DisburseFunds(ctx context.Context, req models.LoanRequest) error {
	time.Sleep(30 * time.Second)
	log.Printf("💸 [Activity: DisburseFunds] -> Transfers loan %.2f฿ to account %s completed", req.PrincipalAmount, req.ShopName)
	return nil
}

func (a *LoanActivities) GetDailySales(ctx context.Context, merchantID string) (float64, error) {
	log.Printf("[Activity: GetDailySales] -> Get daily sales from %s", merchantID)
	mockDailySales := 10000.00
	return mockDailySales, nil
}

func (a *LoanActivities) DeductLoanAmount(ctx context.Context, merchantID string, amount float64) (float64, error) {
	log.Printf("[Activity: DeductRepayment] -> Deduct amount %.2f฿ from %s for deduct loan amount", amount, merchantID)

	if err := a.DB.PingContext(ctx); err != nil {
		log.Printf("❌ [Activity: DeductRepayment] Database connection unavailable! (Ping failed): %v", err)

		return 0, fmt.Errorf("database is down: %w", err)
	}

	query := `
    UPDATE loans 
    SET current_balance = CASE 
        WHEN current_balance - $2 < 0 THEN 0 
        ELSE current_balance - $2 
    END, 
    updated_at = NOW() 
    WHERE merchant_id = $1
    RETURNING current_balance;
	`
	var newBalance float64
	err := a.DB.QueryRowContext(ctx, query, merchantID, amount).Scan(&newBalance)
	return newBalance, err
}
