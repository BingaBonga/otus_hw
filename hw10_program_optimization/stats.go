package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	//nolint:depguard
	jsoniter "github.com/json-iterator/go"
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	json := jsoniter.ConfigFastest
	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			if !strings.Contains(user.Email, "@") {
				return nil, fmt.Errorf("invalid email: %s does not contain @", user.Email)
			}

			foundDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[foundDomain]++
		}
	}
	return result, nil
}
