package service


import (
    "context"
    "github.com/google/uuid"
	pb "github.com/saladin2098/month3/lesson11/public_service/genproto"
	postgres "github.com/saladin2098/month3/lesson11/public_service/storage/postgres"
)
type PublicVoteService struct {
	stg *postgres.Storage
	pb.UnimplementedPublicVoteServiceServer
}
func NewPublicVoteService(stg *postgres.Storage) *PublicVoteService {
    return &PublicVoteService{
        stg: stg,
    }
}
func (ps *PublicVoteService) CreatePublicVote(ctx context.Context,pv *pb.PublicVoteCreate) (*pb.Void, error) {
    id1 := uuid.NewString()
	id2 := uuid.NewString()
    return ps.stg.PublicVoteS.CreatePublicVotes(pv,&id1,&id2)
}
func (ps *PublicVoteService) DeletePublicVote(ctx context.Context,id *pb.ById) (*pb.Void, error) {
    return ps.stg.PublicVoteS.DeletePublicVotes(id)
}
func (ps *PublicVoteService) UpdatePublicVote(ctx context.Context,pv *pb.PublicVote) (*pb.Void, error) {
    return ps.stg.PublicVoteS.UpdatePublicVotes(pv)
}
func (ps *PublicVoteService) GetByIdPublicVote(ctx context.Context,id *pb.ById) (*pb.PublicVote, error) {
    return ps.stg.PublicVoteS.GetByIdPublicVote(id)
}
func (ps *PublicVoteService) GetAllPublicVotes(ctx context.Context,f *pb.Filter) (*pb.GetAllPublicVote, error) {
    return ps.stg.PublicVoteS.GetAllPublicVotes(f)
}
func (ps *PublicVoteService) FindWinner(ctx context.Context,we *pb.WhichElection) (*pb.Winner, error) {
    return ps.stg.PublicVoteS.FindWinner(we)
}