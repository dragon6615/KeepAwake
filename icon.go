//go:build windows

package main

import _ "embed"

// iconData 內嵌由 tools/genicon/main.go 產生的 app.ico（💡 燈泡 emoji）。
// 重新產生方式：go run tools/genicon/main.go
//
//go:embed app.ico
var iconData []byte
