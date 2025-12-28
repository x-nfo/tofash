package seeds

import (
	"fmt"
	"log"
	"tofash/internal/modules/product/model"

	"gorm.io/gorm"
)

func SeedProduct(db *gorm.DB) {
	category := model.Category{
		Name: "Electronics",
		Slug: "electronics",
	}

	if err := db.FirstOrCreate(&category, model.Category{Name: "Electronics"}).Error; err != nil {
		log.Printf("Failed to seed category: %v", err)
	}

	products := []model.Product{
		{
			Name:         "Cotton T-Shirt Red",
			CategorySlug: category.Slug,
			Image:        "tshirt-red.jpg",
			RegulerPrice: 150000,
			SalePrice:    120000,
			Description:  "Premium cotton t-shirt for local fashion dev testing",
			Stock:        50,
			Weight:       250,
			Status:       "ACTIVE",
			Unit:         "pcs",
			Variant:      1,
			// Fashion attributes
			SKU:        "TSHIRT-RED-L-001",
			Size:       "L",
			Color:      "Red",
			Material:   "100% Cotton",
			ImagesJSON: `["tshirt-red-front.jpg","tshirt-red-back.jpg","tshirt-red-detail.jpg"]`,
		},
	}

	for _, p := range products {
		var existing model.Product
		if err := db.Where("sku = ? OR name = ?", p.SKU, p.Name).First(&existing).Error; err != nil {
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Failed to create product %s: %v", p.Name, err)
			} else {
				fmt.Println("Product " + p.Name + " seeded.")
			}
		} else {
			// Update existing product with fashion fields
			existing.SKU = p.SKU
			existing.Size = p.Size
			existing.Color = p.Color
			existing.Material = p.Material
			existing.ImagesJSON = p.ImagesJSON
			existing.Status = "ACTIVE"
			db.Save(&existing)
			fmt.Println("Product " + p.Name + " updated with fashion attributes.")
		}
	}
}
