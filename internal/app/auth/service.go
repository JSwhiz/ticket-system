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
}

type Service struct {
    db     *sqlx.DB
    secret string
    ttl    time.Duration
}

func NewService(db *sqlx.DB, jwtSecret string, ttl time.Duration) *Service {
    return &Service{db: db, secret: jwtSecret, ttl: ttl}
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

func (s *Service) Login(username, password string) (string, *User, error) {
    user, err := s.Authenticate(username, password)
    if err != nil {
        return "", nil, err
    }
    token, err := GenerateToken(s.secret, user.ID, s.ttl)
    if err != nil {
        return "", nil, err
    }
    return token, user, nil
}
