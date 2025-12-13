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
