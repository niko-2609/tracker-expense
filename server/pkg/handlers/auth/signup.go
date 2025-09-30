package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/niko-2609/tracker-expense/database"
	authmodel "github.com/niko-2609/tracker-expense/models/auth"
	apiModel "github.com/niko-2609/tracker-expense/models/common/api"
	txnModel "github.com/niko-2609/tracker-expense/models/transaction"
	"github.com/niko-2609/tracker-expense/pkg/validation"
	"github.com/niko-2609/tracker-expense/utils"
	"gorm.io/gorm"
)

func SignUp(c *fiber.Ctx) error {
	input := new(authmodel.Credentials)

	// Verify raw request
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		log.Error("Bad Request - Invalid Credentials Object")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Invalid Request",
			Data:    nil,
		})
	}

	// Verify request against the defined model
	if err := validation.Validate.Struct(input); err != nil {
		log.Error("Bad Request - Request validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Failed request validation",
			Data:    nil,
		})
	}

	email := input.Email
	password := input.Password

	var userModel *authmodel.User

	// Check if user exists in DB
	userModel, err := utils.GetUserByEmail(email)

	if err != nil {
		// For all other errors that `RecordNotFound`, we return an error.
		// This is done because `RecordNotFound` is not an error in this case and is valid
		// Since we will proceed if a record for the user does not exist in DB.
		if userModel == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// For all other errors, return 500
			log.Error("An unexpected error occurred")
			return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
				Status:  "error",
				Message: "Internal server error",
				Data:    nil,
			})
		}
	}

	// If no errors and a valid user is present, return 409
	if userModel != nil {
		log.Warn("User already exists, should be redirected to /login")
		return c.Status(fiber.StatusConflict).JSON(apiModel.Response{
			Status:  "error",
			Message: "User exists, try signing in",
			Data:    nil,
		})
	}

	if !utils.IsEmail(email) {
		log.Error("Email provided is not valid")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Invalid email or password, please try again",
			Data:    nil,
		})
	}

	//  Encrypt password
	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		log.Error("Error encrypting password")
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: "Unable to sign up user, please try again",
			Data:    nil,
		})
	}

	// // Get Username from Email
	// username, err := utils.ExtractUserName(email)
	// if err != nil {
	// 	log.Error(err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
	// 		Status:  "error",
	// 		Message: "Unable to sign up user, provide a valid email",
	// 		Data:    nil,
	// 	})
	// }

	userModel = &authmodel.User{
		Email:    email,
		Password: hashedPass,
		Username: utils.ExtractUserName(email),
	}

	// Save user to DB
	if err := database.DB.Create(userModel).Error; err != nil {
		log.Error("Unable to create user in DB")
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: "Internal server error, please try again",
			Data:    nil,
		})
	}

	// 1. Insert transaction
	txn := txnModel.Transaction{
		UserID:     1, // Example: get from JWT/session
		Name:       "Test expense",
		Amount:     57,
		TxnType:    "expense",
		CategoryID: uint(8),
		TxnDate:    time.Now(),
		Frequency:  "daily",
	}

	if err := database.DB.Create(&txn).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Debug("Successfully saved transaction")

	err = UpdateDashboardMetrics(uint(1))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(apiModel.Response{
		Status:  "success",
		Message: "Sign up successfull",
		Data:    nil,
	})

}

func UpdateDashboardMetrics(userID uint) error {
	// A. Total Income / Expense / Net Savings
	totalQuery := `
    WITH totals AS (
        SELECT
            SUM(CASE WHEN txn_type='income' THEN amount ELSE 0 END) AS total_income,
            SUM(CASE WHEN txn_type='expense' THEN amount ELSE 0 END) AS total_expense
        FROM transactions
        WHERE user_id = ?
    )
    INSERT INTO user_dashboard_metrics (user_id, total_income, total_expense, net_savings, updated_at)
    SELECT ?, total_income, total_expense, total_income - total_expense, NOW()
    FROM totals
    ON CONFLICT (user_id) DO UPDATE
    SET total_income = EXCLUDED.total_income,
        total_expense = EXCLUDED.total_expense,
        net_savings = EXCLUDED.net_savings,
        updated_at = NOW();
    `
	if err := database.DB.Exec(totalQuery, userID, userID).Error; err != nil {
		return err
	}

	// // B. Monthly Totals (for line chart)
	// monthlyQuery := `
	// WITH monthly AS (
	//     SELECT
	//         TO_CHAR(DATE_TRUNC('month', txn_date), 'YYYY-MM') AS month,
	//         SUM(CASE WHEN txn_type='income' THEN amount ELSE -amount END) AS net
	//     FROM transactions
	//     WHERE user_id = ?
	//     GROUP BY DATE_TRUNC('month', txn_date)
	//     ORDER BY month
	// )
	// UPDATE user_dashboard_metrics
	// SET monthly_totals = (SELECT jsonb_object_agg(month, net) FROM monthly)
	// WHERE user_id = ?;
	// `
	// if err := db.Exec(monthlyQuery, userID, userID).Error; err != nil {
	// 	return err
	// }

	// // C. Top 5 Expense Categories (for pie chart)
	// topCatQuery := `
	// UPDATE user_dashboard_metrics
	// SET top_expense_categories = (
	//     SELECT jsonb_agg(jsonb_build_object('category', c.name, 'amount', SUM(t.amount)))
	//     FROM transactions t
	//     JOIN categories c ON t.category_id = c.id
	//     WHERE t.user_id = ? AND t.txn_type='expense'
	//     GROUP BY c.id, c.name
	//     ORDER BY SUM(t.amount) DESC
	//     LIMIT 5
	// )
	// WHERE user_id = ?;
	// `
	// if err := db.Exec(topCatQuery, userID, userID).Error; err != nil {
	// 	return err
	// }

	return nil
}
