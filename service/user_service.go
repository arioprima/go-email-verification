package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang_email_verification/initializers"
	"golang_email_verification/models"
	"golang_email_verification/repository"
	"golang_email_verification/utils"
	"log"
	"time"
)

type UserService interface {
	Login(ctx context.Context, request models.LoginInput) (models.LoginResponse, error)
	Register(ctx context.Context, request models.RegisterInput) (string, error)
	VerifyEmail(ctx context.Context, request models.VerifyInput) (models.UserResponse, error)
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewUserServiceImpl(userRepository repository.UserRepository, db *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{UserRepository: userRepository, DB: db, Validate: validate}
}

func (service *UserServiceImpl) Login(ctx context.Context, request models.LoginInput) (models.LoginResponse, error) {
	//TODO implement me
	tx, err := service.DB.Begin()
	if err != nil {
		return models.LoginResponse{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			err := tx.Rollback()
			if err != nil {
				log.Println("Error rolling back transaction:", err)
			}
		} else {
			err := tx.Commit()
			if err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	user, err := service.UserRepository.Login(ctx, tx, request.Email)
	if err != nil || user == nil || user.IsVerified != true {
		return models.LoginResponse{}, err
	}

	err = utils.VerifyPassword(user.Password, request.Password)
	if err != nil {
		return models.LoginResponse{}, errors.New("invalid credentials")
	}

	config, _ := initializers.LoadConfig(".")

	//generate jwt token
	tokenPayload := map[string]interface{}{
		"id":        user.ID,
		"email":     user.Email,
		"user_name": user.UserName,
	}

	token, err := utils.GenerateToken(config.TokenExpiresIn, tokenPayload, config.TokenSecret)
	if err != nil {
		return models.LoginResponse{}, errors.New("error generating token")
	}

	return models.LoginResponse{
		ID:        user.ID,
		Email:     user.Email,
		UserName:  user.UserName,
		TokenType: "Bearer",
		Token:     token,
	}, nil
}

func (service *UserServiceImpl) Register(ctx context.Context, request models.RegisterInput) (string, error) {
	// Validasi input
	if err := service.Validate.Struct(request); err != nil {
		return "", fmt.Errorf("kesalahan validasi: %v", err)
	}

	// Mulai transaksi
	tx, err := service.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("kesalahan memulai transaksi: %v", err)
	}

	// Tunda fungsi untuk menangani rollback atau commit transaksi
	defer func() {
		if r := recover(); r != nil {
			// Terjadi kesalahan, rollback transaksi
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Kesalahan rollback transaksi: %v", rollbackErr)
			}
			log.Printf("Panic terjadi: %v", r)
		} else {
			// Tidak ada kesalahan, commit transaksi
			if commitErr := tx.Commit(); commitErr != nil {
				log.Printf("Kesalahan commit transaksi: %v", commitErr)
				// Jika terjadi kesalahan commit, rollback transaksi
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Printf("Kesalahan rollback transaksi setelah kesalahan commit: %v", rollbackErr)
				}
			}
		}
	}()

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return "", fmt.Errorf("kesalahan hashing password: %v", err)
	}

	now := time.Now()
	newUser := models.User{
		ID:         utils.GenerateUUID(),
		UserName:   request.UserName,
		Email:      request.Email,
		Password:   hashedPassword,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Generate dan simpan OTP
	otp := utils.GenerateOTP()

	// Perbarui pengguna di Database
	newUser.VerificationToken = otp
	_, err = service.UserRepository.Register(ctx, tx, &newUser)
	if err != nil {
		return "", fmt.Errorf("kesalahan mendaftarkan pengguna: %v", err)
	}

	// Kirim email
	utils.SendEmail(&newUser, otp)

	return "Cek email Anda untuk kode verifikasi", nil
}

func (service *UserServiceImpl) VerifyEmail(ctx context.Context, request models.VerifyInput) (models.UserResponse, error) {
	log.Println("Start VerifyEmail Function")

	if err := service.Validate.Struct(request); err != nil {
		log.Printf("Validation error: %v", err)
		return models.UserResponse{}, err
	}

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return models.UserResponse{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				log.Printf("Error committing transaction: %v", commitErr)
			}
		}
	}()

	log.Printf("Email: %s, Token: %s", request.Email, request.Token)

	user, err := service.UserRepository.VerifyEmail(ctx, tx, request.Token)
	if err != nil {
		log.Printf("Error verifying email: %v", err)
		return models.UserResponse{}, err
	}

	if user == nil {
		log.Println("User not found")
		return models.UserResponse{}, errors.New("user does not exist")
	}

	// Perbaiki pemanggilan fungsi UpdateUserVerificationStatus untuk mencocokkan perubahan
	err = service.UserRepository.UpdateUserVerificationStatus(ctx, tx, request.Email, request.Token)
	if err != nil {
		log.Printf("Error updating user verification status: %v", err)
		return models.UserResponse{}, err
	}

	log.Println("End VerifyEmail Function")
	return models.UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
