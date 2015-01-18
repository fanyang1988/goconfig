package config

import (
	"testing"
)

func TestConfigFile(t *testing.T) {
	configfile, err := NewConfigFile("test_config_obj.json", true)
	if err != nil {
		t.Error("new err :", err.Error())
	}

	data := configfile.Get()

	t.Log("config : ", data)
	t.Log("config : ", data.Get("info1"))
	t.Log("config : ", data.Get("info2"))
	t.Log("config : ", data.Get("info3"))
	t.Log("config : ", data.Get("info_arr"))
	t.Log("config : ", data.Get("info_obj"))
	t.Log("config : ", data.Get("info_arr").GetIndex(2))
	t.Log("config : ", data.Get("info_obj").Get("o1"))

	t.Log("config : ", data.Get("info_obj").Get("o2").GetIndex(4))
	t.Log("config : ", data.Get("info_obj333").Get("o2").GetIndex(4))

	return
}
