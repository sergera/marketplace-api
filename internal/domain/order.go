package domain

import (
	"errors"
	"regexp"
	"time"
)

type Status int8

const (
	Unknown     Status = 0
	Unconfirmed Status = iota + 1
	InProgress
	Ready
	InTransit
	Delivered
)

func (s Status) String() string {
	switch s {
	case Unconfirmed:
		return "unconfirmed"
	case InProgress:
		return "in_progress"
	case Ready:
		return "ready"
	case InTransit:
		return "in_transit"
	case Delivered:
		return "delivered"
	}
	return "unknown"
}

type OrderModel struct {
	Id     string    `json:"id"`
	Price  uint64    `json:"price"`
	Status string    `json:"Status"`
	Date   time.Time `json:"date"`
}

func (o OrderModel) Validate() error {
	if err := o.ValidateStatus(); err != nil {
		return err
	}
	return nil
}

func (o OrderModel) ValidateStatus() error {
	switch o.Status {
	case Unconfirmed.String():
		return nil
	case InProgress.String():
		return nil
	case Ready.String():
		return nil
	case InTransit.String():
		return nil
	case Delivered.String():
		return nil
	default:
		return errors.New("invalid order status")
	}
}

func (o OrderModel) StatusType() (Status, error) {
	switch o.Status {
	case Unconfirmed.String():
		return Unconfirmed, nil
	case InProgress.String():
		return InProgress, nil
	case Ready.String():
		return Ready, nil
	case InTransit.String():
		return InTransit, nil
	case Delivered.String():
		return Delivered, nil
	default:
		return Unknown, errors.New("invalid order status")
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
