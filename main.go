package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procGetMessageW      = user32.NewProc("GetMessageW")
)

const (
	MOD_ALT  = 0x0001
	MOD_CTRL = 0x0002
	// A 键的虚拟键码
	VK_A      = 0x41
	WM_HOTKEY = 0x0312
)

type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct {
		X int32
		Y int32
	}
}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件路径失败:", err)
		return
	}
	exeDir := fmt.Sprintf("%s/%s", filepath.Dir(exePath), "ScreenCapture.exe")
	hotkeyID := 1
	// 注册热键 Ctrl + Alt + A
	ret, _, err := procRegisterHotKey.Call(
		0,                 // hWnd = NULL
		uintptr(hotkeyID), // 热键ID
		MOD_ALT|MOD_CTRL,  // 修改键
		VK_A,              // 键码
	)
	if ret == 0 {
		fmt.Println("注册热键失败:", err)
		return
	}
	defer procUnregisterHotKey.Call(0, uintptr(hotkeyID))
	var msg MSG
	for {
		// 阻塞直到收到消息
		ret, _, _ := procGetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0,
		)
		if int32(ret) == -1 {
			fmt.Println("GetMessageW 出错")
			time.Sleep(5 * time.Second)
			continue
		}

		if msg.Message == WM_HOTKEY {
			ScreenCapture(exeDir)
		}
	}
}
func ScreenCapture(exePath string) {
	cmd := exec.Command(exePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Start()
	if err != nil {
		return
	}
}
