package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"../font"
	//"../../imutil"
)

var palette2 = []color.Color{
	color.White,
	color.Black,
}

var palette1 = []color.Color{
	color.Black,
	color.White,
	color.RGBA{255, 255, 255, 255},
	color.RGBA{250, 250, 250, 250},
	color.RGBA{245, 245, 245, 245},
	color.RGBA{240, 240, 240, 240},
	color.RGBA{235, 235, 235, 235},
	color.RGBA{230, 230, 230, 230},
	color.RGBA{225, 225, 225, 225},
	color.RGBA{220, 220, 220, 220},
	color.RGBA{215, 215, 215, 215},
	color.RGBA{210, 210, 210, 210},
	color.RGBA{205, 205, 205, 205},
	color.RGBA{200, 200, 200, 200},
	color.RGBA{195, 195, 195, 195},
	color.RGBA{190, 190, 190, 190},
	color.RGBA{185, 185, 185, 185},
	color.RGBA{180, 180, 180, 180},
	color.RGBA{175, 175, 175, 175},
	color.RGBA{170, 170, 170, 170},
	color.RGBA{165, 165, 165, 165},
	color.RGBA{160, 160, 160, 160},
	color.RGBA{155, 155, 155, 155},
	color.RGBA{150, 150, 150, 150},
	color.RGBA{145, 145, 145, 145},
	color.RGBA{140, 140, 140, 140},
	color.RGBA{135, 135, 135, 135},
	color.RGBA{130, 130, 130, 130},
	color.RGBA{125, 125, 125, 125},
	color.RGBA{120, 120, 120, 120},
	color.RGBA{115, 115, 115, 115},
	color.RGBA{110, 110, 110, 110},
	color.RGBA{105, 105, 105, 105},
	color.RGBA{100, 100, 100, 100},
	color.RGBA{95, 95, 95, 95},
	color.RGBA{90, 90, 90, 90},
	color.RGBA{85, 85, 85, 85},
	color.RGBA{80, 80, 80, 80},
	color.RGBA{75, 75, 75, 75},
	color.RGBA{70, 70, 70, 70},
	color.RGBA{65, 65, 65, 65},
	color.RGBA{60, 60, 60, 60},
	color.RGBA{55, 55, 55, 55},
	color.RGBA{50, 50, 50, 50},
	color.RGBA{45, 45, 45, 45},
	color.RGBA{40, 40, 40, 40},
	color.RGBA{35, 35, 35, 35},
	color.RGBA{30, 30, 30, 30},
	color.RGBA{25, 25, 25, 25},
	color.RGBA{20, 20, 20, 20},
	color.RGBA{15, 15, 15, 15},
	color.RGBA{10, 10, 10, 10},
	color.RGBA{5, 5, 5, 5},
	color.RGBA{0, 0, 0, 0},
}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func fileDownloadAndCreateGIF(choice string, start int, end int) {
	argsWithoutProg := os.Args[1:]
	s3link := ""
	outFile := ""
	inFile := ""
	if choice == "slide" {
		s3link = "https://uber-static.s3.amazonaws.com/beacon-gif/33_PICKING_UP_TEST.gif"
		outFile = "outfileSlide.gif"
		inFile = "04_PICKING_UP_TEST.gif"
	} else if choice == "slidemeet" {
		s3link = "https://uber-static.s3.amazonaws.com/beacon-gif/09_MeetAlex.gif"
		outFile = "outfileSlideMeet.gif"
		inFile = "09_MeetAlex.gif"
	}
	res, err := http.Get(s3link)
	if err != nil {
		fmt.Println(err)
	} else {
		fileBytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
		g, err := gif.DecodeAll(bytes.NewReader(fileBytes))
		if err != nil {
			fmt.Println(err)
		}
		arrFadeOut := make([]*image.Paletted, len(g.Image)-end)
		arrFadeDelay := make([]int, (len(g.Image) - end))
		fadeOut := g.Image[len(g.Image)-(len(g.Image)-end):]
		fadeOutDelay := g.Delay[len(g.Delay)-(len(g.Image)-end):]
		copy(arrFadeOut, fadeOut)
		copy(arrFadeDelay, fadeOutDelay)
		if g.Disposal == nil {
			fmt.Println("disposal is nil")
		}
		f, err := font.NewFromFile("../fonts/uber-text-6x6.json")
		if err != nil {
			log.Fatal(err)
		}
		f.Scale(1)
		pt := image.ZP
		x, _ := strconv.Atoi(argsWithoutProg[0])
		y, _ := strconv.Atoi(argsWithoutProg[1])
		pt.X = x
		pt.Y = y
		pt.X = 25
		text := ""
		for i := 3; i < len(argsWithoutProg); i++ {
			text += argsWithoutProg[i]
			//text += " "
		}
		slideAcrossGIF(text, pt, g, f, arrFadeOut, arrFadeDelay, start, end, outFile)
		return
	}
	lol, _ := os.Open(inFile)
	r := bufio.NewReader(lol)
	g, err := gif.DecodeAll(r)
	arrFadeOut := make([]*image.Paletted, len(g.Image)-end)
	arrFadeDelay := make([]int, (len(g.Image) - end))
	fadeOut := g.Image[len(g.Image)-(len(g.Image)-end):]
	fadeOutDelay := g.Delay[len(g.Delay)-(len(g.Image)-end):]
	copy(arrFadeOut, fadeOut)
	copy(arrFadeDelay, fadeOutDelay)
	if g.Disposal == nil {
		fmt.Println("disposal is nil")
	}
	lol.Close()
	f, err := font.NewFromFile("../fonts/uber-text-6x6.json")
	if err != nil {
		log.Fatal(err)
	}
	f.Scale(1)
	pt := image.ZP
	x, _ := strconv.Atoi(argsWithoutProg[0])
	y, _ := strconv.Atoi(argsWithoutProg[1])
	pt.X = x
	pt.Y = y
	pt.X = 25
	text := ""
	for i := 3; i < len(argsWithoutProg); i++ {
		text += argsWithoutProg[i]
		//text += " "
	}
	slideAcrossGIF(text, pt, g, f, arrFadeOut, arrFadeDelay, start, end, outFile)
	return
}

