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
	"github.com/shopspring/decimal"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Payment is an object representing the database table.
type Payment struct {
	ID             int             `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID         int             `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	PaymentMethod  string          `boil:"payment_method" json:"payment_method" toml:"payment_method" yaml:"payment_method"`
	CardNumber     decimal.Decimal `boil:"card_number" json:"card_number" toml:"card_number" yaml:"card_number"`
	ExpirationDate time.Time       `boil:"expiration_date" json:"expiration_date" toml:"expiration_date" yaml:"expiration_date"`
	CreatedAt      time.Time       `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time       `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *paymentR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L paymentL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PaymentColumns = struct {
	ID             string
	UserID         string
	PaymentMethod  string
	CardNumber     string
	ExpirationDate string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "id",
	UserID:         "user_id",
	PaymentMethod:  "payment_method",
	CardNumber:     "card_number",
	ExpirationDate: "expiration_date",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

var PaymentTableColumns = struct {
	ID             string
	UserID         string
	PaymentMethod  string
	CardNumber     string
	ExpirationDate string
	CreatedAt      string
	UpdatedAt      string
}{
	ID:             "payments.id",
	UserID:         "payments.user_id",
	PaymentMethod:  "payments.payment_method",
	CardNumber:     "payments.card_number",
	ExpirationDate: "payments.expiration_date",
	CreatedAt:      "payments.created_at",
	UpdatedAt:      "payments.updated_at",
}

// Generated where

var PaymentWhere = struct {
	ID             whereHelperint
	UserID         whereHelperint
	PaymentMethod  whereHelperstring
	CardNumber     whereHelperdecimal_Decimal
	ExpirationDate whereHelpertime_Time
	CreatedAt      whereHelpertime_Time
	UpdatedAt      whereHelpertime_Time
}{
	ID:             whereHelperint{field: "\"payments\".\"id\""},
	UserID:         whereHelperint{field: "\"payments\".\"user_id\""},
	PaymentMethod:  whereHelperstring{field: "\"payments\".\"payment_method\""},
	CardNumber:     whereHelperdecimal_Decimal{field: "\"payments\".\"card_number\""},
	ExpirationDate: whereHelpertime_Time{field: "\"payments\".\"expiration_date\""},
	CreatedAt:      whereHelpertime_Time{field: "\"payments\".\"created_at\""},
	UpdatedAt:      whereHelpertime_Time{field: "\"payments\".\"updated_at\""},
}

// PaymentRels is where relationship names are stored.
var PaymentRels = struct {
	User           string
	PaymentDetails string
}{
	User:           "User",
	PaymentDetails: "PaymentDetails",
}

// paymentR is where relationships are stored.
type paymentR struct {
	User           *User              `boil:"User" json:"User" toml:"User" yaml:"User"`
	PaymentDetails PaymentDetailSlice `boil:"PaymentDetails" json:"PaymentDetails" toml:"PaymentDetails" yaml:"PaymentDetails"`
}

// NewStruct creates a new relationship struct
func (*paymentR) NewStruct() *paymentR {
	return &paymentR{}
}

func (r *paymentR) GetUser() *User {
	if r == nil {
		return nil
	}
	return r.User
}

func (r *paymentR) GetPaymentDetails() PaymentDetailSlice {
	if r == nil {
		return nil
	}
	return r.PaymentDetails
}

// paymentL is where Load methods for each relationship are stored.
type paymentL struct{}

var (
	paymentAllColumns            = []string{"id", "user_id", "payment_method", "card_number", "expiration_date", "created_at", "updated_at"}
	paymentColumnsWithoutDefault = []string{"user_id", "payment_method", "card_number", "expiration_date"}
	paymentColumnsWithDefault    = []string{"id", "created_at", "updated_at"}
	paymentPrimaryKeyColumns     = []string{"id"}
	paymentGeneratedColumns      = []string{}
)

type (
	// PaymentSlice is an alias for a slice of pointers to Payment.
	// This should almost always be used instead of []Payment.
	PaymentSlice []*Payment
	// PaymentHook is the signature for custom Payment hook methods
	PaymentHook func(context.Context, boil.ContextExecutor, *Payment) error

	paymentQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	paymentType                 = reflect.TypeOf(&Payment{})
	paymentMapping              = queries.MakeStructMapping(paymentType)
	paymentPrimaryKeyMapping, _ = queries.BindMapping(paymentType, paymentMapping, paymentPrimaryKeyColumns)
	paymentInsertCacheMut       sync.RWMutex
	paymentInsertCache          = make(map[string]insertCache)
	paymentUpdateCacheMut       sync.RWMutex
	paymentUpdateCache          = make(map[string]updateCache)
	paymentUpsertCacheMut       sync.RWMutex
	paymentUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var paymentAfterSelectHooks []PaymentHook

var paymentBeforeInsertHooks []PaymentHook
var paymentAfterInsertHooks []PaymentHook

var paymentBeforeUpdateHooks []PaymentHook
var paymentAfterUpdateHooks []PaymentHook

var paymentBeforeDeleteHooks []PaymentHook
var paymentAfterDeleteHooks []PaymentHook

var paymentBeforeUpsertHooks []PaymentHook
var paymentAfterUpsertHooks []PaymentHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Payment) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Payment) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Payment) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Payment) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Payment) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Payment) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Payment) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Payment) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Payment) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range paymentAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPaymentHook registers your hook function for all future operations.
func AddPaymentHook(hookPoint boil.HookPoint, paymentHook PaymentHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		paymentAfterSelectHooks = append(paymentAfterSelectHooks, paymentHook)
	case boil.BeforeInsertHook:
		paymentBeforeInsertHooks = append(paymentBeforeInsertHooks, paymentHook)
	case boil.AfterInsertHook:
		paymentAfterInsertHooks = append(paymentAfterInsertHooks, paymentHook)
	case boil.BeforeUpdateHook:
		paymentBeforeUpdateHooks = append(paymentBeforeUpdateHooks, paymentHook)
	case boil.AfterUpdateHook:
		paymentAfterUpdateHooks = append(paymentAfterUpdateHooks, paymentHook)
	case boil.BeforeDeleteHook:
		paymentBeforeDeleteHooks = append(paymentBeforeDeleteHooks, paymentHook)
	case boil.AfterDeleteHook:
		paymentAfterDeleteHooks = append(paymentAfterDeleteHooks, paymentHook)
	case boil.BeforeUpsertHook:
		paymentBeforeUpsertHooks = append(paymentBeforeUpsertHooks, paymentHook)
	case boil.AfterUpsertHook:
		paymentAfterUpsertHooks = append(paymentAfterUpsertHooks, paymentHook)
	}
}

