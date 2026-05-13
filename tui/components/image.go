package components

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"path"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/draw"
)

// renderImage loads an image from the filesystem and returns it rendered
// as colored half-block characters. maxCols is the maximum display width in
// terminal columns (already accounting for glamour's outer margin).
func renderImage(fsys fs.FS, basePath, imgPath string, maxCols int) (string, error) {
	resolved := imgPath
	if !path.IsAbs(imgPath) {
		resolved = path.Join(basePath, imgPath)
	}

	f, err := fsys.Open(resolved)
	if err != nil {
		return "", err
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	// Leave some room: use at most 3/4 of available width for the image.
	imgMaxW := maxCols * 3 / 4
	if imgMaxW < 20 {
		imgMaxW = 20
	}

	// Scale to fit within imgMaxW width. Each cell = 1 col wide, 2 px tall.
	dstW := imgMaxW
	if srcW < dstW {
		dstW = srcW
	}
	dstH := srcH * dstW / srcW
	// Make height even for half-block pairing.
	if dstH%2 != 0 {
		dstH++
	}
	// Cap height at ~25 rows (50 pixels).
	const maxRows = 25
	if dstH > maxRows*2 {
		dstW = dstW * (maxRows * 2) / dstH
		dstH = maxRows * 2
	}
	if dstW < 1 {
		dstW = 1
	}
	if dstH < 2 {
		dstH = 2
	}

	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	return renderHalfBlocks(dst), nil
}

// renderHalfBlocks converts an RGBA image into terminal text using the
// half-block technique: each character cell represents 2 vertical pixels.
// The upper pixel uses the background color, the lower uses the foreground
// color with the '▄' character. Left padding matches glamour's document margin.
func renderHalfBlocks(img *image.RGBA) string {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	const pad = "  " // match glamour's 2-char document margin
	var sb strings.Builder

	sb.WriteString("\n")
	for y := 0; y < h; y += 2 {
		sb.WriteString(pad)
		for x := 0; x < w; x++ {
			top := img.RGBAAt(x, y)
			var bot color.RGBA
			if y+1 < h {
				bot = img.RGBAAt(x, y+1)
			} else {
				bot = top
			}
			// Use ▄ (lower half block): fg = bottom pixel, bg = top pixel
			fmt.Fprintf(&sb, "\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▄",
				bot.R, bot.G, bot.B,
				top.R, top.G, top.B)
		}
		sb.WriteString("\033[0m\n")
	}
	sb.WriteString("\n")

	return sb.String()
}
