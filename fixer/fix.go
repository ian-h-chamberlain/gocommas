package fixer

import (
	"go/token"
)

func AddMissingCommas(src []byte, commaPositions []token.Position) []byte {
	offset := 0
	for _, commaPos := range commaPositions {
		// insert a comma at the given position
		i := commaPos.Offset + offset

		after := append([]byte(","), src[i:]...)
		src = append(src[:i], after...)
		offset += 1
	}

	return src
}
