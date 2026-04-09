package database

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB inicializa o banco de dados SQLite local
func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		return fmt.Errorf("erro ao abrir banco de dados: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	// Habilita chaves estrangeiras
	_, _ = DB.Exec("PRAGMA foreign_keys = ON;")

	if err = createTables(); err != nil {
		return fmt.Errorf("erro ao criar tabelas: %w", err)
	}

	return nil
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT,
			last_name TEXT,
			email TEXT,
			phone TEXT,
			address TEXT,
			age INTEGER,
			photo TEXT,
			objective TEXT,
			linkedin TEXT,
			github TEXT,
			website TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS education (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id INTEGER,
			institution TEXT,
			course TEXT,
			start_date TEXT,
			end_date TEXT,
			description TEXT,
			FOREIGN KEY(profile_id) REFERENCES profiles(id)
		);`,
		`CREATE TABLE IF NOT EXISTS experiences (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id INTEGER,
			company TEXT,
			position TEXT,
			start_date TEXT,
			end_date TEXT,
			description TEXT,
			FOREIGN KEY(profile_id) REFERENCES profiles(id)
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id INTEGER,
			name TEXT,
			description TEXT,
			url TEXT,
			FOREIGN KEY(profile_id) REFERENCES profiles(id)
		);`,
		`CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			language TEXT DEFAULT 'pt',
			template TEXT DEFAULT 'default',
			show_photo INTEGER DEFAULT 1,
			labels TEXT
		);`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}

	// Migrações manuais para colunas novas (caso o banco já exista)
	_, _ = DB.Exec("ALTER TABLE profiles ADD COLUMN linkedin TEXT")
	_, _ = DB.Exec("ALTER TABLE profiles ADD COLUMN github TEXT")
	_, _ = DB.Exec("ALTER TABLE profiles ADD COLUMN website TEXT")
	_, _ = DB.Exec("ALTER TABLE settings ADD COLUMN labels TEXT")

	return nil
}
