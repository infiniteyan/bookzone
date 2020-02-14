CREATE DATABASE IF NOT EXISTS `bookzone`;
USE `bookzone`;

CREATE TABLE `md_attachment` (
    attachment_id   int(11) NOT NULL AUTO_INCREMENT,
    book_id         int(11) NOT NULL DEFAULT '0',
    document_id     int(11) NOT NULL DEFAULT '0',
    name            varchar(255) NOT NULL DEFAULT '',
    path            varchar(2000) NOT NULL DEFAULT '',
    size            double NOT NULL DEFAULT '0',
    ext             varchar(50) NOT NULL DEFAULT '',
    http_path       varchar(2000) NOT NULL DEFAULT '',
    create_time     datetime NOT NULL,
    create_at       int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (attachment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `md_book_category` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `book_id` int(11) NOT NULL DEFAULT '0',
    `category_id` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `book_id` (`book_id`,`category_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

LOCK TABLES `md_book_category` WRITE;
INSERT INTO `md_book_category` VALUES (1,1,1),(2,1,4);
UNLOCK TABLES;


CREATE TABLE `md_books` (
    `book_id` int(11) NOT NULL AUTO_INCREMENT,
    `book_name` varchar(500) NOT NULL DEFAULT '',
    `identify` varchar(100) NOT NULL DEFAULT '',
    `order_index` int(11) NOT NULL DEFAULT '0',
    `description` varchar(1000) NOT NULL DEFAULT '',
    `cover` varchar(1000) NOT NULL DEFAULT '',
    `editor` varchar(50) NOT NULL DEFAULT '',
    `status` int(11) NOT NULL DEFAULT '0',
    `privately_owned` int(11) NOT NULL DEFAULT '0',
    `private_token` varchar(500) DEFAULT NULL,
    `member_id` int(11) NOT NULL DEFAULT '0',
    `create_time` datetime NOT NULL,
    `modify_time` datetime NOT NULL,
    `release_time` datetime NOT NULL,
    `doc_count` int(11) NOT NULL DEFAULT '0',
    `comment_count` int(11) NOT NULL DEFAULT '0',
    `vcnt` int(11) NOT NULL DEFAULT '0',
    `star` int(11) NOT NULL DEFAULT '0',
    `score` int(11) NOT NULL DEFAULT '40',
    `cnt_score` int(11) NOT NULL DEFAULT '0',
    `cnt_comment` int(11) NOT NULL DEFAULT '0',
    `author` varchar(50) NOT NULL DEFAULT '',
    `author_url` varchar(1000) NOT NULL DEFAULT '',
    PRIMARY KEY (`book_id`),
    UNIQUE KEY `identify` (`identify`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

LOCK TABLES `md_books` WRITE;
INSERT INTO `md_books` VALUES (1,'演示','demo',0,'用于演示的书籍','/static/images/book.png','markdown',0,0,'',1,'2019-12-16 06:16:03','2019-12-16 06:16:03','2019-12-16 06:16:03',1,0,0,0,50,1,0,'','');
UNLOCK TABLES;


CREATE TABLE `md_category` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `pid` int(11) NOT NULL DEFAULT '0',
    `title` varchar(30) NOT NULL DEFAULT '',
    `intro` varchar(255) NOT NULL DEFAULT '',
    `icon` varchar(255) NOT NULL DEFAULT '',
    `cnt` int(11) NOT NULL DEFAULT '0',
    `sort` int(11) NOT NULL DEFAULT '0',
    `status` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `title` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;

LOCK TABLES `md_category` WRITE;
INSERT INTO `md_category` VALUES (1,0,'演示','','',1,0,1),(2,0,'后端','','',0,0,1),(3,0,'前端','','',0,0,1),(4,1,'Demo','','',1,0,1),(5,2,'Go','','',0,0,1),(6,2,'JAVA','','',0,0,1),(7,2,'PHP','','',0,0,1),(8,2,'NET','','',0,0,1),(9,2,'Python','','',0,0,1),(10,3,'HTML','','',0,0,1),(11,3,'CSS','','',0,0,1),(12,3,'JavaScript','','',0,0,1),(13,3,'框架','','',0,0,1);
UNLOCK TABLES;


CREATE TABLE `md_members` (
    `member_id` int(11) NOT NULL AUTO_INCREMENT,
    `account` varchar(30) NOT NULL DEFAULT '',
    `nickname` varchar(30) NOT NULL DEFAULT '',
    `password` varchar(255) NOT NULL DEFAULT '',
    `description` varchar(640) NOT NULL DEFAULT '',
    `email` varchar(100) NOT NULL DEFAULT '',
    `phone` varchar(20) DEFAULT 'null',
    `avatar` varchar(255) NOT NULL DEFAULT '',
    `role` int(11) NOT NULL DEFAULT '1',
    `status` int(11) NOT NULL DEFAULT '0',
    `create_time` datetime NOT NULL,
    `create_at` int(11) NOT NULL DEFAULT '0',
    `last_login_time` datetime DEFAULT NULL,
    PRIMARY KEY (`member_id`),
    UNIQUE KEY `account` (`account`),
    UNIQUE KEY `nickname` (`nickname`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

LOCK TABLES `md_members` WRITE;
INSERT INTO `md_members` VALUES (1,'admin','admin','6fVynJQW4iV-KmCfHPrFucWFxwBKfGB-OY6Gu-9_QsHEFoEqCmgj-M-RwvM6WoIirokO|15|ced0f3c3ba8a223007bd5da110af9c0a3d3985e3c451e80c59789d91|7fec678fcc990d025b378232314a5339e96b26cb55b4ac2b13010f4a8d23c6af','','admin@ziyoubiancheng.com','','/static/images/avatar.png',0,0,'2019-12-16 06:13:31',0,'2019-12-16 14:13:31'),(2,'user1','user1','4mSZoWt1u91t3q6tcSZwFdIMT1wFR9o8Qzo53NRIhmd2FYqschKLYQknxcAADlHdfWLJ|15|98b702a40e8da1402a477983ab3b8fbbf5215b5dc4f5df526af28aa5|7ace8c5c5a49594446197ead810e34c4959e9f72ebcdd64218ecbff23500c5cd','','user1@ziyoubiancheng.com','','/static/images/avatar.png',2,0,'2019-12-19 17:04:26',0,'2019-12-20 01:04:26');
UNLOCK TABLES;


CREATE TABLE `md_comments` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `uid` int(11) NOT NULL DEFAULT '0',
    `book_id` int(11) NOT NULL DEFAULT '0',
    `content` varchar(255) NOT NULL DEFAULT '',
    `time_create` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY `md_comments_uid` (`uid`),
    KEY `md_comments_book_id` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `md_comments_0000` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `uid` int(11) NOT NULL DEFAULT '0',
    `book_id` int(11) NOT NULL DEFAULT '0',
    `content` varchar(255) NOT NULL DEFAULT '',
    `time_create` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY `md_comments_uid` (`uid`),
    KEY `md_comments_book_id` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `md_comments_0001` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `uid` int(11) NOT NULL DEFAULT '0',
    `book_id` int(11) NOT NULL DEFAULT '0',
    `content` varchar(255) NOT NULL DEFAULT '',
    `time_create` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY `md_comments_uid` (`uid`),
    KEY `md_comments_book_id` (`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `md_documents` (
    `document_id` int(11) NOT NULL AUTO_INCREMENT,
    `document_name` varchar(500) NOT NULL DEFAULT '',
    `identify` varchar(100) DEFAULT 'null',
    `book_id` int(11) NOT NULL DEFAULT '0',
    `parent_id` int(11) NOT NULL DEFAULT '0',
    `order_sort` int(11) NOT NULL DEFAULT '0',
    `release` longtext,
    `create_time` datetime NOT NULL,
    `member_id` int(11) NOT NULL DEFAULT '0',
    `modify_time` datetime NOT NULL,
    `modify_at` int(11) NOT NULL DEFAULT '0',
    `version` bigint(20) NOT NULL DEFAULT '0',
    `vcnt` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`document_id`),
    UNIQUE KEY `book_id` (`book_id`,`identify`),
    KEY `md_documents_identify` (`identify`),
    KEY `md_documents_book_id_parent_id_order_sort` (`book_id`,`parent_id`,`order_sort`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

INSERT INTO `md_documents` VALUES (1,'空白文档','blank',1,0,0,'','2019-12-16 14:16:03',1,'2019-12-16 14:16:03',0,0,0);


CREATE TABLE `md_document_store` (
    `document_id` int(11) NOT NULL AUTO_INCREMENT,
    `markdown` longtext NOT NULL,
    `content` longtext NOT NULL,
    PRIMARY KEY (`document_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
INSERT INTO `md_document_store` VALUES (1, '', '');


CREATE TABLE `md_fans` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `member_id` int(11) NOT NULL DEFAULT '0',
    `fans_id` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `member_id` (`member_id`,`fans_id`),
    KEY `md_fans_fans_id` (`fans_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `md_relationship` (
    `relationship_id` int(11) NOT NULL AUTO_INCREMENT,
    `member_id` int(11) NOT NULL DEFAULT '0',
    `book_id` int(11) NOT NULL DEFAULT '0',
    `role_id` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`relationship_id`),
    UNIQUE KEY `member_id` (`member_id`,`book_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

INSERT INTO `md_relationship` VALUES (1,1,1,0);


CREATE TABLE `md_score` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `book_id` int(11) NOT NULL DEFAULT '0',
    `uid` int(11) NOT NULL DEFAULT '0',
    `score` int(11) NOT NULL DEFAULT '0',
    `time_create` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uid` (`uid`,`book_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `md_star` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `member_id` int(11) NOT NULL DEFAULT '0',
    `book_id` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `member_id` (`member_id`,`book_id`),
    KEY `md_star_member_id` (`member_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;