// One returns a single payment record from the query.
func (q paymentQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Payment, error) {
	o := &Payment{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for payments")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Payment records from the query.
func (q paymentQuery) All(ctx context.Context, exec boil.ContextExecutor) (PaymentSlice, error) {
	var o []*Payment

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Payment slice")
	}

	if len(paymentAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Payment records in the query.
func (q paymentQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count payments rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q paymentQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if payments exists")
	}

	return count > 0, nil
}

// User pointed to by the foreign key.
func (o *Payment) User(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	return Users(queryMods...)
}

// PaymentDetails retrieves all the payment_detail's PaymentDetails with an executor.
func (o *Payment) PaymentDetails(mods ...qm.QueryMod) paymentDetailQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"payment_details\".\"payment_id\"=?", o.ID),
	)

	return PaymentDetails(queryMods...)
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (paymentL) LoadUser(ctx context.Context, e boil.ContextExecutor, singular bool, maybePayment interface{}, mods queries.Applicator) error {
	var slice []*Payment
	var object *Payment

	if singular {
		var ok bool
		object, ok = maybePayment.(*Payment)
		if !ok {
			object = new(Payment)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePayment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePayment))
			}
		}
	} else {
		s, ok := maybePayment.(*[]*Payment)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePayment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePayment))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &paymentR{}
		}
		args = append(args, object.UserID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &paymentR{}
			}

			for _, a := range args {
				if a == obj.UserID {
					continue Outer
				}
			}

			args = append(args, obj.UserID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`users`),
		qm.WhereIn(`users.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(paymentAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.User = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.Payments = append(foreign.R.Payments, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID == foreign.ID {
				local.R.User = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.Payments = append(foreign.R.Payments, local)
				break
			}
		}
	}

	return nil
}

// LoadPaymentDetails allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (paymentL) LoadPaymentDetails(ctx context.Context, e boil.ContextExecutor, singular bool, maybePayment interface{}, mods queries.Applicator) error {
	var slice []*Payment
	var object *Payment

	if singular {
		var ok bool
		object, ok = maybePayment.(*Payment)
		if !ok {
			object = new(Payment)
			ok = queries.SetFromEmbeddedStruct(&object, &maybePayment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybePayment))
			}
		}
	} else {
		s, ok := maybePayment.(*[]*Payment)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybePayment)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybePayment))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &paymentR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &paymentR{}
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
		qm.From(`payment_details`),
		qm.WhereIn(`payment_details.payment_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load payment_details")
	}

	var resultSlice []*PaymentDetail
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice payment_details")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on payment_details")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for payment_details")
	}

	if len(paymentDetailAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.PaymentDetails = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &paymentDetailR{}
			}
			foreign.R.Payment = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PaymentID {
				local.R.PaymentDetails = append(local.R.PaymentDetails, foreign)
				if foreign.R == nil {
					foreign.R = &paymentDetailR{}
				}
				foreign.R.Payment = local
				break
			}
		}
	}

	return nil
}

