package services

import (
	"imooc-Product/datamodels"
	"imooc-Product/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	// SubNumberOne 消息队列 后增加 更新商品数量方法
	SubNumberOne(productID int64) error
}

type ProductService struct {
	ProductRepository repositories.IProduct
}

// NewProductService 初始化Product服务函数
func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{ProductRepository: repository}
}

// GetProductByID 通过ID查询商品
func (p *ProductService) GetProductByID(productID int64) (product *datamodels.Product, err error) {
	product, err = p.ProductRepository.SelectByKey(productID)
	return
}

// GetAllProduct 查询所有商品
func (p *ProductService) GetAllProduct() (products []*datamodels.Product, err error) {
	products, err = p.ProductRepository.SelectAll()
	return
}

// DeleteProductByID 商品删除
func (p *ProductService) DeleteProductByID(productID int64) (isDelete bool) {
	isDelete = p.ProductRepository.Delete(productID)
	return
}

// InsertProduct 商品添加
func (p *ProductService) InsertProduct(product *datamodels.Product) (productID int64, err error) {
	productID, err = p.ProductRepository.Insert(product)
	return
}

// UpdateProduct 商品更新
func (p *ProductService) UpdateProduct(product *datamodels.Product) (err error) {
	err = p.ProductRepository.Update(product)
	return
}

// SubNumberOne 消息队列控制商品数量更新
func (p *ProductService) SubNumberOne(productID int64) (err error) {
	err = p.ProductRepository.SubProductNum(productID)
	return
}
