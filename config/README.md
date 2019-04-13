# config
	
	配置文件
	
---

例：
```
	package main
	
	import (
		. "github.com/soekchl/myUtils"
		config "github.com/soekchl/myUtils/config"
	)
	
	func main() {
		config.InitConfig("./config/config.ini", "dev")
		Notice(config.GetConfigBool("test"))
	}
```

config.ini(utf-8 file)
```
	[dev]
	test		=	 true		# desc
```