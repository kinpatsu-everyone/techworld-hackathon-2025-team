package enum

type TrashCategory uint8

const (
	// TrashCategoryNone は指定なし
	TrashCategoryNone TrashCategory = iota
	// TrashCategoryBurnable は燃えるゴミ
	TrashCategoryBurnable
	// TrashCategoryNonBurnable は不燃ごみ
	TrashCategoryNonBurnable
	// TrashCategoryCan は缶ごみ
	TrashCategoryCan
	// TrashCategoryGlassBottle は瓶
	TrashCategoryGlassBottle
	// TrashCategoryPetBottle はペットボトル
	TrashCategoryPetBottle
)

// trashTypeStringToUint8 はゴミ種別の文字列をuint8に変換します
// "燃えるゴミ" -> 1, "不燃ごみ" -> 2, "缶・瓶" -> 3, "ペットボトル" -> 4, "紙" -> 5, その他 -> 1 (デフォルト)
func StringToTrashCategoryEnum(trashType string) TrashCategory {
	switch trashType {
	case "燃えるゴミ":
		return 1
	case "不燃ごみ":
		return 2
	case "缶・瓶", "缶", "瓶":
		return 3
	case "ペットボトル":
		return 4
	case "紙":
		return 5
	default:
		// 不明な場合はデフォルトで燃えるゴミとする
		return 1
	}
}