func main() {
	start := time.Now()
	argsWithoutProg := os.Args[1:]
	if argsWithoutProg[2] == "slide" {
		fileDownloadAndCreateGIF("slide", 99, 120)
		elapsed := time.Since(start)
		fmt.Printf("Time elapsed: %s\n", elapsed)
		fmt.Println("returns here")
		return
	} else if argsWithoutProg[2] == "slidemeet" {
		fileDownloadAndCreateGIF("slidemeet", 8, 93)
		elapsed := time.Since(start)
		fmt.Printf("Time elapsed: %s\n", elapsed)
		fmt.Println("returns here")
		return
	}
}

func wipeUp(text string, pt image.Point, testAnim *gif.GIF, f *font.Font, start int) {
	var colWhite = color.RGBA{50, 50, 50, 255}
	var colBlack = color.RGBA{0, 0, 0, 255}
	for i := 0; i < 10; i++ {
		m := image.NewRGBA(image.Rect(0, 0, 25, 11))
		f.DrawText(m, image.Point{X: 0, Y: -1}, "Meet", color.White)
		if i > 3 {
			f.DrawTextWipeUp(m, pt, colBlack, colWhite, text)
		} else {
			f.DrawTextWipeUp(m, pt, color.Black, colWhite, text)
		}
		r := int(colBlack.R)
		if r+30 > 255 {
			colBlack.R = 255
		} else {
			colBlack.R += 30
		}
		b := int(colBlack.B)
		if b+30 > 255 {
			colBlack.B = 255
		} else {
			colBlack.B += 30
		}
		g := int(colBlack.G)
		if g+30 > 255 {
			colBlack.G = 255
		} else {
			colBlack.G += 30
		}
		wr := int(colWhite.R)
		if wr+25 > 255 {
			colWhite.R = 255
		} else {
			colWhite.R += 25
		}
		wb := int(colWhite.B)
		if wb+25 > 255 {
			colWhite.B = 255
		} else {
			colWhite.B += 25
		}
		wg := int(colWhite.G)
		if wg+25 > 255 {
			colWhite.G = 255
		} else {
			colWhite.G += 25
		}
		var iPalette []color.Color
		for i := 0; i < 25; i++ {
			for j := 0; j < 11; j++ {
				iPalette = AppendIfMissing(iPalette, m.At(i, j))
			}
		}
		z := image.NewPaletted(image.Rect(0, 0, 25, 11), iPalette)
		draw.Draw(z, z.Bounds(), m, image.ZP, draw.Src)
		testAnim.Image[start] = z
		testAnim.Delay[start] = 3
		start++
	}
}

