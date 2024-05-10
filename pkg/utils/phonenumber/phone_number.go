package phonenumber

import (
	"log/slog"

	"github.com/nyaruka/phonenumbers"
)

type RegionCode string

const (
	RegionCodeMalaysia  RegionCode = "MY"
	RegionCodeIndonesia RegionCode = "ID"
)

type PhoneNumber struct {
	number *phonenumbers.PhoneNumber
}

func NewPhoneNumber(number string, regionCode RegionCode) (*PhoneNumber, error) {
	parse, err := phonenumbers.Parse(number, string(regionCode))
	if err != nil {
		slog.Error("failed to parse phone number", "error", err.Error())
		return nil, err
	}

	return &PhoneNumber{number: parse}, nil
}

func (p *PhoneNumber) String() string {
	return phonenumbers.Format(p.number, phonenumbers.E164)[1:]
}

func (p *PhoneNumber) E164() string {
	return phonenumbers.Format(p.number, phonenumbers.E164)
}

func (p *PhoneNumber) National() string {
	return phonenumbers.Format(p.number, phonenumbers.NATIONAL)
}

func (p *PhoneNumber) International() string {
	return phonenumbers.Format(p.number, phonenumbers.INTERNATIONAL)
}

func (p *PhoneNumber) IsValid() bool {
	return phonenumbers.IsValidNumber(p.number)
}
