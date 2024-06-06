package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	pb "github.com/saladin2098/month3/lesson11/public_service/genproto"
)

type ElectionStorage struct {
	db *sql.DB
}
func NewElectionStorage(db *sql.DB) *ElectionStorage {
	return &ElectionStorage{
		db: db,
	}
}
func (es *ElectionStorage) CreateElection(el *pb.Election) (*pb.Void, error) {
	query := `insert into election(id,name,date) values($1, $2, $3)`
	_, err := es.db.Exec(query, el.Id, el.Name, el.Date)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (es *ElectionStorage) DeleteElection(id *pb.ById) (*pb.Void, error) {
	query := `delete from election where id = $1`
	_, err := es.db.Exec(query, id.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (es *ElectionStorage) UpdateElection(el *pb.Election) (*pb.Void, error) {
	query := `update election set name = $1, date = $2 where id = $3`
	_, err := es.db.Exec(query, el.Name, el.Date, el.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Void{}, nil
}
func (es *ElectionStorage) GetByIdElection(id *pb.ById) (*pb.Election, error) {
	query := `select * from election where id = $1`
	row := es.db.QueryRow(query, id.Id)
	el := &pb.Election{}
	err := row.Scan(&el.Id, &el.Name, &el.Date)
	if err != nil {
		return nil, err
	}
	return el, nil
}
func (es ElectionStorage) GetAllElections(filter *pb.Filter) (*pb.GetAllElection, error) {
	var conditions []string
	var args []interface{}

	query := `SELECT id, name, date FROM election`
	if filter.Date != "" {
		conditions = append(conditions, fmt.Sprintf("date = $%d", len(args)+1))
		args = append(args, filter.Date)
	}
	if len(conditions) > 0 {
		query = query + " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := es.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() 

	elections := &pb.GetAllElection{
		Elections: []*pb.Election{},
	}

	for rows.Next() {
		el := &pb.Election{}
		err := rows.Scan(&el.Id, &el.Name, &el.Date)
		if err != nil {
			return nil, err
		}
		elections.Elections = append(elections.Elections, el)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return elections, nil
}
