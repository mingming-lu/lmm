ALTER TABLE `user` MODIFY `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT;
ALTER TABLE `user` MODIFY `name` VARCHAR(31) NOT NULL;
ALTER TABLE `user` MODIFY `token` VARCHAR(63) NOT NULL;