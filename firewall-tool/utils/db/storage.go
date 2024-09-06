package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

func IsSQLite3Installed() bool {
	cmd := exec.Command("which", "sqlite3")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func InstallSQLite3() error {
	pkgManager := getPkgManager()
	cmd := exec.Command("sh", "-c", pkgManager+" install -y sqlite3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getPkgManager() string {
	// Determine the package manager based on the OS
	// For simplicity, assuming Linux with apt-get or yum
	if _, err := exec.LookPath("apt-get"); err == nil {
		return "sudo apt-get"
	}
	if _, err := exec.LookPath("yum"); err == nil {
		return "sudo yum"
	}
	return "sudo apt-get" // Default to apt-get if no other package manager is found
}

func IsSQLite3Running() bool {
	cmd := exec.Command("pgrep", "sqlite3")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func StartSQLite3(dbFile string) error {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			return fmt.Errorf("failed to create database file: %v", err)
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        no INTEGER,
        time TEXT,
        source TEXT,
        destination TEXT,
        protocol TEXT,
        length INTEGER,
        info TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	cmd := exec.Command("sqlite3", dbFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
