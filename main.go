//go:build windows

package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/getlantern/systray"
)

const appVersion = "v1.2"

// AppState 保存應用程式執行狀態。
type AppState struct {
	mu           sync.Mutex
	running      bool
	powerProfile int // 1 or 2
}

func (s *AppState) isRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *AppState) getPowerProfile() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.powerProfile
}

var state = &AppState{powerProfile: 1}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("KeepAwake")
	systray.SetTooltip("KeepAwake - 保持電腦清醒")

	// 版本標題（不可點擊）
	mTitle := systray.AddMenuItem(fmt.Sprintf("保持電腦清醒 %s", appVersion), "")
	mTitle.Disable()

	systray.AddSeparator()

	// 電源設定檔子選單
	mProfileMenu := systray.AddMenuItem("電源設定檔", "切換電源設定檔")
	mProfile1 := mProfileMenu.AddSubMenuItem(powerProfiles[1].Name, powerProfiles[1].Desc)
	mProfile2 := mProfileMenu.AddSubMenuItem(powerProfiles[2].Name, powerProfiles[2].Desc)
	mProfile1.Check() // 預設選第 1 個

	systray.AddSeparator()

	// 狀態顯示（不可點擊）
	mStatus := systray.AddMenuItem("○ 防睡眠已停止", "")
	mStatus.Disable()

	// 切換防睡眠
	mToggle := systray.AddMenuItem("切換防睡眠", "啟動或停止防睡眠功能")

	systray.AddSeparator()

	// 開機自動啟動
	mAutostart := systray.AddMenuItem("開機自動啟動", "設定或取消開機自動啟動")
	if isAutostartEnabled() {
		mAutostart.Check()
	}

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("退出", "退出程式")

	// 啟動時預設開啟防睡眠
	startNosleep(mStatus)

	// 監聽選單事件
	go func() {
		for {
			select {
			case <-mToggle.ClickedCh:
				if state.isRunning() {
					stopNosleep(mStatus)
				} else {
					startNosleep(mStatus)
				}

			case <-mProfile1.ClickedCh:
				switchProfile(1, mProfile1, mProfile2)

			case <-mProfile2.ClickedCh:
				switchProfile(2, mProfile1, mProfile2)

			case <-mAutostart.ClickedCh:
				current := isAutostartEnabled()
				setAutostart(!current)
				if isAutostartEnabled() {
					mAutostart.Check()
				} else {
					mAutostart.Uncheck()
				}

			case <-mQuit.ClickedCh:
				stopNosleep(mStatus)
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	log.Println("[資訊] 程式退出")
}

func startNosleep(mStatus *systray.MenuItem) {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.running {
		return
	}
	setPowerSettings(state.powerProfile)
	setAwake(true)
	state.running = true
	mStatus.SetTitle("● 防睡眠運行中")
	log.Println("[資訊] 防睡眠已啟動")
}

func stopNosleep(mStatus *systray.MenuItem) {
	state.mu.Lock()
	defer state.mu.Unlock()
	if !state.running {
		return
	}
	setAwake(false)
	state.running = false
	mStatus.SetTitle("○ 防睡眠已停止")
	log.Println("[資訊] 防睡眠已停止")
}

func switchProfile(id int, m1, m2 *systray.MenuItem) {
	state.mu.Lock()
	state.powerProfile = id
	running := state.running
	state.mu.Unlock()

	if id == 1 {
		m1.Check()
		m2.Uncheck()
	} else {
		m1.Uncheck()
		m2.Check()
	}

	if running {
		setPowerSettings(id)
	}
	log.Printf("[資訊] 已切換電源設定檔至：%s", powerProfiles[id].Name)
}
