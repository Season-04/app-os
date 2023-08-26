package users

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

type Repository struct {
	usersFilePath string
	usersMutex    *sync.RWMutex
	users         []*User
	maxUserID     uint32
}

func NewRepository(root string) *Repository {
	usersFilePath := filepath.Join(root, "users.yml")
	bytes, err := os.ReadFile(usersFilePath)
	data := &usersFile{}

	if err == nil {
		err = yaml.Unmarshal(bytes, data)
		if err != nil {
			log.Println("Failed to load users", err)
		}
	} else {
		adminPassword, _ := hashPassword("admin")
		data = &usersFile{
			Users: []*User{
				{
					ID:             1,
					Name:           "Admin",
					EmailAddress:   "admin@local",
					HashedPassword: adminPassword,
					Role:           UserRoleAdmin,
				},
			},
		}
		bytes, err = yaml.Marshal(data)
		if err == nil {
			_ = os.WriteFile(usersFilePath, bytes, 0666)
		}
	}
	users := data.Users

	return &Repository{
		usersFilePath: usersFilePath,
		usersMutex:    &sync.RWMutex{},
		users:         users,
		maxUserID:     maxUserID(users),
	}
}

func (r *Repository) GetUserByID(ID uint32) *User {
	r.usersMutex.RLock()
	defer r.usersMutex.RUnlock()

	for _, user := range r.users {
		if user.ID == ID {
			return user
		}
	}
	return nil
}

func (r *Repository) GetUserByEmailAddress(emailAddress string) *User {
	r.usersMutex.RLock()
	defer r.usersMutex.RUnlock()

	for _, user := range r.users {
		if user.EmailAddress == emailAddress {
			return user
		}
	}
	return nil
}

func (r *Repository) GetUserByEmailAddressAndPassword(emailAddress string, password string) *User {
	user := r.GetUserByEmailAddress(emailAddress)

	if checkPassword(password, user.HashedPassword) {
		return user
	}

	return nil
}

func (r *Repository) ListAll() []*User {
	r.usersMutex.RLock()
	defer r.usersMutex.RUnlock()

	return r.users
}

func (r *Repository) CreateUser(user *User, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user.HashedPassword = hashedPassword
	user.ID = r.maxUserID + 1

	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	err = r.saveUsers(append(r.users, user))
	if err != nil {
		return err
	}

	r.maxUserID = user.ID
	return nil
}

func (r *Repository) UpdateUser(user *User) error {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	var existingUser *User = nil
	for _, u := range r.users {
		if u.ID == user.ID {
			existingUser = u
		}
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	existingUser.Name = user.Name
	existingUser.Role = user.Role

	return r.saveUsers(r.users)
}

func (r *Repository) saveUsers(users []*User) error {
	data := &usersFile{
		Users: users,
	}
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(r.usersFilePath, bytes, 0666)
	if err != nil {
		return err
	}

	r.users = data.Users
	return nil
}

type UserRole string

var (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID             uint32     `yaml:"id"`
	Name           string     `yaml:"name"`
	EmailAddress   string     `yaml:"email_address"`
	HashedPassword string     `yaml:"hashed_password"`
	Role           UserRole   `yaml:"role"`
	LastSeenAt     *time.Time `yaml:"last_seen_at"`
}

type usersFile struct {
	Users []*User `yaml:"users"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func checkPassword(password string, hashedPassword string) bool {
	bytes, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(bytes, []byte(password))
	return err == nil
}

func maxUserID(users []*User) uint32 {
	var v uint32 = 0

	for _, user := range users {
		if user.ID > v {
			v = user.ID
		}
	}

	return v
}
