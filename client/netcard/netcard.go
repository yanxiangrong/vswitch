package netcard

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"vswitch/pkg/util/log"
)

func SelectByUser() string {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("在您的电脑上找到以下网卡及对应IP地址")

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error(fmt.Sprint(err))
		os.Exit(-1)
	}

	index := 1
	options := make([]string, 1)
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		fmt.Printf("---- %s %s\n", i.Name, i.HardwareAddr)
		for _, addr := range addrs {
			s := addr.String()
			options = append(options, s[:strings.Index(s, "/")])
			fmt.Printf("%3d. %s\n", index, addr)
			index++
		}
	}
	user := 1
	for {
		fmt.Print("请输入流量出口IP地址的编号: ")

		_, err = fmt.Scanln(&user)
		if err != nil {
			log.Error(fmt.Sprint(err))
			os.Exit(-1)
		}
		if user < 1 || user > len(options)-1 {
			fmt.Println("请输入正确的编号")
			continue
		}
		break
	}

	fmt.Printf("使用IP %s 作为流量出口\n", options[user])
	time.Sleep(100 * time.Millisecond)
	return options[user]
}
