package digestV2

import "golang.org/x/crypto/bcrypt"

type (
	Bcrypter interface {
		Hash() (bytes []byte)
		Check(hash string) bool
	}
	Bcrypt struct{ password []byte }
)

// NewBcrypt 实例化
func NewBcrypt(password string) Bcrypter { return &Bcrypt{password: []byte(password)} }

// Hash 编码
func (my *Bcrypt) Hash() (bytes []byte) {
	bytes, _ = bcrypt.GenerateFromPassword(my.password, bcrypt.DefaultCost)
	return
}

// Check 校验
func (my *Bcrypt) Check(hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), my.password) != nil
}
