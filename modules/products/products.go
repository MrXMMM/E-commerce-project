package products

import (
	"github.com/MrXMMM/E-commerce-Project/modules/appinfo"
	"github.com/MrXMMM/E-commerce-Project/modules/entities"
)

type Product struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    *appinfo.Category `json:"category"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Price       float64           `json:"price"`
	Images      []*entities.Image `json:"images"`
}

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `json:"search"` // title & description
	*entities.PaginationReq
	*entities.SortReq
}
