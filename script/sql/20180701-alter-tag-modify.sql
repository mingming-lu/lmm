ALTER TABLE `tag` DROP COLUMN `user`;
ALTER TABLE `tag` DROP KEY `user`;
ALTER TABLE `tag` ADD UNIQUE `blog_tag` (`blog`, `name`);
ALTER TABLE `tag` ADD INDEX `tag` (`name`);
