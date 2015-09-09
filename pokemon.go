package pokemoncenter

import (
	"strconv"
	"strings"
)

type Pokémon struct {
	PCPokémon
	Level uint8
	Stats
}

type PCPokémon struct {
	Species                            Specie
	CurrHP                             uint16
	PCLevel                            uint8 // Appears to be mostly unused. Should be same as Level
	Condition                          Condition
	Type1, Type2                       Type
	HeldItem                           uint8
	Move1, Move2, Move3, Move4         uint8
	OT                                 uint16
	Experience                         uint24
	EV                                 Stats
	IV                                 uint16
	Move1pp, Move2pp, Move3pp, Move4pp uint8
}

type Stats struct {
	HP, Attack, Defense, Speed, Special uint16
}

type Specie uint8

func (i *Specie) UnmarshalText(data []byte) error {
	*i = Specie(Species[strings.ToLower(string(data))])
	return nil
}

func (i Specie) MarshalText() (data []byte, err error) {
	data = []byte(revSpecies[uint8(i)])
	return
}

type Condition uint8

func (i *Condition) UnmarshalText(data []byte) error {
	*i = Condition(Conditions[strings.ToLower(string(data))])
	return nil
}

func (i Condition) MarshalText() (data []byte, err error) {
	data = []byte(revConditions[uint8(i)])
	return
}

type Type uint8

func (i *Type) UnmarshalText(data []byte) error {
	*i = Type(Types[strings.ToLower(string(data))])
	return nil
}

func (i Type) MarshalText() (data []byte, err error) {
	data = []byte(revTypes[uint8(i)])
	return
}

// Custom type to facilitate the 3-byte width.
// Stores as a 32-bit wide field,
// but writes the lower 24-bits when binary encoded.
type uint24 [3]byte

func (i *uint24) UnmarshalText(data []byte) error {
	num, _ := strconv.Atoi(string(data))
	*i = UInt24(uint32(num))
	return nil
}

func (i uint24) MarshalText() (data []byte, err error) {
	data = []byte(strconv.Itoa(int(i.UInt32())))
	//Safe cast from uint32 to int, since value on occupies lower 3 bytes.
	return
}

// "Casts" the uint24 to a uint32 for usage with outside code.
func (i uint24) UInt32() (val uint32) {
	val |= uint32(i[0]) << 16
	val |= uint32(i[1]) << 8
	val |= uint32(i[2])
	return
}

// "Casts" the uint32 into a uint24.
// Useful for compatability with external code.
// The parameter is truncated to fit the 24 bits.
func UInt24(i uint32) (val uint24) {
	val[0] = byte(i >> 16)
	val[1] = byte(i >> 8)
	val[2] = byte(i)
	return
}

//Convenience Alias
//type Pokemon Pokémon
