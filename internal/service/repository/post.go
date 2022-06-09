package repository

import (
	"context"

	"imman/crud_service/internal/entity"

	"github.com/uptrace/bun"
)

type PostRepository struct {
	db *bun.DB
}

func NewPostRepository(db *bun.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r PostRepository) List(ctx context.Context, limit, offset int) ([]*entity.Post, int, error) {
	posts := make([]*entity.Post, 0)
	count, err := r.db.NewSelect().Model(&posts).Limit(limit).Offset(offset).ScanAndCount(ctx)
	return posts, count, err
}

func (r PostRepository) Detail(ctx context.Context, id int) (*entity.Post, error) {
	var post entity.Post
	err := r.db.NewSelect().Model(&post).Where("id = ?", id).Scan(ctx)
	return &post, err
}

func (r PostRepository) Update(ctx context.Context, post *entity.Post) error {
	_, err := r.db.NewUpdate().Model(post).
		Column("user_id", "title", "body").
		Where("id = ?", post.Id).Exec(ctx)
	return err
}

func (r PostRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.NewDelete().Model((*entity.Post)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
