//go:build ignore

// genicon: 將 PNG 轉換為 64x64 app.ico。
//
// 用法：
//   go run tools/genicon/main.go                     # 從 Twemoji CDN 下載 💡
//   go run tools/genicon/main.go path\to\icon.png    # 使用本地 PNG 檔案
//
// Twemoji 圖形授權：CC BY 4.0 (https://creativecommons.org/licenses/by/4.0/)
package main

import (
"bytes"
"encoding/binary"
"fmt"
"image"
"image/color"
"image/png"
"net/http"
"os"
)

const (
emojiURL   = "https://cdn.jsdelivr.net/gh/twitter/twemoji@latest/assets/72x72/1f4a1.png"
targetSize = 64
)

func main() {
var src image.Image

if len(os.Args) >= 2 {
// 使用本地 PNG 檔案
f, err := os.Open(os.Args[1])
if err != nil {
fatalf("開啟 PNG 失敗: %v", err)
}
defer f.Close()
src, err = png.Decode(f)
if err != nil {
fatalf("PNG 解碼失敗: %v", err)
}
fmt.Printf("使用本地檔案：%s\n", os.Args[1])
} else {
// 從 Twemoji CDN 下載
resp, err := http.Get(emojiURL)
if err != nil {
fatalf("下載 emoji 失敗: %v", err)
}
defer resp.Body.Close()
if resp.StatusCode != 200 {
fatalf("HTTP 狀態碼 %d", resp.StatusCode)
}
src, err = png.Decode(resp.Body)
if err != nil {
fatalf("PNG 解碼失敗: %v", err)
}
fmt.Println("使用 Twemoji CDN")
}

dst := scaleTo(src, targetSize, targetSize)

if err := os.WriteFile("app.ico", toICO(dst, targetSize), 0o644); err != nil {
fatalf("寫入 app.ico 失敗: %v", err)
}
fmt.Println("app.ico 已成功寫入")
}

// scaleTo 以最近鄰法將 src 縮放至 w×h。
func scaleTo(src image.Image, w, h int) *image.NRGBA {
sb := src.Bounds()
sw := sb.Dx()
sh := sb.Dy()
dst := image.NewNRGBA(image.Rect(0, 0, w, h))
for dy := 0; dy < h; dy++ {
sy := sb.Min.Y + dy*sh/h
for dx := 0; dx < w; dx++ {
sx := sb.Min.X + dx*sw/w
rv, gv, bv, av := src.At(sx, sy).RGBA()
if av == 0 {
continue
}
dst.SetNRGBA(dx, dy, color.NRGBA{
R: uint8(rv * 0xff / av),
G: uint8(gv * 0xff / av),
B: uint8(bv * 0xff / av),
A: uint8(av >> 8),
})
}
}
return dst
}

// toICO 將 *image.NRGBA 轉為 32-bit Windows ICO 格式。
func toICO(src *image.NRGBA, size int) []byte {
pix := make([]byte, size*size*4)
for y := 0; y < size; y++ {
for x := 0; x < size; x++ {
c := src.NRGBAAt(x, y)
dstY := size - 1 - y
i := (dstY*size + x) * 4
pix[i+0] = c.B
pix[i+1] = c.G
pix[i+2] = c.R
pix[i+3] = c.A
}
}

rowBytes    := ((size + 31) / 32) * 4
andMask     := make([]byte, rowBytes*size)
imgDataSize := 40 + len(pix) + len(andMask)

buf := &bytes.Buffer{}
le16(buf, 0); le16(buf, 1); le16(buf, 1)
buf.WriteByte(byte(size)); buf.WriteByte(byte(size))
buf.WriteByte(0); buf.WriteByte(0)
le16(buf, 1); le16(buf, 32)
le32(buf, uint32(imgDataSize))
le32(buf, 22)
le32(buf, 40)
lei32(buf, int32(size))
lei32(buf, int32(size*2))
le16(buf, 1); le16(buf, 32)
le32(buf, 0)
le32(buf, uint32(len(pix)))
lei32(buf, 0); lei32(buf, 0)
le32(buf, 0); le32(buf, 0)
buf.Write(pix)
buf.Write(andMask)
return buf.Bytes()
}

func le16(b *bytes.Buffer, v uint16) { binary.Write(b, binary.LittleEndian, v) }
func le32(b *bytes.Buffer, v uint32) { binary.Write(b, binary.LittleEndian, v) }
func lei32(b *bytes.Buffer, v int32) { binary.Write(b, binary.LittleEndian, v) }
func fatalf(format string, args ...any) {
fmt.Fprintf(os.Stderr, "genicon: "+format+"\n", args...)
os.Exit(1)
}
