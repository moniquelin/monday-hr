package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Error for duplicate emails
var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
)

// User struct represents an individual user
type User struct {
	ID        int64     `json:"id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Salary    int64     `json:"salary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int64     `json:"created_by"`
	UpdatedBy int64     `json:"updated_by"`
}

// UserModel struct wraps the connection pool
type UserModel struct {
	DB *sql.DB
}

type Password struct {
	plaintext *string
	hash      []byte
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *Password) Set(plaintextPassword string) error {
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
func (p *Password) Matches(plaintextPassword string) (bool, error) {
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

// Insert new user in the database
func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (role, name, email, password_hash, salary, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

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
	_, err := m.DB.ExecContext(ctx, query,
		user.Role,
		user.Name,
		user.Email,
		user.Password.hash,
		user.Salary,
		createdBy,
		updatedBy)
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

// Get user by email from the database
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
        SELECT id, role, name, email, password_hash, salary, created_at, updated_at, created_by, updated_by
        FROM users
        WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	var createdBy *int64
	var updatedBy *int64

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Role,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Salary,
		&user.CreatedAt,
		&user.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	if createdBy != nil {
		user.CreatedBy = *createdBy
	} else {
		user.CreatedBy = 0
	}

	if updatedBy != nil {
		user.UpdatedBy = *updatedBy
	} else {
		user.UpdatedBy = 0
	}

	return &user, nil
}

// Get user by ID from the database
func (m UserModel) Get(id int64) (*User, error) {
	query := `
        SELECT id, role, name, email, password_hash, salary, created_at, updated_at, created_by, updated_by
        FROM users
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	var createdBy *int64
	var updatedBy *int64

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Role,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Salary,
		&user.CreatedAt,
		&user.UpdatedAt,
		&createdBy,
		&updatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	if createdBy != nil {
		user.CreatedBy = *createdBy
	} else {
		user.CreatedBy = 0
	}

	if updatedBy != nil {
		user.UpdatedBy = *updatedBy
	} else {
		user.UpdatedBy = 0
	}

	return &user, nil
}
