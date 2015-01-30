package goconfig

import (
    "container/list"
    sj "github.com/bitly/go-simplejson"
    "sync"
    "time"
)

const (
    // default watcher check time inc
    DEFAULT_Update_Check_Time = 10 * time.Second
)

// goconfig manage
type Config struct {
    configs   map[string]*configFile
    notifys   *list.List
    is_closed bool
    com_chan  chan int

    mutex sync.RWMutex
}

// create goconfig manage
func New() *Config {
    new_config := &Config{
        configs: make(map[string]*configFile),
        notifys: list.New(),
    }
    new_config.start()
    return new_config
}

// register an config file
//   name: config name
//   path: config file path
//   need_update: is need to upadte when file changed
func (self *Config) Reg(name, path string, need_update bool) error {
    new_config_file, err := newConfigFile(path, need_update)
    if err != nil {
        return err
    }
    self.mutex.Lock()
    self.configs[name] = new_config_file
    self.mutex.Unlock()
    return nil
}

// register a channel to get update notify for what config is changed,
// if "./info.json" named "info"  changed, the notify_chan will recv "info"
func (self *Config) RegNotifyChan(notify_chan chan string) {
    self.mutex.Lock()
    self.notifys.PushBack(notify_chan)
    self.mutex.Unlock()
    return
}

type getDataImpFunc func(*configFile) *sj.Json

func (self *Config) getImp(name string, imp getDataImpFunc) *sj.Json {
    self.mutex.RLock()
    config_file := self.configs[name]
    self.mutex.RUnlock()

    if config_file == nil {
        return nil
    }

    return imp(config_file)
}

// get all data from config
func (self *Config) Get(name string) *sj.Json {
    return self.getImp(name, func(config_file *configFile) *sj.Json {
        re := config_file.get()
        return &re
    })
}

// get data by idx and name
// for xlsx config
func (self *Config) GetByIdx(name, idx string) *sj.Json {
    return self.getImp(name, func(config_file *configFile) *sj.Json {
        re := config_file.getById(idx)
        return &re
    })
}

const (
    command_reload = iota // reload all config data which need update
)

func (self *Config) notify(name string) {
    for e := self.notifys.Front(); e != nil; e = e.Next() {
        channel, ok := e.Value.(chan string)
        if ok {
            channel <- name
        }
    }
}

func (self *Config) closeNotify() {
    for e := self.notifys.Front(); e != nil; e = e.Next() {
        channel, ok := e.Value.(chan string)
        if ok {
            close(channel)
        }
        self.notifys.Remove(e)
    }
}

func (self *Config) reload() {
    self.mutex.RLock()
    defer self.mutex.RUnlock()
    for config_name, config_file := range self.configs {
        update_err, is_update := config_file.update()
        if (update_err == nil) && is_update {
            self.notify(config_name)
        }
    }
}

// start coroutine for update
func (self *Config) start() {
    self.com_chan = make(chan int)
    go func() {
        for {
            time.Sleep(DEFAULT_Update_Check_Time)

            self.mutex.RLock()
            if self.is_closed {
                self.mutex.RUnlock()
                return
            }
            self.mutex.RUnlock()
            self.com_chan <- command_reload
        }
    }()

    go func() {
        for {
            select {
            case command, ok := <-self.com_chan:
                if !ok {
                    return
                }
                switch command {
                case command_reload:
                    self.reload()
                }
            }
        }
    }()
}

// Close Config
func (self *Config) Close() {
    self.mutex.Lock()
    self.is_closed = true
    close(self.com_chan)
    self.closeNotify()
    self.configs = nil
    self.mutex.Unlock()
}
