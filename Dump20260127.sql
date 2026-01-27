-- MySQL dump 10.13  Distrib 8.0.41, for Win64 (x86_64)
--
-- Host: 10.10.12.99    Database: cms_test
-- ------------------------------------------------------
-- Server version	8.0.45

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `academic_details`
--

DROP TABLE IF EXISTS `academic_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `academic_details` (
  `student_id` int DEFAULT NULL,
  `batch` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `year` int DEFAULT NULL,
  `semester` int DEFAULT NULL,
  `degree_level` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `section` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `department` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `student_category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `branch_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `seat_category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `regulation` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `quota` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `university` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `year_of_admission` int DEFAULT NULL,
  `year_of_completion` int DEFAULT NULL,
  `student_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `curriculum_id` int DEFAULT NULL,
  KEY `fk_academic_student` (`student_id`) USING BTREE,
  KEY `fk_academic_curriculum` (`curriculum_id`) USING BTREE,
  CONSTRAINT `fk_academic_curriculum` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_academic_student` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `academic_details`
--

LOCK TABLES `academic_details` WRITE;
/*!40000 ALTER TABLE `academic_details` DISABLE KEYS */;
INSERT INTO `academic_details` VALUES (3,'2024 - 2028',2026,5,'Bachelors','1','CSE','','','','','','',0,0,'',1),(10,'2024-2028',2024,1,NULL,NULL,'Computer Science and Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(11,'2024-2028',2024,1,NULL,NULL,'Computer Science and Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(12,'2024-2028',2024,1,NULL,NULL,'Computer Science and Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(13,'2024-2028',2024,1,NULL,NULL,'Computer Science and Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(14,'2024-2028',2024,1,NULL,NULL,'Computer Science and Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(20,'2024-2028',2024,1,NULL,NULL,'Information Technology',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(21,'2024-2028',2024,1,NULL,NULL,'Information Technology',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(22,'2024-2028',2024,1,NULL,NULL,'Information Technology',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(23,'2024-2028',2024,1,NULL,NULL,'Information Technology',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(24,'2024-2028',2024,1,NULL,NULL,'Information Technology',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(30,'2024-2028',2024,1,NULL,NULL,'Electronics and Communication Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(31,'2024-2028',2024,1,NULL,NULL,'Electronics and Communication Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(32,'2024-2028',2024,1,NULL,NULL,'Electronics and Communication Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1),(33,'2024-2028',2024,1,'wf','e','kmkjmjm','','','','','','',0,0,'',1),(34,'2024-2028',2024,1,NULL,NULL,'Electronics and Communication Engineering',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1);
/*!40000 ALTER TABLE `academic_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `address`
--

DROP TABLE IF EXISTS `address`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `address` (
  `student_id` int DEFAULT NULL,
  `permanent_address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `present_address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `residence_location` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `address_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `address`
--

LOCK TABLES `address` WRITE;
/*!40000 ALTER TABLE `address` DISABLE KEYS */;
INSERT INTO `address` VALUES (3,'48 Yadhavar ','48 Yadhavar ','48 Yadhavar '),(33,'ffv233vsdv','கேடகாரஅள்ளி ,4/178','');
/*!40000 ALTER TABLE `address` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `admission_payment`
--

DROP TABLE IF EXISTS `admission_payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_payment` (
  `student_id` int DEFAULT NULL,
  `dte_register_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `dte_admission_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `receipt_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `receipt_date` date DEFAULT NULL,
  `amount` decimal(10,2) DEFAULT NULL,
  `bank_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `admission_payment_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admission_payment`
--

LOCK TABLES `admission_payment` WRITE;
/*!40000 ALTER TABLE `admission_payment` DISABLE KEYS */;
/*!40000 ALTER TABLE `admission_payment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cluster_departments`
--

DROP TABLE IF EXISTS `cluster_departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cluster_departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cluster_id` int NOT NULL,
  `curriculum_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_department` (`curriculum_id`) USING BTREE,
  KEY `cluster_id` (`cluster_id`) USING BTREE,
  CONSTRAINT `cluster_departments_ibfk_1` FOREIGN KEY (`cluster_id`) REFERENCES `clusters` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cluster_departments`
--

LOCK TABLES `cluster_departments` WRITE;
/*!40000 ALTER TABLE `cluster_departments` DISABLE KEYS */;
/*!40000 ALTER TABLE `cluster_departments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `clusters`
--

DROP TABLE IF EXISTS `clusters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `clusters` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clusters`
--

LOCK TABLES `clusters` WRITE;
/*!40000 ALTER TABLE `clusters` DISABLE KEYS */;
/*!40000 ALTER TABLE `clusters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `co_po_mapping`
--

DROP TABLE IF EXISTS `co_po_mapping`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `co_po_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `co_index` int NOT NULL,
  `po_index` int NOT NULL,
  `mapping_value` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_copo_course` (`course_id`) USING BTREE,
  CONSTRAINT `fk_copo_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=86 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `co_po_mapping`
--

LOCK TABLES `co_po_mapping` WRITE;
/*!40000 ALTER TABLE `co_po_mapping` DISABLE KEYS */;
/*!40000 ALTER TABLE `co_po_mapping` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `co_pso_mapping`
--

DROP TABLE IF EXISTS `co_pso_mapping`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `co_pso_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `co_index` int NOT NULL,
  `pso_index` int NOT NULL,
  `mapping_value` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_copso_course` (`course_id`) USING BTREE,
  CONSTRAINT `fk_copso_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `co_pso_mapping`
--

LOCK TABLES `co_pso_mapping` WRITE;
/*!40000 ALTER TABLE `co_pso_mapping` DISABLE KEYS */;
/*!40000 ALTER TABLE `co_pso_mapping` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `contact_details`
--

DROP TABLE IF EXISTS `contact_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `contact_details` (
  `student_id` int DEFAULT NULL,
  `parent_mobile` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `student_mobile` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `student_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `parent_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `official_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `contact_details_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `contact_details`
--

LOCK TABLES `contact_details` WRITE;
/*!40000 ALTER TABLE `contact_details` DISABLE KEYS */;
INSERT INTO `contact_details` VALUES (3,'0830015502','0830015502','sanjeevirajanramajayam.cs24@bitsathy.ac.in','sanjeevirajanramajayam.cs24@bitsathy.ac.in','sanjeevirajanramajayam.cs24@bitsathy.ac.in'),(33,'0838879911','','mithul0605@gmail.com','mithul0605@gmail.com','mithul0605@gmail.com');
/*!40000 ALTER TABLE `contact_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_allocations`
--

DROP TABLE IF EXISTS `course_allocations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_allocations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `teacher_id` bigint unsigned NOT NULL,
  `academic_year` varchar(50) NOT NULL,
  `semester` int DEFAULT NULL,
  `section` varchar(10) DEFAULT 'A',
  `role` enum('Primary','Assistant') DEFAULT 'Primary',
  `status` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_assignment` (`course_id`,`teacher_id`,`academic_year`,`section`),
  KEY `fk_allocation_teacher` (`teacher_id`),
  CONSTRAINT `fk_allocation_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_allocation_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `teachers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_allocations`
--

LOCK TABLES `course_allocations` WRITE;
/*!40000 ALTER TABLE `course_allocations` DISABLE KEYS */;
INSERT INTO `course_allocations` VALUES (1,1,1,'2025-2026',1,'A','Primary',1,'2026-01-24 09:19:38','2026-01-24 09:19:38'),(2,2,2,'2025-2026',3,'A','Primary',1,'2026-01-24 09:19:38','2026-01-24 09:19:38'),(3,3,3,'2025-2026',1,'A','Primary',1,'2026-01-24 09:19:38','2026-01-24 09:19:38'),(4,4,4,'2025-2026',1,'A','Primary',1,'2026-01-24 09:19:38','2026-01-24 09:19:38'),(5,5,1,'2025-2026',1,'A','Primary',1,'2026-01-24 09:19:38','2026-01-24 09:19:38'),(6,5,3,'2025-2026',1,'A','Assistant',1,'2026-01-24 09:19:38','2026-01-24 09:19:38');
/*!40000 ALTER TABLE `course_allocations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_experiment_topics`
--

DROP TABLE IF EXISTS `course_experiment_topics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_experiment_topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `experiment_id` int NOT NULL,
  `topic_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `topic_order` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_exp_topics` (`experiment_id`) USING BTREE,
  CONSTRAINT `course_experiment_topics_ibfk_1` FOREIGN KEY (`experiment_id`) REFERENCES `course_experiments` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_experiment_topics`
--

LOCK TABLES `course_experiment_topics` WRITE;
/*!40000 ALTER TABLE `course_experiment_topics` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_experiment_topics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_experiments`
--

DROP TABLE IF EXISTS `course_experiments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_experiments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `experiment_number` int NOT NULL,
  `experiment_name` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `hours` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_course_exp` (`course_id`) USING BTREE,
  CONSTRAINT `course_experiments_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_experiments`
--

LOCK TABLES `course_experiments` WRITE;
/*!40000 ALTER TABLE `course_experiments` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_experiments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_objectives`
--

DROP TABLE IF EXISTS `course_objectives`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_objectives` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `objective` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_course_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `course_objectives_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_objectives`
--

LOCK TABLES `course_objectives` WRITE;
/*!40000 ALTER TABLE `course_objectives` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_objectives` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_outcomes`
--

DROP TABLE IF EXISTS `course_outcomes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_outcomes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `outcome` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_course_outcome_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `fk_course_outcomes_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=140 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_outcomes`
--

LOCK TABLES `course_outcomes` WRITE;
/*!40000 ALTER TABLE `course_outcomes` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_outcomes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_prerequisites`
--

DROP TABLE IF EXISTS `course_prerequisites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_prerequisites` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `prerequisite` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_course_prerequisite_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `fk_course_prerequisites_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_prerequisites`
--

LOCK TABLES `course_prerequisites` WRITE;
/*!40000 ALTER TABLE `course_prerequisites` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_prerequisites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_references`
--

DROP TABLE IF EXISTS `course_references`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_references` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `reference_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_course_reference_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `fk_course_references_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_references`
--

LOCK TABLES `course_references` WRITE;
/*!40000 ALTER TABLE `course_references` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_references` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_selflearning`
--

DROP TABLE IF EXISTS `course_selflearning`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_selflearning` (
  `course_id` int NOT NULL,
  `total_hours` int NOT NULL,
  PRIMARY KEY (`course_id`) USING BTREE,
  CONSTRAINT `fk_course_selflearning_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_selflearning`
--

LOCK TABLES `course_selflearning` WRITE;
/*!40000 ALTER TABLE `course_selflearning` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_selflearning` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_selflearning_resources`
--

DROP TABLE IF EXISTS `course_selflearning_resources`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_selflearning_resources` (
  `id` int NOT NULL AUTO_INCREMENT,
  `main_id` int NOT NULL,
  `internal_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_main_position` (`main_id`,`position`) USING BTREE,
  CONSTRAINT `course_selflearning_resources_ibfk_1` FOREIGN KEY (`main_id`) REFERENCES `course_selflearning_topics` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_selflearning_resources`
--

LOCK TABLES `course_selflearning_resources` WRITE;
/*!40000 ALTER TABLE `course_selflearning_resources` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_selflearning_resources` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_selflearning_topics`
--

DROP TABLE IF EXISTS `course_selflearning_topics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_selflearning_topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `main_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_course_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `course_selflearning_topics_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_selflearning_topics`
--

LOCK TABLES `course_selflearning_topics` WRITE;
/*!40000 ALTER TABLE `course_selflearning_topics` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_selflearning_topics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_teamwork`
--

DROP TABLE IF EXISTS `course_teamwork`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_teamwork` (
  `course_id` int NOT NULL,
  `total_hours` int NOT NULL,
  PRIMARY KEY (`course_id`) USING BTREE,
  CONSTRAINT `course_teamwork_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_teamwork`
--

LOCK TABLES `course_teamwork` WRITE;
/*!40000 ALTER TABLE `course_teamwork` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_teamwork` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `course_teamwork_activities`
--

DROP TABLE IF EXISTS `course_teamwork_activities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `course_teamwork_activities` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `activity` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `course_id` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `course_teamwork_activities_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `course_teamwork_activities`
--

LOCK TABLES `course_teamwork_activities` WRITE;
/*!40000 ALTER TABLE `course_teamwork_activities` DISABLE KEYS */;
/*!40000 ALTER TABLE `course_teamwork_activities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `courses`
--

DROP TABLE IF EXISTS `courses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `courses` (
  `course_id` int NOT NULL AUTO_INCREMENT,
  `course_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `course_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `course_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `credit` int DEFAULT NULL,
  `lecture_hrs` int DEFAULT '0',
  `tutorial_hrs` int DEFAULT '0',
  `practical_hrs` int DEFAULT '0',
  `activity_hrs` int DEFAULT '0',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `cia_marks` int DEFAULT '40',
  `see_marks` int DEFAULT '60',
  `total_marks` int GENERATED ALWAYS AS ((`cia_marks` + `see_marks`)) STORED,
  `theory_total_hrs` int DEFAULT '0',
  `tutorial_total_hrs` int DEFAULT '0',
  `practical_total_hrs` int DEFAULT NULL,
  `activity_total_hrs` int DEFAULT '0',
  `tw/sl` int DEFAULT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `curriculum_ref_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `total_hrs` int GENERATED ALWAYS AS ((((`theory_total_hrs` + `activity_total_hrs`) + `tutorial_total_hrs`) + coalesce(`practical_total_hrs`,0))) STORED,
  PRIMARY KEY (`course_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `courses`
--

LOCK TABLES `courses` WRITE;
/*!40000 ALTER TABLE `courses` DISABLE KEYS */;
INSERT INTO `courses` (`course_id`, `course_code`, `course_name`, `course_type`, `credit`, `lecture_hrs`, `tutorial_hrs`, `practical_hrs`, `activity_hrs`, `category`, `cia_marks`, `see_marks`, `theory_total_hrs`, `tutorial_total_hrs`, `practical_total_hrs`, `activity_total_hrs`, `tw/sl`, `visibility`, `source_curriculum_id`, `curriculum_ref_id`, `status`) VALUES (1,'CS101','Programming in C','Theory',3,0,0,0,0,NULL,40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(2,'CS102','Data Structures','Theory',4,0,0,0,0,NULL,40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(3,'IT101','Web Technologies','Theory',3,0,0,0,0,NULL,40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(4,'EC101','Digital Electronics','Theory',3,0,0,0,0,NULL,40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(5,'CS103L','Programming Lab','Practical',2,0,0,0,0,NULL,40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1);
/*!40000 ALTER TABLE `courses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum`
--

DROP TABLE IF EXISTS `curriculum`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `academic_year` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `curriculum_template` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '2026',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `max_credits` int DEFAULT '0',
  `status` tinyint(1) DEFAULT '1',
  `curriculum_ref_id` int DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum`
--

LOCK TABLES `curriculum` WRITE;
/*!40000 ALTER TABLE `curriculum` DISABLE KEYS */;
INSERT INTO `curriculum` VALUES (1,'R2025-2026 B.E CSE','2025-2026','2022','2026-01-23 09:37:12',163,1,NULL);
/*!40000 ALTER TABLE `curriculum` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum_courses`
--

DROP TABLE IF EXISTS `curriculum_courses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum_courses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `semester_id` int NOT NULL,
  `course_id` int NOT NULL,
  `status` tinyint(1) DEFAULT '1',
  `count_towards_limit` tinyint(1) DEFAULT '1' COMMENT 'Whether this course counts towards the curriculum max credit limit',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_rc_regulation` (`curriculum_id`) USING BTREE,
  KEY `fk_rc_semester` (`semester_id`) USING BTREE,
  KEY `fk_rc_course` (`course_id`) USING BTREE,
  CONSTRAINT `fk_rc_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_rc_regulation` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_rc_semester` FOREIGN KEY (`semester_id`) REFERENCES `normal_cards` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=241 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_courses`
--

LOCK TABLES `curriculum_courses` WRITE;
/*!40000 ALTER TABLE `curriculum_courses` DISABLE KEYS */;
/*!40000 ALTER TABLE `curriculum_courses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum_logs`
--

DROP TABLE IF EXISTS `curriculum_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `action` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `changed_by` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'System',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `diff` json DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `curriculum_id` (`curriculum_id`) USING BTREE,
  CONSTRAINT `curriculum_logs_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=339 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_logs`
--

LOCK TABLES `curriculum_logs` WRITE;
/*!40000 ALTER TABLE `curriculum_logs` DISABLE KEYS */;
INSERT INTO `curriculum_logs` VALUES (319,1,'Curriculum Created','Created new curriculum: R2025-2026 B.E CSE (2025-2026)','System','2026-01-23 09:37:12',NULL),(320,1,'Department Overview Created','Created department vision, mission, PEOs, POs, and PSOs','System','2026-01-23 09:38:00',NULL),(321,1,'Mission[0] Deleted','Deleted Mission item at index 0','System','2026-01-23 10:14:32','{\"Mission[0]\": {\"new\": \"\", \"old\": \"mission 1\"}}'),(322,1,'PO[0] Deleted','Deleted PO item at index 0','System','2026-01-23 10:14:32','{\"PO[0]\": {\"new\": \"\", \"old\": \"PO 1\"}}'),(323,1,'PEO[0] Deleted','Deleted PEO item at index 0','System','2026-01-23 10:14:32','{\"PEO[0]\": {\"new\": \"\", \"old\": \"PEO 1\"}}'),(324,1,'PSO[0] Deleted','Deleted PSO item at index 0','System','2026-01-23 10:14:32','{\"PSO[0]\": {\"new\": \"\", \"old\": \"PSO 1\"}}'),(325,1,'Vision Updated','Updated department vision','System','2026-01-23 10:15:12','{\"vision\": {\"new\": \"\", \"old\": \"vision\"}}'),(326,1,'Department Overview Created','Created department vision, mission, PEOs, POs, and PSOs','System','2026-01-23 10:16:00',NULL),(327,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-23 10:17:12','{\"Mission[0]\": {\"new\": \"mission\", \"old\": \"\"}}'),(328,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-23 10:26:27','{\"Mission[0]\": {\"new\": \"mission 1\", \"old\": \"\"}}'),(329,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-23 10:27:05','{\"Mission[0]\": {\"new\": \"Mission 1\", \"old\": \"\"}}'),(330,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-23 10:29:39','{\"Mission[0]\": {\"new\": \"Mission 1\", \"old\": \"\"}}'),(331,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-23 10:31:50','{\"Mission[0]\": {\"new\": \"Mission 1\", \"old\": \"\"}}'),(332,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-23 10:32:51','{\"Mission[0]\": {\"new\": \"Mission 1\", \"old\": \"\"}}'),(333,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-23 10:37:58','{\"Mission[0]\": {\"new\": \"Mission 1\", \"old\": \"\"}}'),(334,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-24 04:07:13','{\"Mission[0]\": {\"new\": \"Mission 1\", \"old\": \"\"}}'),(335,1,'Mission[0] Added','Added Mission item at index 0','System','2026-01-24 05:01:10','{\"Mission[0]\": {\"new\": \"MISSION 1\", \"old\": \"\"}}'),(336,1,'PEO[0] Added','Added PEO item at index 0','System','2026-01-24 05:05:27','{\"PEO[0]\": {\"new\": \"PEO 1\", \"old\": \"\"}}'),(337,1,'PO[0] Added','Added PO item at index 0','System','2026-01-24 05:10:37','{\"PO[0]\": {\"new\": \"PO 1\", \"old\": \"\"}}'),(338,1,'PSO[0] Added','Added PSO item at index 0','System','2026-01-24 05:12:51','{\"PSO[0]\": {\"new\": \"PSOSs\", \"old\": \"\"}}');
/*!40000 ALTER TABLE `curriculum_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum_mission`
--

DROP TABLE IF EXISTS `curriculum_mission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum_mission` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `mission_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_curriculum_id` (`curriculum_id`),
  CONSTRAINT `curriculum_mission_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_mission`
--

LOCK TABLES `curriculum_mission` WRITE;
/*!40000 ALTER TABLE `curriculum_mission` DISABLE KEYS */;
INSERT INTO `curriculum_mission` VALUES (49,1,'mission 1',0,'UNIQUE',NULL,0),(59,1,'MISSION 1',0,'UNIQUE',NULL,1);
/*!40000 ALTER TABLE `curriculum_mission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum_peos`
--

DROP TABLE IF EXISTS `curriculum_peos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum_peos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `peo_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_curriculum_id` (`curriculum_id`),
  CONSTRAINT `curriculum_peos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_peos`
--

LOCK TABLES `curriculum_peos` WRITE;
/*!40000 ALTER TABLE `curriculum_peos` DISABLE KEYS */;
INSERT INTO `curriculum_peos` VALUES (39,1,'PEO 1',0,'UNIQUE',NULL,0),(41,1,'PEO 1',0,'UNIQUE',NULL,1);
/*!40000 ALTER TABLE `curriculum_peos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum_pos`
--

DROP TABLE IF EXISTS `curriculum_pos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum_pos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `po_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_curriculum_id` (`curriculum_id`),
  CONSTRAINT `curriculum_pos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_pos`
--

LOCK TABLES `curriculum_pos` WRITE;
/*!40000 ALTER TABLE `curriculum_pos` DISABLE KEYS */;
INSERT INTO `curriculum_pos` VALUES (50,1,'PO 1',0,'UNIQUE',NULL,0),(52,1,'PO 1',0,'UNIQUE',NULL,1);
/*!40000 ALTER TABLE `curriculum_pos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum_psos`
--

DROP TABLE IF EXISTS `curriculum_psos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum_psos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `pso_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_curriculum_id` (`curriculum_id`),
  CONSTRAINT `curriculum_psos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_psos`
--

LOCK TABLES `curriculum_psos` WRITE;
/*!40000 ALTER TABLE `curriculum_psos` DISABLE KEYS */;
INSERT INTO `curriculum_psos` VALUES (23,1,'PSO 1',0,'UNIQUE',NULL,0),(25,1,'PSOSs',0,'UNIQUE',NULL,1);
/*!40000 ALTER TABLE `curriculum_psos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `curriculum_vision`
--

DROP TABLE IF EXISTS `curriculum_vision`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `curriculum_vision` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `vision` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` int DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `curriculum_vision_ibfk_1` (`curriculum_id`),
  CONSTRAINT `curriculum_vision_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_vision`
--

LOCK TABLES `curriculum_vision` WRITE;
/*!40000 ALTER TABLE `curriculum_vision` DISABLE KEYS */;
INSERT INTO `curriculum_vision` VALUES (16,1,'vision 1',1);
/*!40000 ALTER TABLE `curriculum_vision` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `department_curriculum`
--

DROP TABLE IF EXISTS `department_curriculum`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `department_curriculum` (
  `id` int NOT NULL AUTO_INCREMENT,
  `department_id` int NOT NULL,
  `curriculum_id` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_department_curriculum` (`department_id`,`curriculum_id`),
  KEY `idx_department` (`department_id`),
  KEY `idx_curriculum` (`curriculum_id`),
  CONSTRAINT `fk_dc_curriculum` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `fk_dc_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `department_curriculum`
--

LOCK TABLES `department_curriculum` WRITE;
/*!40000 ALTER TABLE `department_curriculum` DISABLE KEYS */;
/*!40000 ALTER TABLE `department_curriculum` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `department_teachers`
--

DROP TABLE IF EXISTS `department_teachers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `department_teachers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `teacher_id` bigint unsigned NOT NULL,
  `department_id` int NOT NULL,
  `role` varchar(100) DEFAULT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_department_teacher` (`teacher_id`,`department_id`),
  KEY `idx_teacher` (`teacher_id`),
  KEY `idx_department` (`department_id`),
  CONSTRAINT `fk_dt_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `fk_dt_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `teachers` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `department_teachers`
--

LOCK TABLES `department_teachers` WRITE;
/*!40000 ALTER TABLE `department_teachers` DISABLE KEYS */;
INSERT INTO `department_teachers` VALUES (5,7,2,NULL,1,'2026-01-24 06:49:33','2026-01-24 06:49:33'),(6,5,2,NULL,1,'2026-01-25 11:42:51','2026-01-25 11:42:51');
/*!40000 ALTER TABLE `department_teachers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `departments`
--

DROP TABLE IF EXISTS `departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `department_code` varchar(30) DEFAULT NULL,
  `department_name` varchar(255) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_departments_code` (`department_code`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `departments`
--

LOCK TABLES `departments` WRITE;
/*!40000 ALTER TABLE `departments` DISABLE KEYS */;
INSERT INTO `departments` VALUES (1,NULL,'Computer Science and Engineering',1,'2026-01-24 09:19:36','2026-01-24 09:19:36'),(2,NULL,'Information Technology',1,'2026-01-23 10:52:23','2026-01-23 10:52:23'),(3,NULL,'Electronics and Communication Engineering',1,'2026-01-24 09:19:36','2026-01-24 09:19:36');
/*!40000 ALTER TABLE `departments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `honour_cards`
--

DROP TABLE IF EXISTS `honour_cards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `honour_cards` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_regulation` (`curriculum_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `honour_cards`
--

LOCK TABLES `honour_cards` WRITE;
/*!40000 ALTER TABLE `honour_cards` DISABLE KEYS */;
/*!40000 ALTER TABLE `honour_cards` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `honour_vertical_courses`
--

DROP TABLE IF EXISTS `honour_vertical_courses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `honour_vertical_courses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `honour_vertical_id` int NOT NULL,
  `course_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_course_vertical` (`honour_vertical_id`,`course_id`) USING BTREE,
  KEY `course_id` (`course_id`) USING BTREE,
  KEY `idx_vertical` (`honour_vertical_id`) USING BTREE,
  CONSTRAINT `honour_vertical_courses_ibfk_1` FOREIGN KEY (`honour_vertical_id`) REFERENCES `honour_verticals` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `honour_vertical_courses_ibfk_2` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `honour_vertical_courses`
--

LOCK TABLES `honour_vertical_courses` WRITE;
/*!40000 ALTER TABLE `honour_vertical_courses` DISABLE KEYS */;
/*!40000 ALTER TABLE `honour_vertical_courses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `honour_verticals`
--

DROP TABLE IF EXISTS `honour_verticals`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `honour_verticals` (
  `id` int NOT NULL AUTO_INCREMENT,
  `honour_card_id` int NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_honour_card` (`honour_card_id`) USING BTREE,
  CONSTRAINT `honour_verticals_ibfk_1` FOREIGN KEY (`honour_card_id`) REFERENCES `honour_cards` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `honour_verticals`
--

LOCK TABLES `honour_verticals` WRITE;
/*!40000 ALTER TABLE `honour_verticals` DISABLE KEYS */;
/*!40000 ALTER TABLE `honour_verticals` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `hostel_details`
--

DROP TABLE IF EXISTS `hostel_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hostel_details` (
  `student_id` int DEFAULT NULL,
  `hosteller_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `hostel_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `room_no` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `room_capacity` int DEFAULT NULL,
  `room_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `floor_no` int DEFAULT NULL,
  `warden_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `alternate_warden` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `class_advisor` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `status` int DEFAULT '1',
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `hostel_details_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hostel_details`
--

LOCK TABLES `hostel_details` WRITE;
/*!40000 ALTER TABLE `hostel_details` DISABLE KEYS */;
/*!40000 ALTER TABLE `hostel_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `insurance_details`
--

DROP TABLE IF EXISTS `insurance_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `insurance_details` (
  `student_id` int DEFAULT NULL,
  `nominee_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `relationship` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `nominee_age` int DEFAULT NULL,
  `status` int DEFAULT '1',
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `insurance_details_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `insurance_details`
--

LOCK TABLES `insurance_details` WRITE;
/*!40000 ALTER TABLE `insurance_details` DISABLE KEYS */;
/*!40000 ALTER TABLE `insurance_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `normal_cards`
--

DROP TABLE IF EXISTS `normal_cards`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `normal_cards` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `semester_number` int DEFAULT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `card_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'semester',
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_semester_regulation` (`curriculum_id`) USING BTREE,
  CONSTRAINT `fk_semester_regulation` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `normal_cards`
--

LOCK TABLES `normal_cards` WRITE;
/*!40000 ALTER TABLE `normal_cards` DISABLE KEYS */;
/*!40000 ALTER TABLE `normal_cards` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `peo_po_mapping`
--

DROP TABLE IF EXISTS `peo_po_mapping`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `peo_po_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `peo_index` int NOT NULL,
  `po_index` int NOT NULL,
  `mapping_value` int NOT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_peopo_reg` (`curriculum_id`) USING BTREE,
  CONSTRAINT `fk_peopo_reg` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=738 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `peo_po_mapping`
--

LOCK TABLES `peo_po_mapping` WRITE;
/*!40000 ALTER TABLE `peo_po_mapping` DISABLE KEYS */;
/*!40000 ALTER TABLE `peo_po_mapping` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `regulation_clause_history`
--

DROP TABLE IF EXISTS `regulation_clause_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `regulation_clause_history` (
  `id` int NOT NULL AUTO_INCREMENT,
  `clause_id` int NOT NULL,
  `old_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `new_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `changed_by` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `changed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `change_reason` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `status` int DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `clause_id` (`clause_id`) USING BTREE,
  CONSTRAINT `regulation_clause_history_ibfk_1` FOREIGN KEY (`clause_id`) REFERENCES `regulation_clauses` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `regulation_clause_history`
--

LOCK TABLES `regulation_clause_history` WRITE;
/*!40000 ALTER TABLE `regulation_clause_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `regulation_clause_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `regulation_clauses`
--

DROP TABLE IF EXISTS `regulation_clauses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `regulation_clauses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `regulation_id` int NOT NULL,
  `section_no` int NOT NULL,
  `clause_no` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `regulation_id` (`regulation_id`) USING BTREE,
  CONSTRAINT `regulation_clauses_ibfk_1` FOREIGN KEY (`regulation_id`) REFERENCES `regulations` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `regulation_clauses`
--

LOCK TABLES `regulation_clauses` WRITE;
/*!40000 ALTER TABLE `regulation_clauses` DISABLE KEYS */;
/*!40000 ALTER TABLE `regulation_clauses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `regulation_sections`
--

DROP TABLE IF EXISTS `regulation_sections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `regulation_sections` (
  `id` int NOT NULL AUTO_INCREMENT,
  `regulation_id` int NOT NULL,
  `section_no` int NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `display_order` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_section` (`regulation_id`,`section_no`) USING BTREE,
  KEY `idx_regulation` (`regulation_id`) USING BTREE,
  KEY `idx_order` (`regulation_id`,`display_order`) USING BTREE,
  CONSTRAINT `regulation_sections_ibfk_1` FOREIGN KEY (`regulation_id`) REFERENCES `regulations` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `regulation_sections`
--

LOCK TABLES `regulation_sections` WRITE;
/*!40000 ALTER TABLE `regulation_sections` DISABLE KEYS */;
INSERT INTO `regulation_sections` VALUES (1,1,1,'ADMISSION',1,'2025-12-29 04:27:34','2025-12-29 04:27:34');
/*!40000 ALTER TABLE `regulation_sections` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `regulations`
--

DROP TABLE IF EXISTS `regulations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `regulations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` enum('DRAFT','PUBLISHED','LOCKED') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'DRAFT',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `regulations`
--

LOCK TABLES `regulations` WRITE;
/*!40000 ALTER TABLE `regulations` DISABLE KEYS */;
INSERT INTO `regulations` VALUES (1,'R2022','Academic Regulation 2022','DRAFT','2025-12-27 10:20:35','2025-12-27 10:20:35');
/*!40000 ALTER TABLE `regulations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `research_profiles`
--

DROP TABLE IF EXISTS `research_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `research_profiles` (
  `student_id` int DEFAULT NULL,
  `scopus_link` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `google_scholar_link` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `researchgate_link` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `orcid_link` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `h_index` int DEFAULT NULL,
  `status` int DEFAULT '1',
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `research_profiles_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `research_profiles`
--

LOCK TABLES `research_profiles` WRITE;
/*!40000 ALTER TABLE `research_profiles` DISABLE KEYS */;
/*!40000 ALTER TABLE `research_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `school_details`
--

DROP TABLE IF EXISTS `school_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_details` (
  `id` int NOT NULL AUTO_INCREMENT,
  `student_id` int DEFAULT NULL,
  `school_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `board` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `year_of_pass` int DEFAULT NULL,
  `state` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `tc_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `tc_date` date DEFAULT NULL,
  `total_marks` decimal(6,2) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `school_details_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_details`
--

LOCK TABLES `school_details` WRITE;
/*!40000 ALTER TABLE `school_details` DISABLE KEYS */;
INSERT INTO `school_details` VALUES (6,3,'KINGS SCHOOL','',0,'','','2026-01-23',100.00,1);
/*!40000 ALTER TABLE `school_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sharing_tracking`
--

DROP TABLE IF EXISTS `sharing_tracking`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sharing_tracking` (
  `id` int NOT NULL AUTO_INCREMENT,
  `source_curriculum_id` int NOT NULL,
  `target_curriculum_id` int NOT NULL,
  `item_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `source_item_id` int NOT NULL,
  `copied_item_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_source` (`source_curriculum_id`,`item_type`,`source_item_id`) USING BTREE,
  KEY `idx_target` (`target_curriculum_id`,`item_type`) USING BTREE,
  KEY `idx_copied` (`copied_item_id`,`item_type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sharing_tracking`
--

LOCK TABLES `sharing_tracking` WRITE;
/*!40000 ALTER TABLE `sharing_tracking` DISABLE KEYS */;
/*!40000 ALTER TABLE `sharing_tracking` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `students`
--

DROP TABLE IF EXISTS `students`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `students` (
  `student_id` int NOT NULL AUTO_INCREMENT,
  `enrollment_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `register_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `dte_reg_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `application_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `admission_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `student_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `gender` enum('Male','Female','Other') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `dob` date DEFAULT NULL,
  `age` int DEFAULT NULL,
  `father_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `mother_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `guardian_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `religion` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `nationality` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `community` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `mother_tongue` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `blood_group` varchar(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `aadhar_no` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `parent_occupation` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `designation` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `place_of_work` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `parent_income` decimal(10,2) DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  `department_id` int DEFAULT NULL,
  PRIMARY KEY (`student_id`) USING BTREE,
  KEY `fk_students_department` (`department_id`),
  CONSTRAINT `fk_students_department` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `students`
--

LOCK TABLES `students` WRITE;
/*!40000 ALTER TABLE `students` DISABLE KEYS */;
INSERT INTO `students` VALUES (3,'745','654','45623','456','456','test','Male','1998-02-22',27,'Ramajayam Velukonar','234234324','sdfsd','Muslim','Indian','OC','english','B+','223423423434','234234','sfsf','23423',23423494.00,1,NULL),(10,'CSE24001',NULL,NULL,NULL,NULL,'Alice Johnson','Female','2005-01-15',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,1),(11,'CSE24002',NULL,NULL,NULL,NULL,'Bob Smith','Male','2005-02-20',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,1),(12,'CSE24003',NULL,NULL,NULL,NULL,'Charlie Brown','Male','2005-03-25',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,1),(13,'CSE24004',NULL,NULL,NULL,NULL,'Diana Prince','Female','2005-04-30',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,1),(14,'CSE24005',NULL,NULL,NULL,NULL,'Evan Wright','Male','2005-05-05',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,1),(20,'IT24001',NULL,NULL,NULL,NULL,'Fiona Gallagher','Female','2005-06-10',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,2),(21,'IT24002',NULL,NULL,NULL,NULL,'George Miller','Male','2005-07-15',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,2),(22,'IT24003',NULL,NULL,NULL,NULL,'Hannah Abbott','Female','2005-08-20',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,2),(23,'IT24004',NULL,NULL,NULL,NULL,'Ian Malcolm','Male','2005-09-25',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,2),(24,'IT24005',NULL,NULL,NULL,NULL,'Julia Roberts','Female','2005-10-30',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,2),(30,'ECE24001',NULL,NULL,NULL,NULL,'Kevin Hart','Male','2005-11-05',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,3),(31,'ECE24002',NULL,NULL,NULL,NULL,'Laura Croft','Female','2005-12-10',19,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,3),(32,'ECE24003',NULL,NULL,NULL,NULL,'Mike Ross','Male','2006-01-15',18,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,3),(33,'ECE24004','234','12','211','21','Nina Simone','Female','2006-02-20',18,'qwd','weff','wefwef','Hindu','Indian','SCC','12','A-','211223423442','eddwe','wefwefwewef','weefe',323.00,1,3),(34,'ECE24005',NULL,NULL,NULL,NULL,'Oscar Wilde','Male','2006-03-25',18,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,1,3);
/*!40000 ALTER TABLE `students` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `syllabus`
--

DROP TABLE IF EXISTS `syllabus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `syllabus` (
  `id` int NOT NULL AUTO_INCREMENT,
  `model_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `position` int DEFAULT '0',
  `course_id` int NOT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `syllabus_models_fk_courses` (`course_id`) USING BTREE,
  CONSTRAINT `syllabus_models_fk_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `syllabus`
--

LOCK TABLES `syllabus` WRITE;
/*!40000 ALTER TABLE `syllabus` DISABLE KEYS */;
/*!40000 ALTER TABLE `syllabus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `syllabus_titles`
--

DROP TABLE IF EXISTS `syllabus_titles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `syllabus_titles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `model_id` int NOT NULL,
  `title_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `hours` int DEFAULT '0',
  `title` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `model_id` (`model_id`) USING BTREE,
  CONSTRAINT `syllabus_titles_ibfk_1` FOREIGN KEY (`model_id`) REFERENCES `syllabus` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `syllabus_titles`
--

LOCK TABLES `syllabus_titles` WRITE;
/*!40000 ALTER TABLE `syllabus_titles` DISABLE KEYS */;
/*!40000 ALTER TABLE `syllabus_titles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `syllabus_topics`
--

DROP TABLE IF EXISTS `syllabus_topics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `syllabus_topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title_id` int NOT NULL,
  `topic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `title_id` (`title_id`) USING BTREE,
  CONSTRAINT `syllabus_topics_ibfk_1` FOREIGN KEY (`title_id`) REFERENCES `syllabus_titles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `syllabus_topics`
--

LOCK TABLES `syllabus_topics` WRITE;
/*!40000 ALTER TABLE `syllabus_topics` DISABLE KEYS */;
/*!40000 ALTER TABLE `syllabus_topics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teachers`
--

DROP TABLE IF EXISTS `teachers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teachers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `email` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `phone` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `profile_img` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `dept` varchar(100) DEFAULT NULL,
  `desg` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `last_login` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE,
  UNIQUE KEY `email` (`email`) USING BTREE,
  KEY `fk_teachers_dept` (`dept`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teachers`
--

LOCK TABLES `teachers` WRITE;
/*!40000 ALTER TABLE `teachers` DISABLE KEYS */;
INSERT INTO `teachers` VALUES (1,'Dr. Alan Turing','alan@cse.edu','9876543210',NULL,'1','Professor','2026-01-24 09:19:36',1),(2,'Prof. Ada Lovelace','ada@cse.edu','9876543211',NULL,'1','Associate Professor','2026-01-24 09:19:36',1),(3,'Dr. Grace Hopper','grace@it.edu','9876543212',NULL,'2','Assistant Professor','2026-01-24 09:19:36',1),(4,'Prof. Claude Shannon','claude@ece.edu','9876543213',NULL,'3','Professor','2026-01-24 09:19:36',1),(5,'Dr. John von Neumann','john@ece.edu','9876543214','/uploads/teachers/teacher_1769341370.jpg','2','Associate Professor','2026-01-24 09:19:36',1),(7,'Teacher','TestTeacher@gmail.com','234234234','/uploads/teachers/teacher_1769225921.png','2','Lab Assistant','2026-01-23 10:52:23',1);
/*!40000 ALTER TABLE `teachers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `full_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `role` enum('admin','user') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'user',
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `last_login` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `username` (`username`) USING BTREE,
  UNIQUE KEY `email` (`email`) USING BTREE,
  KEY `idx_username` (`username`) USING BTREE,
  KEY `idx_email` (`email`) USING BTREE,
  KEY `idx_role` (`role`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin','$2a$10$H4BOz6nXYrnGQVYwC0eAQul.YF3LyhWTpb7xUf1HAKPT8y18DYDaq','System Administrator','admin@example.com','admin',1,'2026-01-07 05:52:45','2026-01-27 03:58:39','2026-01-27 09:28:38');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-01-27  9:38:14
