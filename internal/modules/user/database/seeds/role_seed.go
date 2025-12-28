package seeds

import (
	"fmt"
	"log"
	"tofash/internal/modules/user/model"

	"gorm.io/gorm"
)

func SeedRole(db *gorm.DB) {
	roles := []model.Role{
		{
			Name: "Admin",
		},
		{
			Name: "Customer",
		},
	}

	for _, v := range roles {
		var role model.Role
		err := db.Where("name = ?", v.Name).First(&role).Error
		if err != nil {
			err := db.Create(&v).Error
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Role " + v.Name + " created")
		}
	}
}
