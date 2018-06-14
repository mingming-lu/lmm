ALTER TABLE photo ADD deleted tinyint NOT NULL DEFAULT 0;
ALTER TABLE photo DROP INDEX `user`;
ALTER TABLE photo ADD UNIQUE `image` (`image`);
