package repository

import (
	"context"
	"database/sql"
	"errors"
	"golang_email_verification/models"
	"golang_email_verification/utils"
	"log"
)

type UserRepository interface {
	Login(ctx context.Context, tx *sql.Tx, email string) (*models.User, error)
	Register(ctx context.Context, tx *sql.Tx, user *models.User) (*models.User, error)
	VerifyEmail(ctx context.Context, tx *sql.Tx, token string) (*models.User, error)
	UpdateUserVerificationStatus(ctx context.Context, tx *sql.Tx, email, token string) error
}

type userRepositoryImpl struct {
	DB *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

func (u *userRepositoryImpl) Login(ctx context.Context, tx *sql.Tx, email string) (*models.User, error) {
	//TODO implement me
	SQL := `SELECT * FROM users WHERE email = $1 and is_verified = true`
	row := tx.QueryRowContext(ctx, SQL, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.VerificationToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepositoryImpl) Register(ctx context.Context, tx *sql.Tx, user *models.User) (*models.User, error) {
	log.Println("Start Register Function")

	// Langkah 1: Check if email already exists
	log.Println("Langkah 1: Memeriksa email yang sudah ada")
	SQL := `SELECT * FROM users WHERE email = $1`
	row := tx.QueryRowContext(ctx, SQL, user.Email)

	var existingUser models.User
	err := row.Scan(
		&existingUser.ID,
		&existingUser.UserName,
		&existingUser.Email,
		&existingUser.Password,
		&existingUser.IsVerified,
		&existingUser.VerificationToken,
		&existingUser.CreatedAt,
		&existingUser.UpdatedAt,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Kesalahan memeriksa email yang sudah ada: %v", err)
			return nil, err
		}
	}

	if existingUser.ID != "" {
		log.Println("Email sudah ada")
		return nil, errors.New("email already exists")
	}

	// Langkah 2: Hash password dan menyimpan pengguna baru
	log.Println("Langkah 2: Hash password dan menyimpan pengguna baru")

	// Set nilai ID sebelum operasi INSERT
	user.ID = utils.GenerateUUID()

	// Gunakan ExecContext untuk INSERT
	SQL = `INSERT INTO users (id, user_name, email, password, is_verified, verification_token, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.ExecContext(ctx, SQL, user.ID, user.UserName, user.Email, user.Password, user.IsVerified, user.VerificationToken, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Printf("Kesalahan menyimpan pengguna baru: %v", err)
		return nil, err
	}

	log.Println("Pengguna berhasil disimpan")
	log.Println("End Register Function")

	return user, nil
}

// Pada fungsi VerifyEmail di userRepositoryImpl
func (u *userRepositoryImpl) VerifyEmail(ctx context.Context, tx *sql.Tx, token string) (*models.User, error) {
	log.Println("Start VerifyEmail Repository Function")

	SQL := `SELECT * FROM users WHERE verification_token = $1`
	row := tx.QueryRowContext(ctx, SQL, token)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.VerificationToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("User not found for token: %s", token)
			return nil, nil
		}

		log.Printf("Error scanning user: %v", err)
		return nil, err
	}

	log.Printf("User found: %+v", user)

	if user.ID == "" {
		log.Println("User ID is empty")
		return nil, errors.New("user ID is empty")
	}

	log.Println("End VerifyEmail Repository Function")
	return &user, nil
}

func (u *userRepositoryImpl) UpdateUserVerificationStatus(ctx context.Context, tx *sql.Tx, email, token string) error {
	// TODO: Implement me
	SQL := `UPDATE users SET is_verified = true WHERE email = $1 AND verification_token = $2`
	result, err := tx.ExecContext(ctx, SQL, email, token)
	if err != nil {
		log.Printf("Error updating user verification status: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	log.Printf("Rows affected: %d", rowsAffected)

	if rowsAffected == 0 {
		log.Println("No rows updated")
		return errors.New("no rows updated")
	}

	return nil
}
