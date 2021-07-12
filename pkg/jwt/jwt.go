package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/smiletrl/micro_ecommerce/pkg/constants"
)

type Provider interface {
	ParseCustomerToken(c echo.Context) (customerID int64, err error)
	NewCustomerToken(customerID int64) (token string, err error)
}

type provider struct {
	JwtSecret string
}

func NewProvider(secret string) Provider {
	return provider{secret}
}

// Authorization header
var authScheme string = "Bearer"

// ParseToken get jwt token from request header, and then get user id signed inside the token
func (p provider) ParseCustomerToken(c echo.Context) (int64, error) {
	var (
		token *jwt.Token
		err   error
	)
	auth := c.Request().Header.Get("Authorization")
	l := len(authScheme)
	if len(auth) <= l+1 || auth[:l] != authScheme {
		return 0, errors.New("Missing auth bearer token in header")
	}
	tokenString := auth[l+1:]
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != constants.AlgorithmHS256 {
			return nil, errors.Errorf("header sign incorrect: %v", t.Header["alg"])
		}
		return []byte(p.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("incorrect or outdated auth token")
	}

	claims := token.Claims.(jwt.MapClaims)
	customerID, ok := claims[constants.AuthCustomerID].(float64)
	if !ok {
		return 0, errors.Errorf("unrecognized customer id: %v", claims[constants.AuthCustomerID])
	}

	return int64(customerID), nil
}

// NewToken from user id. `user` is `staff` in this case.
// We may want to add other info into jwt token for other purposes.
func (p provider) NewCustomerToken(customerID int64) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims[constants.AuthCustomerID] = customerID

	// Set the token to be valid for 3 days.
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(p.JwtSecret))
	return tokenString, err
}

type mockProvider struct{}

func NewMockProvider() Provider {
	return mockProvider{}
}

func (m mockProvider) ParseCustomerToken(c echo.Context) (int64, error) {
	return int64(0), nil
}

func (m mockProvider) NewCustomerToken(userID int64) (string, error) {
	return "secret_token", nil
}
