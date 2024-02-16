package postgres

import (
	"context"
	"fmt"
	"github.com/ciazhar/golang-example/generated/db"
	logger "github.com/ciazhar/golang-example/pkg/log"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func Init() (*pgxpool.Pool, *db.Queries) {
	urlExample := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?search_path=%s",
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.schema"))

	c, err := pgxpool.ParseConfig(urlExample)
	if err != nil {
		panic(err.Error())
	}

	if viper.GetString("profile") == "debug" || viper.GetString("profile") == "test" {
		c.ConnConfig.Logger = zapadapter.NewLogger(logger.Logger)
	}
	c.ConnConfig.Logger = zapadapter.NewLogger(logger.Logger)

	conn, err := pgxpool.ConnectConfig(context.Background(), c)
	if err != nil {
		panic(err.Error())
	}

	q := db.New(conn)

	return conn, q
}
