CREATE TABLE `User` (
    `UserId` varchar(36) NOT NULL comment 'ユーザーID(UUID)',
    `Nickname` varchar(50) NOT NULL comment 'ニックネーム',
    `CreatedAt` datetime NOT NULL default CURRENT_TIMESTAMP comment '作成日時',
    `UpdatedAt` datetime NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP comment '更新日時',
    PRIMARY KEY (`UserId`)
) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT 'ユーザーの基本情報';
