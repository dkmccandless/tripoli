// Code generated by "stringer -type=Suit"; DO NOT EDIT.

package card

import "fmt"

const _Suit_name = "ClubsDiamondsSpadesHearts"

var _Suit_index = [...]uint8{0, 5, 13, 19, 25}

func (i Suit) String() string {
	if i < 0 || i >= Suit(len(_Suit_index)-1) {
		return fmt.Sprintf("Suit(%d)", i)
	}
	return _Suit_name[_Suit_index[i]:_Suit_index[i+1]]
}
