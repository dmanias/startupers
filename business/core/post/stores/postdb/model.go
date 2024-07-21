package postdb

import (
	"github.com/dmanias/startupers/business/core/post"
	"time"

	"github.com/google/uuid"
)

type dbPost struct {
	ID          uuid.UUID `db:"id"`
	IdeaID      uuid.UUID `db:"idea_id"`
	AuthorID    uuid.UUID `db:"author_id"`
	Content     string    `db:"content"`
	OwnerType   string    `db:"owner_type"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBPost(post post.Post) dbPost {
	return dbPost{
		ID:          post.ID,
		IdeaID:      post.IdeaID,
		AuthorID:    post.AuthorID,
		Content:     post.Content,
		OwnerType:   post.OwnerType,
		DateCreated: post.DateCreated.UTC(),
		DateUpdated: post.DateUpdated.UTC(),
	}
}

func toCorePost(dbPost dbPost) post.Post {
	return post.Post{
		ID:          dbPost.ID,
		IdeaID:      dbPost.IdeaID,
		AuthorID:    dbPost.AuthorID,
		Content:     dbPost.Content,
		OwnerType:   dbPost.OwnerType,
		DateCreated: dbPost.DateCreated.In(time.Local),
		DateUpdated: dbPost.DateUpdated.In(time.Local),
	}
}

func toCorePostSlice(dbPosts []dbPost) []post.Post {
	posts := make([]post.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = toCorePost(dbPost)
	}
	return posts
}
