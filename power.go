//go:build windows

package main

import (
	"log"
	"os/exec"
	"syscall"
)

const (
	esContinuous     = 0x80000000
	esSystemRequired = 0x00000001
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	setThreadExecutionState = kernel32.NewProc("SetThreadExecutionState")
)

// setAwake 呼叫 SetThreadExecutionState 防止系統睡眠。
func setAwake(on bool) {
	var flags uintptr
	if on {
		flags = esContinuous | esSystemRequired
	} else {
		flags = esContinuous
	}
	ret, _, _ := setThreadExecutionState.Call(flags)
	if ret == 0 {
		log.Println("[警告] SetThreadExecutionState 呼叫失敗")
	}
}

// PowerProfile 描述一個電源設定檔。
type PowerProfile struct {
	Name      string
	ACMonitor string // AC 螢幕逾時（分鐘）
	Desc      string
}

// powerProfiles 兩個電源設定檔。
var powerProfiles = map[int]PowerProfile{
	1: {
		Name:      "AC螢幕15分鐘",
		ACMonitor: "15",
		Desc:      "AC：螢幕15分鐘/永不睡眠/蓋上不動作　DC：螢幕15分鐘/睡眠30分鐘/蓋上睡眠",
	},
	2: {
		Name:      "AC螢幕60分鐘",
		ACMonitor: "60",
		Desc:      "AC：螢幕60分鐘/永不睡眠/蓋上不動作　DC：螢幕15分鐘/睡眠30分鐘/蓋上睡眠",
	},
}

const (
	guidButtonLid = "4f971e89-eebd-4455-a8de-9e59040e7347"
	guidLidAction = "5ca83367-6e45-459f-a27b-476b1d01c936"
)

// setPowerSettings 執行 powercfg 命令套用指定電源設定檔。
func setPowerSettings(profileID int) {
	profile, ok := powerProfiles[profileID]
	if !ok {
		log.Printf("[警告] 未知電源設定檔 ID: %d", profileID)
		return
	}

	cmds := [][]string{
		{"powercfg", "/change", "monitor-timeout-ac", profile.ACMonitor},
		{"powercfg", "/change", "standby-timeout-ac", "0"},
		{"powercfg", "/change", "monitor-timeout-dc", "15"},
		{"powercfg", "/change", "standby-timeout-dc", "30"},
		{"powercfg", "/setacvalueindex", "SCHEME_CURRENT", guidButtonLid, guidLidAction, "0"},
		{"powercfg", "/setdcvalueindex", "SCHEME_CURRENT", guidButtonLid, guidLidAction, "1"},
		{"powercfg", "/setactive", "SCHEME_CURRENT"},
	}

	failed := 0
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := cmd.Run(); err != nil {
			log.Printf("[警告] 命令失敗 %v: %v", args, err)
			failed++
		}
	}

	if failed > 0 {
		log.Printf("[警告] 部分電源設定命令失敗（%d/%d）", failed, len(cmds))
	} else {
		log.Printf("[資訊] 已套用電源設定檔 [%s]：%s", profile.Name, profile.Desc)
	}
}
