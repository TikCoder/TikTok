package service

import (
	"TikTok/utils/db"
)

type RelationService struct {
}

// IsFollow 查询是否关注
func (r *RelationService) IsFollow(followId, followerId int64) (bool, error) {
	return db.IsFollow(db.DB, followId, followerId)
}
