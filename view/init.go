package view

import (
	"github.com/Bit0r/online-store/conf"
)

var (
	tplRoot = conf.Get("website", "online_store_website", "template_dir").(string) + "/"
)
