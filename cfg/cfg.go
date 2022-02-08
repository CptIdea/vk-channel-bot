package cfg

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Token  string
	Admins []int
}

func ParseEnv() (*Config, error) {
	var admins []int
	adminsString := strings.Split(os.Getenv("ADMINS"), ",")
	for _, s := range adminsString {
		id, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("error parse admin id(%s): %s ", s, err)
		}
		admins = append(
			admins,
			id,
		)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, errors.New("TOKEN must not be empty")
	}

	return &Config{
		Token:  token,
		Admins: admins,
	}, nil
}
