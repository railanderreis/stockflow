package customer

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Modelo de Domínio e Contrato do Repositório
type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Document  string    `json:"document"` // CPF ou CNPJ
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   Address   `json:"address"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Address struct {
	Street  string `json:"street"`
	Number  string `json:"number"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}

type CustomerRepository interface {
	Save(ctx context.Context, customer *Customer) error
	FindByID(ctx context.Context, id string) (*Customer, error)
	FindByDocument(ctx context.Context, document string) (*Customer, error)
	List(ctx context.Context, search string) ([]*Customer, error)
}

// DTOs (Data Transfer Objects)
type CreateCustomerInput struct {
	Name     string  `json:"name"`
	Document string  `json:"document"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Address  Address `json:"address"`
}

type CustomerOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Document  string    `json:"document"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   Address   `json:"address"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

// Caso de Uso
type CustomerUseCase struct {
	customerRepo CustomerRepository
}

func NewCustomerUseCase(customerRepo CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{
		customerRepo: customerRepo,
	}
}

// Criar novo Cliente
func (uc *CustomerUseCase) Create(ctx context.Context, input CreateCustomerInput) (*CustomerOutput, error) {
	if input.Name == "" || input.Document == "" {
		return nil, errors.New("nome e documento (CPF/CNPJ) são obrigatórios")
	}

	// Regra de negócio: não permite documentos duplicados
	existing, _ := uc.customerRepo.FindByDocument(ctx, input.Document)
	if existing != nil {
		return nil, errors.New("já existe um cliente cadastrado com este documento")
	}

	now := time.Now()
	customer := &Customer{
		ID:        fmt.Sprintf("cli_%d", now.UnixNano()),
		Name:      input.Name,
		Document:  input.Document,
		Email:     input.Email,
		Phone:     input.Phone,
		Address:   input.Address,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.customerRepo.Save(ctx, customer); err != nil {
		return nil, fmt.Errorf("erro ao salvar cliente: %w", err)
	}

	return toOutput(customer), nil
}

// Listar Clientes
func (uc *CustomerUseCase) List(ctx context.Context, search string) ([]*CustomerOutput, error) {
	customers, err := uc.customerRepo.List(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar clientes: %w", err)
	}

	var outputs []*CustomerOutput
	for _, c := range customers {
		outputs = append(outputs, toOutput(c))
	}

	return outputs, nil
}

// Obter Cliente por ID
func (uc *CustomerUseCase) GetByID(ctx context.Context, id string) (*CustomerOutput, error) {
	customer, err := uc.customerRepo.FindByID(ctx, id)
	if err != nil || customer == nil {
		return nil, errors.New("cliente não encontrado")
	}

	return toOutput(customer), nil
}

// Mapeador auxiliar interno
func toOutput(c *Customer) *CustomerOutput {
	return &CustomerOutput{
		ID:        c.ID,
		Name:      c.Name,
		Document:  c.Document,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		Active:    c.Active,
		CreatedAt: c.CreatedAt,
	}
}
