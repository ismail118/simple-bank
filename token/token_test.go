package token

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

func Test_errors(t *testing.T) {
	errs := fmt.Errorf("%s:%s", sql.ErrNoRows, ErrExpiredToken)

	if strings.Contains(errs.Error(), ErrExpiredToken.Error()) {
		fmt.Println("lulus")
	}
}
