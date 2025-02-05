package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"finances/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Deposit(c *gin.Context)
	Transfer(c *gin.Context)
	GetLastTransactions(c *gin.Context)
}

type userControllerImpl struct {
	userService usecase.UserService
}

func NewUserController(userService usecase.UserService) UserController {
	return &userControllerImpl{
		userService: userService,
	}
}

func (u *userControllerImpl) Deposit(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var depositRequest struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&depositRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = u.userService.Deposit(context.Background(), userID, depositRequest.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}

func (u *userControllerImpl) Transfer(c *gin.Context) {
    fromUserIDStr := c.Param("userID") // Изменено с "fromUserID" на "userID"
    fromUserID, err := strconv.ParseInt(fromUserIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from user ID"})
        return
    }

    // Получаем ID получателя из параметра пути
    toUserIDStr := c.Param("toUserID")
    toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to user ID"})
        return
    }

    // Парсим тело запроса с суммой перевода
    var transferRequest struct {
        Amount float64 `json:"amount"`
    }
    if err := c.ShouldBindJSON(&transferRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Вызываем сервис для выполнения перевода
    err = u.userService.Transfer(context.Background(), fromUserID, toUserID, transferRequest.Amount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

func (u *userControllerImpl) GetLastTransactions(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	transactions, err := u.userService.GetLastTransactions(context.Background(), userID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
