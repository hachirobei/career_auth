package database

import (
    "fmt"
    "log"
    "os"

	"github.com/joho/godotenv"
    "rawgo.com/m/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type Dbinstance struct {
    Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
    dsn := fmt.Sprintf(
        "host=%s  user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Kuala_Lumpur",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        log.Fatal("Failed to connect to database. \n", err)
        os.Exit(1)
    }

    log.Println("connected")
    db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
    db.AutoMigrate(&models.Candidate{})


    DB = Dbinstance{
        Db: db,
    }
}