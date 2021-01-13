package services

var (
	// ItemsService should be used in the controller in order to use the itemsServiceInterface methods
	ItemsService itemsServiceInterface = &itemsService{}
)

type itemsService struct {
}

type itemsServiceInterface interface {
	GetItem()
	SaveItem()
}

func (s *itemsService) GetItem() {}

func (s *itemsService) SaveItem() {}
