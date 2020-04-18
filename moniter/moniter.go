// 专门用于监控硬件状态
package moniter

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	. "github.com/soekchl/myUtils"
)

var (
	cpuInfo cpu.InfoStat
)

func init() {
	c, err := cpu.Info()
	if err != nil {
		Error(err)
		return
	}
	for i, v := range c {
		if i == 0 {
			cpuInfo = v
		}
	}
}

func GetUseCpuPercent() ([]float64, error) {
	return cpu.Percent(0, false)
}

func GetMemInfo() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

func GetDiskInfo() (us []*disk.UsageStat, err error) {
	d, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	for _, v := range d {
		t, err := disk.Usage(v.Device)
		if err != nil {
			return nil, err
		}
		us = append(us, t)
	}
	return

}

func GetHtml() (str string) {

	str = `
<!DOCTYPE html>
<head>
<meta charset="UTF-8">
<title>服务器状态</title>


</head>
<body>
`

	str += "系统：" + runtime.GOOS + "</br></br>"
	v, err := mem.VirtualMemory()
	if err != nil {
		Error(err)
		return
	}
	str = fmt.Sprint(str, "总内存：", v.Total/1000000, "</br>")
	str = fmt.Sprint(str, "已缓存：", v.Used/1000000, "</br>")
	str = fmt.Sprint(str, "可用：", v.Available/1000000, "</br>")
	str = fmt.Sprint(str, "使用率：", v.UsedPercent, "%</br>")

	str += `<canvas id="memmery" width="250" height="250" ></canvas>`

	d, err := disk.Partitions(false)
	if err != nil {
		Error(err)
		return
	}
	str += "</br>---------------------------------------------------</br>"
	diskShow := ""
	tmpStr := ""
	for i, v := range d {
		dd, err := disk.Usage(v.Device)
		if err != nil {
			Error(err)
			return
		}
		str = fmt.Sprintf("%s<b>%2v</b> 总共：%5v G  | 已使用：%4.1f %%  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;", str, dd.Path, dd.Total/1024/1024/1024, dd.UsedPercent)
		tmpStr += fmt.Sprintf("<canvas id=\"disk%d\" width=\"250\" height=\"250\" ></canvas>", i+1)
		diskShow += fmt.Sprintf("show('disk%d', %v);\n", i+1, int(dd.UsedPercent))
	}
	str = fmt.Sprint(str, "</br>", tmpStr, "</br>")
	str += "</br>---------------------------------------------------</br>"

	str = fmt.Sprint(str, "CPU信息：", cpuInfo.ModelName, "</br>")
	str = fmt.Sprint(str, "CPU核数：", cpuInfo.Cores, "</br>")
	str = fmt.Sprint(str, "</br>CPU使用率：</br>")

	cp, err := cpu.Percent(0, false)
	if err != nil {
		Error(err)
		return
	}
	tmpStr = ""
	cpuShow := ""
	for i, v := range cp {
		tmpStr += fmt.Sprintf("<canvas id=\"cpu%d\" width=\"250\" height=\"250\" ></canvas>", i+1)
		cpuShow += fmt.Sprintf("show('cpu%d', %v);\n", i+1, v)
	}
	str = fmt.Sprint(str, "</br>", tmpStr, "</br>")

	str += `
<script>
	show('memmery', ` + fmt.Sprint(v.UsedPercent) + `);
	` + diskShow + cpuShow + `
	function show(name, persent){
		var canvas = document.getElementById(name),  //获取canvas元素
			context = canvas.getContext('2d'),  //获取画图环境，指明为2d
			centerX = canvas.width/2,   //Canvas中心点x轴坐标
			centerY = canvas.height/2,  //Canvas中心点y轴坐标
			rad = Math.PI*2/100, //将360度分成100份，那么每一份就是rad度
			speed = 0.1; //加载的快慢就靠它了 
		//绘制蓝色外圈
		function blueCircle(n){
			context.save();
			context.strokeStyle = "#000"; //设置描边样式
			context.lineWidth = 5; //设置线宽
			context.beginPath(); //路径开始
			context.arc(centerX, centerY, 100 , -Math.PI/2, -Math.PI/2 +n*rad, false); //用于绘制圆弧context.arc(x坐标，y坐标，半径，起始角度，终止角度，顺时针/逆时针)
			context.stroke(); //绘制
			context.closePath(); //路径结束
			context.restore();
		}
		//绘制白色外圈
		function whiteCircle(){
			context.save();
			context.beginPath();
			context.strokeStyle = "black";
			context.arc(centerX, centerY, 100 , 0, Math.PI*2, false);
			context.stroke();
			context.closePath();
			context.restore();
		}  
		//百分比文字绘制
		function text(n){
			context.save(); //save和restore可以保证样式属性只运用于该段canvas元素
			context.strokeStyle = "#000"; //设置描边样式
			context.font = "40px Arial"; //设置字体大小和字体
			//绘制字体，并且指定位置
			context.strokeText(n.toFixed(0)+"%", centerX-25, centerY+10);
			context.stroke(); //执行绘制
			context.restore();
		} 
		//动画循环
		(function drawFrame(){
			context.clearRect(0, 0, canvas.width, canvas.height);
			whiteCircle();
			text(persent);
			blueCircle(persent);
			if(speed > 100) speed = 0;
			speed += 0.1;
		}());
	}
</script>

<div style="text-align:center;clear:both">
<script src="/gg_bd_ad_720x90.js" type="text/javascript"></script>
<script src="/follow.js" type="text/javascript"></script>
</div>

</body>
</html>
	`
	return
}
