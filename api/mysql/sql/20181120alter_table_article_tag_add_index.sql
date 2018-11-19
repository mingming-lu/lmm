ALTER TABLE `article_tag` DROP INDEX `article_name`;
ALTER TABLE `article_tag` ADD INDEX `name` (`name`);
ALTER TABLE `article_tag` ADD INDEX `article_tag` (`article`, `sort`, `name`);
