ALTER TABLE `category` MODIFY `name` VARCHAR(31) NOT NULL;
ALTER TABLE `category` DROP COLUMN `user`;
ALTER TABLE `category` DROP KEY `user`;
ALTER TABLE `category` ADD UNIQUE `name` (`name`);
