package sprite

import (
	"os"
	"strings"
	"image/png"

	//palette "github.com/pdevine/go-asciisprite/palette"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

// Blocks provides a map of bits to Unicode block character runes.
var Blocks = map[int]rune{
   0: ' ',
   1: '▘',
   2: '▝' ,
   3: '▀',
   4: '▖',
   5: '▌',
   6: '▞',
   7: '▛',
   8: '▗',
   9: '▚',
  10: '▐',
  11: '▜',
  12: '▄',
  13: '▙',
  14: '▟',
  15: '█',
}

type Surface struct {
	Blocks [][]rune
	Width  int
	Height int
	Alpha  bool
}

// ColorMap provides an interpolation map of characters to termbox colors.
var ColorMap = map[rune]tm.Attribute{
	'R': tm.ColorRed,
	'b': tm.Attribute(53),
	't': tm.Attribute(180),
	'Y': tm.ColorYellow,
	'N': tm.ColorBlack,
	'B': tm.ColorBlue,
	'o': tm.Attribute(209),
	'O': tm.Attribute(167),
	'w': tm.ColorWhite,
	'X': tm.ColorWhite,
	'g': tm.ColorGreen,
	'G': tm.Attribute(35),
}

// Convert is a convenience function to create a 1 bit Costume from a string
func Convert(s string) Costume {
	sf := NewSurfaceFromString(s, false)
	return sf.ConvertToCostume()
}

// ColorConvert convenience function to create a color Costume from a string
func ColorConvert(s string, bg tm.Attribute) Costume {
	sf := NewSurfaceFromString(s, false)
	return sf.ConvertToColorCostume(bg)
}


// NewSurface creates a Surface
func NewSurface(width, height int, alpha bool) Surface {
	blocks := make([][]rune, height, height)
	for cnt := 0; cnt < height; cnt++ {
		blocks[cnt] = make([]rune, width, width)
	}

	s := Surface{
		Blocks: blocks,
		Width:  width,
		Height: height,
		Alpha:  alpha,
	}
	return s
}

// NewSurfaceFromString creates a Surface which can be converted to a Costume
func NewSurfaceFromString(s string, alpha bool) Surface {
	l := strings.Split(s, "\n")
	maxR := len(l) + len(l)%2

	// all block sprites must be even
	m := make([][]rune, maxR, maxR)

	var maxC int
	for _, r := range l {
		maxC = max(maxC, len(r) + len(r)%2)
	}

	for rcnt, r := range l {
		m[rcnt] = make([]rune, maxC, maxC)
		for ccnt, c := range r {
			if c == ' ' {
				if !alpha {
					m[rcnt][ccnt] = 0
				} else {
					continue
				}
			} else {
				m[rcnt][ccnt] = c
			}
		}
	}

	// make certain we make a row for any added space
	if len(l) < maxR {
		m[maxR-1] = make([]rune, maxC, maxC)
	}
	sf := Surface{
		Blocks: m,
		Width:  maxC,
		Height: maxR,
		Alpha:  alpha,
	}
	return sf
}

func NewSurfaceFromPng(fn string) Surface {
	f, err := os.Open(fn)
	if err != nil {
		//
	}

	img, err := png.Decode(f)
	if err != nil {
		//
	}

	b := img.Bounds()
	maxR := (b.Max.Y-b.Min.Y) + (b.Max.Y-b.Min.Y)%2
	maxC := (b.Max.X-b.Min.X) + (b.Max.X-b.Min.X)%2

	// all block sprites must be even
	m := make([][]rune, maxR, maxR)

	for y := 0; y < b.Max.Y-b.Min.Y; y++ {
		m[y] = make([]rune, maxC, maxC)
		for x := 0;  x < b.Max.X-b.Min.X; x++ {
			c := img.At(x+b.Min.X, y+b.Min.Y)
			//m[y][x] = rune(palette.Index(c))
			r, g, b, _ := c.RGBA()
			if r != 0 || g != 0 || b != 0 {
				m[y][x] = 'X'
			}
		}
	}

	sf := Surface{
		Blocks: m,
		Width:  maxC,
		Height: maxR,
		Alpha:  false,
	}
	return sf
}

func (s *Surface) Clear() {
	blocks := make([][]rune, s.Height, s.Height)
	for cnt := 0; cnt < s.Height; cnt++ {
		blocks[cnt] = make([]rune, s.Width, s.Width)
	}
	s.Blocks = blocks
}

// ConvertToCostume converts a Surface into a Costume usable in a Sprite
func (s Surface) ConvertToCostume() Costume {
	blocks := []*Block{}

	for rcnt := 0; rcnt < len(s.Blocks); rcnt+=2 {
		// XXX - needs to be max(len(m[rcnt]), len(m[rcnt+1]))
		// for ccnt := 0; ccnt < max(len(m[rcnt]), len(m[rcnt+1])); ccnt+=2 {
		for ccnt := 0; ccnt < len(s.Blocks[rcnt]); ccnt+=2 {
			c := 0
			if s.Blocks[rcnt][ccnt] != 0 {
				c += 1
			}
			if len(s.Blocks[rcnt]) > ccnt+1 && s.Blocks[rcnt][ccnt+1] != 0 {
				c += 2
			}
			if len(s.Blocks) > rcnt+1 && s.Blocks[rcnt+1][ccnt] != 0 {
				c += 4
			}
			if len(s.Blocks) > rcnt+1 && len(s.Blocks[rcnt]) > ccnt+1 && s.Blocks[rcnt+1][ccnt+1] == 'X' {
				c += 8
			}

			if (s.Alpha && c > 0) || (!s.Alpha) {
				b := &Block{
					Char: Blocks[c],
					X:    ccnt/2,
					Y:    rcnt/2,
				}
				blocks = append(blocks, b)
			}
		}
	}
	return Costume{Blocks: blocks, Width: s.Width/2}
}

// ConvertToColorCostume converts a Surface into a color Costume usable in a Sprite
func (s Surface) ConvertToColorCostume(bg tm.Attribute) Costume {
	blocks := []*Block{}

	for rcnt := 0; rcnt < len(s.Blocks); rcnt+=2 {
		for ccnt := 0; ccnt < len(s.Blocks[rcnt]); ccnt+=2 {
			var fg tm.Attribute
			obg := bg

			runes := []rune{
				s.Blocks[rcnt][ccnt],
				s.Blocks[rcnt][ccnt+1],
				s.Blocks[rcnt+1][ccnt],
				s.Blocks[rcnt+1][ccnt+1],
			}

			for _, b := range runes {
				if b > 0 && fg == 0 {
					fg = ColorMap[b]
				} else if b != 0 && ColorMap[b] != fg {
					obg = ColorMap[b]
				}
			}

			// if we didn't set a foreground, just skip the block
			if fg == 0 {
				continue
			}

			c := 0
			for cnt, b := range runes {
				if ColorMap[b] == fg {
					c += int(uint(1) << uint(cnt))
				}
			}

			blk := &Block{
				Char: Blocks[c],
				X:    ccnt/2,
				Y:    rcnt/2,
				Fg:   tm.Attribute(fg),
				Bg:   tm.Attribute(obg),
			}
			blocks = append(blocks, blk)
		}
	}

	costume := Costume{Blocks: blocks}

	return costume
}

// Blit a Surface onto a Surface
func (s Surface) Blit(t Surface, x, y int) error {
	for rcnt, r := range t.Blocks {
		for ccnt, c := range r {
			if rcnt + y < 0 || rcnt + y >= s.Height || ccnt + x < 0 || ccnt + x >= s.Width {
				continue
			}
			if c <= 16 {
				if !t.Alpha {
					s.Blocks[rcnt+y][ccnt+x] = c
				} else {
					s.Blocks[rcnt+y][ccnt+x] |= c
				}
			} else {
				s.Blocks[rcnt+y][ccnt+x] = c
			}
		}
	}
	return nil
}

// Draw a line between two points on a Surface
func (s Surface) Line(x0, y0, x1, y1 int, ch rune) error {
	if x0 >= s.Width || x1 >= s.Width {
		// XXX - put a real error here
		return nil
	}
	if y0 >= s.Height || y1 >= s.Height {
		return nil
	}

	dx := abs(x1-x0)
	dy := -abs(y1-y0)
	err := dx+dy
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	for {
		// draw at x0, y0
		s.Blocks[y0][x0] = ch
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
	return nil
}

func (s Surface) Rectangle(x0, y0, x1, y1 int, ch rune) error {
	if x0 >= s.Width || x1 >= s.Width {
		// XXX - put a real error here
		return nil
	}
	if y0 >= s.Height || y1 >= s.Height {
		return nil
	}
	s.Line(x0, y0, x1, y0, ch)
	s.Line(x1, y0, x1, y1, ch)
	s.Line(x0, y0, x0, y1, ch)
	s.Line(x0, y1, x1, y1, ch)
	return nil
}

