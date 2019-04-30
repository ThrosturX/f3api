package f3api

import (
	"errors"
	"fmt"
	"sync"
)

type NotFoundError error

// Interface for stable storage
type ApiStore interface {
	// Add a payment to the stable storage
	// Precondition: The payment must not exist
	AddPayment(Payment) error

	// Update an existing payment in the stable storage
	// Precondition: The payment must already exist
	UpdatePayment(Payment) error

	// Creates or updates a payment in the stable storage
	// Replaces if it already existed
	StorePayment(Payment) error

	// Delete a payment from the stable storage
	// Precondition: The payment must already exist
	DeletePayment(string) error

	// Fetch a specific payment from the stable storage
	// Precondition: The payment must already exist
	GetPayment(string) (Payment, error)

	// Fetch a list of all payments from the stable storage
	GetAllPayments() ([]Payment, error)

	// Possible TODO: Pagination method
}

type InMemStore struct {
	payments map[string]Payment
	sync.RWMutex
}

func NewInMemStore() *InMemStore {
	store := InMemStore{
		payments: make(map[string]Payment),
	}
	return &store
}

func (s *InMemStore) AddPayment(p Payment) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.payments[p.ID]; ok {
		return errors.New("Cannot add an already existing resource")
	}

	s.payments[p.ID] = p
	return nil
}

func (s *InMemStore) UpdatePayment(p Payment) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.payments[p.ID]; !ok {
		return errors.New("Cannot update a non-existing resource")
	}

	s.payments[p.ID] = p
	return nil
}

func (s *InMemStore) StorePayment(p Payment) error {
	s.Lock()
	defer s.Unlock()

	s.payments[p.ID] = p
	return nil
}

func (s *InMemStore) DeletePayment(id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.payments[id]; !ok {
		return errors.New("Cannot delete a non-existing resource")
	}

	delete(s.payments, id)
	return nil
}

func (s *InMemStore) GetPayment(id string) (Payment, error) {
	var (
		p   Payment
		err error
		ok  bool
	)

	s.RLock()
	defer s.RUnlock()

	if p, ok = s.payments[id]; !ok {
		err = errors.New(fmt.Sprintf("No resource with ID %v", id))
	}

	return p, err
}

func (s *InMemStore) GetAllPayments() ([]Payment, error) {
	var ps []Payment

	s.RLock()
	defer s.RUnlock()

	for _, val := range s.payments {
		ps = append(ps, val)
	}

	return ps, nil
}
