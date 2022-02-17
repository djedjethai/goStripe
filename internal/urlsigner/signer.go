package urlsigner

import (
	"fmt"
	"strings"
	"time"

	goalone "github.com/bwmarrin/go-alone"
)

type Signer struct {
	Secret []byte
}

func (s *Signer) GenerateTokenFromString(data string) string {
	var urlToSign string

	// we gonna add the hash as url params
	// so lets make sure there not one already
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))
	// cast the token
	token := string(tokenBytes)

	return token
}

// verify that the link people click on has not been changed
func (s *Signer) VerifyToken(token string) bool {

	crypt := goalone.New(s.Secret, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(token))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (s *Signer) Expired(token string, minutesUntilExpire int) bool {

	crypt := goalone.New(s.Secret, goalone.Timestamp)
	ts := crypt.Parse([]byte(token))

	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
