package productsrepositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/entities"
	filesusecases "github.com/MrXMMM/E-commerce-Project/modules/files/filesUsecases"
	"github.com/MrXMMM/E-commerce-Project/modules/products"
	productspatterns "github.com/MrXMMM/E-commerce-Project/modules/products/productsPatterns"
	"github.com/jmoiron/sqlx"
)

type IProductsRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) ([]*products.Product, int)
	InsertProduct(req *products.Product) (*products.Product, error)
	UpdateProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
}

type productsRepository struct {
	db           *sqlx.DB
	cfg          config.IConfig
	filesUsecase filesusecases.IFilesUsecase
}

func ProductsRepository(db *sqlx.DB, cfg config.IConfig, filesUsecase filesusecases.IFilesUsecase) IProductsRepository {
	return &productsRepository{
		db:           db,
		cfg:          cfg,
		filesUsecase: filesUsecase,
	}
}

func (r *productsRepository) FindOneProduct(productId string) (*products.Product, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			(
				SELECT
					to_jsonb("ct")
				FROM (
					select
						"c"."id",
						"c"."title"
					FROM "categories" "c"
					LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
					WHERE "pc"."product_id" = "p"."id"
				) AS "ct"
			) AS "category",
			"p"."created_at",
			"p"."updated_at",
			"p"."price",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (	
					SELECT 
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "images" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM "products" "p"
		WHERE "p"."id" = $1
		LIMIT 1
	) AS "t";`

	productBytes := make([]byte, 0)
	product := &products.Product{
		Images: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productBytes, query, productId); err != nil {
		return nil, fmt.Errorf("get products failed: %v", err)
	}
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}
	return product, nil
}

func (r *productsRepository) FindProduct(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productspatterns.FindProductBuilder(r.db, req)
	engineer := productspatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()
	engineer.CountProduct().PrintQuery()

	return result, count
}

func (r *productsRepository) InsertProduct(req *products.Product) (*products.Product, error) {
	builder := productspatterns.InsertProductBuilder(r.db, req)
	engineer := productspatterns.InsertProductEngineer(builder)

	productId, err := engineer.InsertProduct()
	if err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productsRepository) UpdateProduct(req *products.Product) (*products.Product, error) {
	builder := productspatterns.UpdateProductBuilder(r.db, req, r.filesUsecase)
	engineer := productspatterns.UpdateProductEngineer(builder)

	if err := engineer.UpdateProduct(); err != nil {
		return nil, err
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productsRepository) DeleteProduct(productId string) error {

	query := `DELETE FROM "products" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(context.Background(), query, productId); err != nil {
		return fmt.Errorf("delete product failed: %v", err)
	}

	return nil
}
