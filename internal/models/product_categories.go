// Code generated by SQLBoiler 4.13.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ProductCategory is an object representing the database table.
type ProductCategory struct {
	ID          int       `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name        string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	Description string    `boil:"description" json:"description" toml:"description" yaml:"description"`
	CreatedAt   time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *productCategoryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L productCategoryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ProductCategoryColumns = struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}{
	ID:          "id",
	Name:        "name",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

var ProductCategoryTableColumns = struct {
	ID          string
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}{
	ID:          "product_categories.id",
	Name:        "product_categories.name",
	Description: "product_categories.description",
	CreatedAt:   "product_categories.created_at",
	UpdatedAt:   "product_categories.updated_at",
}

// Generated where

var ProductCategoryWhere = struct {
	ID          whereHelperint
	Name        whereHelperstring
	Description whereHelperstring
	CreatedAt   whereHelpertime_Time
	UpdatedAt   whereHelpertime_Time
}{
	ID:          whereHelperint{field: "\"product_categories\".\"id\""},
	Name:        whereHelperstring{field: "\"product_categories\".\"name\""},
	Description: whereHelperstring{field: "\"product_categories\".\"description\""},
	CreatedAt:   whereHelpertime_Time{field: "\"product_categories\".\"created_at\""},
	UpdatedAt:   whereHelpertime_Time{field: "\"product_categories\".\"updated_at\""},
}

// ProductCategoryRels is where relationship names are stored.
var ProductCategoryRels = struct {
	CategoryProducts string
}{
	CategoryProducts: "CategoryProducts",
}

// productCategoryR is where relationships are stored.
type productCategoryR struct {
	CategoryProducts ProductSlice `boil:"CategoryProducts" json:"CategoryProducts" toml:"CategoryProducts" yaml:"CategoryProducts"`
}

// NewStruct creates a new relationship struct
func (*productCategoryR) NewStruct() *productCategoryR {
	return &productCategoryR{}
}

func (r *productCategoryR) GetCategoryProducts() ProductSlice {
	if r == nil {
		return nil
	}
	return r.CategoryProducts
}

// productCategoryL is where Load methods for each relationship are stored.
type productCategoryL struct{}

var (
	productCategoryAllColumns            = []string{"id", "name", "description", "created_at", "updated_at"}
	productCategoryColumnsWithoutDefault = []string{"name", "description"}
	productCategoryColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	productCategoryPrimaryKeyColumns     = []string{"id"}
	productCategoryGeneratedColumns      = []string{}
)

type (
	// ProductCategorySlice is an alias for a slice of pointers to ProductCategory.
	// This should almost always be used instead of []ProductCategory.
	ProductCategorySlice []*ProductCategory
	// ProductCategoryHook is the signature for custom ProductCategory hook methods
	ProductCategoryHook func(context.Context, boil.ContextExecutor, *ProductCategory) error

	productCategoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	productCategoryType                 = reflect.TypeOf(&ProductCategory{})
	productCategoryMapping              = queries.MakeStructMapping(productCategoryType)
	productCategoryPrimaryKeyMapping, _ = queries.BindMapping(productCategoryType, productCategoryMapping, productCategoryPrimaryKeyColumns)
	productCategoryInsertCacheMut       sync.RWMutex
	productCategoryInsertCache          = make(map[string]insertCache)
	productCategoryUpdateCacheMut       sync.RWMutex
	productCategoryUpdateCache          = make(map[string]updateCache)
	productCategoryUpsertCacheMut       sync.RWMutex
	productCategoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var productCategoryAfterSelectHooks []ProductCategoryHook

var productCategoryBeforeInsertHooks []ProductCategoryHook
var productCategoryAfterInsertHooks []ProductCategoryHook

var productCategoryBeforeUpdateHooks []ProductCategoryHook
var productCategoryAfterUpdateHooks []ProductCategoryHook

var productCategoryBeforeDeleteHooks []ProductCategoryHook
var productCategoryAfterDeleteHooks []ProductCategoryHook

var productCategoryBeforeUpsertHooks []ProductCategoryHook
var productCategoryAfterUpsertHooks []ProductCategoryHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ProductCategory) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ProductCategory) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ProductCategory) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ProductCategory) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ProductCategory) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ProductCategory) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ProductCategory) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ProductCategory) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ProductCategory) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range productCategoryAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddProductCategoryHook registers your hook function for all future operations.
func AddProductCategoryHook(hookPoint boil.HookPoint, productCategoryHook ProductCategoryHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		productCategoryAfterSelectHooks = append(productCategoryAfterSelectHooks, productCategoryHook)
	case boil.BeforeInsertHook:
		productCategoryBeforeInsertHooks = append(productCategoryBeforeInsertHooks, productCategoryHook)
	case boil.AfterInsertHook:
		productCategoryAfterInsertHooks = append(productCategoryAfterInsertHooks, productCategoryHook)
	case boil.BeforeUpdateHook:
		productCategoryBeforeUpdateHooks = append(productCategoryBeforeUpdateHooks, productCategoryHook)
	case boil.AfterUpdateHook:
		productCategoryAfterUpdateHooks = append(productCategoryAfterUpdateHooks, productCategoryHook)
	case boil.BeforeDeleteHook:
		productCategoryBeforeDeleteHooks = append(productCategoryBeforeDeleteHooks, productCategoryHook)
	case boil.AfterDeleteHook:
		productCategoryAfterDeleteHooks = append(productCategoryAfterDeleteHooks, productCategoryHook)
	case boil.BeforeUpsertHook:
		productCategoryBeforeUpsertHooks = append(productCategoryBeforeUpsertHooks, productCategoryHook)
	case boil.AfterUpsertHook:
		productCategoryAfterUpsertHooks = append(productCategoryAfterUpsertHooks, productCategoryHook)
	}
}

// One returns a single productCategory record from the query.
func (q productCategoryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ProductCategory, error) {
	o := &ProductCategory{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for product_categories")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ProductCategory records from the query.
func (q productCategoryQuery) All(ctx context.Context, exec boil.ContextExecutor) (ProductCategorySlice, error) {
	var o []*ProductCategory

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ProductCategory slice")
	}

	if len(productCategoryAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ProductCategory records in the query.
func (q productCategoryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count product_categories rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q productCategoryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if product_categories exists")
	}

	return count > 0, nil
}

// CategoryProducts retrieves all the product's Products with an executor via category_id column.
func (o *ProductCategory) CategoryProducts(mods ...qm.QueryMod) productQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"products\".\"category_id\"=?", o.ID),
	)

	return Products(queryMods...)
}

// LoadCategoryProducts allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (productCategoryL) LoadCategoryProducts(ctx context.Context, e boil.ContextExecutor, singular bool, maybeProductCategory interface{}, mods queries.Applicator) error {
	var slice []*ProductCategory
	var object *ProductCategory

	if singular {
		var ok bool
		object, ok = maybeProductCategory.(*ProductCategory)
		if !ok {
			object = new(ProductCategory)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeProductCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeProductCategory))
			}
		}
	} else {
		s, ok := maybeProductCategory.(*[]*ProductCategory)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeProductCategory)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeProductCategory))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &productCategoryR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &productCategoryR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`products`),
		qm.WhereIn(`products.category_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load products")
	}

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice products")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on products")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for products")
	}

	if len(productAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.CategoryProducts = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &productR{}
			}
			foreign.R.Category = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.CategoryID {
				local.R.CategoryProducts = append(local.R.CategoryProducts, foreign)
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.Category = local
				break
			}
		}
	}

	return nil
}

