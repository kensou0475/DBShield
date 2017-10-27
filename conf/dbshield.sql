-- MySQL dump 10.13  Distrib 5.7.20, for osx10.13 (x86_64)
--
-- Host: 127.0.0.1    Database: dbshield
-- ------------------------------------------------------
-- Server version	5.7.18

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `pattern`
--

DROP TABLE IF EXISTS `pattern`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pattern` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` longtext,
  `value` longtext,
  `uuid` varchar(36) NOT NULL DEFAULT '',
  `example_value` longtext,
  `enable` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pattern`
--

LOCK TABLES `pattern` WRITE;
/*!40000 ALTER TABLE `pattern` DISABLE KEYS */;
INSERT INTO `pattern` VALUES (1,'0000e0040000002a0000e008746573740000003b','select * from test;','',NULL,0),(2,'0000e0040000002a0000e008746573740000003b5f757365725f74657374','11','',NULL,0),(3,'0000e0040000002a0000e008746573740000003b5f636c69656e745f3132372e302e302e31','11','',NULL,0),(4,'0000e0040000002a0000e0087465737445786973740000003b','select * from testExist;','',NULL,0),(5,'0000e0040000002a0000e0087465737445786973740000003b5f757365725f74657374','11','',NULL,0),(6,'0000e0040000002a0000e0087465737445786973740000003b5f636c69656e745f3132372e302e302e31','11','',NULL,0),(7,'0000e004404076657273696f6e5f636f6d6d656e740000e00e0000e033','select @@version_comment limit 1','1',NULL,0),(8,'0000e004404076657273696f6e5f636f6d6d656e740000e00e0000e0335f757365725f726f6f74','11','1',NULL,0),(9,'0000e004404076657273696f6e5f636f6d6d656e740000e00e0000e0335f636c69656e745f00000000000000000000ffff7f000001','11','1',NULL,0),(10,'0000e0040000002a0000e00866697273740000e0096e616d650000003c0000e033','select * from first where name<1719','1',NULL,0),(11,'0000e0040000002a0000e00866697273740000e0096e616d650000003c0000e0335f757365725f726f6f74','11','1',NULL,0),(12,'0000e0040000002a0000e00866697273740000e0096e616d650000003c0000e0335f636c69656e745f00000000000000000000ffff7f000001','11','1',NULL,0),(13,'0000e0670000e09b','show databases','1',NULL,0),(14,'0000e0670000e09b5f757365725f726f6f74','11','1',NULL,0),(15,'0000e0670000e09b5f636c69656e745f00000000000000000000ffff7f000001','11','1',NULL,0),(16,'0000e0040000002a0000e0086669727374','select * from first','1',NULL,0),(17,'0000e0040000002a0000e00866697273745f757365725f726f6f74','11','1',NULL,0),(18,'0000e0040000002a0000e00866697273745f636c69656e745f00000000000000000000ffff7f000001','11','1',NULL,0),(19,'424547494e','BEGIN','1',NULL,0),(20,'424547494e','BEGIN','1',NULL,0),(21,'424547494e','BEGIN','1',NULL,0),(22,'424547494e5f757365725f726f6f74','11','1',NULL,0),(23,'424547494e5f757365725f726f6f74','11','1',NULL,0),(24,'424547494e5f757365725f726f6f74','11','1',NULL,0),(25,'424547494e5f757365725f726f6f74','11','1',NULL,0),(26,'424547494e5f757365725f726f6f74','11','1',NULL,0),(27,'424547494e5f636c69656e745f00000000000000000000ffff7f000001','11','1',NULL,0),(28,'424547494e5f757365725f726f6f74','11','1',NULL,0),(29,'424547494e5f757365725f726f6f74','11','1',NULL,0),(30,'424547494e5f636c69656e745f00000000000000000000ffff7f000001','11','1',NULL,0),(31,'0000e0670000e09c','show tables','1','',0),(32,'0000e0670000e09c5f757365725f726f6f74','11','1','',0),(33,'0000e0670000e09c5f636c69656e745f00000000000000000000ffff7f000001','11','1','',0),(34,'0000e0040000002a0000e00866697273740000e0096e616d650000003c0000e0330000e03f69640000003d0000e033','select * from first where name<1719 and id=23','1','',0),(35,'0000e0040000002a0000e00866697273740000e0096e616d650000003c0000e0330000e03f69640000003d0000e0335f757365725f726f6f74','11','1','',0),(36,'0000e0040000002a0000e00866697273740000e0096e616d650000003c0000e0330000e03f69640000003d0000e0335f636c69656e745f00000000000000000000ffff7f000001','11','1','',0),(37,'0000e0040000e0a30000002800000029','SELECT DATABASE()','1','',1),(38,'0000e0040000e0a300000028000000295f757365725f726f6f74','11','1','',1),(39,'0000e0040000e0a300000028000000295f636c69656e745f00000000000000000000ffff7f000001','11','1','',1);
/*!40000 ALTER TABLE `pattern` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `db` varchar(128) DEFAULT NULL,
  `user` varchar(128) DEFAULT NULL,
  `client` varchar(128) DEFAULT NULL,
  `table` varchar(128) DEFAULT NULL,
  `permission` longtext NOT NULL,
  `enable` tinyint(1) NOT NULL DEFAULT '0',
  `uuid` varchar(36) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permission`
