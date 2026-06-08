package loan

import (
	"fmt"
	"log"
	"time"

	"poc-mini-loan/models"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MerchantLoanWorkflow(ctx workflow.Context, req models.LoanRequest) (string, error) {

	options := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var acts *LoanActivities
	logger := workflow.GetLogger(ctx)
	logger.Info("[Workflow] Starting merchant loan workflow process...")

	_ = workflow.ExecuteActivity(ctx, acts.SaveToDatabase, req, "PENDING_PRE_SCREENING").Get(ctx, nil)

	var isPassed bool
	err := workflow.ExecuteActivity(ctx, acts.PreScreening, req).Get(ctx, &isPassed)
	if err != nil {
		return "PRE_SCREENING_ERROR", err
	}

	if !isPassed {
		_ = workflow.ExecuteActivity(ctx, acts.SaveToDatabase, req, "REJECTED").Get(ctx, nil)
		return "LOAN_REJECTED_BY_SYSTEM", nil
	}

	_ = workflow.ExecuteActivity(ctx, acts.SaveToDatabase, req, "WAITING_FOR_APPROVAL").Get(ctx, nil)
	logger.Info("[Workflow] Pre-screen passed. Workflow is sleeping, waiting for manager approval via Signal...")

	var signalResult string
	signalChan := workflow.GetSignalChannel(ctx, "loan-approval-signal")

	signalChan.Receive(ctx, &signalResult)

	logger.Info("[Workflow] Woke up! Received signal outcome: " + signalResult)

	if signalResult != "APPROVED" {
		_ = workflow.ExecuteActivity(ctx, acts.SaveToDatabase, req, "REJECTED_BY_MANAGER").Get(ctx, nil)
		return "LOAN_REJECTED_BY_MANAGER", nil
	}

	err = workflow.ExecuteActivity(ctx, acts.DisburseFunds, req).Get(ctx, nil)
	if err != nil {
		return "DISBURSEMENT_FAILED", err
	}

	_ = workflow.ExecuteActivity(ctx, acts.SaveToDatabase, req, "ACTIVE").Get(ctx, nil)

	remainingBalance := req.PrincipalAmount
	logger.Info("🔄 [Workflow] Entering simulation repayment loop (Deducting every 1 minute)...")

	var round int
	for remainingBalance > 0 {
		logger.Info("[Workflow] Sleeping until next day cycle...")
		_ = workflow.Sleep(ctx, 15*time.Second)

		if round >= 3 {
			return "FAILED_BY_BUSINESS_RULE", fmt.Errorf("auto repayment limit exceeded 3 rounds, shifting to manual process")
		}

		var dailySales float64
		_ = workflow.ExecuteActivity(ctx, acts.GetDailySales, req.MerchantID).Get(ctx, &dailySales)

		deductLoan := dailySales * 0.10

		if deductLoan > remainingBalance {
			deductLoan = remainingBalance
		}

		var newDBBalance float64
		_ = workflow.ExecuteActivity(ctx, acts.DeductLoanAmount, req.MerchantID, deductLoan).Get(ctx, &newDBBalance)

		remainingBalance = newDBBalance

		if remainingBalance <= 0 {
			break
		}
		round++

		log.Printf("Remaining balance left in state memory: %.2f฿", remainingBalance)
	}

	logger.Info("[Workflow] Loan fully paid back! Closing loan application contract...")
	_ = workflow.ExecuteActivity(ctx, acts.SaveToDatabase, req, "CLOSED").Get(ctx, nil)

	return "LOAN_FULLY_REPAYED_AND_CLOSED", nil
}
