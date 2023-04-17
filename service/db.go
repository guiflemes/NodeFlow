package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"flowChart/settings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	User               string
	Password           string
	Host               string
	Port               string
	Database           string
	IdleConnection     int
	OpenConnection     int
	ConnectionLifeTime int
	ConnectionIdleTime int
	ReadTimeout        int
	WriteTimeout       int
	Timeout            int
}

func (conf *DatabaseConfig) Parse() {
	conf.User = settings.GETENV("POSTGRES_USER")
	conf.Password = settings.GETENV("POSTGRES_PASSWORD")
	conf.Host = settings.GETENV("POSTGRES_HOST")
	conf.Port = settings.GETENV("POSTGRES_PORT")
	conf.Database = settings.GETENV("POSTGRES_DB_NAME")
}

func NewPostgresDb(conf *DatabaseConfig) *sqlx.DB {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s connect_timeout=30 sslmode=disable",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Database,
	)

	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err)
	}

	mydb := sqlx.NewDb(db, "postgres")

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	mydb.SetMaxOpenConns(conf.OpenConnection)
	mydb.SetMaxIdleConns(conf.IdleConnection)
	mydb.SetConnMaxLifetime(time.Duration(30) * time.Millisecond)
	mydb.SetConnMaxIdleTime(time.Duration(30) * time.Millisecond)

	return mydb
}
