package mysql

// TrashCategoryToString はゴミ種別のIDを文字列に変換します
// 0:指定なし, 1:燃えるゴミ, 2:不燃ごみ, 3:缶, 4:瓶, 5:ペットボトル
func TrashCategoryToString(category uint8) string {
	switch category {
	case 0:
		return "指定なし"
	case 1:
		return "燃えるゴミ"
	case 2:
		return "不燃ごみ"
	case 3:
		return "缶"
	case 4:
		return "瓶"
	case 5:
		return "ペットボトル"
	default:
		return "不明"
	}
}

