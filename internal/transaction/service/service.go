package service

import (
	"boiler-plate/app/appconf"
	"boiler-plate/internal/transaction/domain"
	"boiler-plate/internal/transaction/repository"
	usersRepo "boiler-plate/internal/users/repository"
	walletRepo "boiler-plate/internal/wallet/repository"
	"boiler-plate/pkg/exception"
	"context"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strconv"
)

// NewService creates new user service
func NewService(
	config *appconf.Config, repo repository.TransactionRepository, users usersRepo.UsersRepository,
	wallet walletRepo.WalletRepository, db *gorm.DB,
	validate *validator.Validate,
) Service {
	return &service{
		config: config, TransactionRepo: repo, WalletRepo: wallet,
		UserRepo: users, validate: validate, DB: db,
	}
}

type service struct {
	DB              *gorm.DB
	config          *appconf.Config
	TransactionRepo repository.TransactionRepository
	UserRepo        usersRepo.UsersRepository
	WalletRepo      walletRepo.WalletRepository
	validate        *validator.Validate
}

func (s service) Credit(
	ctx context.Context, userid string, amount float64,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(userid)
	if err != nil {
		return exception.PermissionDenied("Input of userid must be integer")
	}
	if amount < 1 {
		return exception.PermissionDenied("Input of amount must be greater than zero")
	}
	user, err := s.UserRepo.Detail(ctx, tx, idInt)
	if err != nil {
		return exception.Internal("failed getting sender detail", err)
	}
	if user == nil {
		return exception.NotFound("sender detail not found")
	}
	userTransaction := &domain.Transaction{
		Type:    "Credit",
		Amount:  user.Wallet.Balance + amount,
		Message: "Credit of " + strconv.FormatFloat(amount, 'f', -1, 64),
		UserId:  idInt,
	}
	if err := s.TransactionRepo.Create(ctx, tx, userTransaction); err != nil {
		return exception.Internal("failed creating transaction", err)
	}
	user.Wallet.Increase(amount)
	if err := s.WalletRepo.Update(ctx, tx, user.Wallet.ID, user.Wallet); err != nil {
		return exception.Internal("failed updating wallet", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Transfer(
	ctx context.Context, senderId, receiverId string, amount float64,
) *exception.Exception {
	tx := s.DB.Begin()
	defer tx.Rollback()
	senderIdInt, err := strconv.Atoi(senderId)
	if err != nil {
		return exception.PermissionDenied("Input of senderId must be integer")
	}
	receiverIdInt, err := strconv.Atoi(receiverId)
	if err != nil {
		return exception.PermissionDenied("Input of receiverId must be integer")
	}
	if amount < 1 {
		return exception.PermissionDenied("Input of amount must be greater than zero")
	}
	sender, err := s.UserRepo.Detail(ctx, tx, senderIdInt)
	if err != nil {
		return exception.Internal("failed getting sender detail", err)
	}
	if sender == nil {
		return exception.NotFound("sender detail not found")
	}
	receiver, err := s.UserRepo.Detail(ctx, tx, receiverIdInt)
	if err != nil {
		return exception.Internal("failed getting receiver detail", err)
	}
	if receiver == nil {
		return exception.NotFound("receiver detail not found")
	}
	if sender.Wallet.Balance < amount {
		return exception.PermissionDenied(sender.Name + " does not have enough balance. Balance: " + strconv.FormatFloat(sender.Wallet.Balance, 'f', -1, 64))
	}
	senderTransaction := &domain.Transaction{
		Type:    "Debit",
		Amount:  sender.Wallet.Balance - amount,
		Message: "Transfer to: " + receiver.Name,
		UserId:  senderIdInt,
	}
	if err := s.TransactionRepo.Create(ctx, tx, senderTransaction); err != nil {
		return exception.Internal("failed creating transaction", err)
	}
	sender.Wallet.Decrease(amount)
	if err := s.WalletRepo.Update(ctx, tx, sender.Wallet.ID, sender.Wallet); err != nil {
		return exception.Internal("failed updating wallet", err)
	}
	receiverTransaction := &domain.Transaction{
		Type:    "Credit",
		Amount:  receiver.Wallet.Balance + amount,
		Message: "Transfer from: " + sender.Name,
		UserId:  receiverIdInt,
	}
	if err := s.TransactionRepo.Create(ctx, tx, receiverTransaction); err != nil {
		return exception.Internal("failed creating transaction", err)
	}
	receiver.Wallet.Increase(amount)
	if err := s.WalletRepo.Update(ctx, tx, receiver.Wallet.ID, receiver.Wallet); err != nil {
		return exception.Internal("failed updating wallet", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s service) Find(ctx context.Context, limit string, page string, userid string) (
	*FindResponse, *exception.Exception,
) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	var limitInt, pageInt int

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, exception.PermissionDenied("Input of limit must be integer")
	}
	pageInt, err = strconv.Atoi(page)
	if err != nil {
		return nil, exception.PermissionDenied("Input of page must be integer")
	}
	IdInt, err := strconv.Atoi(userid)
	if err != nil {
		return nil, exception.PermissionDenied("Input of userId must be integer")
	}
	user, err := s.UserRepo.Detail(ctx, tx, IdInt)
	if err != nil {
		return nil, exception.Internal("error getting users", err)
	}
	if user == nil {
		return nil, exception.NotFound("user not found")
	}
	result, pagination, err := s.TransactionRepo.Find(ctx, tx, limitInt, pageInt, IdInt)
	if err != nil {
		return nil, exception.Internal("failed finding transaction", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	finalResponse := &FindResponse{
		Pagination: *pagination,
		Users:      *user,
		Data:       *result,
	}
	return finalResponse, nil
}

func (s service) Detail(ctx context.Context, id string) (*domain.Transaction, *exception.Exception) {
	tx := s.DB.Begin()
	defer tx.Rollback()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, exception.PermissionDenied("Input of id must be integer")
	}
	result, err := s.TransactionRepo.Detail(ctx, tx, idInt)
	if err != nil {
		return nil, exception.Internal("error getting detail users", err)
	}
	if result == nil {
		return nil, exception.NotFound("detail not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return result, nil
}
