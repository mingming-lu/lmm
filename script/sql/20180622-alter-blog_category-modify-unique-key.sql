ALTER TABLE `blog_category` DROP KEY `blog`;
ALTER TABLE `blog_category` DROP KEY `category`;
ALTER TABLE `blog_category` ADD UNIQUE KEY `blog_category` (`blog`, `category`);
-- re-alter by following sql
ALTER TABLE `blog_category` DROP KEY `blog_category`;
ALTER TABLE `blog_category` ADD UNIQUE KEY `blog` (`blog`);
ALTER TABLE `blog_category` ADD INDEX `blog` (`blog`);
