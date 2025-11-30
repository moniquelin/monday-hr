package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/moniquelin/monday-hr/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// Error for duplicate emails
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// User struct represents an individual user
type User struct {
	ID        int64     `json:"id"`
	IsAdmin   bool      `json:"is_admin"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Salary    int64     `json:"salary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

type password struct {
	plaintext *string
	hash      []byte
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// Validates email
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// Validates password
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// Validates user
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user).
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

// UserModel struct wraps the connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert new user in the database
func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (is_admin, name, email, password_hash, salary, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Set created_by = null if the value is 0
	var createdBy any
	if user.CreatedBy == 0 {
		createdBy = nil
	} else {
		createdBy = user.CreatedBy
	}

	// Set updated_by = null if the value is 0
	var updatedBy any
	if user.UpdatedBy == 0 {
		updatedBy = nil
	} else {
		updatedBy = user.UpdatedBy
	}

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE users email constraint
	err := m.DB.QueryRowContext(ctx, query,
		user.IsAdmin,
		user.Name,
		user.Email,
		user.Password.hash,
		user.Salary,
		createdBy,
		updatedBy).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == `pq: duplicate key value violates unique constraint "users_email_key"` { // unique_violation
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
        SELECT id, is_admin, name, email, password_hash, salary, created_at, updated_at, created_by, updated_by
        FROM users
        WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.IsAdmin,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Salary,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}
