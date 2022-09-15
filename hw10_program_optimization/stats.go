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
	u, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return countDomains(u)
}

type users []User

func getUsers(r io.Reader, domain string) (result users, err error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		if !bytes.Contains(scanner.Bytes(), []byte("."+domain)) {
			continue
		}

		user := *userPool.Get().(*User)
		userPool.Put(&user)

		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}

		result = append(result, user)
	}

	return
}

func countDomains(u users) (DomainStat, error) {
	result := make(DomainStat, len(u))

	for _, user := range u {
		at := strings.LastIndex(user.Email, "@")

		if at == 0 {
			continue
		}

		nextDomainLvl := strings.ToLower(user.Email[at+1:])

		result[nextDomainLvl]++
	}
	return result, nil
}
