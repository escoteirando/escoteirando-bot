package consts

var emojisTipoSecao = map[int]string{
	1: MelhorPossivel,
	2: FlorDeLis,
	3: FlorDeLis,
	4: FlorDeLis,
}

func EmojiTipoSecao(tipoSecao int) string {
	if val, ok := emojisTipoSecao[tipoSecao]; ok {
		return val
	}
	return FlorDeLis
}
