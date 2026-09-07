package services

import (
	"os"
	"strconv"
	"strings"
)

func GetPort() string {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		port = 7010
	}
	return strconv.Itoa(port)
}

func GetEnvValues() (url string, k8Cmd string, k8Params []string, timeout int) {
	url = os.Getenv("WEBPAGE_URL")
	if url == "" {
		url = "http://localhost"
	}

	k8Cmd = os.Getenv("K8_COMMAND")
	if k8Cmd == "" {
		k8Cmd = "kubectl"
	}

	k8Prms := os.Getenv("K8_PARAMS")
	if k8Prms != "" {
		k8Params = strings.Split(k8Prms, ",")
	} else {
		k8Params = []string{"get", "nodes"}
	}

	timeout, err := strconv.Atoi(os.Getenv("TIME_OUT_WORKER"))
	if err != nil {
		timeout = 12
	}

	return
}
