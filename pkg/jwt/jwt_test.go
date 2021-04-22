package jwt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCustomerToken(t *testing.T) {
	s := NewService("secret")
	token, err := s.NewCustomerToken(int64(15))
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	fmt.Println(token)
}
