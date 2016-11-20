package auth

import (
	"sync"

	"golang.org/x/oauth2"
)

type callbackTokenSource struct {
	src      oauth2.TokenSource
	callback func(*oauth2.Token) error
	t        *oauth2.Token
	mu       sync.Mutex
}

// CallbackTokenSource returns a TokenSource that calls callback whenever its
// Token method is called and the underlying TokenSource's Token method returns
// a different token from its last returned token. This occurs when the
// underlying TokenSource has refreshed the token.
func CallbackTokenSource(t *oauth2.Token, src oauth2.TokenSource, callback func(*oauth2.Token) error) oauth2.TokenSource {
	return &callbackTokenSource{
		src:      src,
		callback: callback,
		t:        t,
	}
}

func (s *callbackTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, err := s.src.Token()
	if err != nil {
		return nil, err
	}

	if !tokensEqual(t, s.t) {
		err := s.callback(t)
		if err != nil {
			return nil, err
		}
	}

	s.t = t
	return t, nil
}

func tokensEqual(t1, t2 *oauth2.Token) bool {
	return (t1 == t2 ||
		(t1 != nil && t2 != nil &&
			t1.AccessToken == t2.AccessToken &&
			t1.TokenType == t2.TokenType &&
			t1.RefreshToken == t2.RefreshToken &&
			t1.Expiry == t2.Expiry))
}
