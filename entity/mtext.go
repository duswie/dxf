package entity

import (
	"github.com/duswie/dxf/format"
	"github.com/duswie/dxf/table"
)

const MTEXT_CHUNK_SIZE = 250

// Attachment Point
const (
	ATT_TOP_LEFT = iota + 1
	ATT_TOP_CENTER
	ATT_TOP_RIGHT
	ATT_MIDDLE_LEFT
	ATT_MIDDLE_CENTER
	ATT_MIDDLE_RIGHT
	ATT_BOTTOM_LEFT
	ATT_BOTTOM_CENTER
	ATT_BOTTOM_RIGHT
)

// Mtext represents MTEXT Entity.
type Mtext struct {
	*entity
	Coord1           []float64    // 10, 20, 30
	Coord2           []float64    // 11, 21, 31
	Height           float64      // 40 Nominal (initial) text height
	AttachmentPoint  int          // 71
	DrawingDirection int          // 72
	Style            *table.Style // 7
	Rotation         float64      // 50
	RectangleWidth   float64      // 41
	Value            string       // 1,3

}

// IsEntity is for Entity interface.
func (t *Mtext) IsEntity() bool {
	return true
}

// NewText creates a new Text.
func NewMtext() *Mtext {
	t := &Mtext{
		entity:           NewEntity(MTEXT),
		Coord1:           []float64{0.0, 0.0, 0.0},
		Coord2:           []float64{0.0, 0.0, 0.0},
		Height:           1.0,
		Value:            "",
		Style:            table.ST_STANDARD,
		AttachmentPoint:  ATT_TOP_LEFT,
		DrawingDirection: 4,
	}
	return t
}

// Format writes data to formatter.
func (t *Mtext) Format(f format.Formatter) {
	t.entity.Format(f)
	f.WriteString(100, "AcDbMText")
	for i := 0; i < 3; i++ {
		f.WriteFloat((i+1)*10, t.Coord1[i])
	}

	f.WriteFloat(40, t.Height)
	f.WriteFloat(41, t.RectangleWidth)
	f.WriteInt(71, t.AttachmentPoint)
	f.WriteInt(72, t.DrawingDirection)

	//devide the text into chunks, write first chunk with code 1 and following with 3
	for i := 0; i < len(t.Value); i += MTEXT_CHUNK_SIZE {
		end := i + MTEXT_CHUNK_SIZE
		num := 3

		if i > len(t.Value) {
			end = len(t.Value)
		}

		if i == 0 {
			num = 1
		}
		f.WriteString(num, t.Value[i:end])
	}
	f.WriteString(7, t.Style.Name())
	f.WriteFloat(50, t.Rotation)

}

// String outputs data using default formatter.
func (t *Mtext) String() string {
	f := format.NewASCII()
	return t.FormatString(f)
}

// FormatString outputs data using given formatter.
func (t *Mtext) FormatString(f format.Formatter) string {
	t.Format(f)
	return f.Output()
}

func (t *Mtext) BBox() ([]float64, []float64) {
	// TODO: text length, anchor point
	mins := []float64{t.Coord1[0], t.Coord1[1], t.Coord1[2]}
	maxs := []float64{t.Coord1[0], t.Coord1[1] + t.Height, t.Coord1[2]}
	return mins, maxs
}
