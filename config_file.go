package goconfig

import (
	sj "github.com/bitly/go-simplejson"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// ConfugFile
// ConfugFile is a file to watch and read
type ConfigFile struct {
	path  string
	mutex sync.RWMutex
	data  *sj.Json

	time_updated   int64
	is_need_update bool
}

//  new ConfugFile
func newConfigFile(path string, is_need_update bool) (*ConfigFile, error) {
	config_file := &ConfigFile{
		path:           path,
		is_need_update: is_need_update,
	}
	err := config_file.init()
	return config_file, err
}

func (self *ConfigFile) init() error {
	return self.updateConfig()
}

func (self *ConfigFile) get() sj.Json {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	// sj.Json is value
	return *self.data
}

func (self *ConfigFile) getById(idx string) sj.Json {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return *self.data.Get(idx)
}

func (self *ConfigFile) update() (error, bool) {
	is_need, err := self.isNeedUpdate()
	if is_need && err == nil {
		err := self.updateConfig()
		return err, true
	}

	return nil, false
}

func (self *ConfigFile) updateConfig() error {
	bytes, err := ioutil.ReadFile(self.path)
	if err != nil {
		return err
	}

	js, js_err := sj.NewJson(bytes)
	if js_err != nil {
		return js_err
	}

	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.data = js
	self.time_updated = time.Now().Unix()

	return nil
}

func (self *ConfigFile) isNeedUpdate() (bool, error) {
	if !self.is_need_update {
		return false, nil
	}

	file_info, err := os.Stat(self.path)

	if err != nil {
		return false, err
	}

	file_time := file_info.ModTime().Unix()

	self.mutex.RLock()
	defer self.mutex.RUnlock()
	if file_time > self.time_updated {
		return true, nil
	} else {
		return false, nil
	}
}
