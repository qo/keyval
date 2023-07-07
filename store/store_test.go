package store

import (
	"errors"
	"testing"
)

func TestCreateStore(t *testing.T) {
	s, err := CreateStore()
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	expectedLen := 0
	actualLen := len(s.m)
	if actualLen != expectedLen {
		t.Errorf("expected %v, got %v", expectedLen, actualLen)
	}
}

func TestStorePut(t *testing.T) {
	s, _ := CreateStore()
	err := s.Put("name", "john")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	expectedName := "john"
	actualName := s.m["name"]
	if actualName != expectedName {
		t.Errorf("expected %v, got %v", expectedName, actualName)
	}
}

func TestStoreGet(t *testing.T) {
	s, _ := CreateStore()
	age, err := s.Get("age")
	expectedErr := ErrNoSuchKey
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
	expectedAge := ""
	if age != expectedAge {
		t.Errorf("expected %v, got %v", expectedAge, age)
	}
	_ = s.Put("name", "john")
	name, err := s.Get("name")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	expectedName := "john"
	if name != expectedName {
		t.Errorf("expected %v, got %v", expectedName, name)
	}
}

func TestStoreDelete(t *testing.T) {
	s, _ := CreateStore()
	err := s.Delete("age")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	_ = s.Put("age", "30")
	ageBeforeDelete, err := s.Get("age")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	expectedAgeBeforeDelete := "30"
	if ageBeforeDelete != expectedAgeBeforeDelete {
		t.Errorf("expected %v, got %v", expectedAgeBeforeDelete, ageBeforeDelete)
	}
	err = s.Delete("age")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	ageAfterDelete, err := s.Get("age")
	expectedErr := ErrNoSuchKey
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
	expectedAgeAfterDelete := ""
	if ageAfterDelete != expectedAgeAfterDelete {
		t.Errorf("expected %v, got %v", expectedAgeAfterDelete, ageAfterDelete)
	}
}
