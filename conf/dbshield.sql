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
  `example_value` longtext,
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `uuid` varchar(36) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

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
  `enable` tinyint(1) NOT NULL DEFAULT '1',
  `uuid` varchar(36) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `query_action`
--

DROP TABLE IF EXISTS `query_action`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `query_action` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `flow_id` varchar(32) DEFAULT NULL,
  `flow_info` longtext,
  `query` longtext,
  `user` varchar(128) DEFAULT NULL,
  `client_ip` varchar(39) DEFAULT NULL,
  `client_pm` varchar(128) DEFAULT NULL,
  `server_ip` varchar(39) DEFAULT NULL,
  `server_port` int(11) DEFAULT NULL,
  `db` varchar(128) DEFAULT NULL,
  `tables` longtext,
  `time` datetime NOT NULL,
  `duration` bigint(20) NOT NULL DEFAULT '0',
  `query_result` tinyint(1) NOT NULL DEFAULT '1',
  `is_abnormal` tinyint(1) NOT NULL DEFAULT '0',
  `abnormal_type` varchar(32) NOT NULL DEFAULT 'none',
  `action` varchar(36) NOT NULL DEFAULT '',
  `is_alarm` tinyint(1) NOT NULL DEFAULT '0',
  `analysed` tinyint(1) NOT NULL DEFAULT '0',
  `uuid` varchar(36) NOT NULL DEFAULT '',
  `sql_type` varchar(32) DEFAULT NULL,
  `tool` varchar(32) DEFAULT NULL,
  `pattern` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

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
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-10-30 14:24:57
