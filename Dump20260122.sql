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
  `residence_location` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `address_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `address`
--

LOCK TABLES `address` WRITE;
/*!40000 ALTER TABLE `address` DISABLE KEYS */;
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
INSERT INTO `cluster_departments` VALUES (11,1,4,'2026-01-21 22:55:11',1),(12,1,14,'2026-01-21 22:59:21',1),(13,1,15,'2026-01-21 23:00:01',1),(14,1,10,'2026-01-21 23:00:16',1);
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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clusters`
--

LOCK TABLES `clusters` WRITE;
/*!40000 ALTER TABLE `clusters` DISABLE KEYS */;
INSERT INTO `clusters` VALUES (1,'computer cluster','cse departments','2025-12-25 05:07:49',1),(2,'mechanical cluster','mechanical departments','2025-12-25 05:30:51',1);
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
INSERT INTO `co_po_mapping` VALUES (1,1,0,1,3),(2,1,0,2,2),(3,1,0,3,1),(4,2,0,1,3),(5,2,0,2,2),(6,2,0,3,1),(7,2,0,4,1),(8,2,0,5,2),(9,2,0,9,1),(10,2,0,10,1),(11,2,0,12,2),(12,2,1,1,2),(13,2,1,2,3),(14,2,1,3,2),(15,2,1,4,2),(16,2,1,5,1),(17,2,1,9,1),(18,2,1,10,1),(19,2,1,12,3),(20,2,2,1,1),(21,2,2,2,2),(22,2,2,3,3),(23,2,2,4,2),(24,2,2,5,3),(25,2,2,7,1),(26,2,2,9,1),(27,2,2,12,3),(28,4,0,1,1),(29,4,0,2,2),(30,4,0,3,1),(31,4,0,4,3),(32,4,0,5,1),(33,4,0,6,2),(34,4,0,7,1),(35,4,0,10,2),(36,4,0,11,3),(37,4,0,12,1),(38,4,1,1,2),(39,4,1,2,2),(40,4,1,3,1),(41,4,1,4,3),(42,4,1,6,1),(43,4,1,9,2),(44,4,1,10,3),(45,4,1,11,1),(46,4,1,12,2);
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
INSERT INTO `co_pso_mapping` VALUES (1,1,0,1,3),(2,1,0,2,2),(3,2,0,1,3),(4,2,0,2,2),(5,2,0,3,1),(6,2,1,1,2),(7,2,1,2,3),(8,2,1,3,1),(9,2,2,1,2),(10,2,2,2,3),(11,2,2,3,2),(12,4,0,1,2),(13,4,0,2,3),(14,4,0,3,1),(15,4,1,1,2),(16,4,1,2,1),(17,4,1,3,3);
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
/*!40000 ALTER TABLE `contact_details` ENABLE KEYS */;
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
INSERT INTO `course_experiment_topics` VALUES (12,4,'Assess the physical parameters of different materials for engineering applications like radius, thickness and\ndiameter to design the electrical wires, bridges and clothes',0,'2026-01-07 04:25:43'),(13,4,'Evaluate the elastic nature of different solid materials for modern industrial applications like shock absorbers\nof vehicles',1,'2026-01-07 04:25:43');
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
INSERT INTO `course_experiments` VALUES (4,83,1,'Experiment 1',3,'2026-01-07 04:16:38','2026-01-07 04:25:43'),(5,83,2,'Experiment 2',7,'2026-01-07 04:20:01','2026-01-07 04:20:01'),(6,83,3,'Experiment 3',4,'2026-01-07 04:20:11','2026-01-07 04:20:11'),(10,83,4,'Experiment 4',8,'2026-01-07 04:30:23','2026-01-07 04:30:23');
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
INSERT INTO `course_objectives` VALUES (4,83,'To impart mathematical modeling to describe and explore real-world phenomena and data.',0),(5,83,'To provide basic understanding on Linear, quadratic, power and polynomial, exponential, and multi variable models',1),(6,83,'Summarize and apply the methodologies involved in framing the real world problems related to fundamental principles of polynomial equations',2);
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
INSERT INTO `course_outcomes` VALUES (136,83,'Implement the concepts of mathematical modeling based on linear functions in Engineering.',0),(137,83,'Formulate the real-world problems as a quadratic function model',1),(138,83,'Demonstrate the real-world phenomena and data into Power and Polynomial functions',2),(139,83,'Apply the concept of mathematical modeling of exponential functions in Engineering',3);
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
INSERT INTO `course_selflearning` VALUES (22,0),(83,0);
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
INSERT INTO `course_teamwork` VALUES (22,0),(83,0);
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
INSERT INTO `course_teamwork_activities` VALUES (27,22,'hello',0);
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
) ENGINE=InnoDB AUTO_INCREMENT=104 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `courses`
--

