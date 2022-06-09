package service

import (
	"context"

	"imman/crud_service/internal/entity"
	"imman/crud_service/protos/protos/crud_pb"
)

type PostRepository interface {
	List(ctx context.Context, limit, offset int) ([]*entity.Post, int, error)
	Detail(ctx context.Context, id int) (*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id int) error
}

type PostService struct {
	repo PostRepository
	crud_pb.UnimplementedPostCRUDServiceServer
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) List(ctx context.Context, req *crud_pb.ListRequest) (*crud_pb.ListResponse, error) {
	limit, offset := int(req.GetLimit()), int(req.GetOffset())
	posts, count, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return &crud_pb.ListResponse{}, err
	}

	var postResponse []*crud_pb.Post
	for i := range posts {
		postResponse = append(postResponse, &crud_pb.Post{
			Id:     int64(posts[i].Id),
			UserId: int64(posts[i].UserId),
			Title:  posts[i].Title,
			Body:   posts[i].Body,
		})
	}

	return &crud_pb.ListResponse{
		Posts: postResponse,
		Count: int64(count),
	}, nil
}

func (s *PostService) Detail(ctx context.Context, req *crud_pb.DetailRequest) (*crud_pb.DetailResponse, error) {
	id := int(req.GetId())
	post, err := s.repo.Detail(ctx, id)
	if err != nil {
		return &crud_pb.DetailResponse{}, err
	}

	return &crud_pb.DetailResponse{
		Post: &crud_pb.Post{
			Id:     int64(post.Id),
			UserId: int64(post.UserId),
			Title:  post.Title,
			Body:   post.Body,
		},
	}, nil
}

func (s *PostService) Update(ctx context.Context, req *crud_pb.UpdateRequest) (*crud_pb.UpdateResponse, error) {
	post := entity.Post{
		Id:     int(req.Post.GetId()),
		UserId: int(req.Post.GetUserId()),
		Title:  req.Post.GetTitle(),
		Body:   req.Post.GetBody(),
	}

	err := s.repo.Update(ctx, &post)
	return &crud_pb.UpdateResponse{}, err
}

func (s *PostService) Delete(ctx context.Context, req *crud_pb.DeleteRequest) (*crud_pb.DeleteResponse, error) {
	id := int(req.GetId())

	err := s.repo.Delete(ctx, id)
	return &crud_pb.DeleteResponse{}, err
}
