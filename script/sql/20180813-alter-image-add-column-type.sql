ALTER TABLE `image` ADD `type` tinyint NOT NULL DEFAULT 0;
ALTER TABLE `image` ADD INDEX `type` (`type`, `created_at`);
