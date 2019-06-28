package font

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"unicode/utf8"

	"../bitmap"
)

type ColorSetter interface {
	Set(x, y int, c color.Color)
}

type Font struct {
	size  Size
	scale int
	c     color.Color
	m     map[rune]*bitmap.Bitmap
	s     map[rune]Size
}

func New(fi FontInfo) (*Font, error) {
	m, err := newMapFromFontInfo(fi)
	if err != nil {
		return nil, err
	}
	s, err1 := newMapSizeFromFontInfo(fi)
	if err1 != nil {
		return nil, err
	}
	return &Font{
		size:  fi.Size,
		scale: 1,
		c:     color.White,
		m:     m,
		s:     s,
	}, nil
}

func newMapSizeFromFontInfo(fi FontInfo) (map[rune]Size, error) {
	s := make(map[rune]Size)
	for _, ri := range fi.CharSet {
		r := rune(ri.Character)
		if _, ok := s[r]; ok {
			return nil, fmt.Errorf("duplicate rune '%c'", r)
		}
		//x := ri.Size.Dx
		//y := ri.Size.Dy
		//fmt.Printf("x: %d, y: %d\n", x, y)
		s[r] = ri.Size
	}
	return s, nil
}

func NewFromFile(fileName string) (*Font, error) {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var fi FontInfo

	err = json.Unmarshal(data, &fi)
	if err != nil {
		return nil, err
	}

	return New(fi)
}

func newMapFromFontInfo(fi FontInfo) (map[rune]*bitmap.Bitmap, error) {
	m := make(map[rune]*bitmap.Bitmap)
	for _, ri := range fi.CharSet {
		r := rune(ri.Character)
		if _, ok := m[r]; ok {
			return nil, fmt.Errorf("duplicate rune '%c'", r)
		}
		bm, err := newBitmapFromLines(ri.Bitmap, rune(fi.TargetChar), fi.AnchorPos, ri.Size)
		if err != nil {
			return nil, err
		}
		m[r] = bm
		//fmt.Println(bm)
	}
	return m, nil
}

func newBitmapFromLines(lines []string, target rune, pos image.Point, size Size) (*bitmap.Bitmap, error) {
	var (
		nX = size.Dx
		nY = size.Dy
	)
	bm := bitmap.New(nX * nY)
	for iY := 0; iY < nY; iY++ {
		if iY >= len(lines) {
			break
		}
		data := []byte(lines[iY])
		for iX := 0; iX < nX; iX++ {
			r, size := utf8.DecodeRune(data)
			if size == 0 {
				break
			}
			data = data[size:]
			if r == target {
				var (
					x = pos.X + iX
					y = pos.Y + iY
				)
				if (x >= 0) && (x < nX) {
					if (y >= 0) && (y < nY) {
						bm.Set(y*nX+x, true)
					}
				}
			}
		}
	}
	return bm, nil
}

func (f *Font) Scale(scale int) {
	if scale > 0 {
		f.scale = scale
	}
}

func (f *Font) GetScale() int {
	return f.scale
}

func (f *Font) GetSizeMap() map[rune]Size {
	return f.s
}

func (f *Font) Color(c color.Color) {
	f.c = c
}

func (f *Font) GetColor() color.Color {
	return f.c
}

func (f *Font) GetRuneBounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{
			X: f.size.Dx * f.scale,
			Y: f.size.Dy * f.scale,
		},
	}
}

func (f *Font) GetTextBounds(text string) image.Rectangle {

	data := []byte(text)

	maxY := 1
	var x, maxX int

	for {
		r, size := utf8.DecodeRune(data)
		if size == 0 {
			break
		}
		data = data[size:]

		if r == '\n' {
			x = 0
			maxY++
			continue
		}
		if r == '\t' {
			x += 4
			if maxX < x {
				maxX = x
			}
			continue
		}
		x++
		if maxX < x {
			maxX = x
		}
	}

	var (
		x1 = f.size.Dx * f.scale * maxX
		y1 = f.size.Dy * f.scale * maxY
	)

	return image.Rect(0, 0, x1, y1)
}

func (f *Font) GetUberTextBounds(text string) image.Rectangle {
	i := 0
	maxY := 1
	var x, maxX int
	for {
		if i >= len(text) {
			break
		}
		r := []rune(text)[i]
		/*fmt.Println("x at beginning of loop")
		fmt.Println(x)
		fmt.Printf("Rune: %c\n", r)*/
		//r, size := utf8.DecodeRune(data)
		a := f.s[r]
		z := a
		fX := z.Dx
		//fX := f.s[r].Dx
		//fY := f.s[r].Dy
		//fmt.Println("size of RUne")
		//fmt.Println(fX)
		//data = data[size:]

		if r == '\n' {
			x = 0
			maxY++
			continue
		}
		if r == '\t' {
			x += 4
			if maxX < x {
				maxX = x
			}
			continue
		}
		x += fX
		//fmt.Println("position of Point X")
		//fmt.Println(p.X)
		//fmt.Println("position of Y")
		//fmt.Println(p.Y)
		//fmt.Println("New position")
		//fmt.Println(x)
		i++
	}
	var (
		x1 = x
		y1 = f.size.Dy * f.scale * maxY
	)
	return image.Rect(0, 0, x1, y1)
}

