package jwt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCustomerToken(t *testing.T) {
	s := NewProvider("secret")
	token, err := s.NewCustomerToken(int64(12))
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	fmt.Println(token)
}
