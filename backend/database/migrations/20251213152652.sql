-- Create "MonsterTrashCategory" table
CREATE TABLE `MonsterTrashCategory` (
  `MonsterTrashCategoryId` varchar(36) NOT NULL COMMENT "モンスターゴミ種別ID(UUID)",
  `MonsterId` varchar(36) NOT NULL COMMENT "モンスターID(UUID)",
  `TrashCategory` tinyint unsigned NOT NULL COMMENT "ゴミ種別(0:指定なし, 1:燃えるゴミ, 2:不燃ごみ, 3:缶, 4:瓶, 5:ペットボトル)",
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "作成日時",
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新日時",
  PRIMARY KEY (`MonsterTrashCategoryId`),
  INDEX `idx_monster_id` (`MonsterId`),
  UNIQUE INDEX `idx_monster_trash_unique` (`MonsterId`, `TrashCategory`),
  INDEX `idx_trash_category` (`TrashCategory`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin COMMENT "モンスターのゴミ種別(多対多)";
