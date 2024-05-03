package usecase

import (
	"github.com/devfullcycle/20-CleanArch/internal/entity"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	repo entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: repo,
	}
}

func (u *ListOrdersUseCase) Execute() ([]OrderOutputDTO, error) {

	orders, err := u.OrderRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var ordersDTO []OrderOutputDTO

	for _, order := range orders {
		ordersDTO = append(ordersDTO, OrderOutputDTO{
			order.ID,
			order.Price,
			order.Tax,
			order.FinalPrice,
		})
	}

	return ordersDTO, nil
}
