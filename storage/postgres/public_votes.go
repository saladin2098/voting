package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	pb "github.com/saladin2098/month3/lesson11/public_service/genproto"
)
type PublicVotesStorage struct {
	db *sql.DB
}

func NewPublicVotesStorage(db *sql.DB) *PublicVotesStorage {
	return &PublicVotesStorage{
		db: db,
	}
}
func (ps *PublicVotesStorage) CreatePublicVotes(pv *pb.PublicVoteCreate, id *string, id2 *string) (*pb.Void, error) {
	tx, err := ps.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO public_vote(id, election_id, public_id) VALUES($1, $2, $3)`
	_, err = tx.Exec(query, id, pv.ElectionId, pv.PublicId)
	if err != nil {
		return nil, err
	}

	query2 := `INSERT INTO vote(id, election_id, candidate_id) VALUES($1, $2, $3)`
	_, err = tx.Exec(query2, id2, pv.ElectionId, pv.CandidateId)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}
func (ps *PublicVotesStorage) DeletePublicVotes(pv *pb.ById) (*pb.Void, error) {
	query := `delete from public_vote where id = $1`
	_, err := ps.db.Exec(query, pv.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (ps *PublicVotesStorage) UpdatePublicVotes(pv *pb.PublicVote) (*pb.Void, error) {
	query := `update public_vote set `
	var conditions []string
	var args []interface{}
	if pv.ElectionId != "" {
		conditions = append(conditions, fmt.Sprintf("election_id = $%d", len(args)+1))
		args = append(args, pv.ElectionId)
	}
	if pv.PublicId != "" {
		conditions = append(conditions, fmt.Sprintf("public_id = $%d", len(args)+1))
		args = append(args, pv.PublicId)
	}
	query += strings.Join(conditions, ", ")
	query += fmt.Sprintf("where id = $%d", len(args)+1)
	args = append(args, pv.Id)
	_, err := ps.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (ps *PublicVotesStorage) GetByIdPublicVote(id *pb.ById) (*pb.PublicVote, error) {
	query := `select 
        id, 
        election_id, 
        public_id
    from 
        public_vote 
    where 
        id = $1`
	row := ps.db.QueryRow(query, id.Id)
	var pv pb.PublicVote
	err := row.Scan(
		&pv.Id,
		&pv.ElectionId,
		&pv.PublicId)
	if err != nil {
		return nil, err
	}
	return &pv, nil
}
func (ps *PublicVotesStorage) GetAllPublucVotes(filter *pb.Filter) (*pb.GetAllPublicVote, error) {
	query := `select 
        id, 
        election_id, 
        public_id
    from 
        public_vote`
	var conditions []string
	var args []interface{}
	if filter.Election != "" {
		conditions = append(conditions, fmt.Sprintf("election_id = $%d", len(args)+1))
		args = append(args, filter.Election)
	}
	if len(conditions) > 0 {
		query = query + " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out pb.GetAllPublicVote
	for rows.Next() {
		var pv pb.PublicVote
		err := rows.Scan(
			&pv.Id,
			&pv.ElectionId,
			&pv.PublicId)
		if err != nil {
			return nil, err
		}
		out.PublicVotes = append(out.PublicVotes, &pv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &out, nil
}
func (ps *PublicVotesStorage) FindWinner(we *pb.WhichElection) (*pb.Winner, error) {
	query := `
		SELECT election_id, candidate_id, COUNT(*) as vote_count
		FROM votes
		WHERE election_id = $1
		GROUP BY election_id, candidate_id
		ORDER BY vote_count DESC
		LIMIT 1`

	row := ps.db.QueryRow(query, we.ElectionId)

	var winner pb.Winner
	err := row.Scan(&winner.ElectionId, &winner.CandidateId, &winner.Votes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no votes found for the given election")
		}
		return nil, err
	}
	return &winner, nil
}
