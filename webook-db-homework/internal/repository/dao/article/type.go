package dao

type Article struct {
	ID        int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title     string `gorm:"type:varchar(1024)" bson:"title,omitempty"`
	Content   string `gorm:"type:blob" bson:"content,omitempty"`
	AuthorID  int64  `gorm:"index" bson:"author_id,omitempty"`
	CreatedAt int64  `gorm:"column:created_at" bson:"created_at,omitempty"` // 明确指定列名
	UpdatedAt int64  `gorm:"column:updated_at" bson:"updated_at,omitempty"`
	Status    uint8  `gorm:"column:status" bson:"status,omitempty"`
}

type ReaderArticle Article
