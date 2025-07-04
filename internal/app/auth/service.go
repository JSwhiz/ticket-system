package auth

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 			 string `db:"user_id"`
	Username	 string `db:"username"`
    RoleID       string `db:"role_id"`
    DepartmentID int    `db:"department_id"`
	PasswordHash string `db:"password_hash"`
}

type Service struct {
    db     *sqlx.DB
    secret string
    ttl    time.Duration
}

func NewService(db *sqlx.DB, secret string, ttl time.Duration) *Service {
    return &Service{db: db, secret: secret, ttl: ttl}
}

func (s *Service) Authenticate(username, password string) (*User, error) {
    var hash string
    err := s.db.Get(&hash,
        "SELECT password_hash FROM users WHERE username = $1 AND deleted_at IS NULL",
        username,
    )
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
        return nil, errors.New("invalid credentials")
    }
    user := User{}
    err = s.db.Get(&user,
        `SELECT user_id::text, username, role_id::text, department_id
           FROM users
          WHERE username = $1`,
        username,
    )
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *Service) Login(usernameOrEmail, password string) (token string, user User, err error) {
    err = s.db.Get(&user, `
        SELECT user_id, username, role_id, department_id, password_hash
          FROM users
         WHERE username = $1 OR email = $1
    `, usernameOrEmail)
    if err != nil {
        return "", User{}, err
    }
    if bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash),
        []byte(password),
    ) != nil {
        return "", User{}, errors.New("invalid credentials")
    }
    token, err = GenerateToken(s.secret, user.ID, user.RoleID, s.ttl)
    if err != nil {
        return "", User{}, err
    }
    user.PasswordHash = ""
    return token, user, nil
}
