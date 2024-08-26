package entity

import (
	"strconv"
	"strings"

	"github.com/duswie/dxf/format"
)

const THREEDSOLID_LINE_CHUNK_SIZE = 255

var encodeTable map[rune]rune

func createEncodeTable() {
	encodeTable = map[rune]rune{
		' ': ' ', // 0x20
		'_': '@', // 0x40
		'@': '_', // 0x5F
	}
	for c := 0x41; c < 0x5F; c++ {
		encodeTable[rune(c)] = rune(0x5E - (c - 0x41)) // 0x5E->'A', 'B'->0x5D, ...
	}
}

func encode(textLines []string) []string {
	if encodeTable == nil {
		createEncodeTable()
	}
	var result []string
	for _, line := range textLines {
		result = append(result, _encode(line))
	}
	return result
}

func _encode(text string) string {
	var s strings.Builder
	for _, c := range text {
		if val, ok := encodeTable[c]; ok {
			s.WriteRune(val)
			if c == 'A' {
				s.WriteRune(' ') // append a space for an 'A' -> cryptography
			}
		} else {
			s.WriteRune(c ^ 0x5F)
		}
	}
	return s.String()
}

// ThreeDSolid represents 3DFACE Entity.
type ThreeDSolid struct {
	*entity
	data_lines []string
	Bbox       [][]float64
}

// IsEntity is for Entity interface.
func (f *ThreeDSolid) IsEntity() bool {
	return true
}

// New3DFace creates a new ThreeDSolid.
func New3DSolid(lines []string) *ThreeDSolid {
	f := &ThreeDSolid{
		entity:     NewEntity(THREEDSOLID),
		data_lines: lines,
	}
	return f
}

// Format writes data to formatter.
func (f *ThreeDSolid) Format(fm format.Formatter) {

	enc_lines := encode(f.data_lines)

	f.entity.Format(fm)
	fm.WriteString(100, "AcDbModelerGeometry")
	fm.WriteString(70, "1")

	// Write data lines, 255 characters per line. If a line is longer than 255 characters, split it into multiple lines.
	// The first line is written with group code 1, and subsequent lines are written with group code 3.
	// https://documentation.help/AutoCAD-DXF/WS1a9193826455f5ff18cb41610ec0a2e719-7a39.htm
	for _, line := range enc_lines {
		for i := 0; i < len(line); i += THREEDSOLID_LINE_CHUNK_SIZE {
			end := i + THREEDSOLID_LINE_CHUNK_SIZE
			num := 3

			if end > len(line) {
				end = len(line)
			}

			if i == 0 {
				num = 1
			}
			fm.WriteString(num, line[i:end])
		}
	}
}

// String outputs data using default formatter.
func (f *ThreeDSolid) String() string {
	fm := format.NewASCII()
	return f.FormatString(fm)
}

// FormatString outputs data using given formatter.
func (f *ThreeDSolid) FormatString(fm format.Formatter) string {
	f.Format(fm)
	return fm.Output()
}

func (f *ThreeDSolid) BBox() ([]float64, []float64) {

	if f.Bbox != nil {
		return f.Bbox[0], f.Bbox[1]
	}

	//TODO: fix for transformed 3d solids, by find and applying transformation matrix
	mins := make([]float64, 3)
	maxs := make([]float64, 3)

	//find all the points in the data_lines
	for _, line := range f.data_lines {
		//check if line starts with point
		if strings.HasPrefix(line, "point") {
			//extract the coordinates
			coords := strings.Split(line, " ")
			x := coords[4]
			y := coords[5]
			z := coords[6]
			//convert the coordinates to float64
			xf, _ := strconv.ParseFloat(x, 64)
			yf, _ := strconv.ParseFloat(y, 64)
			zf, _ := strconv.ParseFloat(z, 64)

			//check if the coordinates are less than the current minimum
			if xf < mins[0] {
				mins[0] = xf
			}
			if yf < mins[1] {
				mins[1] = yf
			}
			if zf < mins[2] {
				mins[2] = zf
			}
			//check if the coordinates are more than the current maximum
			if xf > maxs[0] {
				maxs[0] = xf
			}
			if yf > maxs[1] {
				maxs[1] = yf
			}
			if zf > maxs[2] {
				maxs[2] = zf
			}
		}
	}
	return mins, maxs
}