// AddCategoryProducts adds the given related objects to the existing relationships
// of the product_category, optionally inserting them as new records.
// Appends related to o.R.CategoryProducts.
// Sets related.R.Category appropriately.
func (o *ProductCategory) AddCategoryProducts(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Product) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.CategoryID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"products\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"category_id"}),
				strmangle.WhereClause("\"", "\"", 2, productPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.CategoryID = o.ID
		}
	}

	if o.R == nil {
		o.R = &productCategoryR{
			CategoryProducts: related,
		}
	} else {
		o.R.CategoryProducts = append(o.R.CategoryProducts, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &productR{
				Category: o,
			}
		} else {
			rel.R.Category = o
		}
	}
	return nil
}

// ProductCategories retrieves all the records using an executor.
func ProductCategories(mods ...qm.QueryMod) productCategoryQuery {
	mods = append(mods, qm.From("\"product_categories\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"product_categories\".*"})
	}

	return productCategoryQuery{q}
}

// FindProductCategory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindProductCategory(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*ProductCategory, error) {
	productCategoryObj := &ProductCategory{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"product_categories\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, productCategoryObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from product_categories")
	}

	if err = productCategoryObj.doAfterSelectHooks(ctx, exec); err != nil {
		return productCategoryObj, err
	}

	return productCategoryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ProductCategory) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no product_categories provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(productCategoryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	productCategoryInsertCacheMut.RLock()
	cache, cached := productCategoryInsertCache[key]
	productCategoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			productCategoryAllColumns,
			productCategoryColumnsWithDefault,
			productCategoryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(productCategoryType, productCategoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(productCategoryType, productCategoryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"product_categories\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"product_categories\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into product_categories")
	}

	if !cached {
		productCategoryInsertCacheMut.Lock()
		productCategoryInsertCache[key] = cache
		productCategoryInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ProductCategory.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ProductCategory) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	productCategoryUpdateCacheMut.RLock()
	cache, cached := productCategoryUpdateCache[key]
	productCategoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			productCategoryAllColumns,
			productCategoryPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update product_categories, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"product_categories\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, productCategoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(productCategoryType, productCategoryMapping, append(wl, productCategoryPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update product_categories row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for product_categories")
	}

	if !cached {
		productCategoryUpdateCacheMut.Lock()
		productCategoryUpdateCache[key] = cache
		productCategoryUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q productCategoryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for product_categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for product_categories")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ProductCategorySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), productCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"product_categories\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, productCategoryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in productCategory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all productCategory")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ProductCategory) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no product_categories provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(productCategoryColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	productCategoryUpsertCacheMut.RLock()
	cache, cached := productCategoryUpsertCache[key]
	productCategoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			productCategoryAllColumns,
			productCategoryColumnsWithDefault,
			productCategoryColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			productCategoryAllColumns,
			productCategoryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert product_categories, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(productCategoryPrimaryKeyColumns))
			copy(conflict, productCategoryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"product_categories\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(productCategoryType, productCategoryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(productCategoryType, productCategoryMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert product_categories")
	}

	if !cached {
		productCategoryUpsertCacheMut.Lock()
		productCategoryUpsertCache[key] = cache
		productCategoryUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single ProductCategory record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ProductCategory) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ProductCategory provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), productCategoryPrimaryKeyMapping)
	sql := "DELETE FROM \"product_categories\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from product_categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for product_categories")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q productCategoryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no productCategoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from product_categories")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for product_categories")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ProductCategorySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(productCategoryBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), productCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"product_categories\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, productCategoryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from productCategory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for product_categories")
	}

	if len(productCategoryAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ProductCategory) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindProductCategory(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ProductCategorySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ProductCategorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), productCategoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"product_categories\".* FROM \"product_categories\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, productCategoryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ProductCategorySlice")
	}

	*o = slice

	return nil
}

// ProductCategoryExists checks if the ProductCategory row exists.
func ProductCategoryExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"product_categories\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if product_categories exists")
	}

	return exists, nil
}