--

LOCK TABLES `permission` WRITE;
/*!40000 ALTER TABLE `permission` DISABLE KEYS */;
INSERT INTO `permission` VALUES (1,'dbshield','root','10.10.10.10','*','SELECT;CREATE',1,'1');
/*!40000 ALTER TABLE `permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `query_action`
--

DROP TABLE IF EXISTS `query_action`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `query_action` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `query` longtext,
  `user` varchar(128) DEFAULT NULL,
  `client` varchar(128) DEFAULT NULL,
  `db` varchar(128) DEFAULT NULL,
  `time` datetime NOT NULL,
  `action` varchar(32) NOT NULL DEFAULT '',
  `uuid` varchar(36) NOT NULL DEFAULT '',
  `duration` bigint(20) NOT NULL DEFAULT '0',
  `client_ip` varchar(39) DEFAULT NULL,
  `client_pm` varchar(128) DEFAULT NULL,
  `tables` longtext,
  `query_result` tinyint(1) NOT NULL DEFAULT '1',
  `is_abnormal` tinyint(1) NOT NULL DEFAULT '0',
  `abnormal_type` varchar(32) NOT NULL DEFAULT 'none',
  `is_alarm` tinyint(1) NOT NULL DEFAULT '0',
  `session_id` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `query_action`
--

LOCK TABLES `query_action` WRITE;
/*!40000 ALTER TABLE `query_action` DISABLE KEYS */;
INSERT INTO `query_action` VALUES (35,'select @@version_comment limit 1','root',NULL,'dbshield','2017-10-27 09:04:34','learning','1',4,'0.0.0.0','','',1,0,'',0,'zsxmRtYLmsHznlejFOAMwurqHBpzXDKl'),(36,'select @@version_comment limit 1','root',NULL,'dbshield','2017-10-27 09:05:24','learning','1',3,'0.0.0.0','','',1,0,'',0,'vRnlXHfaSLAJaMsszVBQNMjNdgttSTNi'),(37,'select @@version_comment limit 1','root',NULL,'dbshield','2017-10-27 09:05:40','learning','1',28,'0.0.0.0','','',1,0,'',0,'LuWSMYTaLckEfbmtaqLpverWdOTVmkBk'),(38,'show databases','root',NULL,'','2017-10-27 09:05:56','learning','1',24,'0.0.0.0','','',1,0,'',0,'tgzsykgYLdojhXqyYfmMBhjplhGfgLhR'),(39,'show databases','root',NULL,'','2017-10-27 09:06:02','learning','1',3,'0.0.0.0','','',1,0,'',0,'JlCDaLHsIasopUfxBbLAfOWCYPbtWvuY'),(40,'show databases','root',NULL,'','2017-10-27 09:06:03','learning','1',3,'0.0.0.0','','',1,0,'',0,'JlCDaLHsIasopUfxBbLAfOWCYPbtWvuY'),(41,'show databases','root',NULL,'','2017-10-27 09:06:04','learning','1',5,'0.0.0.0','','',1,0,'',0,'JlCDaLHsIasopUfxBbLAfOWCYPbtWvuY'),(42,'show databases','root',NULL,'','2017-10-27 09:06:05','learning','1',3,'0.0.0.0','','',1,0,'',0,'JlCDaLHsIasopUfxBbLAfOWCYPbtWvuY'),(43,'show databases','root',NULL,'','2017-10-27 09:06:07','learning','1',10,'0.0.0.0','','',1,0,'',0,'hRzfvMAfjcMhfDMJmvHfDpmxfkOTNtNr'),(44,'show databases','root',NULL,'','2017-10-27 09:06:07','learning','1',3,'0.0.0.0','','',1,0,'',0,'hRzfvMAfjcMhfDMJmvHfDpmxfkOTNtNr'),(45,'show databases','root',NULL,'','2017-10-27 09:06:08','learning','1',6,'0.0.0.0','','',1,0,'',0,'hRzfvMAfjcMhfDMJmvHfDpmxfkOTNtNr'),(46,'show databases','root',NULL,'','2017-10-27 09:06:09','learning','1',3,'0.0.0.0','','',1,0,'',0,'hRzfvMAfjcMhfDMJmvHfDpmxfkOTNtNr');
/*!40000 ALTER TABLE `query_action` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `state`
--

DROP TABLE IF EXISTS `state`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `state` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `key` varchar(5) NOT NULL DEFAULT '',
  `QueryCounter` bigint(20) unsigned NOT NULL DEFAULT '0',
  `AbnormalCounter` bigint(20) unsigned NOT NULL DEFAULT '0',
  `uuid` varchar(36) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `state`
--

LOCK TABLES `state` WRITE;
/*!40000 ALTER TABLE `state` DISABLE KEYS */;
INSERT INTO `state` VALUES (1,'state',13,0,'1');
/*!40000 ALTER TABLE `state` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-10-27 17:28:43
