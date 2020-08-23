package period

type designator byte

const (
	NoFraction designator = 0
	Year       designator = 'Y'
	Month      designator = 'M'
	Week       designator = 'W'
	Day        designator = 'D'
	Hour       designator = 'h'
	Minute     designator = 'm'
	Second     designator = 's'
)

var designatorByte = map[designator]byte{
	Year:   'Y',
	Month:  'M',
	Week:   'W',
	Day:    'D',
	Hour:   'H',
	Minute: 'M',
	Second: 'S',
}

func (d designator) IsAfterT() bool {
	switch d {
	case Hour, Minute, Second:
		return true
	}
	return false
}

func (d designator) Byte() byte {
	return designatorByte[d]
}

func (d designator) IsOneOf(xx ...designator) bool {
	for _, x := range xx {
		if x == d {
			return true
		}
	}
	return false
}

func (d designator) IsNotOneOf(xx ...designator) bool {
	for _, x := range xx {
		if x == d {
			return false
		}
	}
	return true
}
