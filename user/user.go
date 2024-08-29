package user

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:false"`
}
