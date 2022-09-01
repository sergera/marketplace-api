package domain

import (
	"errors"
	"regexp"
	"time"
)

type status int8

const (
	unconfirmed status = iota
	inProgress
	ready
)

func (s status) String() string {
	switch s {
	case unconfirmed:
		return "unconfirmed"
	case inProgress:
		return "in_progress"
	case ready:
		return "ready"
	}
	return "unknown"
}

type OrderModel struct {
	Id     string    `json:"id"`
	Price  string    `json:"price"`
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
}

func (o OrderModel) Validate() error {
	if err := o.ValidatePrice(); err != nil {
		return err
	}
	if err := o.ValidateStatus(); err != nil {
		return err
	}
	return nil
}

func (o OrderModel) ValidatePrice() error {
	pattern := "^[1-9](?:[0-9]+)?$"

	match, err := regexp.MatchString(pattern, o.Price)
	if err != nil {
		return err
	}

	if match {
		return nil
	}

	return errors.New("invalid order price")
}

func (o OrderModel) ValidateStatus() error {
	switch o.Status {
	case unconfirmed.String():
		return nil
	case inProgress.String():
		return nil
	case ready.String():
		return nil
	default:
		return errors.New("invalid order status")
	}
}

type OrderRangeModel struct {
	Start       string
	End         string
	OldestFirst bool
}

func (o OrderRangeModel) Validate() error {
	errorMsg := "invalid range"

	pattern := "^[1-9](?:[0-9]+)?$"

	match, err := regexp.MatchString(pattern, o.Start)
	if err != nil {
		return err
	}

	if !match {
		return errors.New(errorMsg)
	}

	match, err = regexp.MatchString(pattern, o.End)
	if err != nil {
		return err
	}

	if !match {
		return errors.New(errorMsg)
	}

	return nil
}
