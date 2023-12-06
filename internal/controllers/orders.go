package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/qthuy2k1/product-management/internal/utils/email"

	"github.com/qthuy2k1/product-management/internal/repositories"
	"github.com/shopspring/decimal"
	"github.com/signintech/gopdf"
)

type OrderInput struct {
	UserID    int
	Status    string
	Total     decimal.Decimal
	OrderItem []OrderItemInput
}

// CreateOrder creates an order in db given by order model in parameter
func (c *Controller) CreateOrder(ctx context.Context, orderInput OrderInput, orderItemsInput []OrderItemInput) error {
	// check user exists
	user, err := c.Repository.GetUser(ctx, orderInput.UserID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// start a transaction
	tx, err := c.Repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer c.Repository.RollbackTx(tx)

	order, err := c.Repository.CreateOrder(ctx, tx, repositories.Order{
		UserID: orderInput.UserID,
		Status: orderInput.Status,
	})
	if err != nil {
		return err
	}

	order.TotalPrice = decimal.NewNullDecimal(decimal.NewFromFloat(0))
	var oiRepoInputList []repositories.OrderItem
	for _, oi := range orderItemsInput {
		// check product exists
		p, err := c.Repository.GetProduct(ctx, oi.ProductID)
		if err != nil {
			return err
		}

		// check the quantity of product
		if p.Quantity < oi.Quantity {
			return ErrInsufficientQuantity
		}

		oiRepoInputList = append(oiRepoInputList, repositories.OrderItem{
			ProductID: oi.ProductID,
			Quantity:  oi.Quantity,
			Price:     p.Price.Mul(decimal.NewFromInt(int64(p.Quantity))),
		})

		// decrease the quantity of product
		p.Quantity -= oi.Quantity
		if err = c.Repository.UpdateProduct(ctx, tx, p); err != nil {
			return err
		}

		// update order price total
		// TotalPrice + (p.Price * oi.Quantity)
		order.TotalPrice.Decimal = order.TotalPrice.Decimal.Add(p.Price.Mul(decimal.NewFromInt(int64(oi.Quantity))))
	}

	if err = c.Repository.CreateOrderItem(ctx, tx, oiRepoInputList, order); err != nil {
		return err
	}

	// update order price total
	if err = c.Repository.UpdateOrder(ctx, tx, order); err != nil {
		return err
	}

	if err = c.SendEmailOrder(ctx, user.Email, order, oiRepoInputList); err != nil {
		return err
	}

	return tx.Commit()
}

// SendEmailOrder sends an email that contains the order detail to the user
func (c *Controller) SendEmailOrder(ctx context.Context, emailTo string, order models.Order, orderItem []repositories.OrderItem) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	if err := pdf.AddTTFFont("font", "data/fonts/Roboto/Roboto-Regular.ttf"); err != nil {
		return err
	}

	if err := pdf.SetFont("font", "", 18); err != nil {
		log.Print(err.Error())
		return err
	}

	pdf.AddHeader(func() {
		pdf.SetY(5)
		pdf.Cell(nil, "YOUR ORDER")
	})
	pdf.AddPage()

	if err := pdf.SetFont("font", "", 12); err != nil {
		log.Print(err.Error())
		return err
	}
	pdf.Br(30)
	pdf.Text(fmt.Sprintf("Order #%d", order.ID))
	pdf.Br(20)
	pdf.Text(fmt.Sprintf("Date: %s", order.CreatedAt.Format(time.DateTime)))
	pdf.Br(20)
	pdf.Text(fmt.Sprintf("Order Status: %s", order.Status))
	pdf.Br(20)

	user, err := c.Repository.GetUser(ctx, order.UserID)
	if err != nil {
		return err
	}
	pdf.Text(fmt.Sprintf("Name: %s", user.Name))
	pdf.Br(20)
	pdf.Text(fmt.Sprintf("Email: %s", user.Email))
	pdf.Br(20)

	pdf.Line(pdf.GetX(), pdf.GetY(), gopdf.PageSizeA4.W-10, pdf.GetY())
	pdf.Br(20)
	pdf.Text("Product List")
	pdf.Br(20)

	for _, oi := range orderItem {
		p, err := c.Repository.GetProduct(ctx, oi.ProductID)
		if err != nil {
			return err
		}

		pdf.Cell(nil, fmt.Sprintf("Product Name: %s, Quantity: %d, Price: %v", p.Name, oi.Quantity, oi.Price.Mul(decimal.NewFromInt(int64(oi.Quantity)))))
		pdf.Br(20)
	}
	pdf.Br(20)
	pdf.Text(fmt.Sprintf("Total Price: %s", order.TotalPrice.Decimal.String()))

	pdf.WritePdf("test.pdf")

	reader, err := os.Open("test.pdf")
	if err != nil {
		return err
	}
	defer os.Remove("test.pdf")

	sender := email.NewEmailSender()

	m := email.NewMessage("Order Created", "Hi users,\nThis is your order,\nThanks for choosing us!")
	m.To = []string{emailTo}
	m.Attachments = reader
	m.AttachmentName = "order.pdf"

	if err := sender.Send(m); err != nil {
		return err
	}

	return nil
}

