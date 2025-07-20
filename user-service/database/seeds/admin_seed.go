package seeds

import (
	"log"

	"github.com/simonaditia/nyayurin/user-service/internal/core/domain/model"
	"github.com/simonaditia/nyayurin/user-service/utils/conv"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	bytes, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}

	modelRole := model.Role{}
	err = db.Where("name = ?", "Super Admin").First(&modelRole).Error

	admin := model.User{
		Name:       "super admin",
		Email:      "superadmin@gmail.com",
		Password:   bytes,
		IsVerified: true,
		Roles:      []model.Role{modelRole},
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: admin.Email}).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	} else {
		log.Printf("Admin %s seeded successfully", admin.Email)
	}
}
