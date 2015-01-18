goconfig
===================
goconfig is a json config reader and wather  for mobile game server
Current version: 0.0.1

----------

###Getting Started
####To install:
    go get github.com/fanyang1988/goconfig
####Base  usage:
this  is read info

```go
package main

import "github.com/lunny/tango"

func main() {
	config := goconfig.NewConfig()
	config.Reg("info", "../test_config_obj.json", true)
	info_data := config.Get("info")
}
```
this read ../test_config_obj.json as name "info"
To read config data, it need to register file whit a name
use [Reg()](https://gowalker.org/github.com/fanyang1988/goconfig#Config_Reg) function  register "../test_config_obj.json" to config name
And then we can get data by config name

For config in a array like this:
```json
{
    "1" : {
        "Id" : 1,
        "Info" : "Info1"
    },
    "2" : {
        "Id" : 2,
        "Info" : "Info2"
    }
}
``` 
can use [GetByIdx()](https://gowalker.org/github.com/fanyang1988/goconfig#Config_GetByIdx) to get info by  index, this can use less copy

####Notify by Channel when config changed
it can use [RegNotifyChan](https://gowalker.org/github.com/fanyang1988/goconfig#Config_RegNotifyChan) to register  a channel to subscribe config change event, goconfig will send name which changed to channel
```go
func main() {
	notify_chan := make(chan string)

	config := goconfig.NewConfig()
	config.Reg("info", "../test_config_obj.json", true)

	config.RegNotifyChan(notify_chan)

	for {
		select {
		case update_name, ok := <-notify_chan:
			if !ok {
				fmt.Printf("chan closed")
				return
			}

			fmt.Printf("config update %s", update_name)
			new_info_data := config.Get(update_name)
			fmt.Printf("info_data %s", new_info_data.Get("info1"))
			return
		}
	}
}

```

####API
[API Reference](https://gowalker.org/github.com/fanyang1988/goconfig)
