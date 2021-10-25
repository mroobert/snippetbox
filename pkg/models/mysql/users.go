package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/mroobert/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// Define a UserModel type which wraps a sql.DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Method to add a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	sqlInsertQuery := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(sqlInsertQuery, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Method to verify whether a user exists with the provided email address and password.
// This will return the relevant user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	sqlSelectQuery := "SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE"
	row := m.DB.QueryRow(sqlSelectQuery, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// Method to fetch details for a specific user based on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {

	user := models.User{}

	sqlSelectQuery := `SELECT id, name, email, created, active FROM users WHERE id = ?`

	row := m.DB.QueryRow(sqlSelectQuery, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &user, nil
}
