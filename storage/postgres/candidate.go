package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	pb "github.com/saladin2098/month3/lesson11/public_service/genproto"
)
type CandidateStorage struct {
	db *sql.DB
}
func NewCandidateStorage(db *sql.DB) *CandidateStorage {
	return &CandidateStorage{
		db: db,
	}
}
func (cs *CandidateStorage) CreateCandidate(c *pb.CandidateCreate, id *string) (*pb.Void, error) {
	query := `insert into candidate(
		id,
		election_id,
		public_id,
		party_id) values($1, $2, $3, $4)`
	_, err := cs.db.Exec(query, id, c.ElectionId, c.PublicId, c.PartyId)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (cs *CandidateStorage) DeleteCandidate(c *pb.ById) (*pb.Void, error) {
	query := `delete from candidate where id = $1`
	_, err := cs.db.Exec(query, c.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (cs *CandidateStorage) UpdateCandidate(c *pb.Candidate) (*pb.Void, error) {

	var conditions []string
	var args []interface{}
	query := `update candidate set`

	if c.Election!= "" {
		conditions = append(conditions, fmt.Sprintf("election_id = $%d", len(args)+1))
		args = append(args, c.Election)
	}
	if c.Public!= "" {
        conditions = append(conditions, fmt.Sprintf("public_id = $%d", len(args)+1))
        args = append(args, c.Public)
    }
	if c.Party!= "" {
        conditions = append(conditions, fmt.Sprintf("party_id = $%d", len(args)+1))
        args = append(args, c.Party)
    }
	query += strings.Join(conditions, ", ")
	query += fmt.Sprintf("where id = $%d", len(args)+1)
	args = append(args, c.Id)

	_, err := cs.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (cs *CandidateStorage) GetByIdCandidate(c *pb.ById) (*pb.Candidate, error) {
	query := `select 
		id, 
		election_id, 
		public_id, 
		party_id
	from 
		candidate 
	where 
		id = $1`
    row := cs.db.QueryRow(query, c.Id)
    var candidate pb.Candidate
    err := row.Scan(
		&candidate.Id, 
		&candidate.Election, 
		&candidate.Party, 
		&candidate.Public)
    if err!= nil {
        return nil, err
    }
    return &candidate, nil
}
func (cs *CandidateStorage) GetAllCandidates(filter *pb.Filter) (*pb.GetAllCandidate, error) {
	query := `select 
        id, 
        election_id, 
        public_id, 
        party_id
    from 
        candidate`
	var conditions []string
	var args []interface{}	
	if filter.Party != "" {
		conditions = append(conditions, fmt.Sprintf("party_id = $%d", len(args)+1))
        args = append(args, filter.Party)
	}
	if len(conditions) > 0 {
        query = query + " WHERE " + strings.Join(conditions, " AND ")
    }
	rows, err := cs.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out pb.GetAllCandidate
	for rows.Next() {
		var candidate pb.Candidate
        err := rows.Scan(
            &candidate.Id, 
            &candidate.Election, 
            &candidate.Party, 
            &candidate.Public)
        if err!= nil {
            return nil, err
        }
        out.Candidates = append(out.Candidates, &candidate)
	}
	if err = rows.Err(); err!= nil {
        return nil, err
    }
	return &out, nil
}