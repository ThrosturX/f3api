package f3api

import (
	"errors"
	"fmt"
	"sync"
)

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
	// Precondition: A payment with the resource ID already exist
	DeletePayment(id string) error

	// Fetch a specific payment from the stable storage
	// Precondition: A payment with the resource ID already exist
	GetPayment(id string) (Payment, error)

	// Fetch a list of all payments from the stable storage
	// NOTE: Does not paginate!
	GetAllPayments() ([]Payment, error)
}

// Simple in-memory stable storage implementation for testing and demonstration purposes
type InMemStore struct {
	payments map[string]Payment
	sync.RWMutex
}

// Creates a new, blank in-memory ApiStore
func NewInMemStore() *InMemStore {
	store := InMemStore{
		payments: make(map[string]Payment),
	}
	return &store
}

// Add a payment to the stable storage
//
// Precondition: The payment must not exist
func (s *InMemStore) AddPayment(p Payment) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.payments[p.ID]; ok {
		return errors.New("Cannot add an already existing resource")
	}

	s.payments[p.ID] = p
	return nil
}

// Update an existing payment in the stable storage
//
// Precondition: The payment must already exist
func (s *InMemStore) UpdatePayment(p Payment) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.payments[p.ID]; !ok {
		return errors.New("Cannot update a non-existing resource")
	}

	s.payments[p.ID] = p
	return nil
}

// Creates or updates a payment in the stable storage, replacing if necessary/possible
func (s *InMemStore) StorePayment(p Payment) error {
	s.Lock()
	defer s.Unlock()

	s.payments[p.ID] = p
	return nil
}

// Delete a payment from the stable storage
//
// Precondition: The payment must already exist
func (s *InMemStore) DeletePayment(id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.payments[id]; !ok {
		return errors.New("Cannot delete a non-existing resource")
	}

	delete(s.payments, id)
	return nil
}

// Fetch a specific payment from the stable storage
//
// Precondition: A payment with the resource ID must already exist
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

// Fetch a list of all payments from the stable storage
func (s *InMemStore) GetAllPayments() ([]Payment, error) {
	var ps []Payment

	s.RLock()
	defer s.RUnlock()

	for _, val := range s.payments {
		ps = append(ps, val)
	}

	return ps, nil
}
