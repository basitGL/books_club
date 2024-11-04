package config

import (
	"database/sql"
	"errors"
	"os"

	// "errors"
	"fmt"
	"log"

	// "os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ichtrojan/thoth"
	_ "github.com/joho/godotenv/autoload"
)

func Database() *sql.DB {
	logger, _ := thoth.Init("log")
	dbName := "good_reads"
	dbPort := 3306

	user, EXISTS := os.LookupEnv("DB_USER")

	if !EXISTS {
		logger.Log(errors.New("DB_USER not set in .env"))
		log.Fatal("DB_USER not set in .env")
	}

	pass, EXISTS := os.LookupEnv("DB_PASS")

	if !EXISTS {
		logger.Log(errors.New("DB_PASS not set in .env"))
		log.Fatal("DB_PASS not set in .env")
	}

	host, EXISTS := os.LookupEnv("DB_HOST")

	if !EXISTS {
		logger.Log(errors.New("DB_HOST not set in .env"))
		log.Fatal("DB_HOST not set in .env")
	}

	// Format the connection string
	credentials := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, dbPort, dbName)
	fmt.Println(credentials)

	database, err := sql.Open("mysql", credentials)

	if err != nil {
		logger.Log(err)
		log.Fatal(err)
	}

	_, err = database.Exec(`CREATE DATABASE IF NOT EXISTS good_reads`)

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		panic(err.Error())
	} else {
		fmt.Println("Database Created Successfully")
	}

	_, err = database.Exec(`USE good_reads`)

	if err != nil {
		fmt.Println(err)
	}

	return database

}
