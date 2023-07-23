package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Client interface {
	Close()
	Create(tableName string, columns []string, values ...interface{}) (int, error)
	Update(
		tableName string,
		uniqueFieldName string,
		uniqueFieldValue interface{},
		columns []string,
		values ...interface{}) error
	Exists(tableName string, uniqueFiledName string, uniqueFieldValue interface{}) (bool, error)
	FindUnique(resultStruct interface{}, tableName string, uniqueFiledName string, uniqueFieldValue interface{}) error
	FindMany(resultStruct interface{}, tableName string, condition *string, limit *int) error
	Delete(tableName string, uniqueFieldName string, uniqueFieldValue interface{}) error
	DeleteAll(tableName string) error
	GetDB() *sqlx.DB
}

type client struct {
	db *sqlx.DB
}

func NewClient(dsn string) (Client, error) {
	db, err := connect(dsn)
	if err != nil {
		return nil, err
	}
	return &client{db}, nil
}

func connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
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

// Should only be used to execute raw queries in services
func (c *client) GetDB() *sqlx.DB {
	return c.db
}

func (c *client) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *client) Create(tableName string, columns []string, values ...interface{}) (int, error) {
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ","))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.db.ExecContext(ctx, query, values...)
	if err != nil {
		log.Printf("Error %s when inserting row into table", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error %s while getting created row ID", err)
		return 0, err
	}
	log.Printf("Created a new row in table %s", tableName)

	return int(id), nil
}

func (c *client) Update(
	tableName string,
	uniqueFieldName string,
	uniqueFieldValue interface{},
	columns []string,
	values ...interface{}) error {
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = fmt.Sprintf("%s = ?", columns[i])
	}

	query := fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE %s = ?",
		tableName,
		strings.Join(placeholders, ", "),
		uniqueFieldName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	values = append(values, uniqueFieldValue)

	_, err := c.db.ExecContext(ctx, query, values...)
	if err != nil {
		log.Printf("Error %s when inserting row into table", err)
		return err
	}

	log.Printf("Updated row with table %s", tableName)

	return nil
}

func (c *client) Exists(tableName string, columnName string, value interface{}) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM `%s` WHERE %s = ?", tableName, columnName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := c.db.QueryRowContext(ctx, query, value)

	var exists int
	err := row.Scan(&exists)
	if err == sql.ErrNoRows {
		// No rows found, the row doesn't exist
		return false, nil
	} else if err != nil {
		log.Printf("Error %s when querying the database", err)
		return false, err
	}

	return true, nil
}

func (c *client) FindUnique(resultStruct interface{}, tableName string, columnName string, value interface{}) error {
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE %s = ? LIMIT 1", tableName, columnName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	udb := c.db.Unsafe()
	if err := udb.GetContext(ctx, resultStruct, query, value); err != nil {
		log.Printf("Error %s when executing query", err)
		return err
	}

	return nil
}

func (c *client) FindMany(resultStruct interface{}, tableName string, condition *string, limit *int) error {
	query := fmt.Sprintf("SELECT * FROM `%s`", tableName)

	if condition != nil {
		query += fmt.Sprintf(" WHERE %s", *condition)
	}

	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	udb := c.db.Unsafe()
	if err := udb.SelectContext(ctx, resultStruct, query); err != nil {
		log.Printf("Error %s when executing query", err)
		return err
	}

	return nil
}

func (c *client) Delete(tableName string, uniqueFieldName string, uniqueFieldValue interface{}) error {
	return c.delete(tableName, &uniqueFieldName, uniqueFieldValue)
}

func (c *client) DeleteAll(tableName string) error {
	return c.delete(tableName, nil, nil)
}

func (c *client) delete(tableName string, uniqueFieldName *string, uniqueFieldValue interface{}) error {
	query := fmt.Sprintf("DELETE FROM `%s`", tableName)

	if uniqueFieldName != nil && uniqueFieldValue != nil {
		query += fmt.Sprintf(" WHERE %s = ?", *uniqueFieldName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	if uniqueFieldValue != nil {
		_, err = c.db.ExecContext(ctx, query, uniqueFieldValue)
	} else {
		_, err = c.db.ExecContext(ctx, query)
	}

	if err != nil {
		log.Printf("Error %s when executing query", err)
		return err
	}

	return nil
}
