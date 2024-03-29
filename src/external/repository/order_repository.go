package repository

import (
	"errors"
	"log/slog"

	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap/src/entities/enum"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewOrderRepository(db *gorm.DB, logger *slog.Logger) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger,
	}
}

func (o *OrderRepository) Create(order entity.Order) (entity.Order, error) {
	if result := o.db.Create(&order).Model(&order).Preload("Products").Where("id= ?", order.ID).First(&order); result.Error != nil {
		o.logger.Error("result.Error")
		return order, errors.New("create order from repository has failed")
	}

	return order, nil
}

func (o *OrderRepository) GetAll() ([]entity.Order, error) {
	var orders []entity.Order
	results := o.db.Raw("select * from ze_burguer.orders where not status_order = 'FINISHED' order by case when status_order = 'READY' then 1 when status_order = 'PREPARING' then 2 when status_order = 'RECEIVED' then 3 else 4 end asc, created_at asc").Find(&orders)
	if results.Error != nil {
		o.logger.Error("result.Error")
		return orders, errors.New("find orders from repository has failed")
	}

	return orders, nil
}

func (o *OrderRepository) GetById(id int) (*entity.Order, error) {
	var order entity.Order

	result := o.db.Model(&order).Preload("Products").Where("id= ?", id).First(&order)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			o.logger.Error("order not found", slog.Int("id", id))
			return nil, nil
		}

		o.logger.Error("get order by id (%s) from repository has failed", slog.Int("id", id))
		return nil, errors.New("get order by id from repository has failed")
	}

	return &order, nil
}

func (o *OrderRepository) UpdateStatusById(id int, status enum.StatusOrder) error {
	return o.db.Model(&entity.Order{}).Where("id = ?", id).Update("status_order", status).Error
}
