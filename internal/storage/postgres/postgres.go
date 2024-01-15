package postgres

import (
	"database/sql"
	"fmt"
	"gRPCserver/internal/storage"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// TODO: add graceful shutdown
type File struct {
	Name          string
	Creation_date string
	Update_date   string
}

type Storage struct {
	Db *sql.DB
}

func New(host, user, password, dbname string, port int) (*Storage, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("openning db error: %w", err)
	}

	stmt := `CREATE TABLE IF NOT EXISTS files(
        id INTEGER PRIMARY KEY NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
        name VARCHAR (50) NOT NULL UNIQUE,
        creation_date timestamp NOT NULL,
		update_date timestamp NOT NULL);`

	_, err = db.Exec(stmt)
	if err != nil {
		return nil, fmt.Errorf("execution command error: %w", err)
	}

	return &Storage{Db: db}, nil
}

func (s *Storage) SaveFile(filename string) error {
	curr := time.Now()

	dateTime := fmt.Sprintf("%s %s:%s:%s", fmt.Sprint(curr.Date()), strconv.Itoa(curr.Hour()), strconv.Itoa(curr.Minute()), strconv.Itoa(curr.Second()))

	_, err := s.Db.Exec(storage.RequestSaveFile, filename, dateTime, dateTime)
	if err != nil {
		return fmt.Errorf("saving database error: %w", err)
	}

	return nil
}

func (s *Storage) DeleteFile(filename string) error {
	_, err := s.Db.Exec(storage.RequestDeleteFile, filename)
	if err != nil {
		return fmt.Errorf("deleting database error: %w", err)
	}

	return nil
}

func (s *Storage) UpdateFile(filename string) error {
	curr := time.Now()

	dateTime := fmt.Sprintf("%s %s:%s:%s", fmt.Sprint(curr.Date()), strconv.Itoa(curr.Hour()), strconv.Itoa(curr.Minute()), strconv.Itoa(curr.Second()))

	_, err := s.Db.Exec(storage.RequestUpdateFile, dateTime, filename)
	if err != nil {
		return fmt.Errorf("updating database error: %w", err)
	}

	return nil
}

func (s *Storage) GetFullData() ([]File, error) {
	rows, err := s.Db.Query(storage.RequestGetFullData)
	if err != nil {
		return nil, fmt.Errorf("gettingFullData database error: %w", err)
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File
		if err := rows.Scan(&file.Name, &file.Creation_date, &file.Update_date); err != nil {
			return files, fmt.Errorf("scanning data from db error: %w", err)
		}
		file.Creation_date = file.Creation_date[0:10] + " " + file.Creation_date[11:19]
		file.Update_date = file.Update_date[0:10] + " " + file.Update_date[11:19]

		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return files, fmt.Errorf("scanning data from db error: %w", err)
	}
	return files, nil
}
