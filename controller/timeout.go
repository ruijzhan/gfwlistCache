package controller

import (
	"strings"

	"github.com/ruijzhan/routeros"
)

func timeout(domain string) string {
	if strings.HasSuffix(domain, "googlevideo.com") {
		return "3d"
	}
	return routeros.MAX_TIMEOUT
}
