package v1HTTPService

import (
	"strings"

	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/module/httpModule/v1HTTPModule/request"
	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/module/httpModule/v1HTTPModule/response"

	"github.com/gofrs/uuid/v5"
)

type UUIDService struct{}

var UUID UUIDService

// GenerateOne 生成单个UUID
func (UUIDService) GenerateOne(form *request.UUIDGenerateRequest) (response.UUIDResponse, error) {
	var (
		err     error
		u       uuid.UUID
		uuidStr string
	)
	switch form.Version {
	case request.UUIDVersionV1:
		u, err = uuid.NewV1()
	case request.UUIDVersionV4:
		u, err = uuid.NewV4()
	case request.UUIDVersionV6:
		u, err = uuid.NewV6()
	case request.UUIDVersionV7:
		u, err = uuid.NewV7()
	default:
		u, err = uuid.NewV6()
	}

	if err != nil {
		return response.UUIDResponse{}, err
	}

	uuidStr = u.String()

	if form.IsUpper {
		uuidStr = strings.ToUpper(uuidStr)
	}

	if form.NoSubsTractKey {
		uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	}

	return response.UUID.New(uuidStr), nil
}
