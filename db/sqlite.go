package db

import (
	"database/sql"
	"log"
)

type Service interface {
	TextWriter
}

type TextWriter interface {
	AddText(text, hash, url, intent string, prob float64) error
	ContainsHash(hash string) (bool, error)
}

type service struct {
	db *sql.DB
}

func CreateSqliteService() service {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	return service{db}
}

func (s *service) AddLanding(url, hash, name string) error {
	stmt, err := s.db.Prepare(`INSERT into landings (url, hash, name) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(url, hash, name)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddURL(url, hash string) error {
	stmt, err := s.db.Prepare(`INSERT into webpages (url, hash, parsed) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(url, hash, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetURLHash(url string) (string, error) {
	stmt, err := s.db.Prepare(`select hash from webpages where url = ?`)
	if err != nil {
		return "", err
	}
	var hash string
	row := stmt.QueryRow(url)
	err = row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (s *service) SetURLHash(url, hash string) error {
	stmt, err := s.db.Prepare(`update webpages set hash = ? where url = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(hash, url)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SetURLParsed(url string) error {
	stmt, err := s.db.Prepare(`update webpages set parsed = 1 where url = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(url)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SetURLUnparsed(url string) error {
	stmt, err := s.db.Prepare(`update webpages set parsed = 0 where url = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(url)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddText(text, hash, url, intent string, prob float64) error {
	stmt, err := s.db.Prepare(`INSERT into events (url, hash, article, intent, probability) values (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(url, hash, text, intent, prob)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ContainsHash(hash string) (bool, error) {
	stmt, err := s.db.Prepare(`select count(*) from events where hash = ?`)
	if err != nil {
		return false, err
	}
	var c int
	row := stmt.QueryRow(hash)
	err = row.Scan(&c)
	if err != nil {
		return false, err
	}
	return c != 0, nil
}

func (s *service) GetUnparsedURLs() ([]string, error) {
	stmt, err := s.db.Prepare(`select url from webpages where parsed = 0`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []string
	for rows.Next() {
		var url string
		err := rows.Scan(&url)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