func remove(s []color.Color, i int) []color.Color {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func staticTenFrames(text string, pt image.Point, testAnim *gif.GIF, f *font.Font, start int) {
	for i := 0; i < 10; i++ {
		m := image.NewRGBA(image.Rect(0, 0, 25, 11))
		f.DrawText(m, image.Point{X: 0, Y: -1}, "Meet", color.White)
		f.DrawText(m, pt, text, color.White)
		z := image.NewPaletted(image.Rect(0, 0, 25, 11), palette2)
		draw.Draw(z, z.Bounds(), m, image.ZP, draw.Src)
		testAnim.Image[start] = z
		testAnim.Delay[start] = 3
		start++
	}
}

func wipeDown(text string, pt image.Point, testAnim *gif.GIF, f *font.Font, start int) {
	colWhite := color.RGBA{255, 255, 255, 255}
	colBlack := color.RGBA{200, 200, 200, 255}
	for i := 0; i < 10; i++ {
		m := image.NewRGBA(image.Rect(0, 0, 25, 11))
		f.DrawTextWipeDown(m, pt, colWhite, colBlack, text)
		r := int(colBlack.R)
		if r-30 < 0 {
			colBlack.R = 0
		} else {
			colBlack.R -= 30
		}
		b := int(colBlack.B)
		if b-30 < 0 {
			colBlack.B = 0
		} else {
			colBlack.B -= 30
		}
		g := int(colBlack.G)
		if g-30 < 0 {
			colBlack.G = 0
		} else {
			colBlack.G -= 30
		}
		wr := int(colWhite.R)
		if wr-25 < 0 {
			colWhite.R = 0
		} else {
			colWhite.R -= 30
		}
		wb := int(colWhite.B)
		if wb-25 < 0 {
			colWhite.B = 0
		} else {
			colWhite.B -= 30
		}
		wg := int(colWhite.G)
		if wg-25 < 0 {
			colWhite.G = 0
		} else {
			colWhite.G -= 30
		}
		var iPalette []color.Color
		for i := 0; i < 25; i++ {
			for j := 0; j < 11; j++ {
				iPalette = AppendIfMissing(iPalette, m.At(i, j))
			}
		}
		z := image.NewPaletted(image.Rect(0, 0, 25, 11), iPalette)
		draw.Draw(z, z.Bounds(), m, image.ZP, draw.Src)
		testAnim.Image[start] = z
		testAnim.Delay[start] = 3
		start++
	}
}

func AppendIfMissing(slice []color.Color, i color.Color) []color.Color {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func wipeUpGIF(text string, pt image.Point, testAnim *gif.GIF, f *font.Font) *gif.GIF {
	wipeUp(text, pt, testAnim, f, 99)
	staticTenFrames(text, pt, testAnim, f, 109)
	wipeDown(text, pt, testAnim, f, 129)
	fmt.Printf("size of config testAnim: %d\n", unsafe.Sizeof(testAnim.Config.ColorModel))
	fmt.Printf("config width x height: %d x %d \n", testAnim.Config.Width, testAnim.Config.Height)
	for i := 0; i < len(testAnim.Image); i++ {
		fmt.Printf("size of image testAnim: %d\n", unsafe.Sizeof(testAnim.Image[i]))
	}
	if testAnim.Disposal == nil || len(testAnim.Disposal) == 0 {
		fmt.Println("disposal of testAnim is null")
	} else {
		testAnim.Disposal = nil
	}
	gifFile, err := os.OpenFile("outfile123456.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("failed to open GIF file")
	}
	err = gif.EncodeAll(gifFile, testAnim)
	if err != nil {
		fmt.Println("encoding failed")
		fmt.Println(err)
	}
	defer gifFile.Close()
	imgWidth, imgHeight := getGifDimensions(testAnim)
	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), testAnim.Image[0], image.ZP, draw.Src)
	for i, srcImg := range testAnim.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.ZP, draw.Over)

		// save current frame "stack". This will overwrite an existing file with that name
		file, err := os.Create(fmt.Sprintf("%s%d%s", "outputImages/outImageTEST", i, ".png"))
		if err != nil {
			return nil
		}
		//err = jpeg.Encode(file, overpaintImage, nil)
		err = png.Encode(file, overpaintImage)
		if err != nil {
			return nil
		}

		file.Close()
	}
	return testAnim
}

