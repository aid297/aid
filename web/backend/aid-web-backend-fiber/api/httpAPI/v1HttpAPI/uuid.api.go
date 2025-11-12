package v1HttpAPI

import (
	"strings"

	"hr-fiber/global"
	"hr-fiber/module/httpRequest/v1HTTPRequest"
	"hr-fiber/module/httpResponse"
	"hr-fiber/module/httpResponse/v1HTTPResponse"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
)

type UUIDAPI struct{}

var UUID UUIDAPI

// Generate 批量生成uuid
// * URL POST /api/v1/uuid/generate
func (*UUIDAPI) Generate(c *fiber.Ctx) error {
	var (
		title   = "批量生成uuid"
		err     error
		form    v1HTTPRequest.UUIDGenerateRequest
		uuidStr string
		uuids   []string
		uuidRes string
	)

	if form, err = v1HTTPRequest.UUIDGenerate.Bind(c); err != nil {
		global.LOGGER.Error(title, zap.String("表单验证", err.Error()))
		return httpResponse.NewForbidden(httpResponse.Errorf("表单验证失败：%w", err)).JSON(c)
	}

	uuids = make([]string, 0, form.Number)

	generateUuid := func(form *v1HTTPRequest.UUIDGenerateRequest) (string, error) {
		var u uuid.UUID
		switch form.Version {
		case v1HTTPRequest.UUIDVersionV1:
			u, err = uuid.NewV1()
		case v1HTTPRequest.UUIDVersionV4:
			u, err = uuid.NewV4()
		case v1HTTPRequest.UUIDVersionV6:
			u, err = uuid.NewV6()
		case v1HTTPRequest.UUIDVersionV7:
			u, err = uuid.NewV7()
		default:
			u, err = uuid.NewV6()
		}

		if err != nil {
			return "", err
		}

		uuidStr = u.String()
		if form.IsUpper {
			uuidStr = strings.ToUpper(uuidStr)
		}
		if form.NoSubsTractKey {
			uuidStr = strings.ReplaceAll(uuidStr, "-", "")
		}

		return uuidStr, nil
	}

	for range form.Number {
		if uuidRes, err = generateUuid(&form); err != nil {
			global.LOGGER.Error(title, zap.String("生成uuid失败", err.Error()))
			return httpResponse.NewForbidden(httpResponse.Errorf("生成uuid失败：%w", err)).JSON(c)
		}
		uuids = append(uuids, uuidRes)
	}

	return httpResponse.NewOK(httpResponse.Data(v1HTTPResponse.UUIDGenerateResponse{UUIDs: uuids})).JSON(c)
}

func (*UUIDAPI) Versions(c *fiber.Ctx) error {
	return httpResponse.NewOK(httpResponse.Data(
		v1HTTPResponse.UUIDVersionsResponse{Versions: map[string]v1HTTPRequest.UUIDVersion{
			string(v1HTTPRequest.UUIDVersionV1): v1HTTPRequest.UUIDVersionV1,
			string(v1HTTPRequest.UUIDVersionV4): v1HTTPRequest.UUIDVersionV4,
			string(v1HTTPRequest.UUIDVersionV6): v1HTTPRequest.UUIDVersionV6,
			string(v1HTTPRequest.UUIDVersionV7): v1HTTPRequest.UUIDVersionV7,
		}},
	)).JSON(c)
}
