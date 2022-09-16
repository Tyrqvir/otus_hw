package hw10programoptimization

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/goccy/go-json"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var userPool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat, err := getDomainStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return domainStat, nil
}

func getDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)

	result := make(DomainStat)

	for scanner.Scan() {
		if !bytes.Contains(scanner.Bytes(), []byte("."+domain)) {
			continue
		}

		user := *userPool.Get().(*User)

		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		userPool.Put(&user)

		at := strings.LastIndex(user.Email, "@")

		if at == 0 {
			continue
		}

		result[strings.ToLower(user.Email[at+1:])]++
	}

	return result, nil
}
