package handlers

import (
	"fmt"
	"time"

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
	addTransactionReq := new(transactionModels.AddTransactionRequest)

	// Validate incoming request
	if errs, err := validation.ValidateRequest(c, addTransactionReq); err != nil {
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
		Name:        addTransactionReq.Name,
		Frequency:   addTransactionReq.Frequency,
		Amount:      addTransactionReq.Amount,
		CategoryID:  addTransactionReq.CategoryID,
		TxnType:     addTransactionReq.TxnType,
		TxnDate:     time.Now(),
		Description: addTransactionReq.Description,
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
		Data:    addTransactionReq,
	})

}

func UpdateTransactionHandler(c *fiber.Ctx) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(apiModel.Response{
			Status:  "error",
			Message: "User id is required for the transaction",
			Data:    nil,
		})
	}

	transactionID := c.Params("id")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Cannot delete transaction: invalid item request",
			Data:    nil,
		})
	}

	patchTransactionReq := new(transactionModels.UpdateTransactionRequest)

	// Validate incoming request
	if errs, err := validation.ValidateRequest(c, patchTransactionReq); err != nil {
		log.Error(err.Error())
		errMsg := validation.CheckErrors(c, errs, err)
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: errMsg,
			Data:    nil,
		})
	}

	patchMap := buildPatchMap(patchTransactionReq)
	if len(patchMap) == 0 {
		log.Error("No items in PATCH request")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Atleast 1 items is required for PATCH",
			Data:    nil,
		})
	}

	tx := database.DB.Model(&transactionModels.Transaction{}).Where("id = ? AND user_id = ?", transactionID, userID).Updates(patchMap)
	if tx.Error != nil {
		log.Error(tx.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: fmt.Sprintf("Cannot update transaction: %s", tx.Error.Error()),
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(apiModel.Response{
		Status:  "success",
		Message: "Transaction updated",
		Data:    patchMap,
	})
}

func buildPatchMap(patchReq *transactionModels.UpdateTransactionRequest) map[string]any {
	patchMap := make(map[string]any)
	if patchReq.Name != nil {
		patchMap["name"] = patchReq.Name
	}
	if patchReq.Amount != nil {
		patchMap["amount"] = patchReq.Amount
	}
	if patchReq.Frequency != nil {
		patchMap["frequency"] = patchReq.Frequency
	}

	if patchReq.TxnType != nil {
		patchMap["txn_type"] = patchReq.TxnType
	}
	if patchReq.CategoryID != nil {
		patchMap["category_id"] = patchReq.CategoryID
	}
	if patchReq.Description != nil {
		patchMap["description"] = patchReq.Description
	}

	return patchMap
}

func DeleteTransactionHandler(c *fiber.Ctx) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(apiModel.Response{
			Status:  "error",
			Message: "User id is required for the transaction",
			Data:    nil,
		})
	}

	transactionID := c.Params("id")
	if transactionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Cannot delete transaction: invalid item request",
			Data:    nil,
		})
	}

	deleteTX := database.DB.Delete(&transactionModels.Transaction{}, userID, transactionID)
	if deleteTX.Error != nil {
		log.Error(deleteTX.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: fmt.Sprintf("Cannot delete transaction: %s", deleteTX.Error.Error()),
			Data:    nil,
		})
	}

	// Update dashboard metrics
	utils.UpdateDashboardMetrics(userID)

	return c.SendStatus(fiber.StatusOK)
}