func slideAcrossGIF(text string, pt image.Point, testAnim *gif.GIF, f *font.Font, fadeOut []*image.Paletted, fadeDelay []int, start int, end int, file string) *gif.GIF {
	slideAcross(text, pt, testAnim, f, start, end)
	for i := 0; i < len(fadeOut); i++ {
		testAnim.Image = append(testAnim.Image, fadeOut[i])
	}
	for i := 0; i < len(fadeDelay); i++ {
		testAnim.Delay = append(testAnim.Delay, fadeDelay[i])
	}
	testAnim.Disposal = nil
	//testAnim.Config.ColorModel = nil
	if len(testAnim.Image) != len(testAnim.Delay) {
		fmt.Println("lengths not equal")
	}
	gifFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("failed to open GIF file")
	}
	err = gif.EncodeAll(gifFile, testAnim)
	if err != nil {
		fmt.Println("encoding failed")
		fmt.Println(err)
	}
	defer gifFile.Close()
	return testAnim
}

func slideAcross(text string, pt image.Point, testAnim *gif.GIF, f *font.Font, start int, end int) {
	//p1 := testAnim.Image[0].Palette
	var whiteI int
	/*for i := 0; i < len(testAnim.Config.ColorModel.(color.Palette)); i++ {
		r, g, b, a := testAnim.Config.ColorModel.(color.Palette)[i].RGBA()
		fmt.Printf("r: %d, g: %d, b: %d, a: %d\n", r, g, b, a)
	}*/
	for i := 0; i < len(testAnim.Config.ColorModel.(color.Palette)); i++ {
		r, g, b, a := testAnim.Config.ColorModel.(color.Palette)[i].RGBA()
		if r == 65535 && b == 65535 && g == 65535 && a == 65535 {
			whiteI = i
			break
		}
	}
	pt = image.Point{X: 26, Y: 5}
	r := f.GetUberTextBounds(text)
	dx := r.Dx()
	first := testAnim.Image[0]
	b := first.Bounds()
	third := image.NewPaletted(image.Rect(0, 0, 25, 11), testAnim.Config.ColorModel.(color.Palette)) //image1
	for i := 1; i < len(testAnim.Image) && i < start; i++ {
		image1 := testAnim.Image[i-1]
		second := testAnim.Image[i]
		b = image1.Bounds()
		draw.Draw(third, b, image1, image.ZP, draw.Src)
		draw.Draw(third, second.Bounds(), second, image.ZP, draw.Over)
		/*file := fmt.Sprintf("%s%d%s", "../outputImages/hello-world", i, ".png")
		if err := imutil.ImageSaveToPNG(file, third); err != nil {
			log.Fatal(err)
		}*/
	}
	for i := 0; ; i++ {
		m := image.NewPaletted(r, testAnim.Config.ColorModel.(color.Palette))                        //image2
		z := image.NewPaletted(image.Rect(0, 0, 25, 11), testAnim.Config.ColorModel.(color.Palette)) //image3
		b := third.Bounds()
		f.DrawText(m, image.ZP, text, testAnim.Config.ColorModel.(color.Palette)[whiteI])
		draw.Draw(z, b, third, image.ZP, draw.Src)
		draw.Draw(z, m.Bounds().Add(pt), m, image.ZP, draw.Over)
		if i%2 == 0 {
			pt.X--
		}
		if start >= len(testAnim.Image) {
			testAnim.Image = append(testAnim.Image, z)
			testAnim.Delay = append(testAnim.Delay, 4)
		} else {
			testAnim.Image[start] = z
			testAnim.Delay[start] = 4
		}
		start++
		if pt.X < (-1 * dx) {
			break
		}
	}
	if start > end {
		testAnim.Image = testAnim.Image[0:start]
		testAnim.Delay = testAnim.Delay[0:start]
	}
	testAnim.Disposal = nil
}

