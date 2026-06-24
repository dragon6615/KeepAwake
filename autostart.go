//go:build windows

package main

import (
	"log"
	"os"

	"golang.org/x/sys/windows/registry"
)

const (
	appName      = "KeepAwake"
	registryPath = `Software\Microsoft\Windows\CurrentVersion\Run`
)

// isAutostartEnabled 檢查登錄檔是否已設定開機自動啟動。
func isAutostartEnabled() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	_, _, err = key.GetStringValue(appName)
	return err == nil
}

// setAutostart 設定或取消開機自動啟動。
func setAutostart(enable bool) {
	key, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		log.Printf("[警告] 開啟登錄檔失敗: %v", err)
		return
	}
	defer key.Close()

	if enable {
		exePath, err := os.Executable()
		if err != nil {
			log.Printf("[警告] 取得執行檔路徑失敗: %v", err)
			return
		}
		if err := key.SetStringValue(appName, exePath); err != nil {
			log.Printf("[警告] 設定開機啟動失敗: %v", err)
			return
		}
		log.Println("[資訊] 已啟用開機自動啟動")
	} else {
		if err := key.DeleteValue(appName); err != nil && err != registry.ErrNotExist {
			log.Printf("[警告] 刪除開機啟動失敗: %v", err)
			return
		}
		log.Println("[資訊] 已停用開機自動啟動")
	}
}
