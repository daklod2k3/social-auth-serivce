package interfaces

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"shared/entity"
)

type FeedService interface {
	UpdateFeedFromNewPost(postId *primitive.ObjectID, userId *uuid.UUID, postType entity.PostType) error
	GetFeedForUser(userId *uuid.UUID) (interface{}, error)
}
