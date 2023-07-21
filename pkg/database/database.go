package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Client interface {
	Close()
	Create(tableName string, columns []string, values ...interface{}) (int, error)
}

type client struct {
	db *sql.DB
}

func NewClient(dsn string) (Client, error) {
	db, err := connect(dsn)
	if err != nil {
		return nil, err
	}
	return &client{db}, nil
}

func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Error %s when pinging DB", err)
		return nil, err
	}

	log.Printf("Connected to %s successfully\n", dsn)
	return db, nil
}

func (c *client) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *client) Create(tableName string, columns []string, values ...interface{}) (int, error) {
	query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
		tableName, strings.Join(columns, ", "), strings.Repeat("?", len(columns)))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		log.Printf("Error %s when inserting row into table", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error %s while getting created row ID", err)
		return 0, err
	}
	log.Printf("Created a new row")

	return int(id), nil
}
