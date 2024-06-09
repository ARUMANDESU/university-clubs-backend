package domain

import (
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"time"
)

type Post struct {
	ID            string
	Club          EventClub
	Title         string
	Description   string
	Tags          []string
	CoverImages   []CoverImage
	AttachedFiles []EventFile
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func PostFromPb(post *postv1.PostObject) *Post {
	if post == nil {
		return nil
	}

	return &Post{
		ID:            post.Id,
		Club:          *ProtoToEventClub(post.Club),
		Title:         post.Title,
		Description:   post.Description,
		Tags:          post.Tags,
		CoverImages:   ProtoToCoverImageArr(post.CoverImages),
		AttachedFiles: ProtoToEventFileArr(post.AttachedFiles),
		CreatedAt:     post.CreatedAt.AsTime(),
		UpdatedAt:     post.UpdatedAt.AsTime(),
	}
}

func PostsFromPb(posts []*postv1.PostObject) []Post {
	if posts == nil {
		return nil
	}

	pbPosts := make([]Post, 0, len(posts))
	for _, post := range posts {
		pbPosts = append(pbPosts, *PostFromPb(post))
	}

	return pbPosts
}
