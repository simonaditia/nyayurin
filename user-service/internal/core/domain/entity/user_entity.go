package entity

type UserEntity struct {
	ID         int64 `gorm:"primaryKey"`
	Name       string
	Email      string
	Password   string
	RoleName   string
	Address    string
	Lat        string
	Lng        string
	Phone      string
	Photo      string
	IsVerified bool
}
