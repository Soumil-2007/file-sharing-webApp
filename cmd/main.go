package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Soumil-2007/file-sharing-webApp/cmd/api"
	"github.com/Soumil-2007/file-sharing-webApp/configs"
	"github.com/Soumil-2007/file-sharing-webApp/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// Query and delete expired files
				rows, err := db.Query("SELECT id, path FROM files WHERE expires_at IS NOT NULL AND expires_at < NOW()")
				if err != nil {
					continue
				}
				for rows.Next() {
					var id int
					var path string
					rows.Scan(&id, &path)
					os.Remove(path) // Delete from disk
					db.Exec("DELETE FROM files WHERE id = ?", id)
				}
				rows.Close()
			}
		}
	}()

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
