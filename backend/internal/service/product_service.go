package service

import (
	"errors"
	"strings"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(onlyActive bool) ([]models.Product, error) {
	return s.repo.FindAll(onlyActive)
}

func (s *ProductService) GetByID(id uint) (*models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Create(p *models.Product) error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return errors.New("nama produk diperlukan")
	}
	if p.Price <= 0 {
		return errors.New("harga produk mestilah lebih daripada 0")
	}
	if p.Weight <= 0 {
		p.Weight = 0.5 // Default 500g
	}
	return s.repo.Create(p)
}

func (s *ProductService) Update(id uint, input *models.Product) (*models.Product, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak dijumpai")
	}

	if input.Name != "" {
		existing.Name = strings.TrimSpace(input.Name)
	}
	if input.SKU != "" {
		existing.SKU = input.SKU
	}
	if input.Description != "" {
		existing.Description = input.Description
	}
	if input.Price > 0 {
		existing.Price = input.Price
	}
	if input.Stock >= 0 {
		existing.Stock = input.Stock
	}
	if input.Weight > 0 {
		existing.Weight = input.Weight
	}
	if input.ImageURL != "" {
		existing.ImageURL = input.ImageURL
	}
	existing.IsActive = input.IsActive

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}
