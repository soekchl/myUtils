# myUtils
	
	主要模仿 Beego Log  在修改一些自己想要的东东
	
	SetShowLeave 设定输出Log级别
	
	SetOutputFileLog 设定输出名和是否输出到文件

---

例：
```
	SetShowLeave(LeaveWarning)
	SetOutputFileLog("test", true)
		设定后调用
	Notice("Notice")
	Warn("Warn")
```
	显示的时候只显示 Warn 文件记录的时候 Notice,Warn 全部记录
---

0.2 版本
	把所有日志放到 名字_log 目录里
	增加 log  细分存放