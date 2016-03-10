package config

import (
	"fmt"
	"testing"
)

func Test_NewConfig(t *testing.T) {
	config, err := NewConfig("app.conf")
	if err != nil {
		t.Error(err)
		return
	}

	vString, err := config.GetString("httpport")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(vString)

	vInt64, err := config.GetInt64("httpport")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(vInt64)

	vFloat64, err := config.GetFloat64("httpport")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(vFloat64)
}
