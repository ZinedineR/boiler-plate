package responsehelper

import (
	"boiler-plate/internal/base/domain"
)

func GetStatusResponse(statusCode int, description string) *domain.Status {
	switch statusCode {
	case 200:
		return &domain.Status{
			ResponseCode:    statusCode,
			ResponseMessage: "Successful",
		}

	case 500:
		return &domain.Status{
			ResponseCode:    statusCode,
			ResponseMessage: "service unavailable",
		}

	}

	return &domain.Status{
		ResponseCode:    statusCode,
		ResponseMessage: description,
	}
}