func (f *Font) DrawRune(cs ColorSetter, pos image.Point, r rune) {

	bm, ok := f.m[r]
	if !ok {
		return
	}

	f.drawBitmap(cs, pos, bm, r)
	/*if m, ok := cs.(*image.RGBA); ok {
		f.drawBitmapRGBA(m, pos, bm)
	} else {
		f.drawBitmap(cs, pos, bm, r)
	}*/
}

func (f *Font) DrawRuneWipeUp(cs ColorSetter, pos image.Point, r rune, startColor color.RGBA, endColor color.RGBA) {

	bm, ok := f.m[r]
	if !ok {
		return
	}

	f.drawBitmapWipeUp(cs, pos, bm, r, startColor, endColor)
	/*if m, ok := cs.(*image.RGBA); ok {
		f.drawBitmapRGBA(m, pos, bm)
	} else {
		f.drawBitmap(cs, pos, bm, r)
	}*/
}

func (f *Font) DrawRuneWipeDown(cs ColorSetter, pos image.Point, r rune, startColor color.RGBA, endColor color.RGBA) {

	bm, ok := f.m[r]
	if !ok {
		return
	}

	f.drawBitmapWipeDown(cs, pos, bm, r, startColor, endColor)
	/*if m, ok := cs.(*image.RGBA); ok {
		f.drawBitmapRGBA(m, pos, bm)
	} else {
		f.drawBitmap(cs, pos, bm, r)
	}*/
}

func (f *Font) DrawText(cs ColorSetter, pos image.Point, text string, c color.Color) {

	//data := []byte(text)

	//sizeX := f.size.Dx * f.scale
	//sizeY := f.size.Dy * f.scale
	f.Color(c)
	x := 0
	y := 0
	i := 0
	for {
		if i >= len(text) {
			break
		}
		r := []rune(text)[i]
		/*fmt.Println("x at beginning of loop")
		fmt.Println(x)
		fmt.Printf("Rune: %c\n", r)*/
		//r, size := utf8.DecodeRune(data)
		a := f.s[r]
		z := a
		fX := z.Dx
		fY := z.Dx
		//fX := f.s[r].Dx
		//fY := f.s[r].Dy
		//fmt.Println("size of RUne")
		//fmt.Println(fX)
		//data = data[size:]

		if r == '\n' {
			x = 0
			y++
			continue
		}
		if r == '\t' {
			x += 4
			continue
		}

		p := image.Point{
			X: pos.X + x,
			Y: pos.Y + y*fY,
		}
		//fmt.Println("position of Point X")
		//fmt.Println(p.X)
		//fmt.Println("position of Y")
		//fmt.Println(p.Y)
		f.DrawRune(cs, p, r)

		x += fX
		//fmt.Println("New position")
		//fmt.Println(x)
		i++
	}
}

func (f *Font) DrawTextFade(cs ColorSetter, pos image.Point, c color.Color, text string) {
	//data := []byte(text)
	f.Color(c)

	//sizeX := f.size.Dx * f.scale
	//sizeY := f.size.Dy * f.scale

	x := 0
	y := 0
	i := 0

	for {
		if i >= len(text) {
			break
		}
		r := []rune(text)[i]
		a := f.s[r]
		z := a
		fX := z.Dx
		fY := z.Dy

		if r == '\n' {
			x = 0
			y++
			continue
		}
		if r == '\t' {
			x += 4
			continue
		}

		p := image.Point{
			X: pos.X + x,
			Y: pos.Y + y*fY,
		}

		f.DrawRune(cs, p, r)
		x += fX
		i++
	}
}

func (f *Font) DrawTextWipeUp(cs ColorSetter, pos image.Point, startColor color.Color, endColor color.Color, text string) {
	//data := []byte(text)
	f.Color(startColor)
	sr, sg, sb, sa := startColor.RGBA()
	sRGBA := &color.RGBA{
		R: uint8(sr),
		G: uint8(sg),
		B: uint8(sb),
		A: uint8(sa),
	}
	er, eg, eb, ea := endColor.RGBA()
	eRGBA := &color.RGBA{
		R: uint8(er),
		G: uint8(eg),
		B: uint8(eb),
		A: uint8(ea),
	}

	//sizeX := f.size.Dx * f.scale
	//sizeY := f.size.Dy * f.scale

	x := 0
	y := 0
	i := 0

	for {
		if i >= len(text) {
			break
		}
		r := []rune(text)[i]
		a := f.s[r]
		z := a
		fX := z.Dx
		fY := z.Dy

		if r == '\n' {
			x = 0
			y++
			continue
		}
		if r == '\t' {
			x += 4
			continue
		}

		p := image.Point{
			X: pos.X + x,
			Y: pos.Y + y*fY,
		}

		f.DrawRuneWipeUp(cs, p, r, *sRGBA, *eRGBA)
		x += fX
		i++
	}
}