// SetUser of the payment to the related item.
// Sets o.R.User to related.
// Adds o to related.R.Payments.
func (o *Payment) SetUser(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"payments\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, paymentPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID = related.ID
	if o.R == nil {
		o.R = &paymentR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			Payments: PaymentSlice{o},
		}
	} else {
		related.R.Payments = append(related.R.Payments, o)
	}

	return nil
}

// AddPaymentDetails adds the given related objects to the existing relationships
// of the payment, optionally inserting them as new records.
// Appends related to o.R.PaymentDetails.
// Sets related.R.Payment appropriately.
func (o *Payment) AddPaymentDetails(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*PaymentDetail) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PaymentID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"payment_details\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"payment_id"}),
				strmangle.WhereClause("\"", "\"", 2, paymentDetailPrimaryKeyColumns),
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

			rel.PaymentID = o.ID
		}
	}

	if o.R == nil {
		o.R = &paymentR{
			PaymentDetails: related,
		}
	} else {
		o.R.PaymentDetails = append(o.R.PaymentDetails, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &paymentDetailR{
				Payment: o,
			}
		} else {
			rel.R.Payment = o
		}
	}
	return nil
}

// Payments retrieves all the records using an executor.
func Payments(mods ...qm.QueryMod) paymentQuery {
	mods = append(mods, qm.From("\"payments\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"payments\".*"})
	}

	return paymentQuery{q}
}

// FindPayment retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPayment(ctx context.Context, exec boil.ContextExecutor, iD int, selectCols ...string) (*Payment, error) {
	paymentObj := &Payment{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"payments\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, paymentObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from payments")
	}

	if err = paymentObj.doAfterSelectHooks(ctx, exec); err != nil {
		return paymentObj, err
	}

	return paymentObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Payment) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no payments provided for insertion")
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

	nzDefaults := queries.NonZeroDefaultSet(paymentColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	paymentInsertCacheMut.RLock()
	cache, cached := paymentInsertCache[key]
	paymentInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			paymentAllColumns,
			paymentColumnsWithDefault,
			paymentColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(paymentType, paymentMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(paymentType, paymentMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"payments\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"payments\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into payments")
	}

	if !cached {
		paymentInsertCacheMut.Lock()
		paymentInsertCache[key] = cache
		paymentInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Payment.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Payment) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	paymentUpdateCacheMut.RLock()
	cache, cached := paymentUpdateCache[key]
	paymentUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			paymentAllColumns,
			paymentPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update payments, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"payments\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, paymentPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(paymentType, paymentMapping, append(wl, paymentPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update payments row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for payments")
	}

	if !cached {
		paymentUpdateCacheMut.Lock()
		paymentUpdateCache[key] = cache
		paymentUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q paymentQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for payments")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for payments")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PaymentSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"payments\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, paymentPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in payment slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all payment")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Payment) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no payments provided for upsert")
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

	nzDefaults := queries.NonZeroDefaultSet(paymentColumnsWithDefault, o)

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

	paymentUpsertCacheMut.RLock()
	cache, cached := paymentUpsertCache[key]
	paymentUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			paymentAllColumns,
			paymentColumnsWithDefault,
			paymentColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			paymentAllColumns,
			paymentPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert payments, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(paymentPrimaryKeyColumns))
			copy(conflict, paymentPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"payments\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(paymentType, paymentMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(paymentType, paymentMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert payments")
	}

	if !cached {
		paymentUpsertCacheMut.Lock()
		paymentUpsertCache[key] = cache
		paymentUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Payment record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Payment) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Payment provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), paymentPrimaryKeyMapping)
	sql := "DELETE FROM \"payments\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from payments")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for payments")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q paymentQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no paymentQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from payments")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for payments")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PaymentSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(paymentBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"payments\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, paymentPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from payment slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for payments")
	}

	if len(paymentAfterDeleteHooks) != 0 {
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
func (o *Payment) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPayment(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PaymentSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PaymentSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"payments\".* FROM \"payments\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, paymentPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PaymentSlice")
	}

	*o = slice

	return nil
}

// PaymentExists checks if the Payment row exists.
func PaymentExists(ctx context.Context, exec boil.ContextExecutor, iD int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"payments\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if payments exists")
	}

	return exists, nil
}
