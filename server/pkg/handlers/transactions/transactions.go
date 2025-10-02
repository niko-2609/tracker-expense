package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/niko-2609/tracker-expense/database"
	apiModel "github.com/niko-2609/tracker-expense/models/common/api"
	transactionModels "github.com/niko-2609/tracker-expense/models/transaction"
	"github.com/niko-2609/tracker-expense/pkg/validation"
	"github.com/niko-2609/tracker-expense/utils"
)

// Fetch all transactions for a given user from the DB
func GetTransactionsHandler(c *fiber.Ctx) error {

	userID, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(apiModel.Response{
			Status:  "error",
			Message: "User id is required for the transaction",
			Data:    nil,
		})
	}
	// Local variable to handle request validation.
	var transactions []transactionModels.Transaction

	// Fetch transactions from DB
	result := database.DB.Model(&transactionModels.Transaction{}).Where(&transactionModels.Transaction{
		UserID: userID,
	}).Find(&transactions)
	// If error, return no data
	if result.Error != nil {
		log.Error(result.Error)
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: fmt.Sprintf("Cannot fetch transactions: %v", result.Error),
			Data:    nil,
		})
	}

	// Log success message and return the list of transactions.
	log.Debug("Retrieved transactions successfully:", transactions)
	return c.Status(fiber.StatusOK).JSON(apiModel.Response{
		Status:  "success",
		Message: "Operation successfull",
		Data:    transactions,
	})
}

// Adds a new transaction for the respective user to the database.
func AddTransactionHandler(c *fiber.Ctx) error {

	userID, err := utils.GetUserId(c)
	if err != nil {
		log.Error(err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(apiModel.Response{
			Status:  "error",
			Message: "User id is required for the transaction",
			Data:    nil,
		})
	}
	transactionReq := new(transactionModels.TransactionRequest)

	// Validate incoming request
	if errs, err := validation.ValidateRequest(c, transactionReq); err != nil {
		log.Error(err.Error())
		errMsg := validation.CheckErrors(c, errs, err)
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: errMsg,
			Data:    nil,
		})
	}

	// Create a new transaction object
	transaction := &transactionModels.Transaction{
		UserID:      userID,
		Name:        transactionReq.Name,
		Frequency:   transactionReq.Frequency,
		Amount:      transactionReq.Amount,
		CategoryID:  transactionReq.CategoryID,
		TxnType:     transactionReq.TxnType,
		TxnDate:     transactionReq.TxnDate,
		Description: transactionReq.Description,
	}

	// Add transaction to database
	if err := database.DB.Create(transaction).Error; err != nil {
		log.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: "Unable to add transaction, please try again",
			Data:    nil,
		})
	}

	// Update dashboard metrics
	utils.UpdateDashboardMetrics(userID)

	return c.Status(fiber.StatusAccepted).JSON(apiModel.Response{
		Status:  "success",
		Message: "Transaction added successfully",
		Data:    transactionReq,
	})

}
