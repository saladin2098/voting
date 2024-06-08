package postgres_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	genproto "github.com/saladin2098/month3/lesson11/public_service/genproto"
	postgres "github.com/saladin2098/month3/lesson11/public_voting/storage/postgres"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	return db, mock
}


func TestCreateCandidate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	candidateStorage := postgres.NewCandidateStorage(db)

	candidateCreate := &genproto.CandidateCreate{
		ElectionId: "1",
		PublicId:   "2",
		PartyId:    "3",
	}
	id := "1"
	query := `INSERT INTO candidate \(id, election_id, public_id, party_id\) VALUES \(\$1, \$2, \$3, \$4\)`
	mock.ExpectExec(query).WithArgs(id, candidateCreate.ElectionId, candidateCreate.PublicId, candidateCreate.PartyId).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := candidateStorage.CreateCandidate(candidateCreate, &id); err != nil {
		t.Errorf("error creating candidate: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteCandidate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	candidateStorage := postgres.NewCandidateStorage(db)

	id := "1"
	query := `DELETE FROM candidate WHERE id = \$1`
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := candidateStorage.DeleteCandidate(&genproto.ById{Id: id}); err != nil {
		t.Errorf("error deleting candidate: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateCandidate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	candidateStorage := postgres.NewCandidateStorage(db)

	candidate := &genproto.Candidate{
		Id:       "1",
		Election: "1",
		Public:   "2",
		Party:    "3",
	}

	query := `UPDATE candidate SET election_id = \$1, public_id = \$2, party_id = \$3 WHERE id = \$4`
	mock.ExpectExec(query).WithArgs(candidate.Election, candidate.Public, candidate.Party, candidate.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err := candidateStorage.UpdateCandidate(candidate); err != nil {
		t.Errorf("error updating candidate: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
func TestGetByIdCandidate(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	candidateStorage := postgres.NewCandidateStorage(db)

	id := "1"
	query := `SELECT id, election_id, public_id, party_id FROM candidate WHERE id = \$1`
	rows := sqlmock.NewRows([]string{"id", "election_id", "public_id", "party_id"}).AddRow("1", "1", "2", "3")
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(rows)

	got, err := candidateStorage.GetByIdCandidate(&genproto.ById{Id: id})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want := &genproto.Candidate{
		Id:       "1",
		Election: "1",
		Public:   "2",
		Party:    "3",
	}
	if *got != *want {
		t.Errorf("got %v, want %v", got, want)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}


func TestGetAllCandidates(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	candidateStorage := postgres.NewCandidateStorage(db)

	query := `SELECT id, election_id, public_id, party_id FROM candidate`
	rows := sqlmock.NewRows([]string{"id", "election_id", "public_id", "party_id"}).
		AddRow("1", "1", "2", "3").
		AddRow("2", "1", "3", "4")
	mock.ExpectQuery(query).WillReturnRows(rows)

	got, err := candidateStorage.GetAllCandidates(&genproto.Filter{})
	if err != nil {
		t.Fatalf("failed to get all candidates: %v", err)
	}

	want := &genproto.GetAllCandidate{
		Candidates: []*genproto.Candidate{
			{Id: "1", Election: "1", Public: "2", Party: "3"},
			{Id: "2", Election: "1", Public: "3", Party: "4"},
		},
	}
	if len(got.Candidates) != len(want.Candidates) {
		t.Errorf("got %d candidates, want %d", len(got.Candidates), len(want.Candidates))
	}
	for i, gotCandidate := range got.Candidates {
		wantCandidate := want.Candidates[i]
		if *gotCandidate != *wantCandidate {
			t.Errorf("got %v, want %v", gotCandidate, wantCandidate)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

