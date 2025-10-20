package dbhrm

import (
	"context"
	"europm/internal/config"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

var (
	databaseUrl         string
	databaseUser        string
	databasePassword    string
	databaseMaxPoolSize int
	databaseMinPoolSize int
	databaseSSLMode     string
	Enviroment          string

	Pool *pgxpool.Pool
)

func Init() (err error) {
	//config DB
	databaseUrl = config.GetString("db.postgres.url")
	databaseUser = config.GetString("db.postgres.user")
	databasePassword = config.GetString("db.postgres.password")
	databaseMaxPoolSize = config.GetInt("db.postgres.maxPoolSize")
	databaseMinPoolSize = config.GetInt("db.postgres.minPoolSize")
	databaseSSLMode = viper.GetString("db.postgres.sslMode")
	databaseName := config.GetString("db.postgres.database")
	Enviroment = config.GetString("enviroment")
	log.Println("Init db with host: ", databaseUrl)
	var psqlInfo string
	if strings.Compare("prod", Enviroment) == 0 {
		psqlInfo = "postgresql://" + databaseUrl + "/" + databaseName + "?sslmode=" + databaseSSLMode +
			"&user=" + databaseUser +
			"&pool_max_conns=" + strconv.Itoa(databaseMaxPoolSize) +
			"&pool_min_conns=" + strconv.Itoa(databaseMinPoolSize) +
			"&password=" + databasePassword +
			"&target_session_attrs=read-write"
	} else {
		psqlInfo = "postgresql://" + databaseUrl + "/" + databaseName + "?sslmode=" + databaseSSLMode +
			"&user=" + databaseUser +
			"&pool_max_conns=" + strconv.Itoa(databaseMaxPoolSize) +
			"&pool_min_conns=" + strconv.Itoa(databaseMinPoolSize) +
			"&password=" + databasePassword
	}
	Pool, err = pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		log.Fatal("", err)
		return err
	}
	err = Pool.Ping(context.Background())
	if err != nil {
		return err
	}
	log.Println("Database connection established successfully")
	return nil
}
