package db

import (
	"database/sql"

	"codeexec/internal/config"

	log "github.com/sirupsen/logrus"
)

var db *DB

type DB struct {
	conn *sql.DB
}

func InitDB(cfg *config.Config) *DB {
	dsn := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " password=" + cfg.DBPass + " dbname=" + cfg.DBName + " sslmode=disable"
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}
	db = &DB{conn: dbConn}
	return db
}

func GetDB() *DB {
	if db == nil {
		log.Fatalf("db not initialized")
	}
	return db
}

func (d *DB) Close() {
	d.conn.Close()
}

func (d *DB) GetQueries() *Queries {
	return New(d.conn)
}

func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func NullFloatToFloat(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}
