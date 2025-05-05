package repository

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rudyjcruz831/affordAbode_backend/model"
	"github.com/rudyjcruz831/affordAbode_backend/util/errors"
	"gorm.io/gorm"
)

// PGUserRepository is data/repository implementation
// of service layer UserRepository
type pGUserRepository struct {
	DB *gorm.DB
}

type Metadata struct {
	UserID    string     `json:"user_id"`
	CreatedOn time.Time  `json:"created_on"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
type UserCreate struct {
	ID        string `json:"user_id" gorm:"column:user_id;primaryKey"`
	Email     string `json:"email" binding:"required,email" gorm:"column:email;unique;not null"`
	Username  string `json:"username" binding:"required" gorm:"column:username;unique;not null"`
	Password  string `json:"password" binding:"required,gte=6,lte=30" gorm:"column:password"`
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
	Role      string `json:"role" gorm:"column:user_role"`
}

func (UserCreate) TableName() string {
	return "users"
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &pGUserRepository{
		DB: db,
	}
}

// FindByID fetches user by id
func (r *pGUserRepository) FindByID(ctx context.Context, id string) (*model.Users, *errors.AffordAbodeError) {
	// panic("Create function in Pg user repository")
	// we storing ID into model this allows GORM to use the PRIMARY KEY in model
	u := &model.Users{
		ID: id,
	}
	result := r.DB.First(&u)
	if result.Error != nil {
		log.Printf("Error: %v\n", result.Error)
		affordAbodeErr := errors.NewNotFound("id", id)
		return nil, affordAbodeErr
	}

	return u, nil
}

// Create reaches out to database postrges using gorm api
func (r *pGUserRepository) Create(ctx context.Context, u *model.Users) *errors.AffordAbodeError {
	// panic("Create function in Pg user repository")
	uid, _ := uuid.NewRandom()
	// query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"
	u.CreatedOn = time.Now()
	u.UpdatedAt = time.Now()
	u.ID = uid.String()
	u.Role = "user"

	// if result := r.DB.FirstOrCreate(&u, u); result.Error != nil {
	// 	log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, result.Error)
	// 	affordAbodeErr := errors.NewConflict("email", u.Email)
	// 	return affordAbodeErr
	// }

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		userCreate := &UserCreate{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			Password:  u.Password,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      u.Role,
		}

		if err := tx.Create(userCreate).Error; err != nil {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
			// affordAbodeErr := errors.NewConflict("email", u.Email)
			return err
		}

		metadata := &Metadata{
			UserID:    u.ID,
			CreatedOn: u.CreatedOn,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: u.DeletedAt,
		}

		if err := tx.Create(&metadata).Error; err != nil {
			log.Printf("Could not create metadata for user: %v. Reason: %v\n", u.Email, err)
			return err
		}
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v\n", err)
		return errors.NewInternalServerError("Could not create user")
	}

	return nil
}

// FindByEmail fetches user by email
func (r *pGUserRepository) FindByEmail(ctx context.Context, email string) (*model.Users, *errors.AffordAbodeError) {
	// panic("FindByEmail in pGUserRepository")

	u := &model.Users{}

	// using gorm to hit postgresDB using email
	if result := r.DB.Where("email = ?", email).First(u); result.Error != nil {
		log.Printf("Db error : %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			affordAbodeErr := errors.NewNotFound("email", email)
			return nil, affordAbodeErr
		} else {
			affordAbodeErr := errors.NewInternalServerError("")
			return nil, affordAbodeErr
		}

	}

	return u, nil
}

// Update updates user information
func (r *pGUserRepository) Update(ctx context.Context, u *model.Users) *errors.AffordAbodeError {
	userInRepo := &model.Users{}
	result := r.DB.Where("email = ?", u.Email).First(userInRepo)
	if result.Error != nil {
		log.Printf("Db error: %v", result.Error)
		affordAbodeErr := errors.NewNotFound("email", u.Email)
		return affordAbodeErr
	}

	userInRepo.FirstName = u.FirstName
	userInRepo.LastName = u.LastName

	result = r.DB.Save(userInRepo)
	if result.Error != nil {
		log.Printf("Db error %v", result.Error)
		affordAbodeErr := errors.NewInternalServerError("")
		return affordAbodeErr
	}

	return nil
}

// Delete deletes user information
func (r *pGUserRepository) Delete(ctx context.Context, id string) *errors.AffordAbodeError {
	u := &model.Users{}
	results := r.DB.Delete(u, id)
	if results.Error != nil {
		return errors.NewInternalServerError("")
	}

	return nil
}

// func (r *pGUserRepository) UpdateImage(ctx context.Context, u *model.Users, imageURL string) (*model.Users, *errors.AffordAbodeError) {

// 	// must be instantiated to scan into ref using `GetContext`
// 	userInRepo := &model.Users{}

// 	result := r.DB.Where("email = ?", u.Email).First(userInRepo)
// 	if result.Error != nil {
// 		log.Printf("Db error: %v", result.Error)
// 		affordAbodeErr := errors.NewNotFound("email", u.Email)
// 		return nil, affordAbodeErr
// 	}
// 	//update the user Image
// 	userInRepo.Image = imageURL
// 	result = r.DB.Save(userInRepo)
// 	if result.Error != nil {
// 		log.Printf("Db error %v", result.Error)
// 		affordAbodeErr := errors.NewInternalServerError("")
// 		return nil, affordAbodeErr
// 	}
// 	return userInRepo, nil
// }
