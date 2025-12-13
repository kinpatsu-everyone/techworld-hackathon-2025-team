CREATE TABLE `MonsterTrashCategory` (
    `MonsterTrashCategoryId` varchar(36) NOT NULL comment 'モンスターゴミ種別ID(UUID)',
    `MonsterId` varchar(36) NOT NULL comment 'モンスターID(UUID)',
    `TrashCategory` TINYINT UNSIGNED NOT NULL comment 'ゴミ種別(0:指定なし, 1:燃えるゴミ, 2:不燃ごみ, 3:缶, 4:瓶, 5:ペットボトル)',
    `CreatedAt` datetime NOT NULL default CURRENT_TIMESTAMP comment '作成日時',
    `UpdatedAt` datetime NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新日時',
    PRIMARY KEY (`MonsterTrashCategoryId`),
    UNIQUE INDEX `idx_monster_trash_unique` (`MonsterId`, `TrashCategory`),
    INDEX `idx_monster_id` (`MonsterId`),
    INDEX `idx_trash_category` (`TrashCategory`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT 'モンスターのゴミ種別(多対多)';