func (f *Font) DrawTextWipeDown(cs ColorSetter, pos image.Point, startColor color.Color, endColor color.Color, text string) {
	//data := []byte(text)
	f.Color(startColor)
	sr, sg, sb, sa := startColor.RGBA()
	sRGBA := &color.RGBA{
		R: uint8(sr),
		G: uint8(sg),
		B: uint8(sb),
		A: uint8(sa),
	}
	er, eg, eb, ea := endColor.RGBA()
	eRGBA := &color.RGBA{
		R: uint8(er),
		G: uint8(eg),
		B: uint8(eb),
		A: uint8(ea),
	}

	//sizeX := f.size.Dx * f.scale
	//sizeY := f.size.Dy * f.scale

	x := 0
	y := 0
	i := 0

	for {
		if i >= len(text) {
			break
		}
		r := []rune(text)[i]
		a := f.s[r]
		z := a
		fX := z.Dx
		fY := z.Dy

		if r == '\n' {
			x = 0
			y++
			continue
		}
		if r == '\t' {
			x += 4
			continue
		}

		p := image.Point{
			X: pos.X + x,
			Y: pos.Y + y*fY,
		}

		f.DrawRuneWipeDown(cs, p, r, *sRGBA, *eRGBA)
		x += fX
		i++
	}
}

func (f *Font) drawBitmap(cs ColorSetter, pos image.Point, bm *bitmap.Bitmap, r rune) {

	var (
		nX = f.s[r].Dx
		nY = f.s[r].Dy
	)
	//fmt.Printf("rune being printed: %c\n", r)
	if f.scale == 1 {
		y := pos.Y
		for iY := 0; iY < nY; iY++ {
			x := pos.X
			for iX := 0; iX < nX; iX++ {
				if bit, _ := bm.Get(iY*nX + iX); bit {
					cs.Set(x, y, f.c)
				}
				x++
			}
			y++
		}
	}
}

func (f *Font) drawBitmapWipeUp(cs ColorSetter, pos image.Point, bm *bitmap.Bitmap, r rune, startColor color.RGBA, endColor color.RGBA) {
	var (
		nX = f.s[r].Dx
		nY = f.s[r].Dy
	)
	var u uint8
	u = endColor.R - startColor.R
	//fmt.Printf("rune being printed: %c\n", r)
	if f.scale == 1 {
		y := pos.Y
		for iY := 0; iY < nY; iY++ {
			x := pos.X
			for iX := 0; iX < nX; iX++ {
				if bit, _ := bm.Get(iY*nX + iX); bit {
					cs.Set(x, y, startColor)
					//u = endColor.R - startColor.R
				}
				x++
			}
			y++
			startColor.R += (u / 6)
			startColor.G += (u / 6)
			startColor.B += (u / 6)
		}
	}
}

func (f *Font) drawBitmapWipeDown(cs ColorSetter, pos image.Point, bm *bitmap.Bitmap, r rune, startColor color.RGBA, endColor color.RGBA) {
	var (
		nX = f.s[r].Dx
		nY = f.s[r].Dy
	)
	var u uint8
	u = startColor.R - endColor.R
	/*fmt.Printf("value of u: %d\n", u)
	fmt.Printf("startColorR Value: %d\n", startColor.R)
	fmt.Printf("startColorR Value: %d\n", endColor.R)*/
	//fmt.Printf("rune being printed: %c\n", r)
	if f.scale == 1 {
		y := pos.Y
		for iY := 0; iY < nY; iY++ {
			x := pos.X
			for iX := 0; iX < nX; iX++ {
				if bit, _ := bm.Get(iY*nX + iX); bit {
					cs.Set(x, y, startColor)
					//u = endColor.R - startColor.R
				}
				x++
			}
			y++
			startColor.R -= (u / 6)
			startColor.G -= (u / 6)
			startColor.B -= (u / 6)
		}
	}
}

func pixelRect(x, y int, scale int) image.Rectangle {
	return image.Rect(x, y, x+scale, y+scale)
}

func fillRect(cs ColorSetter, r image.Rectangle, c color.Color) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cs.Set(x, y, c)
		}
	}
}
