ALTER TABLE `blog_category` DROP KEY `blog`;
ALTER TABLE `blog_category` DROP KEY `category`;
ALTER TABLE `blog_category` ADD UNIQUE KEY `blog_category` (`blog`, `category`);
