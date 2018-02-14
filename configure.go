package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var instance configuration
var once sync.Once

func configInstance() configuration {
	once.Do(func() {
		instance = getNewConfig()
	})

	return instance
}

type configuration struct {
	LogName                     string
	LogMaxSizeBeforeArchiveInMB int
	LogMaxAgeBeforeDeleteInDays int
	LogUseLocalTime             bool
	LogCompress                 bool

	ApplicationPort string

	PrivateKey string
	IvKey      string
	Database   databaseConfiguration
}

type databaseConfiguration struct {
	DatabaseHost string
	DatabasePort string

	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
}

func getNewConfig() configuration {
	env := os.Getenv("APP_ENVIRONMENT")
	c := getConfiguration(env)
	c.Database = getDatabaseConfiguration(env, c.PrivateKey, c.IvKey)
	return c
}

func getConfiguration(env string) configuration {
	fmt.Println("[I] Loading config for APP_ENVIRONMENT: " + env)
	propertiesFileName := getConfigurationFileName(env)

	fmt.Println("[L] Reading... " + propertiesFileName)
	confFile, confErr := ioutil.ReadFile(propertiesFileName)
	if confErr != nil {
		panic(confErr)
	}

	var c configuration
	json.Unmarshal(confFile, &c)
	panicIfConfigurationContainEmpty(c, propertiesFileName)
	defaultValueIfValueIsZero(c)

	return c
}

func getDatabaseConfiguration(env, privateKey, ivKey string) databaseConfiguration {
	databasePropertiesFileName := getDatabaseConfigurationFileName(env)
	fmt.Println("[L] Reading... " + databasePropertiesFileName)
	dbFile, dbErr := ioutil.ReadFile(databasePropertiesFileName)
	if dbErr != nil {
		panic(dbErr)
	}

	var db databaseConfiguration
	json.Unmarshal(dbFile, &db)
	panicIfDatabaseConfigurationContainEmpty(db, databasePropertiesFileName)

	fmt.Println("[P] Decrypting file...")
	db = decryptNecessaryFields(db, privateKey, ivKey)
	fmt.Println("[F] File decrypted!")

	return db
}

func getConfigurationFileName(env string) string {
	switch env {
	case "DEVELOPMENT":
		return "/conf/application-dev.json"
	case "UAT":
		return "/conf/application-uat.json"
	case "PRODUCTION":
		return "/conf/application-production.json"
	case "LOCAL":
		return "conf/application-local.json"
	default:
		panic("[E] APP_ENVIRONMENT not set please set one of 'DEVELOPMENT', 'UAT', 'PRODUCTION', 'LOCAL'")
	}
}

func getDatabaseConfigurationFileName(env string) string {
	switch env {
	case "DEVELOPMENT":
		return "/conf/db-con-dev.json"
	case "UAT":
		return "/conf/db-con-uat.json"
	case "PRODUCTION":
		return "/conf/db-con-production.json"
	case "LOCAL":
		return "conf/db-con-local.json"
	default:
		panic("[E] APP_ENVIRONMENT not set please set one of 'DEVELOPMENT', 'UAT', 'PRODUCTION', 'LOCAL'")
	}
}

func decryptNecessaryFields(c databaseConfiguration, privateKey, ivKey string) databaseConfiguration {
	decryptedDbUsername, _ := decrypt(c.DatabaseUsername, privateKey, ivKey)
	c.DatabaseUsername = decryptedDbUsername

	decryptedDbPassword, _ := decrypt(c.DatabasePassword, privateKey, ivKey)
	c.DatabasePassword = decryptedDbPassword

	return c
}

func panicIfConfigurationContainEmpty(c configuration, fileName string) {
	if c.ApplicationPort == "" {
		panic("[!] ApplicationPort shold not be nil in " + fileName)
	} else if c.LogName == "" {
		panic("[!] LogName shold not be nil in " + fileName)
	} else if c.PrivateKey == "" {
		panic("[!] PrivateKey shold not be nil in " + fileName)
	}
}

func panicIfDatabaseConfigurationContainEmpty(c databaseConfiguration, fileName string) {
	if c.DatabaseHost == "" {
		panic("[!] DatabaseHost shold not be nil in " + fileName)
	} else if c.DatabasePort == "" {
		panic("[!] DatabasePort shold not be nil in " + fileName)
	} else if c.DatabaseUsername == "" {
		panic("[!] DatabaseUsername shold not be nil in " + fileName)
	} else if c.DatabasePassword == "" {
		panic("[!] DatabasePassword shold not be nil in " + fileName)
	}
}

func defaultValueIfValueIsZero(c configuration) configuration {
	if c.ApplicationPort == "" {
		fmt.Println("[!] ApplicationPort not found, set it to 8080")
		c.ApplicationPort = "8080"
	}

	if c.LogMaxSizeBeforeArchiveInMB == 0 {
		fmt.Println("[!] LogMaxSizeBeforeArchiveInMB not found, set it to 100")
		c.LogMaxSizeBeforeArchiveInMB = 100
	}

	if c.LogMaxAgeBeforeDeleteInDays == 0 {
		fmt.Println("[!] LogMaxAgeBeforeDeleteInDays not found, set it to 366")
		c.LogMaxAgeBeforeDeleteInDays = 366
	}

	return c
}
