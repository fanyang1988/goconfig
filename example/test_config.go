package main

import (
	"fmt"
	"github.com/fanyang1988/goconfig"
)

func main() {
	fmt.Printf("This is an Test\n")
	notify_chan := make(chan string)
	config := goconfig.NewConfig()
	config.Start()
	config.Reg("info", "../test_config_obj.json", true)
	config.RegNotifyChan(notify_chan)

	info_data := config.Get("info")
	fmt.Printf("info_data %s", info_data.Get("info1"))

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