func generatePairingFinalGIF(text string, pt image.Point, anim *gif.GIF, fo *font.Font) *gif.GIF {
	hiString := "HI"
	r := []rune(hiString)
	hString := strings.SplitAfter(hiString, "")[0]
	iString := strings.SplitAfter(hiString, "")[1]
	h := r[0]
	x := fo.GetSizeMap()[h].Dx
	pz := image.ZP
	pz.Y = -1
	var colWhite = color.RGBA{255, 255, 255, 255}
	var colBlack = color.RGBA{0, 0, 0, 255}
	for j := 0; j < 5; j++ {
		m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
		fo.DrawTextFade(m, pz, colBlack, hString)
		colBlack.R += (255 / 5)
		colBlack.B += (255 / 5)
		colBlack.G += (255 / 5)
		colBlack.A += (255 / 5)
		anim.Image = append(anim.Image, m)
		anim.Delay = append(anim.Delay, 3)
	}
	pz.X += x
	colBlack = color.RGBA{0, 0, 0, 255}
	pzH := image.ZP
	pzH.Y = -1
	for j := 0; j < 5; j++ {
		m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
		fo.DrawText(m, pzH, hString, color.White)
		fo.DrawTextFade(m, pz, colBlack, iString)
		colBlack.R += (255 / 5)
		colBlack.B += (255 / 5)
		colBlack.G += (255 / 5)
		colBlack.A += (255 / 5)
		anim.Image = append(anim.Image, m)
		anim.Delay = append(anim.Delay, 3)
	}
	ptLetter := &image.Point{
		X: 0,
		Y: 5,
	}
	ptHI := &image.Point{
		X: 0,
		Y: -1,
	}
	ptArray := make([]image.Point, len(text))
	counter := 0
	ptArray[0] = *ptLetter
	textArray := strings.SplitAfter(text, "")
	colBlack = color.RGBA{0, 0, 0, 255}
	for j := 0; j < len(text); j++ {
		ri := []rune(text)
		letter := ri[counter]
		letterX := fo.GetSizeMap()[letter].Dx
		for k := 0; k < 5; k++ {
			m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
			draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
			fo.DrawText(m, *ptHI, hiString, color.White)
			if j >= 1 {
				substring := text[0:j]
				fo.DrawText(m, ptArray[0], substring, color.White)
			}
			fo.DrawTextFade(m, *ptLetter, colBlack, textArray[counter])
			colBlack.R += (255 / 5)
			colBlack.B += (255 / 5)
			colBlack.G += (255 / 5)
			colBlack.A += (255 / 5)
			anim.Image = append(anim.Image, m)
			anim.Delay = append(anim.Delay, 3)
		}
		colBlack = color.RGBA{0, 0, 0, 255}
		ptArray[counter] = *ptLetter
		ptLetter.X += letterX
		counter++
	}
	ptLetter = &image.Point{
		X: 0,
		Y: 5,
	}
	for j := 0; j < 20; j++ {
		m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
		fo.DrawText(m, *ptHI, hiString, color.White)
		fo.DrawText(m, *ptLetter, text, color.White)
		anim.Image = append(anim.Image, m)
		anim.Delay = append(anim.Delay, 3)
	}
	pzH = image.ZP
	pzH.Y = -1
	pzH.X += x
	pz = image.ZP
	pz.Y = -1
	for j := 0; j < 5; j++ {
		m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
		fo.DrawTextFade(m, pz, colWhite, hString)
		colWhite.R -= (255 / 5)
		colWhite.B -= (255 / 5)
		colWhite.G -= (255 / 5)
		colWhite.A -= (255 / 5)
		fo.DrawText(m, pzH, iString, color.White)
		fo.DrawText(m, *ptLetter, text, color.White)
		anim.Image = append(anim.Image, m)
		anim.Delay = append(anim.Delay, 3)
	}
	for j := 0; j < 1; j++ {
		m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
		fo.DrawText(m, pzH, iString, color.White)
		fo.DrawText(m, *ptLetter, text, color.White)
		anim.Image = append(anim.Image, m)
		anim.Delay = append(anim.Delay, 3)
	}
	colWhite = color.RGBA{255, 255, 255, 255}
	for j := 0; j < 5; j++ {
		m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
		fo.DrawTextFade(m, pzH, colWhite, iString)
		fo.DrawText(m, *ptLetter, text, color.White)
		colWhite.R -= (255 / 5)
		colWhite.B -= (255 / 5)
		colWhite.G -= (255 / 5)
		colWhite.A -= (255 / 5)
		anim.Image = append(anim.Image, m)
		anim.Delay = append(anim.Delay, 3)
	}
	for j := 0; j < 1; j++ {
		m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
		fo.DrawText(m, *ptLetter, text, color.White)
		anim.Image = append(anim.Image, m)
		anim.Delay = append(anim.Delay, 3)
	}
	counter = 0
	for j := 0; j < len(text); j++ {
		colWhite = color.RGBA{255, 255, 255, 255}
		ri := []rune(text)
		letter := ri[counter]
		letterX := fo.GetSizeMap()[letter].Dx
		for k := 0; k < 5; k++ {
			m := image.NewPaletted(image.Rect(0, 0, 25, 11), palette1)
			draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Over)
			if j < len(text)-1 {
				substring := text[j+1 : len(text)]
				subPoint := &image.Point{
					X: ptLetter.X + letterX,
					Y: ptLetter.Y,
				}
				fo.DrawText(m, *subPoint, substring, color.White)
			}
			fo.DrawTextFade(m, *ptLetter, colWhite, textArray[counter])
			colWhite.R -= (255 / 5)
			colWhite.B -= (255 / 5)
			colWhite.G -= (255 / 5)
			colWhite.A -= (255 / 5)
			anim.Image = append(anim.Image, m)
			anim.Delay = append(anim.Delay, 3)
		}
		ptArray[counter] = *ptLetter
		ptLetter.X += letterX
		counter++
	}
	gifFile, err := os.OpenFile("outfile1234.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("failed to open GIF file")
	}
	err = gif.EncodeAll(gifFile, anim)
	if err != nil {
		fmt.Println("encoding failed")
		fmt.Println(err)
	}
	gifFile.Close()
	imgWidth, imgHeight := getGifDimensions(anim)
	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(overpaintImage, overpaintImage.Bounds(), anim.Image[0], image.ZP, draw.Src)
	for i, srcImg := range anim.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), srcImg, image.ZP, draw.Over)

		// save current frame "stack". This will overwrite an existing file with that name
		file, err := os.Create(fmt.Sprintf("%s%d%s", "outputImages/outImage", i, ".png"))
		if err != nil {
			return nil
		}

		err = png.Encode(file, overpaintImage)
		if err != nil {
			return nil
		}

		file.Close()
	}
	return anim
}

func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}
