package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"gitlab-beta.engr.illinois.edu/sp-box/boxsync/auth/mocks"
)

var (
	now  = time.Now()
	tok1 = oauth2.Token{
		AccessToken:  "foo",
		TokenType:    "bearer",
		RefreshToken: "bar",
		Expiry:       now,
	}
	tok2 = tok1
	tok3 = oauth2.Token{
		AccessToken:  "baz",
		TokenType:    "bearer",
		RefreshToken: "quux",
		Expiry:       now.Add(1 * time.Hour),
	}
)

func TestTokensEqual(t *testing.T) {
	assert.True(t, tokensEqual(nil, nil))
	assert.True(t, tokensEqual(&tok1, &tok1))
	assert.True(t, tokensEqual(&tok1, &tok2))
	assert.False(t, tokensEqual(&tok1, &tok3))
	assert.False(t, tokensEqual(nil, &tok1))
	assert.False(t, tokensEqual(&tok1, nil))
}

func TestSameToken(t *testing.T) {
	src := new(mocks.TokenSource)
	src.On("Token").Return(&tok2, nil).Once()
	defer src.AssertExpectations(t)

	count := 0
	cts := CallbackTokenSource(&tok1, src, func(t *oauth2.Token) error {
		count++
		return nil
	})

	retTok, err := cts.Token()
	assert.NoError(t, err, "Token method should not return an error")

	assert.Equal(t, 0, count, "Callback should not have been called")
	assert.True(t, &tok2 == retTok, "Token pointers should be equal")
}

func TestDifferentToken(t *testing.T) {
	src := new(mocks.TokenSource)
	src.On("Token").Return(&tok3, nil).Twice()
	defer src.AssertExpectations(t)

	count := 0
	cts := CallbackTokenSource(&tok1, src, func(t *oauth2.Token) error {
		count++
		return nil
	})

	retTok, err := cts.Token()
	assert.NoError(t, err, "Token method should not return an error")
	assert.Equal(t, 1, count, "Callback should have been called")
	assert.True(t, &tok3 == retTok, "Token pointers should be equal")

	retTok, err = cts.Token()
	assert.NoError(t, err, "Token method should not return an error")
	assert.Equal(t, 1, count, "Callback should not have been called again")
	assert.True(t, &tok3 == retTok, "Token pointers should be equal")
}
