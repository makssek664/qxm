package main;
import (
	"time"
	"gorm.io/gorm"
)
type Event struct {
	gorm.Model
	UserID uint
	Date time.Time
	Name string
	Desc string
	Color uint32
	ID uint
}

// Helper to convert seperate RGB pairs into a u32 word.
func (e *Event)RGBToColor(r uint8, g uint8, b uint8) {
	e.Color = (r << 24 | g << 16 | b << 8)
}
