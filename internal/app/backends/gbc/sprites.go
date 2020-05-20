package gbc

import (
	"sort"

	"github.com/omstrumpf/goemu/internal/app/log"
)

// sprite represents a gameboy sprite
type sprite struct {
	yPos    byte
	xPos    byte
	tileNum byte

	priority    bool
	yFlip       bool
	xFlip       bool
	paletteFlag bool
	tileBank    bool
	paletteNum  byte
}

// oam is the Object Attribute Memory, which holds the sprite data
type oam struct {
	sprites [40]sprite
}

func (oam *oam) Read(addr uint16) byte {
	spriteNum := addr >> 2
	if spriteNum >= 40 {
		log.Warningf("Encountered overflowed OAM read")
		return 0xFF
	}

	sprite := oam.sprites[spriteNum]

	switch addr & 0x3 {
	case 0:
		return sprite.yPos
	case 1:
		return sprite.xPos
	case 2:
		return sprite.tileNum
	case 3:
		var ret byte
		if sprite.priority {
			ret |= 0x80
		}
		if sprite.yFlip {
			ret |= 0x40
		}
		if sprite.xFlip {
			ret |= 0x20
		}
		if sprite.paletteFlag {
			ret |= 0x10
		}
		if sprite.tileBank {
			ret |= 0x08
		}
		ret |= (sprite.paletteNum & 0x07)
		return ret
	}

	log.Warningf("Unexpected OAM read fallthrough")
	return 0xFF
}

func (oam *oam) Write(addr uint16, val byte) {
	spriteNum := addr >> 2
	if spriteNum >= 40 {
		log.Warningf("Encountered overflowed OAM write")
	}

	sprite := &oam.sprites[spriteNum]

	switch addr & 0x3 {
	case 0:
		sprite.yPos = val
	case 1:
		sprite.xPos = val
	case 2:
		sprite.tileNum = val
	case 3:
		sprite.priority = (val&0x80 != 0)
		sprite.yFlip = (val&0x40 != 0)
		sprite.xFlip = (val&0x20 != 0)
		sprite.paletteFlag = (val&0x10 != 0)
		sprite.tileBank = (val&0x08 != 0)
		sprite.paletteNum = val & 0x07
	}
}

func (oam *oam) VisibleSpritesOnLine(line byte, tallSprites bool) []*sprite {
	var ret []*sprite

	height := byte(8)
	if tallSprites {
		height = 16
	}

	// Collect all sprites on the line
	for i := range oam.sprites {
		sprite := &oam.sprites[i]
		if sprite.yPos != 0 && line+16 >= sprite.yPos && line+16-height < sprite.yPos {
			ret = append(ret, sprite)
		}
	}

	// Sort by x position
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].xPos < ret[j].xPos
	})

	// Keep the first 10 from the left
	if len(ret) > 10 {
		ret = ret[:10]
	}

	// Reverse
	for l, r := 0, len(ret)-1; l < r; l, r = l+1, r-1 {
		ret[l], ret[r] = ret[r], ret[l]
	}

	return ret
}