// UpdateOrder updates an order in db given by order model in parameter
func (c *Controller) UpdateOrder(ctx context.Context, orderID int, orderInput OrderInput) error {
	order, err := c.Repository.GetOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, repositories.ErrOrderNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	// check user exists
	if _, err := c.Repository.GetUser(ctx, orderInput.UserID); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	tx, err := c.Repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer c.Repository.RollbackTx(tx)

	if orderInput.OrderItem != nil {
		// init new total price
		order.TotalPrice = decimal.NewNullDecimal(decimal.NewFromFloat(0))
		for _, oi := range orderInput.OrderItem {
			// check product exists
			p, err := c.Repository.GetProduct(ctx, oi.ProductID)
			if err != nil {
				if errors.Is(err, repositories.ErrProductNotFound) {
					return ErrProductNotFound
				}
				return err
			}

			// check the product quantity
			if p.Quantity < oi.Quantity {
				return ErrInsufficientQuantity
			}

			if _, err = c.Repository.GetOrderItem(ctx, oi.ID); err != nil {
				if errors.Is(err, repositories.ErrOrderItemNotFound) {
					return ErrOrderItemNotFound
				}
			}

			if err = c.Repository.UpdateOrderItem(ctx, tx, oi.ID, repositories.OrderItem{
				OrderID:   order.ID,
				ProductID: oi.ProductID,
				Quantity:  oi.Quantity,
				Price:     p.Price,
			}); err != nil {
				return err
			}

			// adjust the quantity of product
			p.Quantity -= oi.Quantity
			if err = c.Repository.UpdateProduct(ctx, tx, p); err != nil {
				return err
			}

			// update order price total
			order.TotalPrice.Decimal = order.TotalPrice.Decimal.Add(p.Price.Mul(decimal.NewFromInt(int64(oi.Quantity))))
		}
	}

	order.UserID = orderInput.UserID
	order.Status = orderInput.Status

	if err = c.Repository.UpdateOrder(ctx, tx, order); err != nil {
		return err
	}

	return tx.Commit()
}

type OrderOutputGraph struct {
	ID         int
	UserName   string
	UserEmail  string
	Status     string
	TotalPrice decimal.Decimal
	CreatedAt  time.Time
	Items      []OrderItemOutput
}

type OrderFilterCtrl struct {
	Sorting    []Sorting
	Pagination Pagination
	FilterDate FilterDate
}

type Pagination struct {
	Limit int
	Page  int
}

type Sorting struct {
	ColumnName string
	Desc       bool
}

type FilterDate struct {
	StartDate string
	EndDate   string
}

// GetOrders retrieves all the orders in db
func (c *Controller) GetOrders(ctx context.Context, filter OrderFilterCtrl) ([]OrderOutputGraph, int64, error) {
	filterOrderRepo := repositories.OrderFilterRepo{
		Pagination: repositories.Pagination{
			Limit: filter.Pagination.Limit,
			Page:  filter.Pagination.Page,
		},
		FilterDate: repositories.FilterDate{
			StartDate: filter.FilterDate.StartDate,
			EndDate:   filter.FilterDate.EndDate,
		},
	}

	for _, o := range filter.Sorting {
		var sortOrder string
		if o.Desc {
			sortOrder = "desc"
		} else {
			sortOrder = "asc"
		}
		filterOrderRepo.Sorting = append(filterOrderRepo.Sorting, repositories.Sorting{
			ColumnName: o.ColumnName,
			SortOrder:  sortOrder,
		})
	}

	orders, count, err := c.Repository.GetOrders(ctx, filterOrderRepo)
	if err != nil {
		return nil, 0, err
	}

	var ordersOutput []OrderOutputGraph
	for _, o := range orders {
		order := OrderOutputGraph{
			ID:         o.ID,
			UserName:   o.UserName,
			UserEmail:  o.UserEmail,
			Status:     o.Status,
			TotalPrice: o.TotalPrice,
			CreatedAt:  o.CreatedAt,
		}

		// remove the {} and split to a slice
		itemsID := strings.Split(strings.Trim(o.ItemID, `{}`), ",")
		productsName := strings.Split(o.ProductName, ",")
		itemsQuantity := strings.Split(strings.Trim(o.Quantity, "{}"), ",")
		itemsPrice := strings.Split(strings.Trim(o.ItemPrice, "{}"), ",")

		for index := range itemsID {
			// convert string to specific values
			id, err := strconv.Atoi(itemsID[index])
			if err != nil {
				return nil, 0, err
			}
			quantity, err := strconv.Atoi(itemsQuantity[index])
			if err != nil {
				return nil, 0, err
			}
			price, err := strconv.ParseFloat(itemsPrice[index], 64)
			if err != nil {
				return nil, 0, err
			}

			order.Items = append(order.Items, OrderItemOutput{
				ID:          id,
				ProductName: strings.Trim(productsName[index], `{}"`),
				Quantity:    quantity,
				Price:       decimal.NewFromFloat(price),
			})
		}

		ordersOutput = append(ordersOutput, order)
	}

	return ordersOutput, count, nil
}
