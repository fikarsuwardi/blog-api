package database

import (
	"log"

	"blog-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func Migrate() error {
	log.Println("Running database migration...")

	// Auto migrate semua models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)

	if err != nil {
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}

func SeedData() error {
	log.Println("Seeding initial data...")

	// Cek apakah sudah ada user
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count > 0 {
		log.Println("Data already exists, skipping seed")
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Buat contoh user
	users := []models.User{
		{
			Email:    "john@example.com",
			Password: string(hashedPassword),
			Name:     "John Doe",
		},
		{
			Email:    "jane@example.com",
			Password: string(hashedPassword),
			Name:     "Jane Smith",
		},
	}

	result := DB.Create(&users)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Seeded %d users successfully", len(users))
	return nil
}
