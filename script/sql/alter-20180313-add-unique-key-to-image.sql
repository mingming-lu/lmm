ALTER TABLE image DROP INDEX `url`;
ALTER TABLE image ADD UNIQUE (`user`, `name`);
