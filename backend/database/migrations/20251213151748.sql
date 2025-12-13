-- Create "Monster" table
CREATE TABLE `Monster` (
  `MonsterId` varchar(36) NOT NULL COMMENT "モンスターID(UUID)",
  `Nickname` varchar(50) NOT NULL COMMENT "ニックネーム",
  `ImageUrl` text NOT NULL COMMENT "画像のURL",
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "作成日時",
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新日時",
  PRIMARY KEY (`MonsterId`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin COMMENT "モンスターの基本情報";
-- Create "MonsterAttribute" table
CREATE TABLE `MonsterAttribute` (
  `MonsterId` varchar(36) NOT NULL COMMENT "モンスターID(UUID)",
  `AttributeName` varchar(50) NOT NULL COMMENT "属性の名前",
  `ColorCode` varchar(30) NOT NULL COMMENT "カラーコード",
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "作成日時",
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新日時",
  PRIMARY KEY (`MonsterId`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin COMMENT "モンスターの属性情報";
-- Create "MonsterLocation" table
CREATE TABLE `MonsterLocation` (
  `MonsterLocationId` varchar(36) NOT NULL COMMENT "モンスター位置ID(UUID)",
  `MonsterId` varchar(36) NOT NULL COMMENT "モンスターID(UUID)",
  `Latitude` decimal(10,8) NOT NULL COMMENT "緯度(-90.0 ~ 90.0)",
  `Longitude` decimal(11,8) NOT NULL COMMENT "経度(-180.0 ~ 180.0)",
  `RecordedAt` datetime NOT NULL COMMENT "位置情報取得日時",
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "作成日時",
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新日時",
  PRIMARY KEY (`MonsterLocationId`),
  INDEX `idx_location` (`Latitude`, `Longitude`),
  INDEX `idx_monster_id` (`MonsterId`),
  INDEX `idx_recorded_at` (`RecordedAt`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin COMMENT "モンスターの位置情報";
-- Create "User" table
CREATE TABLE `User` (
  `UserId` varchar(36) NOT NULL COMMENT "ユーザーID(UUID)",
  `Nickname` varchar(50) NOT NULL COMMENT "ニックネーム",
  `CreatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "作成日時",
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新日時",
  PRIMARY KEY (`UserId`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin COMMENT "ユーザーの基本情報";
