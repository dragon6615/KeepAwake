# KeepAwake

Windows 系統匣（system tray）小工具，用 Go 寫成，功能是讓電腦保持清醒（不進入睡眠）。

## Build

```bash
go build -ldflags "-H windowsgui" -o KeepAwake.exe
```

`-H windowsgui` 用來抑制黑色 console 視窗（tray app 必須）。

## 發版流程（Release）

採用**手動打 tag 觸發**的方式，刻意不做「改版號自動發版」，以保持可控、避免誤觸發。

`.github/workflows/release.yml` 會在推送 `v*` tag 時自動 build 並發布 GitHub Release（附上可下載的 `KeepAwake.exe`）。

發新版時的完整步驟：

1. 更新 `main.go` 裡的 `const appVersion`（例如改成 `v1.3`）
2. commit + push
3. 打 tag 並推送，版號需與 `appVersion` 對齊：
   ```bash
   git tag v1.3
   git push origin v1.3
   ```

> 給 Claude：當使用者說「發 vX.Y」/「release vX.Y」時，請直接完成上述整個流程。

## 其他 CI

`.github/workflows/build.yml`：每次 push / PR 到 `main` 會自動編譯驗證，並上傳 `KeepAwake.exe` 作為 artifact（暫存）。
