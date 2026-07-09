package digestsV2

import "golang.org/x/crypto/bcrypt"

type (
	BcryptEncoder interface {
		Hash() (bytes []byte)
		Valid(hash string) bool
		Invalid(hash string) bool
	}

	BcryptEncoderImpl struct{ password []byte }
)

// NewBcrypt 实例化
func NewBcrypt(password string) BcryptEncoder { return &BcryptEncoderImpl{password: []byte(password)} }

// Hash 编码
func (my *BcryptEncoderImpl) Hash() (bytes []byte) {
	bytes, _ = bcrypt.GenerateFromPassword(my.password, bcrypt.DefaultCost)
	return
}

// Valid 校验是否通过
func (my *BcryptEncoderImpl) Valid(hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), my.password) == nil
}

// Invalid 校验是否未通过
func (my *BcryptEncoderImpl) Invalid(hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), my.password) != nil
}
