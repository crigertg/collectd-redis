package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type redisInstance struct {
	name string
	host string
	port int
	pw   string
}

func validateConnectionString(arg string) error {
	// regular expression matches <host>:<port>@<password>
	// Host can be a name/domain or IP. Password ist optional
	re := regexp.MustCompile(`^[\w-.]+:[\w.-]+(?:\.[\w\.-]+)*:([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])(:\w+)?$`)

	if !re.MatchString(arg) {
		return fmt.Errorf("%s does not match <name>:<host>:<port>[:<password>] scheme", arg)
	}
	return nil
}

func parseArgToInstance(connStr string) redisInstance {
	errCheckFatal(validateConnectionString(connStr))
	connInfo := strings.Split(connStr, ":")
	name := connInfo[0]
	host := connInfo[1]
	port, err := strconv.Atoi(connInfo[2])
	errCheckFatal(err)
	var pw string
	if len(connInfo) > 3 {
		pw = connInfo[3]
	} else {
		pw = ""
	}

	return redisInstance{
		name: name,
		host: host,
		port: port,
		pw:   pw,
	}
}

func (reIn redisInstance) fetchInfo() (string, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", reIn.host, reIn.port),
		Password: reIn.pw,
	})
	output, err := client.Info().Result()
	return output, err
}
