package token

import (
	"fmt"
	"github.com/ismail118/simple-bank/errors"
	"strings"
	"testing"
)

func Test_errors(t *testing.T) {
	errs := fmt.Errorf("%s:%s", errors.ErrNoRow, ErrExpiredToken)

	if strings.Contains(errs.Error(), ErrExpiredToken.Error()) {
		fmt.Println("lulus")
	}
}
