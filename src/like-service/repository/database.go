package repository

import (
	"errors"
	"fmt"
	"like-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrLikeNotFound = errors.New("no matching record found in database")
)

type Database struct {
	Client *gorm.DB
}

func NewDatabase(host string, dbUser string, password string, dbName string, port string) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &Database{
		Client: db,
	}, nil
}

func (d *Database) FindAllLikes() (*[]models.Like, error) {
	var err error

	var likes []models.Like

	err = d.Client.Order("id").Find(&likes).Error

	if err != nil {
		return &[]models.Like{}, err
	}
	return &likes, nil
}

func (d *Database) AddLike(id int) error {
	tx := d.Client.Exec("UPDATE likes SET count = count + 1 WHERE id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
