package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Config interface {
	GetString(key string) (string, error)
	GetInt64(key string) (int64, error)
	GetFloat64(key string) (float64, error)
}

type config struct {
	container map[string]string
}

func (config *config) GetString(key string) (string, error) {
	fail := ""
	v := config.container[key]
	if fail == v {
		return fail, fmt.Errorf("find no config for key : %s", key)
	}

	return v, nil
}

func (config *config) GetInt64(key string) (int64, error) {

	var fail int64 = 0

	v, err := config.GetString(key)
	if err != nil {
		return fail, err
	}

	vInt64, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fail, err
	}

	return vInt64, nil
}

func (config *config) GetFloat64(key string) (float64, error) {
	var fail float64 = 0

	v, err := config.GetString(key)
	if err != nil {
		return fail, err
	}

	vFloat64, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fail, err
	}

	return vFloat64, nil
}

func NewConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	container := make(map[string]string)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		filteredLine, err := filterSpace(line)
		if err != nil {
			return nil, err
		}

		isAnnotationOrWarp, err := inspectAnnotationOrWarp(line)
		if err != nil {
			return nil, err
		}

		if isAnnotationOrWarp {
			continue
		}

		key, value, err := resolveKeyValue(filteredLine)
		if err != nil {
			return nil, err
		}
		container[key] = value
	}

	return &config{container: container}, nil
}

func filterSpace(content []byte) ([]byte, error) {
	reg, err := regexp.Compile(`\s*`)
	if err != nil {
		return nil, err
	}

	content = reg.ReplaceAll(content, []byte(""))

	return content, nil
}

func inspectAnnotationOrWarp(content []byte) (bool, error) {
	isWarp := 0
	fail := false

	reg, err := regexp.Compile(`^#`)
	if err != nil {
		return fail, err
	}

	return len(content) == isWarp || reg.Match(content), nil
}

func resolveKeyValue(content []byte) (key, value string, err error) {
	v := string(content)
	sep := "="
	keyAndValue := strings.Split(v, sep)

	if isKeyValueStruct(keyAndValue) {
		err = fmt.Errorf("%s is no key-value struct", v)
		return
	}

	key = keyAndValue[0]
	value = keyAndValue[1]
	return
}

func isKeyValueStruct(v []string) bool {
	return len(v) != 2
}
