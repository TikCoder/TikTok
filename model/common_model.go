package model

// Response 通用型响应
type Response struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type User struct {
	// gorm.Model
	Id              int64  `gorm:"column:id; primary_key; auto_increment" json:"id"`
	Username        string `gorm:"column:user_name" json:"username"`
	Password        string `gorm:"column:password" json:"password"`
	Name            string `gorm:"column:name" json:"name"`
	WorkCount       int64  `gorm:"column:work_count" json:"work_count"`
	FollowCount     int64  `gorm:"column:follow_count" json:"follow_count"`
	FollowerCount   int64  `gorm:"column:follower_count" json:"follower_count"`
	Avatar          string `gorm:"column:avatar" json:"avatar"`
	BackgroundImage string `gorm:"column:background_image" json:"background_image"`
	Signature       string `gorm:"column:signature" json:"signature"`
	TotalFav        int64  `gorm:"column:total_favorited" json:"total_favorite"` // 我 获赞的总数
	FavCount        int64  `gorm:"column:favorite_count" json:"favorite_count"`  // 我 点赞的视频总数
}

type Video struct {
	Id            int64  `gorm:"column:id; primary_key; auto_increment" json:"id"`              // video_id
	AuthorId      int64  `gorm:"column:author_id;" json:"author_id,omitempty" json:"author_id"` // 谁发布的
	PlayUrl       string `gorm:"column:play_url;" json:"play_url"`                              // videoURL
	CoverUrl      string `gorm:"column:cover_url;" json:"cover_url"`                            // picURL
	FavoriteCount int64  `gorm:"column:favorite_count;" json:"favorite_count"`                  // 点赞数
	CommentCount  int64  `gorm:"column:comment_count;" json:"comment_count"`                    // 评论数
	PublishTime   int64  `gorm:"column:publish_time;" json:"publish_time,omitempty"`            // 发布时间
	Title         string `gorm:"column:title;" json:"title"`                                    // 标题
	// Author        User   `gorm:"foreignkey:AuthorId"`           // 作者
}

type Favorite struct {
	Id      int64 `gorm:"column:id; primary_key; auto_increment" json:"id"`
	UserId  int64 `gorm:"column:user_id;" json:"user_id"`
	VideoId int64 `gorm:"column:video_id;" json:"video_id"`
}

type Relation struct {
	Id         int64 `gorm:"column:id; primary_key; auto_increment" json:"id"`
	FollowId   int64 `gorm:"column:follow_id;" json:"follow_id"`
	FollowerId int64 `gorm:"column:follower_id;" json:"follower_id"`
}
