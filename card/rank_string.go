// Code generated by "stringer -type=Rank"; DO NOT EDIT.

package card

import "fmt"

const _Rank_name = "TwoThreeFourFiveSixSevenEightNineTenJackQueenKingAce"

var _Rank_index = [...]uint8{0, 3, 8, 12, 16, 19, 24, 29, 33, 36, 40, 45, 49, 52}

func (i Rank) String() string {
	if i < 0 || i >= Rank(len(_Rank_index)-1) {
		return fmt.Sprintf("Rank(%d)", i)
	}
	return _Rank_name[_Rank_index[i]:_Rank_index[i+1]]
}
