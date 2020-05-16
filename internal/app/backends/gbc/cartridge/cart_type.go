package cartridge

// CartType is the type of the cartridge, specifying which memory bank controller and features are present.
type CartType byte

// CartTypes
const (
	ROM CartType = iota
	MBC1
	MBC1RAM
	MBC1RAMBAT
	MBC2
	MBC2BAT
	ROMRAM
	ROMRAMBAT
	MMM01
	MM01RAM
	MM01RAMBAT
	MBC3TIMBAT
	MBC3TIMRAMBAT
	MBC3
	MBC3RAM
	MBC3RAMBAT
	MBC4
	MBC4RAM
	MBC4RAMBAT
	MBC5
	MBC5RAM
	MBC5RAMBAT
	MBC5RUMBLE
	MBC5RUMBLERAM
	MBC5RUMBLERAMBAT
	POCKETCAM
	BANDAITAMA5
	HUC3
	HUC1RAMBAT
)

func (c CartType) String() string {
	switch c {
	case ROM:
		return "ROM"
	case MBC1:
		return "MBC1"
	case MBC1RAM:
		return "MBC1+RAM"
	case MBC1RAMBAT:
		return "MBC1+RAM+BAT"
	case MBC2:
		return "MBC2"
	case MBC2BAT:
		return "MBC2+BAT"
	case ROMRAM:
		return "ROM+RAM"
	case ROMRAMBAT:
		return "ROM+RAM+BAT"
	case MMM01:
		return "MMM01"
	case MM01RAM:
		return "MM01+RAM"
	case MM01RAMBAT:
		return "MM01+RAM+BAT"
	case MBC3TIMBAT:
		return "MBC3+TIM+BAT"
	case MBC3TIMRAMBAT:
		return "MBC3+TIM+RAM+BAT"
	case MBC3:
		return "MBC3"
	case MBC3RAM:
		return "MBC3+RAM"
	case MBC3RAMBAT:
		return "MBC3+RAM+BAT"
	case MBC4:
		return "MBC4"
	case MBC4RAM:
		return "MBC4+RAM"
	case MBC4RAMBAT:
		return "MBC4+RAM+BAT"
	case MBC5:
		return "MBC5"
	case MBC5RAM:
		return "MBC5+RAM"
	case MBC5RAMBAT:
		return "MBC5+RAM+BAT"
	case MBC5RUMBLE:
		return "MBC5RUMBLE"
	case MBC5RUMBLERAM:
		return "MBC5RUMBLE+RAM"
	case MBC5RUMBLERAMBAT:
		return "MBC5RUMBLE+RAM+BAT"
	case POCKETCAM:
		return "POCKETCAM"
	case BANDAITAMA5:
		return "BANDAITAMA5"
	case HUC3:
		return "HUC3"
	case HUC1RAMBAT:
		return "HUC1+RAM+BAT"
	default:
		return "UNKNOWN"
	}
}