LOCK TABLES `courses` WRITE;
/*!40000 ALTER TABLE `courses` DISABLE KEYS */;
INSERT INTO `courses` (`course_id`, `course_code`, `course_name`, `course_type`, `credit`, `lecture_hrs`, `tutorial_hrs`, `practical_hrs`, `activity_hrs`, `category`, `cia_marks`, `see_marks`, `theory_total_hrs`, `tutorial_total_hrs`, `practical_total_hrs`, `activity_total_hrs`, `tw/sl`, `visibility`, `source_curriculum_id`, `curriculum_ref_id`, `status`) VALUES (1,'CS101','Introduction to Programming','Theory',3,0,0,0,0,'Core',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(2,'CS3801','Cloud Computing','Theory',3,0,0,0,0,'Elective',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(3,'CS201','Data Structures','Theory',4,3,1,0,0,'Core',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(4,'CS3501','Database Management Systems','Theory',4,3,1,2,0,'Core',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(17,'22MA101','ENGINEERING MATHEMATICS I ','Theory',4,0,10,3,1,'BS - Basic Sciences',40,60,2,0,NULL,0,NULL,'CLUSTER',NULL,NULL,1),(18,'22PH102 ','ENGINEERING PHYSICS ','Theory',3,0,10,2,3,'BS - Basic Sciences',40,60,1,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(19,'22CSH01 ','EXPLORATORY DATA ANALYSIS ','Theory',3,2,0,2,300,'PE - Professional Elective',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(20,'26MA101','Linear Algebra and Calculus','Theory',4,3,1,0,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(21,'26PH102','Engineering Physics','Theory',3,3,1,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(22,'26CH103','Engineering Chemistry','Theory',2,2,0,0,0,'ES - Engineering Sciences',30,70,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(23,'26GE004','Digital Computer Electronics','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(24,'26GE005','Problem Solving using C','Theory',3,2,0,0,0,'ES - Engineering Sciences',30,70,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(25,'26HS001','Communicative English','Theory',2,2,0,0,0,'HSS - Humanities and Social Sciences',30,70,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(26,'26HS002','தமிழர் மரபு / Heritage of Tamils ','Theory',1,1,0,0,0,'HSS - Humanities and Social Sciences',15,85,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(27,'26PH108','Physical Science Laboratory','Experiment',2,0,0,4,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(28,'26GE006','C Programming Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(29,'26SD001','Skill Development Course  I','Experiment',1,0,0,2,0,'EEC - Employability Enhancement Course',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(30,'26MA201','Differential Equations and Transforms','Theory',4,3,1,0,0,'ES - Engineering Sciences',60,40,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(31,'26PH202','Materials Science','Theory',3,3,0,0,0,'BS - Basic Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(32,'26CS203','Fundamentals of Web Principles','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(33,'26CS204','Computer Organization and Architecture','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(34,'26GE007','Python Programming','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(35,'26HS005','Professional Communication','Theory',2,2,0,0,0,'HSS - Humanities and Social Sciences',30,70,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(36,'26HS006','தமிழரும் தொழில்நுட்பமும் / Tamils and Technology','Theory',1,1,0,0,0,'HSS - Humanities and Social Sciences',15,85,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(37,'26GE008','Python Programming Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(38,'26CS209','Web Principles Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(39,'26SD002','Skill Development Course II','Experiment',1,0,0,2,0,'EEC - Employability Enhancement Course',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(40,'26CS301','Discrete Mathematics','Theory',4,3,1,0,0,'BS - Basic Sciences',60,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(41,'26CS302','Data Structures and Algorithms','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(42,'26CS303','Operating Systems','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(43,'26CS304','Object Oriented Programming with Java','Theory',3,2,0,2,0,'ES',30,70,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(44,'26CS305','Software Engineering','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(45,'26CS306','Database Management Systems','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(46,'26CS307','Standards in Computer Science','Theory',1,1,0,0,0,'ES - Engineering Sciences',15,85,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(47,'26CS308','Data Structures and Algorithms Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(48,'26CS309','Database Management Systems Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(49,'26CS310','Design Thinking and Innovation Laboratory (AICTE, & NEP)','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(50,'26CS401','Probability and Statistics','Theory',4,3,1,0,0,'ES - Engineering Sciences',60,40,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(51,'26CS402','Full Stack Development','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(52,'26CS403','Artificial Intelligence Essentials','Theory',3,3,0,0,0,'ES',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(53,'26CS404','Design and Analysis of Algorithms','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(54,'26CS405','Theory of Computation','Theory',4,3,1,0,0,'ES - Engineering Sciences',60,40,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(55,'26CS406','Computer Networks','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(56,'26HS009','Environmental Sciences and Sustainability ','Theory',2,2,0,0,0,'HSS - Humanities and Social Sciences',30,70,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(57,'26CS408','Full Stack Development Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(58,'26CS409','Computer Networks Laboratory','Lab',1,0,0,2,0,'ES',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(59,'26CS410','Community Engagement Project','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(60,'26CS501','Compiler Design','Theory',4,3,1,0,0,'ES - Engineering Sciences',60,40,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(61,'26CS502','Cloud Infrastructure Services','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(62,'26CS503','Bigdata Analytics','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(63,'26CS504','Machine Learning Essentials','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(64,'26XXIV','Professional Elective IV','Theory',3,0,0,0,0,'ES',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(65,'26CS507','Cloud Infrastructure Services Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(66,'26CS508','Machine Learning Essentials Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(67,'26CS509','Technology Integration Project','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(68,'26CS601','Software Project Management and Quality Assurance','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(69,'26CS602','Deep Learning','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(70,'26CS603','Cryptography and Cyber Security','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(71,'26CS607','Software Project Management and Quality Assurance Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(72,'26CS608','Deep Learning Laboratory ','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(73,'26CS609','Innovation and Product Development Project / Industry Oriented Course / Summer Internship','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(74,'26CS701','Generative AI and Large Language Models','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(75,'26CS702','IoT and Edge Computing','Theory',3,3,0,0,0,'ES - Engineering Sciences',45,55,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(76,'26XXIV','Professional Elective IV','Theory',3,0,0,0,0,'ES',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(77,'26XXV','Professional Elective V','Theory',3,0,0,0,0,'ES',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(78,'26XXVI','Professional Elective VI','Theory',3,0,0,0,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(79,'26CS706','Generative AI and Large Language Models Laboratory','Experiment',1,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(80,'26CS707','Capstone Project work Level I / Internship Pro','Experiment',3,0,0,6,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(81,'26CS801','Capstone Project Work Level II / Internship Project / Startup Product','Experiment',8,0,0,16,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(82,'hello','efvwbg','Theory',2,1,0,2,0,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(83,'wcnejce','ervewvwrv','Theory',3,0,0,0,0,'BS - Basic Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(84,'bdzhgd','bgbdfb','Theory',2,1,134,0,0,'ES - Engineering Sciences',40,60,13,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(85,'CS130','check1','Theory',3,3,15,0,30,'BS - Basic Sciences',40,60,45,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(86,'CS230','check 2','Theory',3,3,1,0,2,'ES - Engineering Sciences',40,60,0,0,NULL,0,NULL,'UNIQUE',NULL,NULL,1),(87,'CS303','check 3','Theory',3,3,1,0,2,'BS - Basic Sciences',40,60,45,15,NULL,30,NULL,'UNIQUE',NULL,NULL,1),(88,'cs102323','check4','Lab',3,0,0,3,0,'BS - Basic Sciences',40,60,0,0,45,0,10,'UNIQUE',NULL,NULL,1),(89,'CS134','check5','Theory&Lab',3,2,1,1,0,'ES - Engineering Sciences',40,60,30,15,15,0,0,'UNIQUE',NULL,NULL,1),(90,'CS140','1check2022','Theory',3,1,2,0,0,'ES - Engineering Sciences',40,60,15,30,0,0,0,'UNIQUE',NULL,NULL,1),(91,'CS120','2Check2022','Lab',32,1,2,1,0,'BS - Basic Sciences',40,60,15,30,15,0,0,'UNIQUE',NULL,NULL,1),(92,'CS150','3maincheck','Lab',3,2,2,4,0,'PE - Professional Elective',40,60,30,30,60,0,0,'UNIQUE',NULL,NULL,1),(93,'CS1','check2026theory','Theory',3,1,2,0,3,'BS - Basic Sciences',40,60,15,30,0,45,0,'UNIQUE',NULL,NULL,1),(94,'CS2','check2026lab','Lab',3,1,2,3,1,'ES - Engineering Sciences',40,60,15,30,45,15,10,'UNIQUE',NULL,NULL,1),(95,'CS3','check 123','Lab',3,1,2,3,3,'ES - Engineering Sciences',40,60,0,0,45,0,10,'UNIQUE',NULL,NULL,1),(96,'CS1234','1234','Theory',3,1,3,3,1,'HSS - Humanities and Social Sciences',40,60,15,45,0,15,0,'UNIQUE',NULL,NULL,1),(97,'CS789','hello','Theory',3,1,2,0,1,'ES - Engineering Sciences',40,60,15,30,0,15,0,'UNIQUE',NULL,NULL,1),(98,'CS190','2026 theory check','Theory',3,1,2,0,1,'ES - Engineering Sciences',40,60,15,30,0,15,0,'UNIQUE',NULL,NULL,1),(99,'CS600','theory check','Theory',3,1,2,0,0,'BS - Basic Sciences',40,60,15,30,0,0,0,'UNIQUE',NULL,NULL,1),(100,'CS601','lathery','Lab',3,0,0,2,0,'ES - Engineering Sciences',40,60,0,0,30,0,0,'UNIQUE',NULL,NULL,1),(101,'CS603','theory&labcheck','Theory&Lab',3,2,2,1,0,'ES - Engineering Sciences',40,60,30,30,15,0,0,'UNIQUE',NULL,NULL,1),(102,'CS103','checkbug','Theory',2,3,3,0,0,'ES - Engineering Sciences',40,60,45,45,0,0,0,'UNIQUE',NULL,NULL,1),(103,'CS121','check1','Theory',110,1,12,0,0,'ES - Engineering Sciences',40,60,15,180,0,0,0,'UNIQUE',NULL,NULL,1);
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
  `template_config` json DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `max_credits` int DEFAULT '0',
  `curriculum_ref_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum`
--

LOCK TABLES `curriculum` WRITE;
/*!40000 ALTER TABLE `curriculum` DISABLE KEYS */;
INSERT INTO `curriculum` VALUES (4,'BE - CSE ','2025-2026','2026',NULL,'2025-11-12 06:09:01',163,NULL,1),(10,'R2026-CSE','2024-2025','2026',NULL,'2026-01-06 04:47:56',162,NULL,1),(14,'check 1','2024-2025','2022',NULL,'2026-01-13 09:56:01',162,NULL,1),(15,'check2','2024-2025','2022',NULL,'2026-01-21 22:59:50',160,NULL,1);
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
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_rc_regulation` (`curriculum_id`) USING BTREE,
  KEY `fk_rc_semester` (`semester_id`) USING BTREE,
  KEY `fk_rc_course` (`course_id`) USING BTREE,
  CONSTRAINT `fk_rc_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_rc_regulation` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_rc_semester` FOREIGN KEY (`semester_id`) REFERENCES `normal_cards` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=232 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_courses`
--

LOCK TABLES `curriculum_courses` WRITE;
/*!40000 ALTER TABLE `curriculum_courses` DISABLE KEYS */;
INSERT INTO `curriculum_courses` VALUES (137,10,48,20,1),(138,10,48,21,1),(139,10,48,22,1),(140,10,48,23,1),(141,10,48,24,1),(142,10,48,25,1),(143,10,48,26,1),(144,10,48,27,1),(145,10,48,28,1),(146,10,48,29,1),(147,10,49,30,1),(148,10,49,31,1),(149,10,49,32,1),(150,10,49,33,1),(151,10,49,34,1),(152,10,49,35,1),(153,10,49,36,1),(154,10,49,37,1),(155,10,49,38,1),(156,10,49,39,1),(157,10,50,40,1),(158,10,50,41,1),(159,10,50,42,1),(160,10,50,43,1),(161,10,50,44,1),(162,10,50,45,1),(163,10,50,46,1),(164,10,50,47,1),(165,10,50,48,1),(166,10,50,49,1),(167,10,51,50,1),(168,10,51,51,1),(169,10,51,52,1),(170,10,51,53,1),(171,10,51,54,1),(172,10,51,55,1),(173,10,51,56,1),(174,10,51,57,1),(175,10,51,58,1),(176,10,51,59,1),(177,10,52,60,1),(178,10,52,61,1),(179,10,52,62,1),(180,10,52,63,1),(181,10,52,64,1),(182,10,52,64,1),(183,10,52,65,1),(184,10,52,66,1),(185,10,52,67,1),(186,10,53,68,1),(187,10,53,69,1),(188,10,53,70,1),(189,10,53,71,1),(190,10,53,72,1),(191,10,53,73,1),(192,10,53,64,1),(193,10,53,64,1),(194,10,53,64,1),(195,10,54,74,1),(196,10,54,75,1),(197,10,54,76,1),(198,10,54,77,1),(199,10,54,78,1),(200,10,54,79,1),(201,10,54,80,1),(202,10,55,81,1),(213,14,71,90,1),(214,14,71,91,1),(215,14,71,92,1),(216,4,3,93,1),(217,4,3,94,1),(218,4,3,95,1),(219,4,3,96,1),(220,4,3,1,1),(221,4,3,1,1),(222,4,3,1,1),(223,4,3,1,1),(224,4,3,97,1),(225,4,3,98,1),(226,14,71,99,1),(227,14,71,100,1),(228,14,71,101,1),(229,14,71,102,1),(230,14,71,1,1),(231,14,71,103,1);
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
) ENGINE=InnoDB AUTO_INCREMENT=271 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_logs`
--

LOCK TABLES `curriculum_logs` WRITE;
/*!40000 ALTER TABLE `curriculum_logs` DISABLE KEYS */;
INSERT INTO `curriculum_logs` VALUES (1,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:21:59',NULL),(2,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:24:07',NULL),(3,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:28:35',NULL),(4,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:58:12','{\"mission\": {\"new\": [\"To impart need based education \", \"To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.\", \"To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.\"], \"old\": [\"To impart need based education to meet the requirements of the industry and society.\", \"To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.\", \"To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.\"]}}'),(5,4,'Vision Updated','Updated department vision','System','2025-11-18 06:18:04','{\"vision\": {\"new\": \"To excel \", \"old\": \"To excel in the field of Computer Science and Engineering\"}}'),(6,4,'Vision Updated','Updated department vision','System','2025-11-18 06:18:25','{\"vision\": {\"new\": \"To excel in the field of Computer Science and Engineering\", \"old\": \"To excel \"}}'),(7,4,'Mission[3] Added','Added Mission item at index 3','System','2025-11-18 06:24:25','{\"Mission[3]\": {\"new\": \"to be good\", \"old\": \"\"}}'),(8,4,'Mission[3] Deleted','Deleted Mission item at index 3','System','2025-11-18 06:24:44','{\"Mission[3]\": {\"new\": \"\", \"old\": \"to be good\"}}'),(10,4,'CO-PO/PSO Mapping Saved','Updated CO-PO and CO-PSO mappings for course: ENGINEERING MATHEMATICS I','System','2025-11-19 04:25:50','{\"co_po_mappings\": {\"new\": {\"CO0-PO1\": 1, \"CO0-PO2\": 2, \"CO1-PO1\": 2, \"CO1-PO2\": 2, \"CO2-PO1\": 2, \"CO2-PO2\": 1, \"CO3-PO1\": 2, \"CO3-PO2\": 2, \"CO4-PO1\": 1, \"CO4-PO2\": 2}, \"old\": {\"CO0-PO1\": 1, \"CO0-PO2\": 2, \"CO1-PO1\": 2, \"CO1-PO2\": 2, \"CO2-PO1\": 2, \"CO2-PO2\": 1, \"CO3-PO1\": 2, \"CO3-PO2\": 2, \"CO4-PO1\": 1, \"CO4-PO2\": 2}}, \"co_pso_mappings\": {\"new\": {\"CO0-PSO1\": 1, \"CO1-PSO2\": 2, \"CO2-PSO2\": 1, \"CO3-PSO1\": 1, \"CO3-PSO3\": 3, \"CO4-PSO2\": 2, \"CO4-PSO3\": 2}, \"old\": {}}}'),(11,4,'Course Updated','Updated course: 22CH103  - ENGINEERING CHEMISTRY I ','System','2025-11-19 10:26:39','{\"credit\": {\"new\": 141, \"old\": 3}, \"category\": {\"new\": \"BS\", \"old\": \"BS - Basic Sciences\"}}'),(12,4,'Course Updated','Updated course: 22CS108  - COMPREHENSIVE WORK','System','2025-11-19 10:27:00','{\"credit\": {\"new\": 4, \"old\": 1}, \"category\": {\"new\": \"PE\", \"old\": \"EEC - Employability Enhancement Course\"}}'),(13,4,'Course Updated','Updated course: 22GE001  - FUNDAMENTALS OF COMPUTING ','System','2025-11-19 10:27:17','{\"credit\": {\"new\": 4, \"old\": 3}, \"category\": {\"new\": \"BS\", \"old\": \"ES - Engineering Sciences\"}}'),(14,4,'Course Updated','Updated course: 22CH103  - ENGINEERING CHEMISTRY I ','System','2025-12-22 04:10:33','{\"credit\": {\"new\": 3, \"old\": 141}}'),(15,4,'Course Updated','Updated course: 22CH103  - ENGINEERING CHEMISTRY I ','System','2025-12-22 04:10:33','{\"credit\": {\"new\": 3, \"old\": 141}}'),(16,4,'Course Updated','Updated course: 22MA101 - ENGINEERING MATHEMATICS I','System','2025-12-22 06:01:23','{\"category\": {\"new\": \"BS\", \"old\": \"BS - Basic Sciences\"}, \"tutorial_hours\": {\"new\": 30, \"old\": 1}}'),(17,4,'Course Updated','Updated course: 22MA101 - ENGINEERING MATHEMATICS I','System','2025-12-22 08:28:03','{\"lecture_hours\": {\"new\": 2, \"old\": 3}, \"practical_hours\": {\"new\": 1, \"old\": 0}}'),(45,4,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 17:28:37','{\"Mission[0]\": {\"new\": \"To impart need based education EXTRA\", \"old\": \"To impart need based education \"}}'),(48,4,'Course Added','Added course 26MA101 - Linear Algebra and Calculus to Semester 3','System','2025-12-26 05:43:55',NULL),(49,4,'CO-PO/PSO Mapping Saved','Updated CO-PO and CO-PSO mappings for course: Linear Algebra and Calculus','System','2025-12-26 05:53:59','{\"co_po_mappings\": {\"new\": {\"CO0-PO1\": 3, \"CO0-PO2\": 1, \"CO0-PO4\": 1, \"CO0-PO5\": 3, \"CO1-PO2\": 2, \"CO1-PO3\": 1, \"CO1-PO4\": 1, \"CO2-PO3\": 3, \"CO3-PO1\": 1}, \"old\": {}}, \"co_pso_mappings\": {\"new\": {}, \"old\": {}}}'),(50,4,'Course Added','Added course 22PH102  - ENGINEERING PHYSICS  to Semester 3','System','2025-12-26 05:55:54',NULL),(51,4,'Course Added','Added course 22CH103  - ENGINEERING CHEMISTRY I  to Semester 3','System','2025-12-26 05:56:36',NULL),(52,4,'Course Added','Added course 22MA101 - ENGINEERING MATHEMATICS I  to Semester 3','System','2025-12-26 06:17:48',NULL),(53,4,'Course Added','Added course 22PH102  - ENGINEERING PHYSICS  to Semester 3','System','2025-12-26 06:18:23',NULL),(55,4,'Semester Added','Added Semester 0','System','2026-01-05 04:25:39',NULL),(56,4,'Semester Added','Added electives','System','2026-01-05 04:28:31',NULL),(57,4,'Semester Added','Added electives','System','2026-01-05 04:32:16',NULL),(58,4,'Semester Added','Added semester','System','2026-01-05 04:35:58',NULL),(59,4,'Semester Added','Added electives','System','2026-01-05 04:36:11',NULL),(60,4,'Semester Added','Added vertical ','System','2026-01-05 04:36:32',NULL),(61,4,'Honour Card Added','Added Honour Card: honour verticals','System','2026-01-05 04:38:03',NULL),(62,4,'Semester Added','Added Elective','System','2026-01-05 04:52:35',NULL),(63,4,'Semester Added','Added Vertical','System','2026-01-05 04:53:39',NULL),(64,4,'Semester Added','Added Vertical','System','2026-01-05 04:53:53',NULL),(65,4,'Semester Added','Added Semester','System','2026-01-05 04:54:01',NULL),(66,4,'Semester Added','Added Vertical','System','2026-01-05 08:18:03',NULL),(67,4,'Course Removed','Removed course ENGINEERING MATHEMATICS I  from Semester 3','System','2026-01-05 09:35:37',NULL),(68,4,'PEO[0] Updated','Updated PEO item at index 0','System','2026-01-06 04:31:33','{\"PEO[0]\": {\"new\": \"Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.\", \"old\": \"Graduates will apply computer science and engineering principles and practices to solvereal- world problems with their technical competence.\"}}'),(69,4,'PEO[1] Updated','Updated PEO item at index 1','System','2026-01-06 04:31:33','{\"PEO[1]\": {\"new\": \"Pursue continuous learning in emerging technologies such as AI, data science, and cybersecurity to remain adaptable professionals.\", \"old\": \"Graduates will have the domain knowledge to pursue higher education and apply cuttingedge research to develop solutions for socially relevant problems.\"}}'),(70,4,'PEO[2] Updated','Updated PEO item at index 2','System','2026-01-06 04:31:33','{\"PEO[2]\": {\"new\": \"Demonstrate leadership, teamwork, and ethical responsibility in developing sustainable software and computing solutions.\", \"old\": \"Graduates will communicate effectively and practice their profession with ethics,integrity, leadership, teamwork, and social responsibility, and pursue lifelong learning throughout their careers.\"}}'),(71,4,'PSO[0] Updated','Updated PSO item at index 0','System','2026-01-06 04:32:14','{\"PSO[0]\": {\"new\": \"Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.\", \"old\": \"Apply suitable algorithmic thinking and data management practices to design develop, and evaluate effective solutions for real-life and research problems.\"}}'),(72,4,'PSO[1] Updated','Updated PSO item at index 1','System','2026-01-06 04:32:14','{\"PSO[1]\": {\"new\": \"Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.\", \"old\": \"Design and develop cost-effective solutions based on cutting-edge hardware and software tools and techniques to meet global requirements.\"}}'),(73,10,'Curriculum Created','Created new curriculum: R2026-CSE (2025-2026)','System','2026-01-06 04:47:56',NULL),(74,10,'Department Overview Created','Created department vision, mission, PEOs, POs, and PSOs','System','2026-01-06 04:49:17',NULL),(75,10,'PSO[1] Added','Added PSO item at index 1','System','2026-01-06 04:49:57','{\"PSO[1]\": {\"new\": \"Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.\", \"old\": \"\"}}'),(76,10,'PSO[0] Added','Added PSO item at index 0','System','2026-01-06 04:49:57','{\"PSO[0]\": {\"new\": \"Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.\", \"old\": \"\"}}'),(77,10,'PO[2] Added','Added PO item at index 2','System','2026-01-06 04:52:59','{\"PO[2]\": {\"new\": \"Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.\", \"old\": \"\"}}'),(78,10,'PO[8] Added','Added PO item at index 8','System','2026-01-06 04:52:59','{\"PO[8]\": {\"new\": \"Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings\", \"old\": \"\"}}'),(79,10,'PO[5] Added','Added PO item at index 5','System','2026-01-06 04:52:59','{\"PO[5]\": {\"new\": \"The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice\", \"old\": \"\"}}'),(80,10,'PO[1] Added','Added PO item at index 1','System','2026-01-06 04:52:59','{\"PO[1]\": {\"new\": \"Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering sciences.\", \"old\": \"\"}}'),(81,10,'PO[10] Added','Added PO item at index 10','System','2026-01-06 04:52:59','{\"PO[10]\": {\"new\": \"Project Management and Finance: Demonstrate knowledge and understanding of the engineering and management principles and apply these to one’s own work, as a member and leader in a team, to manage projects and in multidisciplinary environments.\", \"old\": \"\"}}'),(82,10,'PO[9] Added','Added PO item at index 9','System','2026-01-06 04:52:59','{\"PO[9]\": {\"new\": \"Communication: Communicate effectively on complex engineering activities with the engineering community and with society at large, such as, being able to comprehend and write effective reports and design documentation, make effective presentations, and give and receive clear instructions.\", \"old\": \"\"}}'),(83,10,'PO[3] Added','Added PO item at index 3','System','2026-01-06 04:52:59','{\"PO[3]\": {\"new\": \"Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.\", \"old\": \"\"}}'),(84,10,'PO[4] Added','Added PO item at index 4','System','2026-01-06 04:52:59','{\"PO[4]\": {\"new\": \"Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.\", \"old\": \"\"}}'),(85,10,'PO[0] Added','Added PO item at index 0','System','2026-01-06 04:52:59','{\"PO[0]\": {\"new\": \"Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.\", \"old\": \"\"}}'),(86,10,'PO[11] Added','Added PO item at index 11','System','2026-01-06 04:52:59','{\"PO[11]\": {\"new\": \"Life-long Learning: Recognize the need for, and have the preparation and ability to engage in independent and life-long learning in the broadest context of technological change.\", \"old\": \"\"}}'),(87,10,'PO[7] Added','Added PO item at index 7','System','2026-01-06 04:52:59','{\"PO[7]\": {\"new\": \"Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.\", \"old\": \"\"}}'),(88,10,'PO[6] Added','Added PO item at index 6','System','2026-01-06 04:52:59','{\"PO[6]\": {\"new\": \"Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.\", \"old\": \"\"}}'),(89,10,'Card Added','Added Semester 1','System','2026-01-06 04:57:25',NULL),(90,4,'PEO[0] Updated','Updated PEO item at index 0','System','2026-01-06 05:10:13','{\"PEO[0]\": {\"new\": \"hello\", \"old\": \"Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.\"}}'),(91,4,'PEO[0] Updated','Updated PEO item at index 0','System','2026-01-06 05:10:55','{\"PEO[0]\": {\"new\": \"Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.\", \"old\": \"hello\"}}'),(92,10,'Course Added','Added course 26MA101 - Linear Algebra and Calculus to Semester 48','System','2026-01-06 05:28:05',NULL),(93,10,'Course Added','Added course 26PH102 - Engineering Physics to Semester 48','System','2026-01-06 05:35:35',NULL),(94,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 05:35:43','{\"category\": {\"new\": \"ES\", \"old\": \"BS - Basic Sciences\"}, \"cia_marks\": {\"new\": 40, \"old\": 60}, \"see_marks\": {\"new\": 60, \"old\": 40}}'),(95,10,'Course Added','Added course 26CH103 - Engineering Chemistry to Semester 48','System','2026-01-06 05:37:34',NULL),(96,10,'Course Added','Added course 26GE004 - Digital Computer Electronics to Semester 48','System','2026-01-06 05:39:17',NULL),(97,10,'Course Added','Added course 26GE005 - Problem Solving using C to Semester 48','System','2026-01-06 06:03:03',NULL),(98,10,'Course Added','Added course 26HS001 - Communicative English to Semester 48','System','2026-01-06 06:08:14',NULL),(99,10,'Course Added','Added course 26HS002 - தமிழர் மரபு / Heritage of Tamils  to Semester 48','System','2026-01-06 06:09:10',NULL),(100,10,'Course Added','Added course 26PH108 - Physical Science Laboratory to Semester 48','System','2026-01-06 06:11:53',NULL),(101,10,'Course Added','Added course 26GE006 - C Programming Laboratory to Semester 48','System','2026-01-06 06:12:57',NULL),(102,10,'Course Added','Added course 26SD001 - Skill Development Course  I to Semester 48','System','2026-01-06 06:17:04',NULL),(103,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:20:19','{\"category\": {\"new\": \"BS\", \"old\": \"ES\"}}'),(104,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:20:27','{\"category\": {\"new\": \"PC\", \"old\": \"BS\"}}'),(105,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:20:39','{\"category\": {\"new\": \"BS\", \"old\": \"PC\"}}'),(106,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:21:07','{\"category\": {\"new\": \"ES\", \"old\": \"BS\"}}'),(107,10,'Card Added','Added Semester 2','System','2026-01-06 06:29:27',NULL),(108,10,'Course Added','Added course 26MA201 - Differential Equations and Transforms to Semester 49','System','2026-01-06 06:31:04',NULL),(109,10,'Course Added','Added course 26PH202 - Materials Science to Semester 49','System','2026-01-06 06:33:17',NULL),(110,10,'Course Added','Added course 26CS203 - Fundamentals of Web Principles to Semester 49','System','2026-01-06 06:34:05',NULL),(111,10,'Course Added','Added course 26CS204 - Computer Organization and Architecture to Semester 49','System','2026-01-06 06:34:39',NULL),(112,10,'Course Added','Added course 26GE007 - Python Programming to Semester 49','System','2026-01-06 06:35:17',NULL),(113,10,'Course Added','Added course 26HS005 - Professional Communication to Semester 49','System','2026-01-06 06:36:17',NULL),(114,10,'Course Added','Added course 26HS006 - தமிழரும் தொழில்நுட்பமும் / Tamils and Technology to Semester 49','System','2026-01-06 06:37:28',NULL),(115,10,'Course Added','Added course 26GE008 - Python Programming Laboratory to Semester 49','System','2026-01-06 06:38:21',NULL),(116,10,'Course Added','Added course 26CS209 - Web Principles Laboratory to Semester 49','System','2026-01-06 06:39:13',NULL),(117,10,'Course Added','Added course 26SD002 - Skill Development Course II to Semester 49','System','2026-01-06 06:40:28',NULL),(118,10,'Card Added','Added Semester 3','System','2026-01-06 06:40:47',NULL),(119,10,'Course Added','Added course 26CS301 - Discrete Mathematics to Semester 50','System','2026-01-06 06:42:11',NULL),(120,10,'Course Added','Added course 26CS302 - Data Structures and Algorithms to Semester 50','System','2026-01-06 06:43:30',NULL),(121,10,'Course Added','Added course 26CS303 - Operating Systems to Semester 50','System','2026-01-06 06:44:19',NULL),(122,10,'Course Added','Added course 26CS304 - Object Oriented Programming with Java to Semester 50','System','2026-01-06 06:45:12',NULL),(123,10,'Course Updated','Updated course: 26CS304 - Object Oriented Programming with Java','System','2026-01-06 06:46:09','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"cia_marks\": {\"new\": 40, \"old\": 30}, \"see_marks\": {\"new\": 60, \"old\": 70}, \"practical_hours\": {\"new\": 2, \"old\": 0}}'),(124,10,'Course Updated','Updated course: 26CS304 - Object Oriented Programming with Java','System','2026-01-06 06:46:38','{\"cia_marks\": {\"new\": 30, \"old\": 40}, \"see_marks\": {\"new\": 70, \"old\": 60}}'),(125,10,'Course Added','Added course 26CS305 - Software Engineering to Semester 50','System','2026-01-06 06:47:33',NULL),(126,10,'Course Added','Added course 26CS306 - Database Management Systems to Semester 50','System','2026-01-06 06:48:12',NULL),(127,10,'Course Added','Added course 26CS307 - Standards in Computer Science to Semester 50','System','2026-01-06 06:48:52',NULL),(128,10,'Course Added','Added course 26CS308 - Data Structures and Algorithms Laboratory to Semester 50','System','2026-01-06 06:49:39',NULL),(129,10,'Course Added','Added course 26CS309 - Database Management Systems Laboratory to Semester 50','System','2026-01-06 06:50:43',NULL),(130,10,'Course Added','Added course 26CS310 - Design Thinking and Innovation Laboratory (AICTE, & NEP) to Semester 50','System','2026-01-06 06:52:07',NULL),(131,10,'Card Added','Added Semester 4','System','2026-01-06 09:03:44',NULL),(132,10,'Course Added','Added course 26CS401 - Probability and Statistics to Semester 51','System','2026-01-06 09:04:53',NULL),(133,10,'Course Added','Added course 26CS402 - Full Stack Development to Semester 51','System','2026-01-06 09:05:42',NULL),(134,10,'Course Added','Added course 26CS403 - Artificial Intelligence Essentials to Semester 51','System','2026-01-06 09:06:27',NULL),(135,10,'Course Added','Added course 26CS404 - Design and Analysis of Algorithms to Semester 51','System','2026-01-06 09:07:04',NULL),(136,10,'Course Updated','Updated course: 26CS403 - Artificial Intelligence Essentials','System','2026-01-06 09:07:27','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"see_marks\": {\"new\": 55, \"old\": 54}}'),(137,10,'Course Added','Added course 26CS405 - Theory of Computation to Semester 51','System','2026-01-06 09:08:33',NULL),(138,10,'Course Added','Added course 26CS406 - Computer Networks to Semester 51','System','2026-01-06 09:09:37',NULL),(139,10,'Course Added','Added course 26HS009 - Environmental Sciences and Sustainability  to Semester 51','System','2026-01-06 09:10:29',NULL),(140,10,'Course Added','Added course 26CS408 - Full Stack Development Laboratory to Semester 51','System','2026-01-06 09:11:13',NULL),(141,10,'Course Added','Added course 26CS409 - Computer Networks Laboratory to Semester 51','System','2026-01-06 09:12:17',NULL),(142,10,'Course Added','Added course 26CS410 - Community Engagement Project to Semester 51','System','2026-01-06 09:13:04',NULL),(143,10,'Course Updated','Updated course: 26CS409 - Computer Networks Laboratory','System','2026-01-06 09:13:57','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"course_type\": {\"new\": \"Lab\", \"old\": \"Experiment\"}, \"lecture_hours\": {\"new\": 0, \"old\": 2}, \"practical_hours\": {\"new\": 2, \"old\": 0}}'),(144,10,'Card Added','Added Semester 5','System','2026-01-06 09:14:38',NULL),(145,10,'Course Added','Added course 26CS501 - Compiler Design to Semester 52','System','2026-01-06 09:15:30',NULL),(146,10,'Course Added','Added course 26CS502 - Cloud Infrastructure Services to Semester 52','System','2026-01-06 09:16:38',NULL),(147,10,'Course Added','Added course 26CS503 - Bigdata Analytics to Semester 52','System','2026-01-06 09:17:35',NULL),(148,10,'Course Added','Added course 26CS504 - Machine Learning Essentials to Semester 52','System','2026-01-06 09:18:18',NULL),(149,10,'Course Added','Added course 26XX - Professional Elective I to Semester 52','System','2026-01-06 09:20:24',NULL),(150,10,'Course Added','Added course 26XX - Open Elective I to Semester 52','System','2026-01-06 09:21:09',NULL),(151,10,'Course Added','Added course 26CS507 - Cloud Infrastructure Services Laboratory to Semester 52','System','2026-01-06 09:21:55',NULL),(152,10,'Course Added','Added course 26CS508 - Machine Learning Essentials Laboratory to Semester 52','System','2026-01-06 09:22:33',NULL),(153,10,'Course Added','Added course 26CS509 - Technology Integration Project to Semester 52','System','2026-01-06 09:23:11',NULL),(154,10,'Card Added','Added Semester 6','System','2026-01-06 09:23:53',NULL),(155,10,'Course Added','Added course 26CS601 - Software Project Management and Quality Assurance to Semester 53','System','2026-01-06 09:24:39',NULL),(156,10,'Course Added','Added course 26CS602 - Deep Learning to Semester 53','System','2026-01-06 09:25:17',NULL),(157,10,'Course Added','Added course 26CS603 - Cryptography and Cyber Security to Semester 53','System','2026-01-06 09:25:58',NULL),(158,10,'Course Added','Added course 26CS607 - Software Project Management and Quality Assurance Laboratory to Semester 53','System','2026-01-06 09:26:58',NULL),(159,10,'Course Added','Added course 26CS608 - Deep Learning Laboratory  to Semester 53','System','2026-01-06 09:27:32',NULL),(160,10,'Course Added','Added course 26CS609 - Innovation and Product Development Project / Industry Oriented Course / Summer Internship to Semester 53','System','2026-01-06 09:28:14',NULL),(161,10,'Course Added','Added course 26XX - Professional Elective II to Semester 53','System','2026-01-06 09:29:07',NULL),(162,10,'Course Added','Added course 26XX - Professional Elective III to Semester 53','System','2026-01-06 09:29:43',NULL),(163,10,'Course Added','Added course 26XX - Open Elective II to Semester 53','System','2026-01-06 09:30:24',NULL),(164,10,'Course Updated','Updated course: 26XX - Open Elective II','System','2026-01-06 09:31:05','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"course_name\": {\"new\": \"Open Elective II\", \"old\": \"Professional Elective I\"}}'),(165,10,'Course Updated','Updated course: 26XX - Professional Elective II','System','2026-01-06 09:32:01','{\"course_name\": {\"new\": \"Professional Elective II\", \"old\": \"Open Elective II\"}}'),(166,10,'Card Added','Added Semester 7','System','2026-01-06 09:32:40',NULL),(167,10,'Course Added','Added course 26CS701 - Generative AI and Large Language Models to Semester 54','System','2026-01-06 09:33:11',NULL),(168,10,'Course Added','Added course 26CS702 - IoT and Edge Computing to Semester 54','System','2026-01-06 09:33:49',NULL),(169,10,'Course Added','Added course 26XXIV - Professional Elective IV to Semester 54','System','2026-01-06 09:35:02',NULL),(170,10,'Course Added','Added course 26XXV - Professional Elective V to Semester 54','System','2026-01-06 09:35:40',NULL),(171,10,'Course Updated','Updated course: 26XXIV - Professional Elective IV','System','2026-01-06 09:36:12','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}}'),(172,10,'Course Updated','Updated course: 26XXV - Professional Elective V','System','2026-01-06 09:36:45','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}}'),(173,10,'Course Added','Added course 26XXVI - Professional Elective VI to Semester 54','System','2026-01-06 09:37:35',NULL),(174,10,'Course Added','Added course 26CS706 - Generative AI and Large Language Models Laboratory to Semester 54','System','2026-01-06 09:38:35',NULL),(175,10,'Course Added','Added course 26CS707 - Capstone Project work Level I / Internship Pro to Semester 54','System','2026-01-06 09:39:41',NULL),(176,10,'Course Updated','Updated course: 26XXIII - Professional Elective III','System','2026-01-06 09:40:13','{\"course_code\": {\"new\": \"26XXIII\", \"old\": \"26XX\"}, \"course_name\": {\"new\": \"Professional Elective III\", \"old\": \"Professional Elective II\"}}'),(177,10,'Card Added','Added Semester 8','System','2026-01-06 09:41:03',NULL),(178,10,'Course Added','Added course 26CS801 - Capstone Project Work Level II / Internship Project / Startup Product to Semester 55','System','2026-01-06 09:42:07',NULL),(179,10,'Curriculum Updated','Updated curriculum details','System','2026-01-06 09:47:58','{\"academic_year\": {\"new\": \"\", \"old\": \"2025-2026\"}}'),(180,10,'Course Updated','Updated course: 26XXIV - Professional Elective IV','System','2026-01-06 09:50:30','{\"course_code\": {\"new\": \"26XXIV\", \"old\": \"26XXIII\"}, \"course_name\": {\"new\": \"Professional Elective IV\", \"old\": \"Professional Elective III\"}}'),(204,4,'PO[12] Added','Added PO item at index 12','System','2026-01-13 08:49:03','{\"PO[12]\": {\"new\": \"check 1\", \"old\": \"\"}}'),(205,4,'PEO[3] Added','Added PEO item at index 3','System','2026-01-13 08:49:03','{\"PEO[3]\": {\"new\": \"check 1\", \"old\": \"\"}}'),(213,10,'Curriculum Updated','Updated curriculum details','System','2026-01-13 08:58:03','{\"academic_year\": {\"new\": \"2024-2025\", \"old\": \"\"}}'),(217,10,'Semester Updated','Updated Semester 8 to Semester 7','System','2026-01-13 08:59:19','{\"semester_number\": {\"new\": 7, \"old\": 8}}'),(218,10,'Semester Updated','Updated Semester 7 to Semester 8','System','2026-01-13 08:59:26','{\"semester_number\": {\"new\": 8, \"old\": 7}}'),(219,10,'Semester Updated','Updated Semester 8 to Semester 9','System','2026-01-13 09:04:46','{\"semester_number\": {\"new\": 9, \"old\": 8}}'),(220,10,'Semester Updated','Updated Semester 9 to Semester 8','System','2026-01-13 09:05:00','{\"semester_number\": {\"new\": 8, \"old\": 9}}'),(221,10,'Semester Updated','Updated Semester 8 to Semester 9','System','2026-01-13 09:12:33','{\"semester_number\": {\"new\": 9, \"old\": 8}}'),(222,10,'Semester Updated','Updated Semester 9 to Semester 8','System','2026-01-13 09:12:37','{\"semester_number\": {\"new\": 8, \"old\": 9}}'),(223,10,'Card Added','Added Vertical 1','System','2026-01-13 09:12:44',NULL),(224,10,'Card Added','Added Vertical 2','System','2026-01-13 09:12:57',NULL),(225,10,'Card Added','Added New Card','System','2026-01-13 09:13:32',NULL),(226,10,'Card Added','Added New Card','System','2026-01-13 09:13:46',NULL),(227,10,'Honour Card Added','Added Honour Card: honour vertical *','System','2026-01-13 09:14:16',NULL),(228,10,'Honour Card Added','Added Honour Card: sdknvjnvfdjnvfnj','System','2026-01-13 09:22:50',NULL),(229,10,'Honour Card Added','Added Honour Card: Honour Vertical *','System','2026-01-13 09:28:30',NULL),(230,10,'Card Added','Added Vertical 2','System','2026-01-13 09:42:10',NULL),(231,10,'Honour Card Added','Added Honour Card: Honour vertical *','System','2026-01-13 09:51:43',NULL),(232,4,'Honour Card Added','Added Honour Card: Honour Vertical *','System','2026-01-13 09:52:21',NULL),(233,14,'Curriculum Created','Created new curriculum: check 1 (2024-2025)','System','2026-01-13 09:56:01',NULL),(234,14,'Card Added','Added Semester 1','System','2026-01-13 09:56:11',NULL),(235,14,'Honour Card Added','Added Honour Card: Honour card','System','2026-01-13 09:56:19',NULL),(236,4,'Course Added','Added course CS101 - cue to Semester 3','System','2026-01-13 11:04:23',NULL),(237,4,'Course Added','Added course CS130 - check1 to Semester 3','System','2026-01-13 11:08:19',NULL),(238,4,'Course Added','Added course CS230 - check 2 to Semester 3','System','2026-01-13 12:11:59',NULL),(239,4,'Course Added','Added course CS303 - check 3 to Semester 3','System','2026-01-13 12:24:32',NULL),(240,4,'Course Added','Added course cs102323 - check4 to Semester 3','System','2026-01-19 10:54:51',NULL),(241,4,'Course Added','Added course CS134 - check5 to Semester 3','System','2026-01-19 11:33:08',NULL),(242,14,'Course Added','Added course CS140 - 1check2022 to Semester 71','System','2026-01-19 14:01:55',NULL),(243,14,'Course Added','Added course CS120 - 2Check2022 to Semester 71','System','2026-01-19 14:22:02',NULL),(244,14,'Course Added','Added course CS150 - 3maincheck to Semester 71','System','2026-01-19 14:23:23',NULL),(245,4,'Course Removed','Removed course ENGINEERING PHYSICS  from Semester 3','System','2026-01-19 14:41:58',NULL),(246,4,'Course Removed','Removed course Introduction to Programming from Semester 3','System','2026-01-19 14:42:00',NULL),(247,4,'Course Removed','Removed course check4 from Semester 3','System','2026-01-19 14:42:03',NULL),(248,4,'Course Removed','Removed course check1 from Semester 3','System','2026-01-19 14:42:05',NULL),(249,4,'Course Removed','Removed course check5 from Semester 3','System','2026-01-19 14:42:07',NULL),(250,4,'Course Removed','Removed course check 2 from Semester 3','System','2026-01-19 14:42:11',NULL),(251,4,'Course Removed','Removed course check 3 from Semester 3','System','2026-01-19 14:42:13',NULL),(252,4,'Course Added','Added course CS1 - check2026theory to Semester 3','System','2026-01-19 14:49:31',NULL),(253,4,'Course Added','Added course CS2 - check2026lab to Semester 3','System','2026-01-19 14:52:03',NULL),(254,4,'Course Added','Added course CS3 - check 123 to Semester 3','System','2026-01-19 15:13:06',NULL),(255,4,'Course Added','Added course CS1234 - 1234 to Semester 3','System','2026-01-19 15:13:53',NULL),(256,4,'Course Added','Added course CS101 - theory check to Semester 3','System','2026-01-20 15:53:23',NULL),(257,4,'Course Added','Added course CS101 - theory check to Semester 3','System','2026-01-20 15:55:10',NULL),(258,4,'Course Added','Added course CS101 - hello to Semester 3','System','2026-01-20 16:01:45',NULL),(259,4,'Course Added','Added course CS101 - theroycheck to Semester 3','System','2026-01-20 16:07:41',NULL),(260,4,'Course Added','Added course CS789 - hello to Semester 3','System','2026-01-20 16:09:16',NULL),(261,4,'Course Added','Added course CS190 - 2026 theory check to Semester 3','System','2026-01-20 16:13:40',NULL),(262,14,'Course Added','Added course CS600 - theory check to Semester 71','System','2026-01-21 08:53:37',NULL),(263,14,'Course Added','Added course CS601 - lathery to Semester 71','System','2026-01-21 09:04:50',NULL),(264,14,'Course Added','Added course CS603 - theory&labcheck to Semester 71','System','2026-01-21 09:06:13',NULL),(265,14,'Course Added','Added course CS103 - checkbug to Semester 71','System','2026-01-21 20:57:38',NULL),(266,14,'Course Added','Added course CS101 - hello to Semester 71','System','2026-01-21 22:24:52',NULL),(267,14,'Course Added','Added course CS121 - check1 to Semester 71','System','2026-01-21 22:26:25',NULL),(268,15,'Curriculum Created','Created new curriculum: check2 (2024-2025)','System','2026-01-21 22:59:50',NULL),(269,15,'Mission[0] Added','Added Mission item at index 0','System','2026-01-22 08:52:44','{\"Mission[0]\": {\"new\": \"remove check 1\", \"old\": \"\"}}'),(270,15,'Mission[0] Deleted','Deleted Mission item at index 0','System','2026-01-22 08:54:32','{\"Mission[0]\": {\"new\": \"\", \"old\": \"remove check 1\"}}');
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
  `source_department_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_mission_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_mission`
--

LOCK TABLES `curriculum_mission` WRITE;
/*!40000 ALTER TABLE `curriculum_mission` DISABLE KEYS */;
INSERT INTO `curriculum_mission` VALUES (2,2,'To impart need based education EXTRA',0,'CLUSTER',NULL,NULL,1),(3,2,'To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.',1,'CLUSTER',NULL,NULL,1),(4,2,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',2,'CLUSTER',NULL,NULL,1),(8,4,'hello',0,'UNIQUE',NULL,NULL,1),(9,5,'hello hi',0,'UNIQUE',NULL,NULL,1),(10,5,'hello',1,'CLUSTER',NULL,NULL,1),(11,5,'hi',2,'UNIQUE',NULL,NULL,1),(12,4,'hello',1,'CLUSTER',5,NULL,1),(27,3,'MODIFIED',0,'CLUSTER',NULL,NULL,1),(28,3,'To impart need based education EXTRA',1,'CLUSTER',2,NULL,1),(29,6,'ELIMINATE',0,'UNIQUE',NULL,NULL,1),(30,3,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',2,'CLUSTER',2,NULL,1),(31,6,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',1,'CLUSTER',2,NULL,1),(34,8,'To impart need based education to meet the requirements of the industry and society.',0,'UNIQUE',NULL,NULL,1),(35,8,'To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.',1,'UNIQUE',NULL,NULL,1),(36,8,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',2,'UNIQUE',NULL,NULL,1),(37,7,'MODIFIED',0,'CLUSTER',3,NULL,1),(38,5,'check 1',3,'UNIQUE',NULL,NULL,1),(39,11,'remove check 1',0,'UNIQUE',NULL,NULL,0);
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
  `source_department_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_peos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_peos`
--

LOCK TABLES `curriculum_peos` WRITE;
/*!40000 ALTER TABLE `curriculum_peos` DISABLE KEYS */;
INSERT INTO `curriculum_peos` VALUES (1,2,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'CLUSTER',NULL,NULL,1),(2,2,'Pursue continuous learning in emerging technologies such as AI, data science, and cybersecurity to remain adaptable professionals.',1,'UNIQUE',NULL,NULL,1),(3,2,'Demonstrate leadership, teamwork, and ethical responsibility in developing sustainable software and computing solutions.',2,'UNIQUE',NULL,NULL,1),(6,4,'hound cwvwvvsf',0,'UNIQUE',NULL,NULL,1),(25,3,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'CLUSTER',2,NULL,1),(26,6,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'CLUSTER',2,NULL,1),(27,7,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'UNIQUE',NULL,NULL,1),(28,7,'Pursue continuous learning in emerging technologies such as AI, data science, and cybersecurity to remain adaptable professionals.',1,'UNIQUE',NULL,NULL,1),(29,7,'Demonstrate leadership, teamwork, and ethical responsibility in developing sustainable software and computing solutions.',2,'UNIQUE',NULL,NULL,1),(30,8,'Graduates will apply computer science and engineering principles and practices to solvereal- world problems with their technical competence.',0,'UNIQUE',NULL,NULL,1),(31,8,'Graduates will have the domain knowledge to pursue higher education and apply cuttingedge research to develop solutions for socially relevant problems.',1,'UNIQUE',NULL,NULL,1),(32,8,'Graduates will communicate effectively and practice their profession with ethics,integrity, leadership, teamwork, and social responsibility, and pursue lifelong learning throughout their careers.',2,'UNIQUE',NULL,NULL,1),(33,5,'check 1',0,'UNIQUE',NULL,NULL,1);
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
  `source_department_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_pos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_pos`
--

LOCK TABLES `curriculum_pos` WRITE;
/*!40000 ALTER TABLE `curriculum_pos` DISABLE KEYS */;
INSERT INTO `curriculum_pos` VALUES (1,2,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'CLUSTER',NULL,NULL,1),(2,2,'Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering sciences.',1,'UNIQUE',NULL,NULL,1),(3,2,'Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.',2,'UNIQUE',NULL,NULL,1),(4,2,'Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.',3,'UNIQUE',NULL,NULL,1),(5,2,'Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.',4,'UNIQUE',NULL,NULL,1),(6,2,'The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice',5,'UNIQUE',NULL,NULL,1),(7,2,'Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.',6,'UNIQUE',NULL,NULL,1),(8,2,'Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.',7,'UNIQUE',NULL,NULL,1),(9,2,'Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings.',8,'UNIQUE',NULL,NULL,1),(10,2,'Communication: Communicate effectively on complex engineering activities with the engineering community and with society at large, such as, being able to comprehend and write effective reports and design documentation, make effective presentations, and give and receive clear instructions.',9,'UNIQUE',NULL,NULL,1),(11,2,'Project Management and Finance: Demonstrate knowledge and understanding of the engineering and management principles and apply these to one’s own work, as a member and leader in a team, to manage projects and in multidisciplinary environments.',10,'UNIQUE',NULL,NULL,1),(12,2,'Life-long Learning: Recognize the need for, and have the preparation and ability to engage in independent and life-long learning in the broadest context of technological change.',11,'UNIQUE',NULL,NULL,1),(18,4,'sdcvvsvd',0,'UNIQUE',NULL,NULL,1),(19,4,'huii',1,'UNIQUE',NULL,NULL,1),(21,6,'MODIFIED',0,'UNIQUE',NULL,NULL,1),(22,3,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'CLUSTER',2,NULL,1),(23,7,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'UNIQUE',NULL,NULL,1),(24,7,'Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering sciences.',1,'UNIQUE',NULL,NULL,1),(25,7,'Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.',2,'UNIQUE',NULL,NULL,1),(26,7,'Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.',3,'UNIQUE',NULL,NULL,1),(27,7,'Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.',4,'UNIQUE',NULL,NULL,1),(28,7,'The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice',5,'UNIQUE',NULL,NULL,1),(29,7,'Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.',6,'UNIQUE',NULL,NULL,1),(30,7,'Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.',7,'UNIQUE',NULL,NULL,1),(31,7,'Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings',8,'UNIQUE',NULL,NULL,1),(32,7,'Communication: Communicate effectively on complex engineering activities with the engineering community and with society at large, such as, being able to comprehend and write effective reports and design documentation, make effective presentations, and give and receive clear instructions.',9,'UNIQUE',NULL,NULL,1),(33,7,'Project Management and Finance: Demonstrate knowledge and understanding of the engineering and management principles and apply these to one’s own work, as a member and leader in a team, to manage projects and in multidisciplinary environments.',10,'UNIQUE',NULL,NULL,1),(34,7,'Life-long Learning: Recognize the need for, and have the preparation and ability to engage in independent and life-long learning in the broadest context of technological change.',11,'UNIQUE',NULL,NULL,1),(35,8,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'UNIQUE',NULL,NULL,1),(36,8,'Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering science',1,'UNIQUE',NULL,NULL,1),(37,8,'Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.',2,'UNIQUE',NULL,NULL,1),(38,8,'Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.',3,'UNIQUE',NULL,NULL,1),(39,8,'Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.',4,'UNIQUE',NULL,NULL,1),(40,8,'The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice',5,'UNIQUE',NULL,NULL,1),(41,8,'Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.',6,'UNIQUE',NULL,NULL,1),(42,8,'Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.',7,'UNIQUE',NULL,NULL,1),(43,8,'Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings.',8,'UNIQUE',NULL,NULL,1),(44,5,'check 1',0,'UNIQUE',NULL,NULL,1);
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
  `source_department_id` int DEFAULT NULL,
  `status` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_psos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_psos`
--

LOCK TABLES `curriculum_psos` WRITE;
/*!40000 ALTER TABLE `curriculum_psos` DISABLE KEYS */;
INSERT INTO `curriculum_psos` VALUES (1,2,'Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.',0,'UNIQUE',NULL,NULL,1),(2,2,'Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.',1,'UNIQUE',NULL,NULL,1),(6,4,'s. csdcscsC',0,'CLUSTER',5,NULL,1),(9,5,'s. csdcscsC',0,'UNIQUE',NULL,NULL,1),(10,5,'s. csdcscsC',1,'UNIQUE',NULL,NULL,1),(15,7,'Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.',0,'UNIQUE',NULL,NULL,1),(16,7,'Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.',1,'UNIQUE',NULL,NULL,1),(17,8,'Apply suitable algorithmic thinking and data management practices to design develop, and evaluate effective solutions for real-life and research problems.',0,'UNIQUE',NULL,NULL,1),(18,8,'Design and develop cost-effective solutions based on cutting-edge hardware and software tools and techniques to meet global requirements.',1,'UNIQUE',NULL,NULL,1),(19,5,'check 1',2,'UNIQUE',NULL,NULL,1);
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
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `curriculum_vision`
--

LOCK TABLES `curriculum_vision` WRITE;
/*!40000 ALTER TABLE `curriculum_vision` DISABLE KEYS */;
INSERT INTO `curriculum_vision` VALUES (2,4,'To excel in the field of Computer Science and Engineering',1),(3,6,'',1),(4,8,'',1),(5,7,'ccdcacdccasdcsdcCHECK 1',1),(6,9,'',1),(7,10,'',1),(8,11,'To excel in the field of Computer Science and Engineering, to meet the emerging needsof the\nindustry, society, and beyond.',1),(9,12,'',1),(10,14,'',1),(11,15,'',1);
/*!40000 ALTER TABLE `curriculum_vision` ENABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `honour_cards`
--

LOCK TABLES `honour_cards` WRITE;
/*!40000 ALTER TABLE `honour_cards` DISABLE KEYS */;
INSERT INTO `honour_cards` VALUES (6,4,'Honour Vertical *','2026-01-13 09:52:21','UNIQUE',NULL,1),(7,14,'Honour card','2026-01-13 09:56:19','UNIQUE',NULL,1);
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
INSERT INTO `honour_verticals` VALUES (3,6,'data','2026-01-13 09:52:29',1);
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
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `normal_cards`
--

LOCK TABLES `normal_cards` WRITE;
/*!40000 ALTER TABLE `normal_cards` DISABLE KEYS */;
INSERT INTO `normal_cards` VALUES (3,4,1,'UNIQUE',NULL,'semester',1),(41,4,NULL,'UNIQUE',NULL,'elective',1),(42,4,1,'UNIQUE',NULL,'vertical',1),(43,4,2,'UNIQUE',NULL,'vertical',1),(44,4,2,'UNIQUE',NULL,'semester',1),(47,4,3,'UNIQUE',NULL,'vertical',1),(48,10,1,'UNIQUE',NULL,'semester',1),(49,10,2,'UNIQUE',NULL,'semester',1),(50,10,3,'UNIQUE',NULL,'semester',1),(51,10,4,'UNIQUE',NULL,'semester',1),(52,10,5,'UNIQUE',NULL,'semester',1),(53,10,6,'UNIQUE',NULL,'semester',1),(54,10,7,'UNIQUE',NULL,'semester',1),(55,10,8,'UNIQUE',NULL,'semester',1),(71,14,1,'UNIQUE',NULL,'semester',1);
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
INSERT INTO `peo_po_mapping` VALUES (11,4,1,1,3,1),(12,4,1,2,3,1),(13,4,1,3,3,1),(14,4,1,4,3,1),(15,4,1,5,3,1),(16,4,1,6,3,1),(17,4,1,7,3,1),(18,4,1,11,3,1),(19,4,1,12,3,1),(20,4,2,1,3,1),(21,4,2,2,3,1),(22,4,2,3,3,1),(23,4,2,4,3,1),(24,4,2,5,3,1),(25,4,2,6,3,1),(26,4,2,7,3,1),(27,4,2,10,3,1),(28,4,3,8,3,1),(29,4,3,9,3,1),(30,4,3,10,3,1),(31,4,3,11,3,1),(32,4,3,12,3,1);
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
  `status` int DEFAULT '1',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `student_id` (`student_id`) USING BTREE,
  CONSTRAINT `school_details_ibfk_1` FOREIGN KEY (`student_id`) REFERENCES `students` (`student_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `school_details`
--

LOCK TABLES `school_details` WRITE;
/*!40000 ALTER TABLE `school_details` DISABLE KEYS */;
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
) ENGINE=InnoDB AUTO_INCREMENT=84 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sharing_tracking`
--

LOCK TABLES `sharing_tracking` WRITE;
/*!40000 ALTER TABLE `sharing_tracking` DISABLE KEYS */;
INSERT INTO `sharing_tracking` VALUES (31,2,3,'mission',2,28,'2025-12-25 17:27:27'),(33,2,3,'mission',4,30,'2025-12-25 18:09:15'),(34,2,6,'mission',4,31,'2025-12-25 18:09:16'),(64,2,3,'peos',1,25,'2025-12-26 06:42:35'),(65,2,6,'peos',1,26,'2025-12-26 06:42:35'),(76,2,3,'semester',3,33,'2025-12-26 09:19:28'),(77,2,6,'semester',3,34,'2025-12-26 09:19:34'),(78,2,3,'pos',1,22,'2025-12-26 09:35:55'),(79,2,3,'semester',42,33,'2026-01-05 06:26:00'),(80,2,3,'semester',43,45,'2026-01-05 06:34:07'),(81,2,6,'semester',43,46,'2026-01-05 06:34:12'),(82,2,3,'semester',44,45,'2026-01-05 08:49:05'),(83,2,6,'semester',44,46,'2026-01-05 08:49:06');
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
  `status` int DEFAULT '1',
  PRIMARY KEY (`student_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `students`
--

LOCK TABLES `students` WRITE;
/*!40000 ALTER TABLE `students` DISABLE KEYS */;
INSERT INTO `students` VALUES (1,'2024UCS1999','7376241CS199','4388757858','2211','39889','check1','Male','2006-01-20',20,'checkf','checkm','checkg','Hindu','Indian','OC','tamil','A+','384837463777','checkoccu','checkdesig','check place',50000.00,1),(2,'7376242122','7847VDjh44','7488','38898','188748','ahjhdjn','Male','2006-01-04',20,'vvbdfv','vjdfjvjhadbvbdhfbh','b vhhjdvjd','Hindu','Indian','OC','bsdcsbvsbfv','A-','122212221222','fhbejbfjebrffebrhbf','ehrbbfebfqb','bberbfhjbf',2333333.00,NULL),(3,'745','654','456','456','456','dsf','Male','2025-10-09',0,'Ramajayam Velukonar','234234324','sdfsd','Muslim','Indian','OC','tamil','B+','223423423434','234234','sfsf','23423',234234.00,NULL);
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
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `syllabus`
--

LOCK TABLES `syllabus` WRITE;
/*!40000 ALTER TABLE `syllabus` DISABLE KEYS */;
INSERT INTO `syllabus` VALUES (6,'Module 1','Module 1',0,17,1),(7,'Experiment 2','Experiment 2',1,17,1),(8,'Module 1','Module 1',0,18,1),(11,'Unit 1','Unit 1',0,83,1),(12,'Unit 2','Unit 2',1,83,1);
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
INSERT INTO `syllabus_titles` VALUES (6,6,'Experiment 1',5,'Experiment 1',0),(7,8,'lineAR ',5,'lineAR ',0),(8,11,'MATHEMATICS MODELING OF LINEAR FUNCTIONS',8,'MATHEMATICS MODELING OF LINEAR FUNCTIONS',0);
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
INSERT INTO `syllabus_topics` VALUES (14,6,'Rank of a Matrix','Rank of a Matrix',0),(15,7,'HSFFS','HSFFS',0);
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
  `dept` int DEFAULT NULL,
  `desg` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `theory` int DEFAULT '0',
  `lab` int DEFAULT '0',
  `last_login` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'Active',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE,
  UNIQUE KEY `email` (`email`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teachers`
--

LOCK TABLES `teachers` WRITE;
/*!40000 ALTER TABLE `teachers` DISABLE KEYS */;
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
INSERT INTO `users` VALUES (1,'admin','$2a$10$H4BOz6nXYrnGQVYwC0eAQul.YF3LyhWTpb7xUf1HAKPT8y18DYDaq','System Administrator','admin@example.com','admin',1,'2026-01-07 05:52:45','2026-01-22 08:44:33','2026-01-22 14:14:34');
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

-- Dump completed on 2026-01-22 15:52:55
