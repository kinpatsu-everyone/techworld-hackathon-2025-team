CREATE TABLE `Monster` (
    `MonsterId` varchar(36) NOT NULL comment 'モンスターID(UUID)',
    `Nickname` varchar(50) NOT NULL comment 'ニックネーム',
    `OriginalTrashBinImageUrl` TEXT NOT NULL comment '元のゴミ箱の画像URL',
    `GeneratedMonsterImageUrl` TEXT NOT NULL comment '生成したモンスターの画像URL',
    `Latitude` DECIMAL(10, 8) NULL comment '緯度(-90.0 ~ 90.0)',
    `Longitude` DECIMAL(11, 8) NULL comment '経度(-180.0 ~ 180.0)',
    `CreatedAt` datetime NOT NULL default CURRENT_TIMESTAMP comment '作成日時',
    `UpdatedAt` datetime NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新日時',
    PRIMARY KEY (`MonsterId`),
    INDEX `idx_location` (`Latitude`, `Longitude`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT 'モンスターの基本情報';
