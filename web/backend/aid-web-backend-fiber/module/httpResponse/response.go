package httpResponse

import "github.com/gofiber/fiber/v2"

type (
	response struct {
		code int
		msg  string
		data any
	}
)

func New(attrs ...ResponseAttributer) response {
	res := response{}
	for _, attr := range attrs {
		attr.Register(&res)
	}
	return res
}

func NewOK(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{OK()}
	def = append(def, attrs...)
	return New(def...)
}

func NewCreated(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{Created()}
	def = append(def, attrs...)
	return New(def...)
}

func NewUpdated(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{Updated()}
	def = append(def, attrs...)
	return New(def...)
}

func NewDeleted(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{Deleted()}
	def = append(def, attrs...)
	return New(def...)
}

func NewNoLogin(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{NoLogin()}
	def = append(def, attrs...)
	return New(def...)
}

func NewNoPermission(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{NoPermission()}
	def = append(def, attrs...)
	return New(def...)
}

func NewForbidden(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{Forbidden()}
	def = append(def, attrs...)
	return New(def...)
}

func NewNotFound(attrs ...ResponseAttributer) response {
	var def []ResponseAttributer = []ResponseAttributer{NotFound()}
	def = append(def, attrs...)
	return New(def...)
}

func (my response) JSON(c *fiber.Ctx) error { return c.Status(my.code).JSON(my.data) }
