package goconfig

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func changeConfigFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString(`{
		"info1" : 333,
		"info2" : "aaa2",
		"info3" : "aaa3",
		"info_arr" : ["a1", "a2", "a3"],
		"info_obj" : {
			"o1" : "sss",
			"o2" : ["o21", "o22"]
			}
	}`)
	writer.Flush()

}

func TestConfigImp(t *testing.T) {
	config_mng := NewConfig()
	config_mng.Reg("info", "test_config_obj.json", true)

	data := config_mng.Get("info")

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

	if data.Get("info3").MustString() != "aaa3" {
		t.Error("info error!")
	}

	config_mng.Close()

	return
}

func TestConfigGetByIdx(t *testing.T) {
	config_mng := NewConfig()
	config_mng.Reg("info", "test_config_array.json", true)

	data1 := config_mng.GetByIdx("info", "1")
	data2 := config_mng.GetByIdx("info", "2")
	datanil := config_mng.GetByIdx("info", "9999")

	t.Log("config : ", data1)
	t.Log("config : ", data2)
	t.Log("config : ", datanil)

	if data1.Get("Info").MustString() != "Info1" {
		t.Error("data1 error!")
	}
	if data2.Get("Info").MustString() != "Info2" {
		t.Error("data1 error!")
	}

	config_mng.Close()

	return
}

func TestConfigUpdate(t *testing.T) {
	update_chan := make(chan string)

	config_mng := NewConfig()
	config_mng.Reg("info", "test_config_obj.json", true)
	config_mng.RegNotifyChan(update_chan)

	data := config_mng.Get("info")

	t.Log("config : ", data)

	go func() {
		time.Sleep(2 * time.Second)
		changeConfigFile("test_config_obj.json")
	}()

	for {
		select {
		case update_name, ok := <-update_chan:
			if !ok {
				return
			}
			t.Log("config : ", update_name)
			data = config_mng.Get("info")
			goto RETURN_INFO
		}
	}

RETURN_INFO:
	config_mng.Close()

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
