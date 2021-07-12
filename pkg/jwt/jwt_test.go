package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomerToken(t *testing.T) {
	s := NewProvider("secret")
	token, err := s.NewCustomerToken(int64(12))
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	fmt.Println(token)
}
