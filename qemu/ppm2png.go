//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("screenshot.ppm")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	magic, _ := reader.ReadString('\n')
	if strings.TrimSpace(magic) != "P6" {
		fmt.Printf("Invalid magic: %s\n", magic)
		return
	}

	// Read dimensions (skipping comments)
	readNonComment := func() string {
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				return line
			}
		}
	}

	dims := strings.Fields(readNonComment())
	width, _ := strconv.Atoi(dims[0])
	height, _ := strconv.Atoi(dims[1])
	maxVal, _ := strconv.Atoi(readNonComment())
	_ = maxVal

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, _ := reader.ReadByte()
			g, _ := reader.ReadByte()
			b, _ := reader.ReadByte()
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	out, _ := os.Create("screenshot.png")
	defer out.Close()
	png.Encode(out, img)
	fmt.Printf("Successfully converted PPM to screenshot.png (%dx%d)\n", width, height)
}
