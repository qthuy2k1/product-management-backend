package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Product struct {
	ID          int             `redis:"id"`
	Name        string          `redis:"name"`
	Description string          `redis:"description"`
	Price       decimal.Decimal `redis:"price"`
	Quantity    int             `redis:"quantity"`
	AuthorID    int             `redis:"author_id"`
	CategoryID  int             `redis:"category_id"`
	CreatedAt   time.Time       `redis:"created_at"`
	UpdatedAt   time.Time       `redis:"updated_at"`
}

// CreateProduct creates a product in db given by product model in parameter
func (r *Repository) CreateProduct(ctx context.Context, pReq Product) error {
	product := models.Product{
		Name:        pReq.Name,
		Description: pReq.Description,
		Price:       pReq.Price,
		Quantity:    pReq.Quantity,
		AuthorID:    pReq.AuthorID,
		CategoryID:  pReq.CategoryID,
	}
	if err := product.Insert(ctx, boil.GetContextDB(), boil.Infer()); err != nil {
		return pkgerrors.WithStack(err)
	}
	return nil
}

// GetProduct retrieves a product in db by ID
func (r *Repository) GetProduct(ctx context.Context, id int) (models.Product, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("product:%d", id))
	if len(res.Val()) == 0 {
		product, err := models.FindProduct(ctx, boil.GetContextDB(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.Product{}, ErrProductNotFound
			}
			return models.Product{}, err
		}

		// set cache
		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("product:%d", id), map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price.String(),
			"quantity":    product.Quantity,
			"author_id":   product.AuthorID,
			"category_id": product.CategoryID,
			"created_at":  product.CreatedAt,
			"updated_at":  product.UpdatedAt,
		}); errCache.Err() != nil {
			return models.Product{}, errCache.Err()
		}

		return *product, nil
	}

	var productScan Product
	if err := res.Scan(&productScan); err != nil {
		return models.Product{}, err
	}

	return models.Product{
		ID:          productScan.ID,
		Name:        productScan.Name,
		Description: productScan.Description,
		Price:       productScan.Price,
		Quantity:    productScan.Quantity,
		AuthorID:    productScan.AuthorID,
		CategoryID:  productScan.CategoryID,
		CreatedAt:   productScan.CreatedAt,
		UpdatedAt:   productScan.UpdatedAt,
	}, nil
}

// UpdateProduct updates a product in db given by product model in parameter
func (r *Repository) UpdateProduct(ctx context.Context, tx *sql.Tx, pReq models.Product) error {
	ctxExec := boil.GetContextDB()
	if tx != nil {
		ctxExec = tx
	}

	if _, err := pReq.Update(ctx, ctxExec, boil.Blacklist("id", "created_at")); err != nil {
		return err
	}

	if err := r.Redis.HSet(ctx, fmt.Sprintf("product:%d", pReq.ID), map[string]interface{}{
		"id":          pReq.ID,
		"name":        pReq.Name,
		"description": pReq.Description,
		"price":       pReq.Price.String(),
		"quantity":    pReq.Quantity,
		"author_id":   pReq.AuthorID,
		"category_id": pReq.CategoryID,
		"created_at":  pReq.CreatedAt,
		"updated_at":  pReq.UpdatedAt,
	}); err != nil {
		return err.Err()
	}

	return nil
}

// DeleteProduct deletes a product in db by ID
func (r *Repository) DeleteProduct(ctx context.Context, id int) error {
	product, err := r.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	pModel := models.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		AuthorID:    product.AuthorID,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	if _, err = pModel.Delete(ctx, boil.GetContextDB()); err != nil {
		return err
	}
	return nil
}

type ProductRepoFilter struct {
	Name string
	Date string
}

type ProductOutput struct {
	ID           int             `boil:"id"`
	Name         string          `boil:"name"`
	Description  string          `boil:"description"`
	Price        decimal.Decimal `boil:"price"`
	Quantity     int             `boil:"quantity"`
	AuthorID     int             `boil:"author_id"`
	CategoryName string          `boil:"category"`
	CreatedAt    time.Time       `boil:"created_at"`
	UpdatedAt    time.Time       `boil:"updated_at"`
}

