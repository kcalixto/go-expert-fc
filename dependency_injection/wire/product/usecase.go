package product

type ProductUseCase struct {
	repository ProductRepositoryInterface
}

func NewProductUseCase(repository ProductRepositoryInterface) *ProductUseCase {
	return &ProductUseCase{repository}
}

func (uc *ProductUseCase) GetProductByID(id int) (*Product, error) {
	return uc.repository.GetProductByID(id)
}
