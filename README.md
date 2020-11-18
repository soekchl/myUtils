# myUtils
	
	主要模仿 Beego Log  在修改一些自己想要的东东
	
	SetShowLeave 设定输出Log级别
	
	SetOutPutLeave 设定 输出到文件的Log级别
	
	SetOutputFileLog 设定输出名和是否输出到文件
	
	处理部分log的方法
		grep -r "【E】" xxx.log >>  error.log
	
	完善测试用例
		go test -bench=.*
---

例：
```
	SetShowLeave(LeaveWarning) 
	//SetShowLeave(LeaveNoShow) // 完全不显示-高并发项目中用到比较多
	SetOutputFileLog("test")
		设定后调用
	Notice("Notice")
	Warn("Warn")
```
	显示的时候只显示 Warn 文件记录的时候 Notice,Warn 全部记录
---

0.2 版本
	把所有日志放到 名字_log 目录里

	增加 log  细分存放

1.0 版本
	删除  -> 增加 log  细分存放

	优化log输出文件缓慢

	优化文件显示频繁创建变量

	增加文件输出的时候缓存并且间隔输出

	增加立即文件输出接口	（建议defer使用）
