package helpers

import (
	"SamkoOfMraz/models"
	"fmt"
	"math/rand"
	"net"
)

func GenerateToken() string {
	var token = ""
	var num = 0
	for i := 0; i < 40; i++ {
		num = rand.Intn(90-49+1) + 49

		if (num < 57 && num > 49) || (num < 90 && num > 65) {
			token = token + string(rune(num)) //rune!!!
		} else {
			i--
		}

	}
	return token
}
func IpAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {

				continue
			}

			for _, addr := range addrs {

				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					fmt.Println(ipnet.IP.String())
					return ipnet.IP.String()
				}
			}
		}
	}

	return "localhost"
}
func ContainsValue(m map[models.UserForGet]string, value string) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}
func FindKeyByValue(m map[models.UserForGet]string, value string) models.UserForGet {
	for k, v := range m {
		if v == value {
			return k
		}
	}
	return models.UserForGet{}
}
func GetUserByToken(tokenMap map[models.UserForGet]string, token string) (models.UserForGet, bool) {
	for user, t := range tokenMap {
		if t == token {
			return user, true
		}
	}
	return models.UserForGet{}, false
}
