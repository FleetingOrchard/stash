ALTER TABLE `performers` ADD COLUMN `ignore_auto_tag` boolean not null default '0'; 
ALTER TABLE `studios` ADD COLUMN `ignore_auto_tag` boolean not null default '0';
ALTER TABLE `tags` ADD COLUMN `ignore_auto_tag` boolean not null default '0';

CREATE TABLE `bookmarks` (
  `id` integer not null primary key autoincrement,
  `url` varchar(255) not null,
  `name` varchar(255),
  'position' integer not null
)