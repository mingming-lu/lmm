ALTER TABLE `image` CHANGE `name` `uid` VARCHAR(255) NOT NULL; 
ALTER TABLE `image` DROP KEY `user`;
ALTER TABLE `image` ADD UNIQUE `uid` (`uid`);
ALTER TABLE `image` ADD INDEX `created_at` (`created_at`);
