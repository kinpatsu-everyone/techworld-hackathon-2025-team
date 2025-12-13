CREATE TABLE `MonsterAttribute` (
    `MonsterId` varchar(36) NOT NULL comment 'モンスターID(UUID)',
    `AttributeName` varchar(50) NOT NULL comment '属性の名前',
    `ColorCode` varchar(30) NOT NULL comment 'カラーコード',
    `CreatedAt` datetime NOT NULL default CURRENT_TIMESTAMP comment '作成日時',
    `UpdatedAt` datetime NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新日時',
    PRIMARY KEY (`MonsterId`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT 'モンスターの属性情報';
