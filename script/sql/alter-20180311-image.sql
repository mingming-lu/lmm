ALTER TABLE image DROP COLUMN `type`;
ALTER TABLE image CHANGE `url` `name` varchar(255);
