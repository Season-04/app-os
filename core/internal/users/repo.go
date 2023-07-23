package users

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
	"sync"

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

func (r *Repository) CreateUser(user *User, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user.HashedPassword = hashedPassword
	user.ID = r.maxUserID + 1

	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()
	data := &usersFile{
		Users: append(r.users, user),
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
	r.maxUserID = user.ID
	return nil
}

type User struct {
	ID             uint32 `yaml:"id"`
	Name           string `yaml:"name"`
	EmailAddress   string `yaml:"email_address"`
	HashedPassword string `yaml:"hashed_password"`
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
