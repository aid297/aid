package secret

import (
	"github.com/aid297/aid/secret/asymmetric"
	"github.com/aid297/aid/secret/symmetric"
)

type APP struct {
	Asymmetric asymmetric.APP
	Symmetric  symmetric.APP
}
