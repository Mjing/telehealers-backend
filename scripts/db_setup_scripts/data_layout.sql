-- MySQL dump 10.13  Distrib 8.0.31, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: telehealers
-- ------------------------------------------------------
-- Server version	8.0.30

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `advices`
--

DROP TABLE IF EXISTS `advices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `advices` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqNm` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `advices`
--

LOCK TABLES `advices` WRITE;
/*!40000 ALTER TABLE `advices` DISABLE KEYS */;
/*!40000 ALTER TABLE `advices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `appointment_prescription_map`
--

DROP TABLE IF EXISTS `appointment_prescription_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointment_prescription_map` (
  `appointment_id` int NOT NULL,
  `prescription_id` int DEFAULT NULL,
  PRIMARY KEY (`appointment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `appointment_prescription_map`
--

LOCK TABLES `appointment_prescription_map` WRITE;
/*!40000 ALTER TABLE `appointment_prescription_map` DISABLE KEYS */;
/*!40000 ALTER TABLE `appointment_prescription_map` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `appointments`
--

DROP TABLE IF EXISTS `appointments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `requested_start_time` time DEFAULT NULL,
  `requested_end_time` time DEFAULT NULL,
  `start_time` time DEFAULT NULL,
  `end_time` time DEFAULT NULL,
  `doctor_id` int NOT NULL,
  `patient_id` int NOT NULL,
  `prescription_id` int DEFAULT NULL,
  `patient_health_info_id` int DEFAULT NULL,
  `date` date DEFAULT NULL,
  `last_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_on` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `doctor_id` (`doctor_id`),
  KEY `patient_id` (`patient_id`),
  KEY `prescription_id` (`prescription_id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `appointments`
--

LOCK TABLES `appointments` WRITE;
/*!40000 ALTER TABLE `appointments` DISABLE KEYS */;
INSERT INTO `appointments` VALUES (1,NULL,NULL,NULL,NULL,0,0,NULL,NULL,'2022-11-14','2022-12-01 15:09:00','2022-12-01 15:09:18'),(2,'19:34:53',NULL,NULL,NULL,2,5,NULL,NULL,'2022-11-14','2022-12-01 15:09:00','2022-12-01 15:09:18'),(3,'19:36:07',NULL,NULL,NULL,2,5,NULL,NULL,'2022-11-14','2022-12-01 15:09:00','2022-12-01 15:09:18'),(4,'20:41:17',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(5,'20:47:08',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(6,'20:51:42',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(7,'20:54:13',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(8,'20:56:30',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(9,'20:57:52',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(10,'21:08:30',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(11,'21:10:13',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(12,'21:12:31',NULL,NULL,NULL,74,9,NULL,NULL,'2022-11-15','2022-12-01 15:09:00','2022-12-01 15:09:18'),(13,'21:13:50',NULL,'21:14:50','21:14:53',74,9,12,NULL,'2022-11-15','2022-12-01 18:41:34','2022-12-01 15:09:18'),(14,'14:30:37',NULL,'14:30:46','14:30:50',74,9,13,NULL,'2022-11-16','2022-12-01 18:44:36','2022-12-01 15:09:18'),(15,'22:17:48',NULL,'22:17:52','22:17:56',74,22,NULL,NULL,'2022-11-20','2022-12-01 15:09:00','2022-12-01 15:09:18'),(16,'16:11:57',NULL,NULL,NULL,74,9,NULL,NULL,'2022-12-03','2022-12-03 16:11:57','2022-12-03 16:11:57'),(17,'16:12:26',NULL,'16:12:30',NULL,74,9,NULL,NULL,'2022-12-03','2022-12-03 16:12:30','2022-12-03 16:12:26'),(18,'16:13:05',NULL,'16:13:15','16:13:38',74,9,NULL,NULL,'2022-12-03','2022-12-03 16:13:38','2022-12-03 16:13:05'),(19,'17:07:32',NULL,'17:07:55','17:07:57',74,9,NULL,NULL,'2022-12-03','2022-12-03 17:07:57','2022-12-03 17:07:32'),(20,'10:42:35',NULL,NULL,NULL,74,9,NULL,NULL,'2022-12-26','2022-12-26 10:42:35','2022-12-26 10:42:35'),(21,'10:45:13',NULL,'10:46:07',NULL,74,9,NULL,NULL,'2022-12-26','2022-12-26 10:46:07','2022-12-26 10:45:13'),(22,'14:49:57',NULL,NULL,NULL,74,9,NULL,NULL,'2022-12-26','2022-12-26 14:49:57','2022-12-26 14:49:57');
/*!40000 ALTER TABLE `appointments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `doctor_admission_applications`
--

DROP TABLE IF EXISTS `doctor_admission_applications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `doctor_admission_applications` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `approved` tinyint(1) DEFAULT '0',
  `comments` varchar(255) DEFAULT '',
  `reviewer_comments` varchar(255) DEFAULT '',
  `applied_on` datetime NOT NULL,
  `registration_number` varchar(100) NOT NULL,
  `password` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `registration_number` (`registration_number`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `doctor_admission_applications`
--

LOCK TABLES `doctor_admission_applications` WRITE;
/*!40000 ALTER TABLE `doctor_admission_applications` DISABLE KEYS */;
INSERT INTO `doctor_admission_applications` VALUES (5,'1','1',0,'','','2022-09-26 09:20:46','1','1');
/*!40000 ALTER TABLE `doctor_admission_applications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `doctor_languages`
--

DROP TABLE IF EXISTS `doctor_languages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `doctor_languages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(500) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `doctor_languages`
--

LOCK TABLES `doctor_languages` WRITE;
/*!40000 ALTER TABLE `doctor_languages` DISABLE KEYS */;
/*!40000 ALTER TABLE `doctor_languages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `doctor_specialities`
--

DROP TABLE IF EXISTS `doctor_specialities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `doctor_specialities` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` varchar(500) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `doctor_specialities`
--

LOCK TABLES `doctor_specialities` WRITE;
/*!40000 ALTER TABLE `doctor_specialities` DISABLE KEYS */;
INSERT INTO `doctor_specialities` VALUES (1,'General','General online physical treatment and advises');
/*!40000 ALTER TABLE `doctor_specialities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `doctor_to_specialities_map`
--

DROP TABLE IF EXISTS `doctor_to_specialities_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `doctor_to_specialities_map` (
  `doctor_id` int NOT NULL,
  `speciality_id` int NOT NULL,
  PRIMARY KEY (`doctor_id`,`speciality_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `doctor_to_specialities_map`
--

LOCK TABLES `doctor_to_specialities_map` WRITE;
/*!40000 ALTER TABLE `doctor_to_specialities_map` DISABLE KEYS */;
INSERT INTO `doctor_to_specialities_map` VALUES (74,1);
/*!40000 ALTER TABLE `doctor_to_specialities_map` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `doctors`
--

DROP TABLE IF EXISTS `doctors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `doctors` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT '',
  `email` varchar(255) DEFAULT '',
  `phone` varchar(255) DEFAULT '',
  `about` varchar(500) DEFAULT '',
  `profile_picture` int DEFAULT '0',
  `qualification` varchar(255) DEFAULT '',
  `password` varchar(100) NOT NULL,
  `registration_number` varchar(100) DEFAULT NULL,
  `sign_pic` int DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `registration_number` (`registration_number`),
  UNIQUE KEY `doctors_emails` (`email`),
  UNIQUE KEY `sign_pic` (`sign_pic`),
  UNIQUE KEY `profile_picture` (`profile_picture`)
) ENGINE=InnoDB AUTO_INCREMENT=75 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `doctors`
--

LOCK TABLES `doctors` WRITE;
/*!40000 ALTER TABLE `doctors` DISABLE KEYS */;
INSERT INTO `doctors` VALUES (74,'Dr asd def','1','','Testing profile&nbsp; <font size=\"6\">Updation Feature</font><br>',26,'','1','1',0);
/*!40000 ALTER TABLE `doctors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `file_store`
--

DROP TABLE IF EXISTS `file_store`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `file_store` (
  `id` int NOT NULL AUTO_INCREMENT,
  `created_on` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_id` int NOT NULL,
  `user_type` enum('doctor','patient','admin','guest','organization') NOT NULL,
  `path` varchar(255) NOT NULL,
  `file_tag` enum('profile_pic','sign_pic','else') DEFAULT 'else',
  PRIMARY KEY (`id`),
  UNIQUE KEY `path` (`path`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `file_store`
--

LOCK TABLES `file_store` WRITE;
/*!40000 ALTER TABLE `file_store` DISABLE KEYS */;
INSERT INTO `file_store` VALUES (1,'2022-12-13 17:47:46',23,'doctor','a/b/c','else'),(9,'2022-12-13 18:58:09',123,'guest','example.png','else'),(10,'2022-12-13 19:21:46',123,'guest','asd','else'),(11,'2022-12-17 18:22:10',74,'doctor','.//doctor/74/images.jpeg','else'),(22,'2022-12-21 07:11:04',74,'doctor','.//doctor/74/2831406983images.jpeg','else'),(23,'2022-12-21 08:20:05',74,'doctor','.//doctor/74/1579807097images.jpeg','else'),(24,'2022-12-21 08:43:44',74,'doctor','.//doctor/74/1760821654images.jpeg','else'),(25,'2022-12-21 09:05:43',74,'doctor','.//doctor/74/1612650118images.jpeg','else'),(26,'2022-12-21 09:10:52',74,'doctor','.//doctor/74/1918398626images.jpeg','else'),(27,'2022-12-21 18:58:25',9,'patient','.//patient/9/3722231965images.jpeg','else'),(28,'2022-12-21 18:58:32',9,'patient','.//patient/9/3446029827images.jpeg','else'),(29,'2022-12-21 19:01:28',9,'patient','.//patient/9/2468912083images.jpeg','else');
/*!40000 ALTER TABLE `file_store` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `helpdesk_tickets`
--

DROP TABLE IF EXISTS `helpdesk_tickets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `helpdesk_tickets` (
  `id` int NOT NULL AUTO_INCREMENT,
  `type` enum('med:service:update','query') DEFAULT 'query',
  `status` enum('new','inprogress','completed','declined') DEFAULT 'new',
  `description` varchar(1000) DEFAULT '',
  `created_by` int NOT NULL,
  `creator_type` enum('doctor','patient','admin','guest','organization') NOT NULL,
  `last_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `helpdesk_tickets`
--

LOCK TABLES `helpdesk_tickets` WRITE;
/*!40000 ALTER TABLE `helpdesk_tickets` DISABLE KEYS */;
INSERT INTO `helpdesk_tickets` VALUES (1,'med:service:update','new','Dr. Tets Update Page[id:0] requested specialization update\nSpecialization requested: Not working beyond 10.\nDescription:W/L balance',74,'doctor','2022-12-21 08:32:51'),(2,'med:service:update','new','Dr. Tets Update Page[id:0] requested specialization update\nSpecialization requested: Add 12th pass.\nDescription:Matriculated',74,'doctor','2022-12-21 08:32:51'),(3,'med:service:update','new','1[id:74] requested specialization update\nSpecialization requested: wad.\nDescription:weed',74,'doctor','2022-12-21 08:35:10'),(4,'med:service:update','new','1[id:74] requested specialization update\nSpecialization requested: wad.\nDescription:weed',74,'doctor','2022-12-21 08:36:30');
/*!40000 ALTER TABLE `helpdesk_tickets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `med_tests`
--

DROP TABLE IF EXISTS `med_tests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `med_tests` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqNm` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `med_tests`
--

LOCK TABLES `med_tests` WRITE;
/*!40000 ALTER TABLE `med_tests` DISABLE KEYS */;
INSERT INTO `med_tests` VALUES (17,'test-1[MODIFIED]','testcase 1'),(18,'test-12[MODIFIED]','testcase 12'),(19,'test-2[MODIFIED]','testcase 2'),(20,'test-1','testcase 1'),(21,'test-12','testcase 12'),(22,'test-2','testcase 2');
/*!40000 ALTER TABLE `med_tests` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `medicines`
--

DROP TABLE IF EXISTS `medicines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `medicines` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `description` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniqNm` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=96 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `medicines`
--

LOCK TABLES `medicines` WRITE;
/*!40000 ALTER TABLE `medicines` DISABLE KEYS */;
INSERT INTO `medicines` VALUES (85,'something-a-zole','Updated description'),(86,'Paracetamol','analgesic and antipyretic drug');
/*!40000 ALTER TABLE `medicines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `patient_health_info`
--

DROP TABLE IF EXISTS `patient_health_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `patient_health_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `gender` enum('M','F','NA') DEFAULT NULL,
  `height` varchar(255) DEFAULT NULL,
  `weight` varchar(255) DEFAULT NULL,
  `bp` varchar(255) DEFAULT NULL,
  `health_complaints` varchar(255) DEFAULT NULL,
  `patient_id` int NOT NULL,
  `created_on` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `info_to_patient` (`patient_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `patient_health_info`
--

LOCK TABLES `patient_health_info` WRITE;
/*!40000 ALTER TABLE `patient_health_info` DISABLE KEYS */;
/*!40000 ALTER TABLE `patient_health_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `patients`
--

DROP TABLE IF EXISTS `patients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `patients` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `phone` varchar(50) DEFAULT NULL,
  `date_of_birth` date DEFAULT NULL,
  `password` varchar(100) NOT NULL,
  `profile_picture` int DEFAULT NULL,
  `about` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email_2` (`email`),
  UNIQUE KEY `email` (`email`,`phone`),
  UNIQUE KEY `phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `patients`
--

LOCK TABLES `patients` WRITE;
/*!40000 ALTER TABLE `patients` DISABLE KEYS */;
INSERT INTO `patients` VALUES (9,'Patient MC Pat Face','1','1231241253',NULL,'1',29,NULL),(19,'1','2',NULL,NULL,'1',NULL,NULL),(20,'1','3',NULL,NULL,'1',NULL,NULL),(21,'1','5',NULL,NULL,'1',NULL,NULL),(22,'test','test',NULL,NULL,'test',NULL,NULL);
/*!40000 ALTER TABLE `patients` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prescription_to_advices_map`
--

DROP TABLE IF EXISTS `prescription_to_advices_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prescription_to_advices_map` (
  `prescription_id` int NOT NULL,
  `advice_id` int NOT NULL,
  `description` varchar(255) DEFAULT '',
  `last_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `prescription_id` (`prescription_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prescription_to_advices_map`
--

LOCK TABLES `prescription_to_advices_map` WRITE;
/*!40000 ALTER TABLE `prescription_to_advices_map` DISABLE KEYS */;
/*!40000 ALTER TABLE `prescription_to_advices_map` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prescription_to_medicines_map`
--

DROP TABLE IF EXISTS `prescription_to_medicines_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prescription_to_medicines_map` (
  `prescription_id` int NOT NULL,
  `medicine_id` int NOT NULL,
  `description` varchar(255) DEFAULT '',
  `last_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `prescription_id` (`prescription_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prescription_to_medicines_map`
--

LOCK TABLES `prescription_to_medicines_map` WRITE;
/*!40000 ALTER TABLE `prescription_to_medicines_map` DISABLE KEYS */;
INSERT INTO `prescription_to_medicines_map` VALUES (7,86,'','2022-12-01 16:18:55'),(8,86,'','2022-12-01 16:18:55'),(9,86,'','2022-12-01 16:18:55');
/*!40000 ALTER TABLE `prescription_to_medicines_map` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prescription_to_tests_map`
--

DROP TABLE IF EXISTS `prescription_to_tests_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prescription_to_tests_map` (
  `prescription_id` int NOT NULL,
  `test_id` int NOT NULL,
  `description` varchar(255) DEFAULT '',
  `last_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `prescription_id` (`prescription_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prescription_to_tests_map`
--

LOCK TABLES `prescription_to_tests_map` WRITE;
/*!40000 ALTER TABLE `prescription_to_tests_map` DISABLE KEYS */;
INSERT INTO `prescription_to_tests_map` VALUES (7,86,'','2022-12-01 16:19:14'),(8,86,'','2022-12-01 16:19:14'),(9,86,'','2022-12-01 16:19:14');
/*!40000 ALTER TABLE `prescription_to_tests_map` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prescriptions`
--

DROP TABLE IF EXISTS `prescriptions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prescriptions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `comment_on_medicines` varchar(500) DEFAULT (_utf8mb4''),
  `comment_on_tests` varchar(500) DEFAULT (_utf8mb4''),
  `comment_on_advices` varchar(500) DEFAULT (_utf8mb4''),
  `name` varchar(255) DEFAULT NULL,
  `last_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_on` datetime DEFAULT CURRENT_TIMESTAMP,
  `description` varchar(500) DEFAULT '',
  `created_by` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  FULLTEXT KEY `name_2` (`name`,`description`),
  FULLTEXT KEY `name_3` (`name`,`description`,`comment_on_medicines`,`comment_on_tests`,`comment_on_advices`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prescriptions`
--

LOCK TABLES `prescriptions` WRITE;
/*!40000 ALTER TABLE `prescriptions` DISABLE KEYS */;
INSERT INTO `prescriptions` VALUES (10,'','','ROj kha BC<br>',NULL,'2022-12-04 04:38:56','2022-12-01 16:47:22','',0),(11,'','','asd fsda<br>',NULL,'2022-12-04 04:38:56','2022-12-01 16:50:58','',0),(12,'','zxcv','asd <br>',NULL,'2022-12-04 04:38:56','2022-12-01 18:41:34','',0),(13,'def','asd','zvxz',NULL,'2022-12-04 04:38:56','2022-12-01 18:44:36','',0);
/*!40000 ALTER TABLE `prescriptions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_availability_status`
--

DROP TABLE IF EXISTS `user_availability_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_availability_status` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_type` enum('doctor','patient','admin','guest','organization') NOT NULL,
  `user_id` int NOT NULL,
  `session_id` varchar(50) NOT NULL,
  `status` enum('online','offline') DEFAULT 'offline',
  `last_login` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `session_id` (`session_id`),
  UNIQUE KEY `unique_users` (`user_type`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=133 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_availability_status`
--

LOCK TABLES `user_availability_status` WRITE;
/*!40000 ALTER TABLE `user_availability_status` DISABLE KEYS */;
INSERT INTO `user_availability_status` VALUES (1,'doctor',123,'123','online','2022-12-13 17:58:46'),(2,'doctor',74,'8b69fcab-8531-11ed-966f-0242ac120002','offline','2022-12-26 15:25:35'),(59,'patient',9,'99171e69-8546-11ed-966f-0242ac120002','offline','2022-12-26 17:56:18');
/*!40000 ALTER TABLE `user_availability_status` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-01-01 10:26:34
