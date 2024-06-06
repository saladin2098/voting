package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	"github.com/saladin2098/month3/lesson11/public_service/config"
	storage "github.com/saladin2098/month3/lesson11/public_service/storage"
)

type Storage struct {
	db *sql.DB
	CandidateS storage.CandidateI
	ElectionS storage.ElectionI
	PublicVoteS storage.PublicVoteI
}
func ConnectDB() (*Storage, error) {
	cfg := config.Load()
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase)
    db, err := sql.Open("postgres", dbConn)
	if err!= nil {
        return nil, err
    }
	err = db.Ping()
	if err!= nil {
        return nil, err
    }
	e_repo := NewElectionStorage(db)
	c_repo := NewCandidateStorage(db)
	pv_repo := NewPublicVotesStorage(db)
	return &Storage{
		db: db,
		CandidateS: c_repo,
        ElectionS: e_repo,
        PublicVoteS: pv_repo,
		}, nil
}
func (s *Storage) Election() storage.ElectionI {
	if s.ElectionS == nil {
		s.ElectionS = &ElectionStorage{db: s.db}
	}
	return s.ElectionS
}
func (s *Storage) Candidate() storage.CandidateI {
	if s.CandidateS == nil {
        s.CandidateS = &CandidateStorage{db: s.db}
    }
    return s.CandidateS
}
func (s *Storage) PublicVote() storage.PublicVoteI {
	if s.PublicVoteS == nil {
        s.PublicVoteS = &PublicVotesStorage{db: s.db}
    }
    return s.PublicVoteS
}