package repository

type Favorite struct {
	UserId     int64 `gorm:"primaryKey;autoIncrement:false"`
	VideoId    int64 `gorm:"primaryKey;autoIncrement:false"`
	IsFavorite bool
}
