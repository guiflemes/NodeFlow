package repository

import (
	"flowChart/domain"
	"fmt"
	"time"

	"flowChart/settings"

	"github.com/jmoiron/sqlx"
)

type Database struct {
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

func (db *Database) Parse() {
	db.User = settings.GETENV("POSTGRES_USER")
	db.Password = settings.GETENV("POSTGRES_PASSWORD")
	db.Host = settings.GETENV("POSTGRES_HOST")
	db.Port = settings.GETENV("POSTGRES_PORT")
	db.Database = settings.GETENV("POSTGRES_DB_NAME")
}

type PostgresRepo[T comparable] struct {
	client *sqlx.DB
}

func NewPostgresRepo[T comparable](conf *Database) *PostgresRepo[T] {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Password,
		conf.Database,
	)

	db, err := sqlx.Connect("postgres", dns)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(conf.OpenConnection)
	db.SetMaxIdleConns(conf.IdleConnection)
	db.SetConnMaxLifetime(time.Duration(conf.ConnectionLifeTime) * time.Millisecond)
	db.SetConnMaxIdleTime(time.Duration(conf.IdleConnection) * time.Millisecond)

	return &PostgresRepo[T]{
		client: db,
	}
}

func (r *PostgresRepo[T]) Save(flowChart *domain.FlowChart[T]) {}

func (r *PostgresRepo[T]) saveNode(node *domain.Node[T]) {}
