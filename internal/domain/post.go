package domain

import (
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"time"
)

type Post struct {
	ID            string       `json:"id,omitempty"`
	Club          EventClub    `json:"club"`
	Title         string       `json:"title,omitempty"`
	Description   string       `json:"description,omitempty"`
	Tags          []string     `json:"tags,omitempty"`
	CoverImages   []CoverImage `json:"cover_images,omitempty"`
	AttachedFiles []EventFile  `json:"attached_files,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
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
