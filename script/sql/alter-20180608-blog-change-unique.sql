ALTER TABLE `blog` DROP KEY `user`;
ALTER TABLE `blog` ADD UNIQUE `title` (`title`);