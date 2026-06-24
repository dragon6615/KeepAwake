# KeepAwake（保持電腦清醒）

一個常駐 Windows 系統匣（system tray）的輕量防睡眠小工具，用 Go 撰寫。
啟動後會讓電腦保持清醒、不進入睡眠，並可一鍵切換電源設定檔、設定開機自動啟動。

![圖示](icon.png)

---

## 功能特色

- 🛡️ **防止系統睡眠** — 呼叫 Windows API `SetThreadExecutionState` 讓系統保持喚醒（注意：不阻止螢幕關閉，由電源設定檔控制）。
- ⚡ **兩種電源設定檔** — 透過 `powercfg` 一鍵切換 AC 螢幕逾時等設定。
- 🚀 **開機自動啟動** — 透過登錄檔 `HKCU\...\Run` 設定，可隨時開關。
- 🔔 **系統匣常駐** — 不佔工作列，所有操作集中在右下角圖示選單。
- 🟢 **啟動即生效** — 程式開啟時預設立即啟用防睡眠。

---

## 系統需求

- Windows（程式以 `//go:build windows` 標記，僅支援 Windows）
- 編譯需要 Go 1.26+

> `powercfg` 部分指令可能需要系統管理員權限才能完整套用；若部分指令失敗，主控台會輸出 `[警告]` 訊息，但防睡眠核心功能仍正常運作。

---

## 安裝與執行

### 直接執行

編譯後的執行檔，雙擊即可：

```
KeepAwake.exe
```

執行後會出現在系統匣（右下角），預設**立即啟用防睡眠**。

### 從原始碼編譯

```powershell
# 取得相依套件
go mod download

# 編譯（-H windowsgui 隱藏主控台視窗）
go build -ldflags "-H windowsgui" -o KeepAwake.exe
```

> `rsrc.syso` 為已產生的 Windows 資源檔，內含 exe 圖示，編譯時會自動嵌入。

---

## 使用說明

在系統匣圖示上**點擊**，會展開選單：

| 選單項目 | 說明 |
| --- | --- |
| 保持電腦清醒 v1.2 | 標題，顯示版本（不可點擊） |
| **電源設定檔** | 子選單，切換兩種設定檔（見下方） |
| ○ / ● 狀態列 | 顯示目前防睡眠狀態（不可點擊） |
| **切換防睡眠** | 啟動 / 停止防睡眠功能 |
| **開機自動啟動** | 勾選代表已啟用開機自啟 |
| **退出** | 停止防睡眠並結束程式 |

### 電源設定檔

| 設定檔 | AC（接電源） | DC（電池） |
| --- | --- | --- |
| **設定檔 1：AC螢幕15分鐘**（預設） | 螢幕 15 分鐘 / 永不睡眠 / 蓋上不動作 | 螢幕 15 分鐘 / 睡眠 30 分鐘 / 蓋上睡眠 |
| **設定檔 2：AC螢幕60分鐘** | 螢幕 60 分鐘 / 永不睡眠 / 蓋上不動作 | 螢幕 15 分鐘 / 睡眠 30 分鐘 / 蓋上睡眠 |

切換設定檔時，若防睡眠正在運行，會立即套用新設定。

---

## 程式架構

```
KeepAwake/
├── main.go            # 進入點、系統匣 UI、狀態管理、選單事件迴圈
├── power.go           # 核心防睡眠（SetThreadExecutionState）與 powercfg 電源設定檔
├── autostart.go       # 登錄檔開機自動啟動（HKCU\...\Run）
├── icon.go            # 以 //go:embed 內嵌 app.ico
├── app.ico            # 系統匣圖示（💡 燈泡）
├── rsrc.syso          # 編譯時嵌入的 Windows 資源（exe 圖示）
├── go.mod / go.sum    # 相依套件
└── tools/
    └── genicon/
        └── main.go    # 工具：將 PNG / Twemoji 轉為 64×64 app.ico
```

### 模組職責

- **main.go** — 程式進入點。使用 [`getlantern/systray`](https://github.com/getlantern/systray) 建立系統匣選單，以執行緒安全的 `AppState`（`sync.Mutex`）保存執行狀態與目前設定檔，並在 goroutine 中監聽選單點擊事件。
- **power.go** — 防睡眠核心。`setAwake()` 直接呼叫 `kernel32.dll` 的 `SetThreadExecutionState`（`ES_CONTINUOUS | ES_SYSTEM_REQUIRED`）；`setPowerSettings()` 依設定檔執行一系列 `powercfg` 指令（螢幕/睡眠逾時、蓋上動作），執行時隱藏命令列視窗。
- **autostart.go** — 讀寫 `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run` 下的 `KeepAwake` 值，以執行檔絕對路徑設定開機自啟。
- **icon.go** — 透過 `//go:embed` 將 `app.ico` 編入執行檔。

### 相依套件

- [`github.com/getlantern/systray`](https://github.com/getlantern/systray) — 跨平台系統匣選單
- [`golang.org/x/sys/windows/registry`](https://pkg.go.dev/golang.org/x/sys/windows/registry) — Windows 登錄檔存取

---

## 重新產生圖示

`tools/genicon` 是一個獨立工具（`//go:build ignore`），可將 PNG 或 Twemoji 燈泡轉成 64×64 的 `app.ico`：

```powershell
# 從 Twemoji CDN 下載 💡 燈泡
go run tools/genicon/main.go

# 或使用本地 PNG
go run tools/genicon/main.go path\to\icon.png
```

> Twemoji 圖形授權：CC BY 4.0

---

## 運作原理

1. **防睡眠**：`SetThreadExecutionState` 告知 Windows「目前有應用程式需要系統保持運作」，系統因而不會自動進入睡眠。停止時清除 `ES_SYSTEM_REQUIRED` 旗標即恢復正常。
2. **電源設定檔**：以 `powercfg /change` 與 `/setacvalueindex`、`/setdcvalueindex` 調整螢幕逾時、睡眠逾時及蓋上筆電動作，最後 `/setactive SCHEME_CURRENT` 套用。
3. **開機自啟**：在登錄檔 `Run` 鍵寫入執行檔路徑，登入時由 Windows 自動啟動。
