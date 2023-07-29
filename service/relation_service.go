package service

import "tiktok2023/db"

type RelationService struct {
}

// IsFollow 查询是否关注
func (r *RelationService) IsFollow(followId, followerId int64) (bool, error) {
	return db.IsFollow(db.DB, followId, followerId)
}