// GetProducts retrieves all the products in db
func (r *Repository) GetProducts(ctx context.Context, filter ProductRepoFilter) ([]ProductOutput, error) {
	var queryMod []qm.QueryMod

	queryMod = append(queryMod, qm.Select(fmt.Sprintf("%s.%s as id, %s.%s as name, %s.%s as description, %s, %s, %s, %s.%s as category, %s.%s as created_at, %s.%s as updated_at", models.TableNames.Products, models.ProductColumns.ID, models.TableNames.Products, models.ProductColumns.Name, models.TableNames.Products, models.ProductColumns.Description, models.ProductColumns.Price, models.ProductColumns.Quantity, models.ProductColumns.AuthorID, models.TableNames.ProductCategories, models.ProductCategoryColumns.Name, models.TableNames.Products, models.ProductColumns.CreatedAt, models.TableNames.Products, models.ProductColumns.UpdatedAt)))

	if filter.Name != "" {
		queryMod = append(queryMod, qm.Where(fmt.Sprintf("UPPER(%s.%s) LIKE UPPER(?)", models.TableNames.Products, models.ProductColumns.Name), "%"+filter.Name+"%"))
	}

	if filter.Date != "" {
		queryMod = append(queryMod, qm.Where(fmt.Sprintf("DATE(%s.%s) = ?", models.TableNames.Products, models.ProductColumns.CreatedAt), filter.Date))
	}

	queryMod = append(queryMod, qm.InnerJoin(fmt.Sprintf("%s on %s.%s = %s.%s", models.TableNames.ProductCategories, models.TableNames.Products, models.ProductColumns.CategoryID, models.TableNames.ProductCategories, models.ProductColumns.ID)))

	var products []ProductOutput
	if err := models.Products(queryMod...).Bind(ctx, boil.GetContextDB(), &products); err != nil {
		return nil, err
	}

	return products, nil
}

type GetProductsGraph struct {
	Product models.Product         `boil:"products,bind"`
	User    models.User            `boil:"users,bind"`
	PCate   models.ProductCategory `boil:"product_categories,bind"`
}

// GetProductsGraph retrieves all products in db and is used in GraphQL.
func (r *Repository) GetProductsGraph(ctx context.Context, filter ProductRepoFilter) ([]GetProductsGraph, error) {
	var queryMod []qm.QueryMod

	productsTable := models.TableNames.Products
	usersTable := models.TableNames.Users
	productCategoriesTable := models.TableNames.ProductCategories

	// select all product with inner join
	queryMod = append(queryMod,
		qm.Select(
			productsTable+".id",
			productsTable+".name",
			productsTable+".description",
			productsTable+".price",
			productsTable+".quantity",
			productsTable+".created_at",
			productsTable+".updated_at",
			usersTable+".id",
			usersTable+".name",
			usersTable+".email",
			usersTable+".role",
			usersTable+".status",
			usersTable+".created_at",
			usersTable+".updated_at",
			productCategoriesTable+".id",
			productCategoriesTable+".name",
			productCategoriesTable+".description",
			productCategoriesTable+".created_at",
			productCategoriesTable+".updated_at",
		),
	)

	// filter query
	if filter.Name != "" {
		queryMod = append(queryMod, qm.Where(fmt.Sprintf("UPPER(%s.%s) LIKE UPPER(?)", models.TableNames.Products, models.ProductColumns.Name), "%"+filter.Name+"%"))
	}

	// filter query
	if filter.Date != "" {
		queryMod = append(queryMod, qm.Where(fmt.Sprintf("DATE(%s.%s) = ?", models.TableNames.Products, models.ProductColumns.CreatedAt), filter.Date))
	}

	// inner join product categories table
	queryMod = append(queryMod, qm.InnerJoin(fmt.Sprintf("%s on %s.%s = %s.%s", models.TableNames.ProductCategories, models.TableNames.Products, models.ProductColumns.CategoryID, models.TableNames.ProductCategories, models.ProductCategoryColumns.ID)))

	// inner join users table
	queryMod = append(queryMod, qm.InnerJoin(fmt.Sprintf("%s on %s.%s = %s.%s", models.TableNames.Users, models.TableNames.Products, models.ProductColumns.AuthorID, models.TableNames.Users, models.UserColumns.ID)))

	var products []GetProductsGraph
	if err := models.Products(queryMod...).Bind(ctx, boil.GetContextDB(), &products); err != nil {
		return nil, err
	}

	return products, nil
}

// UpsertProducts updates list of products. If a product already exists in db, updates only changed values instead
func (r *Repository) UpsertProducts(ctx context.Context, products []Product) error {
	// prepare statement
	query := []string{`INSERT INTO products(name, description, price, quantity, category_id, author_id) VALUES `}
	for _, p := range products {
		query = append(query, fmt.Sprintf(`('%s', '%s', %v, %d, %d, %d),`, p.Name, p.Description, p.Price, p.Quantity, p.CategoryID, p.AuthorID))
	}

	// remove the "," character at the end of the last query
	query[len(query)-1] = strings.TrimSuffix(query[len(query)-1], ",")

	// if conflict then update
	query = append(query, ` ON CONFLICT(name) DO UPDATE SET description=EXCLUDED.description, price=EXCLUDED.price, quantity=EXCLUDED.quantity, category_id=EXCLUDED.category_id, author_id=EXCLUDED.author_id `)

	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, strings.Join(query, " "))
	if err != nil {
		return err
	}

	fmt.Println("start query", time.Now())
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
