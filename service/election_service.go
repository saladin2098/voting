package service

import (
	"context"

    "github.com/google/uuid"
    pb "github.com/saladin2098/month3/lesson11/public_service/genproto"
    postgres "github.com/saladin2098/month3/lesson11/public_service/storage/postgres"
)

type ElectionService struct {
	stg *postgres.Storage
	pb.UnimplementedElectionServiceServer
}
func NewElectionService(stg *postgres.Storage) *ElectionService {
    return &ElectionService{
        stg: stg,
    }
}
func (es *ElectionService) CreateElection(ctx context.Context,e *pb.Election) (*pb.Void, error) {
    id := uuid.NewString()
	e.Id = id
    return es.stg.ElectionS.CreateElection(e)
}
func (es *ElectionService) DeleteElection(ctx context.Context,id *pb.ById) (*pb.Void, error) {
    return es.stg.ElectionS.DeleteElection(id)
}
func (es *ElectionService) UpdateElection(ctx context.Context,e *pb.Election) (*pb.Void, error) {
    return es.stg.ElectionS.UpdateElection(e)
}
func (es *ElectionService) GetByIdElection(ctx context.Context,id *pb.ById) (*pb.Election, error) {
    return es.stg.ElectionS.GetByIdElection(id)
}
func (es *ElectionService) GetAllElections(ctx context.Context,f *pb.Filter) (*pb.GetAllElection, error) {
    return es.stg.ElectionS.GetAllElections(f)
}