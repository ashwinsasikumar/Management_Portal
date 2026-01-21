# ************************************************************
# Sequel Ace SQL dump
# Version 20096
#
# https://sequel-ace.com/
# https://github.com/Sequel-Ace/Sequel-Ace
#
# Host: localhost (MySQL 9.5.0)
# Database: cms_local
# Generation Time: 2026-01-19 04:24:06 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table cluster_departments
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cluster_departments`;

CREATE TABLE `cluster_departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `cluster_id` int NOT NULL,
  `curriculum_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_department` (`curriculum_id`) USING BTREE,
  KEY `cluster_id` (`cluster_id`) USING BTREE,
  CONSTRAINT `cluster_departments_ibfk_1` FOREIGN KEY (`cluster_id`) REFERENCES `clusters` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `cluster_departments` WRITE;
/*!40000 ALTER TABLE `cluster_departments` DISABLE KEYS */;

INSERT INTO `cluster_departments` (`id`, `cluster_id`, `curriculum_id`, `created_at`)
VALUES
	(2,1,3,'2025-12-25 05:10:57'),
	(3,2,4,'2025-12-25 05:31:02'),
	(4,2,5,'2025-12-25 05:31:10'),
	(6,1,7,'2026-01-07 08:12:58'),
	(7,1,9,'2026-01-07 08:13:04'),
	(8,1,6,'2026-01-12 11:33:59'),
	(10,1,2,'2026-01-12 12:04:13');

/*!40000 ALTER TABLE `cluster_departments` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table clusters
# ------------------------------------------------------------

DROP TABLE IF EXISTS `clusters`;

CREATE TABLE `clusters` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `clusters` WRITE;
/*!40000 ALTER TABLE `clusters` DISABLE KEYS */;

INSERT INTO `clusters` (`id`, `name`, `description`, `created_at`)
VALUES
	(1,'computer cluster','cse departments','2025-12-25 05:07:49'),
	(2,'mechanical cluster','mechanical departments','2025-12-25 05:30:51');

/*!40000 ALTER TABLE `clusters` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table co_po_mapping
# ------------------------------------------------------------

DROP TABLE IF EXISTS `co_po_mapping`;

CREATE TABLE `co_po_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `co_index` int NOT NULL,
  `po_index` int NOT NULL,
  `mapping_value` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_copo_course` (`course_id`) USING BTREE,
  CONSTRAINT `fk_copo_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=86 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `co_po_mapping` WRITE;
/*!40000 ALTER TABLE `co_po_mapping` DISABLE KEYS */;

INSERT INTO `co_po_mapping` (`id`, `course_id`, `co_index`, `po_index`, `mapping_value`)
VALUES
	(1,1,0,1,3),
	(2,1,0,2,2),
	(3,1,0,3,1),
	(4,2,0,1,3),
	(5,2,0,2,2),
	(6,2,0,3,1),
	(7,2,0,4,1),
	(8,2,0,5,2),
	(9,2,0,9,1),
	(10,2,0,10,1),
	(11,2,0,12,2),
	(12,2,1,1,2),
	(13,2,1,2,3),
	(14,2,1,3,2),
	(15,2,1,4,2),
	(16,2,1,5,1),
	(17,2,1,9,1),
	(18,2,1,10,1),
	(19,2,1,12,3),
	(20,2,2,1,1),
	(21,2,2,2,2),
	(22,2,2,3,3),
	(23,2,2,4,2),
	(24,2,2,5,3),
	(25,2,2,7,1),
	(26,2,2,9,1),
	(27,2,2,12,3),
	(28,4,0,1,1),
	(29,4,0,2,2),
	(30,4,0,3,1),
	(31,4,0,4,3),
	(32,4,0,5,1),
	(33,4,0,6,2),
	(34,4,0,7,1),
	(35,4,0,10,2),
	(36,4,0,11,3),
	(37,4,0,12,1),
	(38,4,1,1,2),
	(39,4,1,2,2),
	(40,4,1,3,1),
	(41,4,1,4,3),
	(42,4,1,6,1),
	(43,4,1,9,2),
	(44,4,1,10,3),
	(45,4,1,11,1),
	(46,4,1,12,2);

/*!40000 ALTER TABLE `co_po_mapping` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table co_pso_mapping
# ------------------------------------------------------------

DROP TABLE IF EXISTS `co_pso_mapping`;

CREATE TABLE `co_pso_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `co_index` int NOT NULL,
  `pso_index` int NOT NULL,
  `mapping_value` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_copso_course` (`course_id`) USING BTREE,
  CONSTRAINT `fk_copso_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `co_pso_mapping` WRITE;
/*!40000 ALTER TABLE `co_pso_mapping` DISABLE KEYS */;

INSERT INTO `co_pso_mapping` (`id`, `course_id`, `co_index`, `pso_index`, `mapping_value`)
VALUES
	(1,1,0,1,3),
	(2,1,0,2,2),
	(3,2,0,1,3),
	(4,2,0,2,2),
	(5,2,0,3,1),
	(6,2,1,1,2),
	(7,2,1,2,3),
	(8,2,1,3,1),
	(9,2,2,1,2),
	(10,2,2,2,3),
	(11,2,2,3,2),
	(12,4,0,1,2),
	(13,4,0,2,3),
	(14,4,0,3,1),
	(15,4,1,1,2),
	(16,4,1,2,1),
	(17,4,1,3,3);

/*!40000 ALTER TABLE `co_pso_mapping` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table course_experiment_topics
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_experiment_topics`;

CREATE TABLE `course_experiment_topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `experiment_id` int NOT NULL,
  `topic_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `topic_order` int DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_exp_topics` (`experiment_id`) USING BTREE,
  CONSTRAINT `course_experiment_topics_ibfk_1` FOREIGN KEY (`experiment_id`) REFERENCES `course_experiments` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `course_experiment_topics` WRITE;
/*!40000 ALTER TABLE `course_experiment_topics` DISABLE KEYS */;

INSERT INTO `course_experiment_topics` (`id`, `experiment_id`, `topic_text`, `topic_order`, `created_at`)
VALUES
	(12,4,'Assess the physical parameters of different materials for engineering applications like radius, thickness and\ndiameter to design the electrical wires, bridges and clothes',0,'2026-01-07 04:25:43'),
	(13,4,'Evaluate the elastic nature of different solid materials for modern industrial applications like shock absorbers\nof vehicles',1,'2026-01-07 04:25:43');

/*!40000 ALTER TABLE `course_experiment_topics` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table course_experiments
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_experiments`;

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
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `course_experiments` WRITE;
/*!40000 ALTER TABLE `course_experiments` DISABLE KEYS */;

INSERT INTO `course_experiments` (`id`, `course_id`, `experiment_number`, `experiment_name`, `hours`, `created_at`, `updated_at`)
VALUES
	(4,83,1,'Experiment 1',3,'2026-01-07 04:16:38','2026-01-07 04:25:43'),
	(5,83,2,'Experiment 2',7,'2026-01-07 04:20:01','2026-01-07 04:20:01'),
	(6,83,3,'Experiment 3',4,'2026-01-07 04:20:11','2026-01-07 04:20:11'),
	(10,83,4,'Experiment 4',8,'2026-01-07 04:30:23','2026-01-07 04:30:23');

/*!40000 ALTER TABLE `course_experiments` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table course_objectives
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_objectives`;

CREATE TABLE `course_objectives` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `objective` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_course_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `course_objectives_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `course_objectives` WRITE;
/*!40000 ALTER TABLE `course_objectives` DISABLE KEYS */;

INSERT INTO `course_objectives` (`id`, `course_id`, `objective`, `position`)
VALUES
	(4,83,'To impart mathematical modeling to describe and explore real-world phenomena and data.',0),
	(5,83,'To provide basic understanding on Linear, quadratic, power and polynomial, exponential, and multi variable models',1),
	(6,83,'Summarize and apply the methodologies involved in framing the real world problems related to fundamental principles of polynomial equations',2);

/*!40000 ALTER TABLE `course_objectives` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table course_outcomes
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_outcomes`;

CREATE TABLE `course_outcomes` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `outcome` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_course_outcome_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `fk_course_outcomes_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=140 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `course_outcomes` WRITE;
/*!40000 ALTER TABLE `course_outcomes` DISABLE KEYS */;

INSERT INTO `course_outcomes` (`id`, `course_id`, `outcome`, `position`)
VALUES
	(136,83,'Implement the concepts of mathematical modeling based on linear functions in Engineering.',0),
	(137,83,'Formulate the real-world problems as a quadratic function model',1),
	(138,83,'Demonstrate the real-world phenomena and data into Power and Polynomial functions',2),
	(139,83,'Apply the concept of mathematical modeling of exponential functions in Engineering',3);

/*!40000 ALTER TABLE `course_outcomes` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table course_prerequisites
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_prerequisites`;

CREATE TABLE `course_prerequisites` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `prerequisite` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_course_prerequisite_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `fk_course_prerequisites_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;



# Dump of table course_references
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_references`;

CREATE TABLE `course_references` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `reference_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_course_reference_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `fk_course_references_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;



# Dump of table course_selflearning
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_selflearning`;

CREATE TABLE `course_selflearning` (
  `course_id` int NOT NULL,
  `total_hours` int NOT NULL,
  PRIMARY KEY (`course_id`) USING BTREE,
  CONSTRAINT `fk_course_selflearning_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `course_selflearning` WRITE;
/*!40000 ALTER TABLE `course_selflearning` DISABLE KEYS */;

INSERT INTO `course_selflearning` (`course_id`, `total_hours`)
VALUES
	(22,0),
	(83,0);

/*!40000 ALTER TABLE `course_selflearning` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table course_selflearning_resources
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_selflearning_resources`;

CREATE TABLE `course_selflearning_resources` (
  `id` int NOT NULL AUTO_INCREMENT,
  `main_id` int NOT NULL,
  `internal_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_main_position` (`main_id`,`position`) USING BTREE,
  CONSTRAINT `course_selflearning_resources_ibfk_1` FOREIGN KEY (`main_id`) REFERENCES `course_selflearning_topics` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;



# Dump of table course_selflearning_topics
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_selflearning_topics`;

CREATE TABLE `course_selflearning_topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `main_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_course_position` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `course_selflearning_topics_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;



# Dump of table course_teamwork
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_teamwork`;

CREATE TABLE `course_teamwork` (
  `course_id` int NOT NULL,
  `total_hours` int NOT NULL,
  PRIMARY KEY (`course_id`) USING BTREE,
  CONSTRAINT `course_teamwork_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `course_teamwork` WRITE;
/*!40000 ALTER TABLE `course_teamwork` DISABLE KEYS */;

INSERT INTO `course_teamwork` (`course_id`, `total_hours`)
VALUES
	(22,0),
	(83,0);

/*!40000 ALTER TABLE `course_teamwork` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table course_teamwork_activities
# ------------------------------------------------------------

DROP TABLE IF EXISTS `course_teamwork_activities`;

CREATE TABLE `course_teamwork_activities` (
  `id` int NOT NULL AUTO_INCREMENT,
  `course_id` int NOT NULL,
  `activity` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `course_id` (`course_id`,`position`) USING BTREE,
  CONSTRAINT `course_teamwork_activities_ibfk_1` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `course_teamwork_activities` WRITE;
/*!40000 ALTER TABLE `course_teamwork_activities` DISABLE KEYS */;

INSERT INTO `course_teamwork_activities` (`id`, `course_id`, `activity`, `position`)
VALUES
	(27,22,'hello',0);

/*!40000 ALTER TABLE `course_teamwork_activities` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table courses
# ------------------------------------------------------------

DROP TABLE IF EXISTS `courses`;

CREATE TABLE `courses` (
  `course_id` int NOT NULL AUTO_INCREMENT,
  `course_code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `course_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `course_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `credit` int DEFAULT NULL,
  `theory_total_hrs` int DEFAULT '0',
  `activity_hrs` int DEFAULT '0',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `lecture_hrs` int DEFAULT '0',
  `tutorial_hrs` int DEFAULT '0',
  `practical_hrs` int DEFAULT '0',
  `cia_marks` int DEFAULT '40',
  `see_marks` int DEFAULT '60',
  `total_marks` int GENERATED ALWAYS AS ((`cia_marks` + `see_marks`)) STORED,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `curriculum_ref_id` int DEFAULT NULL,
  `activity_total_hrs` int DEFAULT '0',
  `tutorial_total_hrs` int DEFAULT '0',
  `total_hrs` int GENERATED ALWAYS AS (((`theory_total_hrs` + `activity_total_hrs`) + `tutorial_total_hrs`)) STORED,
  `practical_total_hrs` int DEFAULT NULL,
  PRIMARY KEY (`course_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=88 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `courses` WRITE;
/*!40000 ALTER TABLE `courses` DISABLE KEYS */;

INSERT INTO `courses` (`course_id`, `course_code`, `course_name`, `course_type`, `credit`, `theory_total_hrs`, `activity_hrs`, `category`, `lecture_hrs`, `tutorial_hrs`, `practical_hrs`, `cia_marks`, `see_marks`, `total_marks`, `visibility`, `source_curriculum_id`, `curriculum_ref_id`, `activity_total_hrs`, `tutorial_total_hrs`, `total_hrs`, `practical_total_hrs`)
VALUES
	(1,'CS101','Introduction to Programming','Theory',3,0,0,'Core',0,0,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(2,'CS3801','Cloud Computing','Theory',3,0,0,'Elective',0,0,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(3,'CS201','Data Structures','Theory',4,0,0,'Core',3,1,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(4,'CS3501','Database Management Systems','Theory',4,0,0,'Core',3,1,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(17,'22MA101','ENGINEERING MATHEMATICS I ','Theory',4,2,1,'BS - Basic Sciences',0,10,3,40,60,100,'CLUSTER',NULL,NULL,0,0,2,NULL),
	(18,'22PH102 ','ENGINEERING PHYSICS ','Theory',3,1,3,'BS - Basic Sciences',0,10,2,40,60,100,'UNIQUE',NULL,NULL,0,0,1,NULL),
	(19,'22CSH01 ','EXPLORATORY DATA ANALYSIS ','Theory',3,0,300,'PE - Professional Elective',2,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(20,'26MA101','Linear Algebra and Calculus','Theory',4,0,0,'ES - Engineering Sciences',3,1,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(21,'26PH102','Engineering Physics','Theory',3,0,0,'ES - Engineering Sciences',3,1,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(22,'26CH103','Engineering Chemistry','Theory',2,0,0,'ES - Engineering Sciences',2,0,0,30,70,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(23,'26GE004','Digital Computer Electronics','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(24,'26GE005','Problem Solving using C','Theory',3,0,0,'ES - Engineering Sciences',2,0,0,30,70,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(25,'26HS001','Communicative English','Theory',2,0,0,'HSS - Humanities and Social Sciences',2,0,0,30,70,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(26,'26HS002','தமிழர் மரபு / Heritage of Tamils ','Theory',1,0,0,'HSS - Humanities and Social Sciences',1,0,0,15,85,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(27,'26PH108','Physical Science Laboratory','Experiment',2,0,0,'ES - Engineering Sciences',0,0,4,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(28,'26GE006','C Programming Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(29,'26SD001','Skill Development Course  I','Experiment',1,0,0,'EEC - Employability Enhancement Course',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(30,'26MA201','Differential Equations and Transforms','Theory',4,0,0,'ES - Engineering Sciences',3,1,0,60,40,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(31,'26PH202','Materials Science','Theory',3,0,0,'BS - Basic Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(32,'26CS203','Fundamentals of Web Principles','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(33,'26CS204','Computer Organization and Architecture','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(34,'26GE007','Python Programming','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(35,'26HS005','Professional Communication','Theory',2,0,0,'HSS - Humanities and Social Sciences',2,0,0,30,70,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(36,'26HS006','தமிழரும் தொழில்நுட்பமும் / Tamils and Technology','Theory',1,0,0,'HSS - Humanities and Social Sciences',1,0,0,15,85,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(37,'26GE008','Python Programming Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(38,'26CS209','Web Principles Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(39,'26SD002','Skill Development Course II','Experiment',1,0,0,'EEC - Employability Enhancement Course',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(40,'26CS301','Discrete Mathematics','Theory',4,0,0,'BS - Basic Sciences',3,1,0,60,60,120,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(41,'26CS302','Data Structures and Algorithms','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(42,'26CS303','Operating Systems','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(43,'26CS304','Object Oriented Programming with Java','Theory',3,0,0,'ES',2,0,2,30,70,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(44,'26CS305','Software Engineering','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(45,'26CS306','Database Management Systems','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(46,'26CS307','Standards in Computer Science','Theory',1,0,0,'ES - Engineering Sciences',1,0,0,15,85,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(47,'26CS308','Data Structures and Algorithms Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(48,'26CS309','Database Management Systems Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(49,'26CS310','Design Thinking and Innovation Laboratory (AICTE, & NEP)','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(50,'26CS401','Probability and Statistics','Theory',4,0,0,'ES - Engineering Sciences',3,1,0,60,40,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(51,'26CS402','Full Stack Development','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(52,'26CS403','Artificial Intelligence Essentials','Theory',3,0,0,'ES',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(53,'26CS404','Design and Analysis of Algorithms','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(54,'26CS405','Theory of Computation','Theory',4,0,0,'ES - Engineering Sciences',3,1,0,60,40,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(55,'26CS406','Computer Networks','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(56,'26HS009','Environmental Sciences and Sustainability ','Theory',2,0,0,'HSS - Humanities and Social Sciences',2,0,0,30,70,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(57,'26CS408','Full Stack Development Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(58,'26CS409','Computer Networks Laboratory','Lab',1,0,0,'ES',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(59,'26CS410','Community Engagement Project','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(60,'26CS501','Compiler Design','Theory',4,0,0,'ES - Engineering Sciences',3,1,0,60,40,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(61,'26CS502','Cloud Infrastructure Services','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(62,'26CS503','Bigdata Analytics','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(63,'26CS504','Machine Learning Essentials','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(64,'26XXIV','Professional Elective IV','Theory',3,0,0,'ES',0,0,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(65,'26CS507','Cloud Infrastructure Services Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(66,'26CS508','Machine Learning Essentials Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(67,'26CS509','Technology Integration Project','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(68,'26CS601','Software Project Management and Quality Assurance','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(69,'26CS602','Deep Learning','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(70,'26CS603','Cryptography and Cyber Security','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(71,'26CS607','Software Project Management and Quality Assurance Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(72,'26CS608','Deep Learning Laboratory ','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(73,'26CS609','Innovation and Product Development Project / Industry Oriented Course / Summer Internship','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(74,'26CS701','Generative AI and Large Language Models','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(75,'26CS702','IoT and Edge Computing','Theory',3,0,0,'ES - Engineering Sciences',3,0,0,45,55,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(76,'26XXIV','Professional Elective IV','Theory',3,0,0,'ES',0,0,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(77,'26XXV','Professional Elective V','Theory',3,0,0,'ES',0,0,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(78,'26XXVI','Professional Elective VI','Theory',3,0,0,'ES - Engineering Sciences',0,0,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(79,'26CS706','Generative AI and Large Language Models Laboratory','Experiment',1,0,0,'ES - Engineering Sciences',0,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(80,'26CS707','Capstone Project work Level I / Internship Pro','Experiment',3,0,0,'ES - Engineering Sciences',0,0,6,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(81,'26CS801','Capstone Project Work Level II / Internship Project / Startup Product','Experiment',8,0,0,'ES - Engineering Sciences',0,0,16,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(82,'hello','efvwbg','Theory',2,0,0,'ES - Engineering Sciences',1,0,2,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(83,'wcnejce','ervewvwrv','Theory',3,0,0,'BS - Basic Sciences',0,0,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(84,'bdzhgd','bgbdfb','Theory',2,13,0,'ES - Engineering Sciences',1,134,0,40,60,100,'UNIQUE',NULL,NULL,0,0,13,NULL),
	(85,'CS130','check1','Theory',3,45,30,'BS - Basic Sciences',3,15,0,40,60,100,'UNIQUE',NULL,NULL,0,0,45,NULL),
	(86,'CS230','check 2','Theory',3,0,2,'ES - Engineering Sciences',3,1,0,40,60,100,'UNIQUE',NULL,NULL,0,0,0,NULL),
	(87,'CS303','check 3','Theory',3,45,2,'BS - Basic Sciences',3,1,0,40,60,100,'UNIQUE',NULL,NULL,30,15,90,NULL);

/*!40000 ALTER TABLE `courses` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum`;

CREATE TABLE `curriculum` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `academic_year` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `curriculum_template` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '2026',
  `template_config` json DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `max_credits` int DEFAULT '0',
  `curriculum_ref_id` int DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum` WRITE;
/*!40000 ALTER TABLE `curriculum` DISABLE KEYS */;

INSERT INTO `curriculum` (`id`, `name`, `academic_year`, `curriculum_template`, `template_config`, `created_at`, `max_credits`, `curriculum_ref_id`)
VALUES
	(4,'BE - CSE ','2025-2026','2026',NULL,'2025-11-12 06:09:01',163,NULL),
	(6,'BT - AIML','2024-2028','2026',NULL,'2025-12-24 10:33:52',164,NULL),
	(7,'ME - MECHANICAL','2024-2028','2026',NULL,'2025-12-25 05:27:44',145,NULL),
	(10,'R2026-CSE','2024-2025','2026',NULL,'2026-01-06 04:47:56',162,NULL),
	(14,'check 1','2024-2025','2022',NULL,'2026-01-13 09:56:01',162,NULL);

/*!40000 ALTER TABLE `curriculum` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum_courses
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum_courses`;

CREATE TABLE `curriculum_courses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `semester_id` int NOT NULL,
  `course_id` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_rc_regulation` (`curriculum_id`) USING BTREE,
  KEY `fk_rc_semester` (`semester_id`) USING BTREE,
  KEY `fk_rc_course` (`course_id`) USING BTREE,
  CONSTRAINT `fk_rc_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_rc_regulation` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_rc_semester` FOREIGN KEY (`semester_id`) REFERENCES `normal_cards` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=211 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum_courses` WRITE;
/*!40000 ALTER TABLE `curriculum_courses` DISABLE KEYS */;

INSERT INTO `curriculum_courses` (`id`, `curriculum_id`, `semester_id`, `course_id`)
VALUES
	(102,4,3,18),
	(133,6,33,17),
	(134,6,33,18),
	(137,10,48,20),
	(138,10,48,21),
	(139,10,48,22),
	(140,10,48,23),
	(141,10,48,24),
	(142,10,48,25),
	(143,10,48,26),
	(144,10,48,27),
	(145,10,48,28),
	(146,10,48,29),
	(147,10,49,30),
	(148,10,49,31),
	(149,10,49,32),
	(150,10,49,33),
	(151,10,49,34),
	(152,10,49,35),
	(153,10,49,36),
	(154,10,49,37),
	(155,10,49,38),
	(156,10,49,39),
	(157,10,50,40),
	(158,10,50,41),
	(159,10,50,42),
	(160,10,50,43),
	(161,10,50,44),
	(162,10,50,45),
	(163,10,50,46),
	(164,10,50,47),
	(165,10,50,48),
	(166,10,50,49),
	(167,10,51,50),
	(168,10,51,51),
	(169,10,51,52),
	(170,10,51,53),
	(171,10,51,54),
	(172,10,51,55),
	(173,10,51,56),
	(174,10,51,57),
	(175,10,51,58),
	(176,10,51,59),
	(177,10,52,60),
	(178,10,52,61),
	(179,10,52,62),
	(180,10,52,63),
	(181,10,52,64),
	(182,10,52,64),
	(183,10,52,65),
	(184,10,52,66),
	(185,10,52,67),
	(186,10,53,68),
	(187,10,53,69),
	(188,10,53,70),
	(189,10,53,71),
	(190,10,53,72),
	(191,10,53,73),
	(192,10,53,64),
	(193,10,53,64),
	(194,10,53,64),
	(195,10,54,74),
	(196,10,54,75),
	(197,10,54,76),
	(198,10,54,77),
	(199,10,54,78),
	(200,10,54,79),
	(201,10,54,80),
	(202,10,55,81),
	(207,4,3,1),
	(208,4,3,85),
	(209,4,3,86),
	(210,4,3,87);

/*!40000 ALTER TABLE `curriculum_courses` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum_logs
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum_logs`;

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
) ENGINE=InnoDB AUTO_INCREMENT=240 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum_logs` WRITE;
/*!40000 ALTER TABLE `curriculum_logs` DISABLE KEYS */;

INSERT INTO `curriculum_logs` (`id`, `curriculum_id`, `action`, `description`, `changed_by`, `created_at`, `diff`)
VALUES
	(1,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:21:59',NULL),
	(2,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:24:07',NULL),
	(3,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:28:35',NULL),
	(4,4,'Department Overview Updated','Updated department vision, mission, PEOs, POs, and PSOs','System','2025-11-18 05:58:12','{\"mission\": {\"new\": [\"To impart need based education \", \"To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.\", \"To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.\"], \"old\": [\"To impart need based education to meet the requirements of the industry and society.\", \"To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.\", \"To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.\"]}}'),
	(5,4,'Vision Updated','Updated department vision','System','2025-11-18 06:18:04','{\"vision\": {\"new\": \"To excel \", \"old\": \"To excel in the field of Computer Science and Engineering\"}}'),
	(6,4,'Vision Updated','Updated department vision','System','2025-11-18 06:18:25','{\"vision\": {\"new\": \"To excel in the field of Computer Science and Engineering\", \"old\": \"To excel \"}}'),
	(7,4,'Mission[3] Added','Added Mission item at index 3','System','2025-11-18 06:24:25','{\"Mission[3]\": {\"new\": \"to be good\", \"old\": \"\"}}'),
	(8,4,'Mission[3] Deleted','Deleted Mission item at index 3','System','2025-11-18 06:24:44','{\"Mission[3]\": {\"new\": \"\", \"old\": \"to be good\"}}'),
	(10,4,'CO-PO/PSO Mapping Saved','Updated CO-PO and CO-PSO mappings for course: ENGINEERING MATHEMATICS I','System','2025-11-19 04:25:50','{\"co_po_mappings\": {\"new\": {\"CO0-PO1\": 1, \"CO0-PO2\": 2, \"CO1-PO1\": 2, \"CO1-PO2\": 2, \"CO2-PO1\": 2, \"CO2-PO2\": 1, \"CO3-PO1\": 2, \"CO3-PO2\": 2, \"CO4-PO1\": 1, \"CO4-PO2\": 2}, \"old\": {\"CO0-PO1\": 1, \"CO0-PO2\": 2, \"CO1-PO1\": 2, \"CO1-PO2\": 2, \"CO2-PO1\": 2, \"CO2-PO2\": 1, \"CO3-PO1\": 2, \"CO3-PO2\": 2, \"CO4-PO1\": 1, \"CO4-PO2\": 2}}, \"co_pso_mappings\": {\"new\": {\"CO0-PSO1\": 1, \"CO1-PSO2\": 2, \"CO2-PSO2\": 1, \"CO3-PSO1\": 1, \"CO3-PSO3\": 3, \"CO4-PSO2\": 2, \"CO4-PSO3\": 2}, \"old\": {}}}'),
	(11,4,'Course Updated','Updated course: 22CH103  - ENGINEERING CHEMISTRY I ','System','2025-11-19 10:26:39','{\"credit\": {\"new\": 141, \"old\": 3}, \"category\": {\"new\": \"BS\", \"old\": \"BS - Basic Sciences\"}}'),
	(12,4,'Course Updated','Updated course: 22CS108  - COMPREHENSIVE WORK','System','2025-11-19 10:27:00','{\"credit\": {\"new\": 4, \"old\": 1}, \"category\": {\"new\": \"PE\", \"old\": \"EEC - Employability Enhancement Course\"}}'),
	(13,4,'Course Updated','Updated course: 22GE001  - FUNDAMENTALS OF COMPUTING ','System','2025-11-19 10:27:17','{\"credit\": {\"new\": 4, \"old\": 3}, \"category\": {\"new\": \"BS\", \"old\": \"ES - Engineering Sciences\"}}'),
	(14,4,'Course Updated','Updated course: 22CH103  - ENGINEERING CHEMISTRY I ','System','2025-12-22 04:10:33','{\"credit\": {\"new\": 3, \"old\": 141}}'),
	(15,4,'Course Updated','Updated course: 22CH103  - ENGINEERING CHEMISTRY I ','System','2025-12-22 04:10:33','{\"credit\": {\"new\": 3, \"old\": 141}}'),
	(16,4,'Course Updated','Updated course: 22MA101 - ENGINEERING MATHEMATICS I','System','2025-12-22 06:01:23','{\"category\": {\"new\": \"BS\", \"old\": \"BS - Basic Sciences\"}, \"tutorial_hours\": {\"new\": 30, \"old\": 1}}'),
	(17,4,'Course Updated','Updated course: 22MA101 - ENGINEERING MATHEMATICS I','System','2025-12-22 08:28:03','{\"lecture_hours\": {\"new\": 2, \"old\": 3}, \"practical_hours\": {\"new\": 1, \"old\": 0}}'),
	(18,6,'Curriculum Created','Created new curriculum: BT - AIML (2024-2028)','System','2025-12-24 10:33:52',NULL),
	(19,7,'Curriculum Created','Created new curriculum: ME - MECH (2024-2028)','System','2025-12-25 05:27:44',NULL),
	(22,7,'Department Overview Created','Created department vision, mission, PEOs, POs, and PSOs','System','2025-12-25 05:30:15',NULL),
	(23,7,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 10:20:02','{\"Mission[0]\": {\"new\": \"hello hi\", \"old\": \"hello\"}}'),
	(26,7,'Mission[1] Added','Added Mission item at index 1','System','2025-12-25 10:42:13','{\"Mission[1]\": {\"new\": \"hello\", \"old\": \"\"}}'),
	(27,7,'Mission[2] Added','Added Mission item at index 2','System','2025-12-25 10:42:13','{\"Mission[2]\": {\"new\": \"hi\", \"old\": \"\"}}'),
	(28,6,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 11:00:48','{\"Mission[0]\": {\"new\": \"To impart need based education ASHWIN\", \"old\": \"To impart need based education \"}}'),
	(29,6,'PO[0] Updated','Updated PO item at index 0','System','2025-12-25 11:21:24','{\"PO[0]\": {\"new\": \"ASHWIN\", \"old\": \"Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.\"}}'),
	(30,6,'PEO[0] Updated','Updated PEO item at index 0','System','2025-12-25 11:26:51','{\"PEO[0]\": {\"new\": \"hello\", \"old\": \"Graduates will apply computer science and engineering principles and practices to solvereal- world problems with their technical competence.\"}}'),
	(31,6,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 16:35:59','{\"Mission[0]\": {\"new\": \"To impart need based education MITHILESH\", \"old\": \"To impart need based education ASHWIN\"}}'),
	(32,6,'Mission[2] Updated','Updated Mission item at index 2','System','2025-12-25 16:35:59','{\"Mission[2]\": {\"new\": \"To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.\", \"old\": \"To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.\"}}'),
	(33,6,'Mission[1] Updated','Updated Mission item at index 1','System','2025-12-25 16:36:00','{\"Mission[1]\": {\"new\": \"To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.\", \"old\": \"To impart need based education \"}}'),
	(34,6,'Mission[3] Deleted','Deleted Mission item at index 3','System','2025-12-25 16:36:00','{\"Mission[3]\": {\"new\": \"\", \"old\": \"To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.\"}}'),
	(35,6,'PEO[0] Updated','Updated PEO item at index 0','System','2025-12-25 16:42:33','{\"PEO[0]\": {\"new\": \"removed\", \"old\": \"Graduates will apply computer science and engineering principles and practices to solvereal- world problems with their technical competence.\"}}'),
	(36,6,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 16:49:41','{\"Mission[0]\": {\"new\": \"MODIFIED\", \"old\": \"To impart need based education \"}}'),
	(37,6,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 16:55:23','{\"Mission[0]\": {\"new\": \"MODIFIED\", \"old\": \"To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.\"}}'),
	(38,6,'PEO[0] Updated','Updated PEO item at index 0','System','2025-12-25 17:00:59','{\"PEO[0]\": {\"new\": \"MODIFIED\", \"old\": \"Graduates will apply computer science and engineering principles and practices to solvereal- world problems with their technical competence.\"}}'),
	(39,6,'Mission[0] Deleted','Deleted Mission item at index 0','System','2025-12-25 17:21:44','{\"Mission[0]\": {\"new\": \"\", \"old\": \"MODIFIED\"}}'),
	(40,6,'PEO[0] Deleted','Deleted PEO item at index 0','System','2025-12-25 17:21:59','{\"PEO[0]\": {\"new\": \"\", \"old\": \"MODIFIED\"}}'),
	(41,6,'PSO[0] Deleted','Deleted PSO item at index 0','System','2025-12-25 17:21:59','{\"PSO[0]\": {\"new\": \"\", \"old\": \"Apply suitable algorithmic thinking and data management practices to design develop, and evaluate effective solutions for real-life and research problems.\"}}'),
	(42,6,'PSO[1] Deleted','Deleted PSO item at index 1','System','2025-12-25 17:21:59','{\"PSO[1]\": {\"new\": \"\", \"old\": \"Apply suitable algorithmic thinking and data management practices to design develop, and evaluate effective solutions for real-life and research problems.\"}}'),
	(43,6,'PSO[2] Deleted','Deleted PSO item at index 2','System','2025-12-25 17:21:59','{\"PSO[2]\": {\"new\": \"\", \"old\": \"Design and develop cost-effective solutions based on cutting-edge hardware and software tools and techniques to meet global requirements.\"}}'),
	(44,6,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 17:26:38','{\"Mission[0]\": {\"new\": \"MODIFIED\", \"old\": \"To impart need based education \"}}'),
	(45,4,'Mission[0] Updated','Updated Mission item at index 0','System','2025-12-25 17:28:37','{\"Mission[0]\": {\"new\": \"To impart need based education EXTRA\", \"old\": \"To impart need based education \"}}'),
	(48,4,'Course Added','Added course 26MA101 - Linear Algebra and Calculus to Semester 3','System','2025-12-26 05:43:55',NULL),
	(49,4,'CO-PO/PSO Mapping Saved','Updated CO-PO and CO-PSO mappings for course: Linear Algebra and Calculus','System','2025-12-26 05:53:59','{\"co_po_mappings\": {\"new\": {\"CO0-PO1\": 3, \"CO0-PO2\": 1, \"CO0-PO4\": 1, \"CO0-PO5\": 3, \"CO1-PO2\": 2, \"CO1-PO3\": 1, \"CO1-PO4\": 1, \"CO2-PO3\": 3, \"CO3-PO1\": 1}, \"old\": {}}, \"co_pso_mappings\": {\"new\": {}, \"old\": {}}}'),
	(50,4,'Course Added','Added course 22PH102  - ENGINEERING PHYSICS  to Semester 3','System','2025-12-26 05:55:54',NULL),
	(51,4,'Course Added','Added course 22CH103  - ENGINEERING CHEMISTRY I  to Semester 3','System','2025-12-26 05:56:36',NULL),
	(52,4,'Course Added','Added course 22MA101 - ENGINEERING MATHEMATICS I  to Semester 3','System','2025-12-26 06:17:48',NULL),
	(53,4,'Course Added','Added course 22PH102  - ENGINEERING PHYSICS  to Semester 3','System','2025-12-26 06:18:23',NULL),
	(55,4,'Semester Added','Added Semester 0','System','2026-01-05 04:25:39',NULL),
	(56,4,'Semester Added','Added electives','System','2026-01-05 04:28:31',NULL),
	(57,4,'Semester Added','Added electives','System','2026-01-05 04:32:16',NULL),
	(58,4,'Semester Added','Added semester','System','2026-01-05 04:35:58',NULL),
	(59,4,'Semester Added','Added electives','System','2026-01-05 04:36:11',NULL),
	(60,4,'Semester Added','Added vertical ','System','2026-01-05 04:36:32',NULL),
	(61,4,'Honour Card Added','Added Honour Card: honour verticals','System','2026-01-05 04:38:03',NULL),
	(62,4,'Semester Added','Added Elective','System','2026-01-05 04:52:35',NULL),
	(63,4,'Semester Added','Added Vertical','System','2026-01-05 04:53:39',NULL),
	(64,4,'Semester Added','Added Vertical','System','2026-01-05 04:53:53',NULL),
	(65,4,'Semester Added','Added Semester','System','2026-01-05 04:54:01',NULL),
	(66,4,'Semester Added','Added Vertical','System','2026-01-05 08:18:03',NULL),
	(67,4,'Course Removed','Removed course ENGINEERING MATHEMATICS I  from Semester 3','System','2026-01-05 09:35:37',NULL),
	(68,4,'PEO[0] Updated','Updated PEO item at index 0','System','2026-01-06 04:31:33','{\"PEO[0]\": {\"new\": \"Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.\", \"old\": \"Graduates will apply computer science and engineering principles and practices to solvereal- world problems with their technical competence.\"}}'),
	(69,4,'PEO[1] Updated','Updated PEO item at index 1','System','2026-01-06 04:31:33','{\"PEO[1]\": {\"new\": \"Pursue continuous learning in emerging technologies such as AI, data science, and cybersecurity to remain adaptable professionals.\", \"old\": \"Graduates will have the domain knowledge to pursue higher education and apply cuttingedge research to develop solutions for socially relevant problems.\"}}'),
	(70,4,'PEO[2] Updated','Updated PEO item at index 2','System','2026-01-06 04:31:33','{\"PEO[2]\": {\"new\": \"Demonstrate leadership, teamwork, and ethical responsibility in developing sustainable software and computing solutions.\", \"old\": \"Graduates will communicate effectively and practice their profession with ethics,integrity, leadership, teamwork, and social responsibility, and pursue lifelong learning throughout their careers.\"}}'),
	(71,4,'PSO[0] Updated','Updated PSO item at index 0','System','2026-01-06 04:32:14','{\"PSO[0]\": {\"new\": \"Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.\", \"old\": \"Apply suitable algorithmic thinking and data management practices to design develop, and evaluate effective solutions for real-life and research problems.\"}}'),
	(72,4,'PSO[1] Updated','Updated PSO item at index 1','System','2026-01-06 04:32:14','{\"PSO[1]\": {\"new\": \"Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.\", \"old\": \"Design and develop cost-effective solutions based on cutting-edge hardware and software tools and techniques to meet global requirements.\"}}'),
	(73,10,'Curriculum Created','Created new curriculum: R2026-CSE (2025-2026)','System','2026-01-06 04:47:56',NULL),
	(74,10,'Department Overview Created','Created department vision, mission, PEOs, POs, and PSOs','System','2026-01-06 04:49:17',NULL),
	(75,10,'PSO[1] Added','Added PSO item at index 1','System','2026-01-06 04:49:57','{\"PSO[1]\": {\"new\": \"Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.\", \"old\": \"\"}}'),
	(76,10,'PSO[0] Added','Added PSO item at index 0','System','2026-01-06 04:49:57','{\"PSO[0]\": {\"new\": \"Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.\", \"old\": \"\"}}'),
	(77,10,'PO[2] Added','Added PO item at index 2','System','2026-01-06 04:52:59','{\"PO[2]\": {\"new\": \"Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.\", \"old\": \"\"}}'),
	(78,10,'PO[8] Added','Added PO item at index 8','System','2026-01-06 04:52:59','{\"PO[8]\": {\"new\": \"Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings\", \"old\": \"\"}}'),
	(79,10,'PO[5] Added','Added PO item at index 5','System','2026-01-06 04:52:59','{\"PO[5]\": {\"new\": \"The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice\", \"old\": \"\"}}'),
	(80,10,'PO[1] Added','Added PO item at index 1','System','2026-01-06 04:52:59','{\"PO[1]\": {\"new\": \"Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering sciences.\", \"old\": \"\"}}'),
	(81,10,'PO[10] Added','Added PO item at index 10','System','2026-01-06 04:52:59','{\"PO[10]\": {\"new\": \"Project Management and Finance: Demonstrate knowledge and understanding of the engineering and management principles and apply these to one’s own work, as a member and leader in a team, to manage projects and in multidisciplinary environments.\", \"old\": \"\"}}'),
	(82,10,'PO[9] Added','Added PO item at index 9','System','2026-01-06 04:52:59','{\"PO[9]\": {\"new\": \"Communication: Communicate effectively on complex engineering activities with the engineering community and with society at large, such as, being able to comprehend and write effective reports and design documentation, make effective presentations, and give and receive clear instructions.\", \"old\": \"\"}}'),
	(83,10,'PO[3] Added','Added PO item at index 3','System','2026-01-06 04:52:59','{\"PO[3]\": {\"new\": \"Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.\", \"old\": \"\"}}'),
	(84,10,'PO[4] Added','Added PO item at index 4','System','2026-01-06 04:52:59','{\"PO[4]\": {\"new\": \"Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.\", \"old\": \"\"}}'),
	(85,10,'PO[0] Added','Added PO item at index 0','System','2026-01-06 04:52:59','{\"PO[0]\": {\"new\": \"Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.\", \"old\": \"\"}}'),
	(86,10,'PO[11] Added','Added PO item at index 11','System','2026-01-06 04:52:59','{\"PO[11]\": {\"new\": \"Life-long Learning: Recognize the need for, and have the preparation and ability to engage in independent and life-long learning in the broadest context of technological change.\", \"old\": \"\"}}'),
	(87,10,'PO[7] Added','Added PO item at index 7','System','2026-01-06 04:52:59','{\"PO[7]\": {\"new\": \"Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.\", \"old\": \"\"}}'),
	(88,10,'PO[6] Added','Added PO item at index 6','System','2026-01-06 04:52:59','{\"PO[6]\": {\"new\": \"Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.\", \"old\": \"\"}}'),
	(89,10,'Card Added','Added Semester 1','System','2026-01-06 04:57:25',NULL),
	(90,4,'PEO[0] Updated','Updated PEO item at index 0','System','2026-01-06 05:10:13','{\"PEO[0]\": {\"new\": \"hello\", \"old\": \"Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.\"}}'),
	(91,4,'PEO[0] Updated','Updated PEO item at index 0','System','2026-01-06 05:10:55','{\"PEO[0]\": {\"new\": \"Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.\", \"old\": \"hello\"}}'),
	(92,10,'Course Added','Added course 26MA101 - Linear Algebra and Calculus to Semester 48','System','2026-01-06 05:28:05',NULL),
	(93,10,'Course Added','Added course 26PH102 - Engineering Physics to Semester 48','System','2026-01-06 05:35:35',NULL),
	(94,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 05:35:43','{\"category\": {\"new\": \"ES\", \"old\": \"BS - Basic Sciences\"}, \"cia_marks\": {\"new\": 40, \"old\": 60}, \"see_marks\": {\"new\": 60, \"old\": 40}}'),
	(95,10,'Course Added','Added course 26CH103 - Engineering Chemistry to Semester 48','System','2026-01-06 05:37:34',NULL),
	(96,10,'Course Added','Added course 26GE004 - Digital Computer Electronics to Semester 48','System','2026-01-06 05:39:17',NULL),
	(97,10,'Course Added','Added course 26GE005 - Problem Solving using C to Semester 48','System','2026-01-06 06:03:03',NULL),
	(98,10,'Course Added','Added course 26HS001 - Communicative English to Semester 48','System','2026-01-06 06:08:14',NULL),
	(99,10,'Course Added','Added course 26HS002 - தமிழர் மரபு / Heritage of Tamils  to Semester 48','System','2026-01-06 06:09:10',NULL),
	(100,10,'Course Added','Added course 26PH108 - Physical Science Laboratory to Semester 48','System','2026-01-06 06:11:53',NULL),
	(101,10,'Course Added','Added course 26GE006 - C Programming Laboratory to Semester 48','System','2026-01-06 06:12:57',NULL),
	(102,10,'Course Added','Added course 26SD001 - Skill Development Course  I to Semester 48','System','2026-01-06 06:17:04',NULL),
	(103,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:20:19','{\"category\": {\"new\": \"BS\", \"old\": \"ES\"}}'),
	(104,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:20:27','{\"category\": {\"new\": \"PC\", \"old\": \"BS\"}}'),
	(105,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:20:39','{\"category\": {\"new\": \"BS\", \"old\": \"PC\"}}'),
	(106,10,'Course Updated','Updated course: 26MA101 - Linear Algebra and Calculus','System','2026-01-06 06:21:07','{\"category\": {\"new\": \"ES\", \"old\": \"BS\"}}'),
	(107,10,'Card Added','Added Semester 2','System','2026-01-06 06:29:27',NULL),
	(108,10,'Course Added','Added course 26MA201 - Differential Equations and Transforms to Semester 49','System','2026-01-06 06:31:04',NULL),
	(109,10,'Course Added','Added course 26PH202 - Materials Science to Semester 49','System','2026-01-06 06:33:17',NULL),
	(110,10,'Course Added','Added course 26CS203 - Fundamentals of Web Principles to Semester 49','System','2026-01-06 06:34:05',NULL),
	(111,10,'Course Added','Added course 26CS204 - Computer Organization and Architecture to Semester 49','System','2026-01-06 06:34:39',NULL),
	(112,10,'Course Added','Added course 26GE007 - Python Programming to Semester 49','System','2026-01-06 06:35:17',NULL),
	(113,10,'Course Added','Added course 26HS005 - Professional Communication to Semester 49','System','2026-01-06 06:36:17',NULL),
	(114,10,'Course Added','Added course 26HS006 - தமிழரும் தொழில்நுட்பமும் / Tamils and Technology to Semester 49','System','2026-01-06 06:37:28',NULL),
	(115,10,'Course Added','Added course 26GE008 - Python Programming Laboratory to Semester 49','System','2026-01-06 06:38:21',NULL),
	(116,10,'Course Added','Added course 26CS209 - Web Principles Laboratory to Semester 49','System','2026-01-06 06:39:13',NULL),
	(117,10,'Course Added','Added course 26SD002 - Skill Development Course II to Semester 49','System','2026-01-06 06:40:28',NULL),
	(118,10,'Card Added','Added Semester 3','System','2026-01-06 06:40:47',NULL),
	(119,10,'Course Added','Added course 26CS301 - Discrete Mathematics to Semester 50','System','2026-01-06 06:42:11',NULL),
	(120,10,'Course Added','Added course 26CS302 - Data Structures and Algorithms to Semester 50','System','2026-01-06 06:43:30',NULL),
	(121,10,'Course Added','Added course 26CS303 - Operating Systems to Semester 50','System','2026-01-06 06:44:19',NULL),
	(122,10,'Course Added','Added course 26CS304 - Object Oriented Programming with Java to Semester 50','System','2026-01-06 06:45:12',NULL),
	(123,10,'Course Updated','Updated course: 26CS304 - Object Oriented Programming with Java','System','2026-01-06 06:46:09','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"cia_marks\": {\"new\": 40, \"old\": 30}, \"see_marks\": {\"new\": 60, \"old\": 70}, \"practical_hours\": {\"new\": 2, \"old\": 0}}'),
	(124,10,'Course Updated','Updated course: 26CS304 - Object Oriented Programming with Java','System','2026-01-06 06:46:38','{\"cia_marks\": {\"new\": 30, \"old\": 40}, \"see_marks\": {\"new\": 70, \"old\": 60}}'),
	(125,10,'Course Added','Added course 26CS305 - Software Engineering to Semester 50','System','2026-01-06 06:47:33',NULL),
	(126,10,'Course Added','Added course 26CS306 - Database Management Systems to Semester 50','System','2026-01-06 06:48:12',NULL),
	(127,10,'Course Added','Added course 26CS307 - Standards in Computer Science to Semester 50','System','2026-01-06 06:48:52',NULL),
	(128,10,'Course Added','Added course 26CS308 - Data Structures and Algorithms Laboratory to Semester 50','System','2026-01-06 06:49:39',NULL),
	(129,10,'Course Added','Added course 26CS309 - Database Management Systems Laboratory to Semester 50','System','2026-01-06 06:50:43',NULL),
	(130,10,'Course Added','Added course 26CS310 - Design Thinking and Innovation Laboratory (AICTE, & NEP) to Semester 50','System','2026-01-06 06:52:07',NULL),
	(131,10,'Card Added','Added Semester 4','System','2026-01-06 09:03:44',NULL),
	(132,10,'Course Added','Added course 26CS401 - Probability and Statistics to Semester 51','System','2026-01-06 09:04:53',NULL),
	(133,10,'Course Added','Added course 26CS402 - Full Stack Development to Semester 51','System','2026-01-06 09:05:42',NULL),
	(134,10,'Course Added','Added course 26CS403 - Artificial Intelligence Essentials to Semester 51','System','2026-01-06 09:06:27',NULL),
	(135,10,'Course Added','Added course 26CS404 - Design and Analysis of Algorithms to Semester 51','System','2026-01-06 09:07:04',NULL),
	(136,10,'Course Updated','Updated course: 26CS403 - Artificial Intelligence Essentials','System','2026-01-06 09:07:27','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"see_marks\": {\"new\": 55, \"old\": 54}}'),
	(137,10,'Course Added','Added course 26CS405 - Theory of Computation to Semester 51','System','2026-01-06 09:08:33',NULL),
	(138,10,'Course Added','Added course 26CS406 - Computer Networks to Semester 51','System','2026-01-06 09:09:37',NULL),
	(139,10,'Course Added','Added course 26HS009 - Environmental Sciences and Sustainability  to Semester 51','System','2026-01-06 09:10:29',NULL),
	(140,10,'Course Added','Added course 26CS408 - Full Stack Development Laboratory to Semester 51','System','2026-01-06 09:11:13',NULL),
	(141,10,'Course Added','Added course 26CS409 - Computer Networks Laboratory to Semester 51','System','2026-01-06 09:12:17',NULL),
	(142,10,'Course Added','Added course 26CS410 - Community Engagement Project to Semester 51','System','2026-01-06 09:13:04',NULL),
	(143,10,'Course Updated','Updated course: 26CS409 - Computer Networks Laboratory','System','2026-01-06 09:13:57','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"course_type\": {\"new\": \"Lab\", \"old\": \"Experiment\"}, \"lecture_hours\": {\"new\": 0, \"old\": 2}, \"practical_hours\": {\"new\": 2, \"old\": 0}}'),
	(144,10,'Card Added','Added Semester 5','System','2026-01-06 09:14:38',NULL),
	(145,10,'Course Added','Added course 26CS501 - Compiler Design to Semester 52','System','2026-01-06 09:15:30',NULL),
	(146,10,'Course Added','Added course 26CS502 - Cloud Infrastructure Services to Semester 52','System','2026-01-06 09:16:38',NULL),
	(147,10,'Course Added','Added course 26CS503 - Bigdata Analytics to Semester 52','System','2026-01-06 09:17:35',NULL),
	(148,10,'Course Added','Added course 26CS504 - Machine Learning Essentials to Semester 52','System','2026-01-06 09:18:18',NULL),
	(149,10,'Course Added','Added course 26XX - Professional Elective I to Semester 52','System','2026-01-06 09:20:24',NULL),
	(150,10,'Course Added','Added course 26XX - Open Elective I to Semester 52','System','2026-01-06 09:21:09',NULL),
	(151,10,'Course Added','Added course 26CS507 - Cloud Infrastructure Services Laboratory to Semester 52','System','2026-01-06 09:21:55',NULL),
	(152,10,'Course Added','Added course 26CS508 - Machine Learning Essentials Laboratory to Semester 52','System','2026-01-06 09:22:33',NULL),
	(153,10,'Course Added','Added course 26CS509 - Technology Integration Project to Semester 52','System','2026-01-06 09:23:11',NULL),
	(154,10,'Card Added','Added Semester 6','System','2026-01-06 09:23:53',NULL),
	(155,10,'Course Added','Added course 26CS601 - Software Project Management and Quality Assurance to Semester 53','System','2026-01-06 09:24:39',NULL),
	(156,10,'Course Added','Added course 26CS602 - Deep Learning to Semester 53','System','2026-01-06 09:25:17',NULL),
	(157,10,'Course Added','Added course 26CS603 - Cryptography and Cyber Security to Semester 53','System','2026-01-06 09:25:58',NULL),
	(158,10,'Course Added','Added course 26CS607 - Software Project Management and Quality Assurance Laboratory to Semester 53','System','2026-01-06 09:26:58',NULL),
	(159,10,'Course Added','Added course 26CS608 - Deep Learning Laboratory  to Semester 53','System','2026-01-06 09:27:32',NULL),
	(160,10,'Course Added','Added course 26CS609 - Innovation and Product Development Project / Industry Oriented Course / Summer Internship to Semester 53','System','2026-01-06 09:28:14',NULL),
	(161,10,'Course Added','Added course 26XX - Professional Elective II to Semester 53','System','2026-01-06 09:29:07',NULL),
	(162,10,'Course Added','Added course 26XX - Professional Elective III to Semester 53','System','2026-01-06 09:29:43',NULL),
	(163,10,'Course Added','Added course 26XX - Open Elective II to Semester 53','System','2026-01-06 09:30:24',NULL),
	(164,10,'Course Updated','Updated course: 26XX - Open Elective II','System','2026-01-06 09:31:05','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}, \"course_name\": {\"new\": \"Open Elective II\", \"old\": \"Professional Elective I\"}}'),
	(165,10,'Course Updated','Updated course: 26XX - Professional Elective II','System','2026-01-06 09:32:01','{\"course_name\": {\"new\": \"Professional Elective II\", \"old\": \"Open Elective II\"}}'),
	(166,10,'Card Added','Added Semester 7','System','2026-01-06 09:32:40',NULL),
	(167,10,'Course Added','Added course 26CS701 - Generative AI and Large Language Models to Semester 54','System','2026-01-06 09:33:11',NULL),
	(168,10,'Course Added','Added course 26CS702 - IoT and Edge Computing to Semester 54','System','2026-01-06 09:33:49',NULL),
	(169,10,'Course Added','Added course 26XXIV - Professional Elective IV to Semester 54','System','2026-01-06 09:35:02',NULL),
	(170,10,'Course Added','Added course 26XXV - Professional Elective V to Semester 54','System','2026-01-06 09:35:40',NULL),
	(171,10,'Course Updated','Updated course: 26XXIV - Professional Elective IV','System','2026-01-06 09:36:12','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}}'),
	(172,10,'Course Updated','Updated course: 26XXV - Professional Elective V','System','2026-01-06 09:36:45','{\"category\": {\"new\": \"ES\", \"old\": \"ES - Engineering Sciences\"}}'),
	(173,10,'Course Added','Added course 26XXVI - Professional Elective VI to Semester 54','System','2026-01-06 09:37:35',NULL),
	(174,10,'Course Added','Added course 26CS706 - Generative AI and Large Language Models Laboratory to Semester 54','System','2026-01-06 09:38:35',NULL),
	(175,10,'Course Added','Added course 26CS707 - Capstone Project work Level I / Internship Pro to Semester 54','System','2026-01-06 09:39:41',NULL),
	(176,10,'Course Updated','Updated course: 26XXIII - Professional Elective III','System','2026-01-06 09:40:13','{\"course_code\": {\"new\": \"26XXIII\", \"old\": \"26XX\"}, \"course_name\": {\"new\": \"Professional Elective III\", \"old\": \"Professional Elective II\"}}'),
	(177,10,'Card Added','Added Semester 8','System','2026-01-06 09:41:03',NULL),
	(178,10,'Course Added','Added course 26CS801 - Capstone Project Work Level II / Internship Project / Startup Product to Semester 55','System','2026-01-06 09:42:07',NULL),
	(179,10,'Curriculum Updated','Updated curriculum details','System','2026-01-06 09:47:58','{\"academic_year\": {\"new\": \"\", \"old\": \"2025-2026\"}}'),
	(180,10,'Course Updated','Updated course: 26XXIV - Professional Elective IV','System','2026-01-06 09:50:30','{\"course_code\": {\"new\": \"26XXIV\", \"old\": \"26XXIII\"}, \"course_name\": {\"new\": \"Professional Elective IV\", \"old\": \"Professional Elective III\"}}'),
	(202,7,'PEO[0] Added','Added PEO item at index 0','System','2026-01-13 08:47:52','{\"PEO[0]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(203,7,'PO[0] Added','Added PO item at index 0','System','2026-01-13 08:47:52','{\"PO[0]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(204,4,'PO[12] Added','Added PO item at index 12','System','2026-01-13 08:49:03','{\"PO[12]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(205,4,'PEO[3] Added','Added PEO item at index 3','System','2026-01-13 08:49:03','{\"PEO[3]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(206,7,'Mission[3] Added','Added Mission item at index 3','System','2026-01-13 08:50:10','{\"Mission[3]\": {\"new\": \"CHECK 1\", \"old\": \"\"}}'),
	(207,7,'Vision Updated','Updated department vision','System','2026-01-13 08:50:10','{\"vision\": {\"new\": \"ccdcacdccasdcsdcCHECK 1\", \"old\": \"ccdcacdccasdcsdc\"}}'),
	(208,7,'PSO[2] Added','Added PSO item at index 2','System','2026-01-13 08:50:10','{\"PSO[2]\": {\"new\": \"CHECK 1\", \"old\": \"\"}}'),
	(209,7,'PEO[0] Added','Added PEO item at index 0','System','2026-01-13 08:54:06','{\"PEO[0]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(210,7,'PSO[2] Added','Added PSO item at index 2','System','2026-01-13 08:54:06','{\"PSO[2]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(211,7,'PO[0] Added','Added PO item at index 0','System','2026-01-13 08:54:06','{\"PO[0]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(212,7,'Mission[3] Added','Added Mission item at index 3','System','2026-01-13 08:54:06','{\"Mission[3]\": {\"new\": \"check 1\", \"old\": \"\"}}'),
	(213,10,'Curriculum Updated','Updated curriculum details','System','2026-01-13 08:58:03','{\"academic_year\": {\"new\": \"2024-2025\", \"old\": \"\"}}'),
	(214,7,'Curriculum Updated','Updated curriculum details','System','2026-01-13 08:58:30','{\"name\": {\"new\": \"ME - MECHANICAL\", \"old\": \"ME - MECH\"}}'),
	(215,7,'Curriculum Updated','Updated curriculum details','System','2026-01-13 08:58:38','{\"max_credits\": {\"new\": 145, \"old\": 143}}'),
	(216,7,'PEO-PO Mapping Saved','Updated PEO-PO mappings for the curriculum','System','2026-01-13 08:58:50',NULL),
	(217,10,'Semester Updated','Updated Semester 8 to Semester 7','System','2026-01-13 08:59:19','{\"semester_number\": {\"new\": 7, \"old\": 8}}'),
	(218,10,'Semester Updated','Updated Semester 7 to Semester 8','System','2026-01-13 08:59:26','{\"semester_number\": {\"new\": 8, \"old\": 7}}'),
	(219,10,'Semester Updated','Updated Semester 8 to Semester 9','System','2026-01-13 09:04:46','{\"semester_number\": {\"new\": 9, \"old\": 8}}'),
	(220,10,'Semester Updated','Updated Semester 9 to Semester 8','System','2026-01-13 09:05:00','{\"semester_number\": {\"new\": 8, \"old\": 9}}'),
	(221,10,'Semester Updated','Updated Semester 8 to Semester 9','System','2026-01-13 09:12:33','{\"semester_number\": {\"new\": 9, \"old\": 8}}'),
	(222,10,'Semester Updated','Updated Semester 9 to Semester 8','System','2026-01-13 09:12:37','{\"semester_number\": {\"new\": 8, \"old\": 9}}'),
	(223,10,'Card Added','Added Vertical 1','System','2026-01-13 09:12:44',NULL),
	(224,10,'Card Added','Added Vertical 2','System','2026-01-13 09:12:57',NULL),
	(225,10,'Card Added','Added New Card','System','2026-01-13 09:13:32',NULL),
	(226,10,'Card Added','Added New Card','System','2026-01-13 09:13:46',NULL),
	(227,10,'Honour Card Added','Added Honour Card: honour vertical *','System','2026-01-13 09:14:16',NULL),
	(228,10,'Honour Card Added','Added Honour Card: sdknvjnvfdjnvfnj','System','2026-01-13 09:22:50',NULL),
	(229,10,'Honour Card Added','Added Honour Card: Honour Vertical *','System','2026-01-13 09:28:30',NULL),
	(230,10,'Card Added','Added Vertical 2','System','2026-01-13 09:42:10',NULL),
	(231,10,'Honour Card Added','Added Honour Card: Honour vertical *','System','2026-01-13 09:51:43',NULL),
	(232,4,'Honour Card Added','Added Honour Card: Honour Vertical *','System','2026-01-13 09:52:21',NULL),
	(233,14,'Curriculum Created','Created new curriculum: check 1 (2024-2025)','System','2026-01-13 09:56:01',NULL),
	(234,14,'Card Added','Added Semester 1','System','2026-01-13 09:56:11',NULL),
	(235,14,'Honour Card Added','Added Honour Card: Honour card','System','2026-01-13 09:56:19',NULL),
	(236,4,'Course Added','Added course CS101 - cue to Semester 3','System','2026-01-13 11:04:23',NULL),
	(237,4,'Course Added','Added course CS130 - check1 to Semester 3','System','2026-01-13 11:08:19',NULL),
	(238,4,'Course Added','Added course CS230 - check 2 to Semester 3','System','2026-01-13 12:11:59',NULL),
	(239,4,'Course Added','Added course CS303 - check 3 to Semester 3','System','2026-01-13 12:24:32',NULL);

/*!40000 ALTER TABLE `curriculum_logs` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum_mission
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum_mission`;

CREATE TABLE `curriculum_mission` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `mission_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `source_department_id` int DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_mission_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum_mission` WRITE;
/*!40000 ALTER TABLE `curriculum_mission` DISABLE KEYS */;

INSERT INTO `curriculum_mission` (`id`, `curriculum_id`, `mission_text`, `position`, `visibility`, `source_curriculum_id`, `source_department_id`)
VALUES
	(2,2,'To impart need based education EXTRA',0,'CLUSTER',NULL,NULL),
	(3,2,'To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.',1,'CLUSTER',NULL,NULL),
	(4,2,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',2,'CLUSTER',NULL,NULL),
	(8,4,'hello',0,'UNIQUE',NULL,NULL),
	(9,5,'hello hi',0,'UNIQUE',NULL,NULL),
	(10,5,'hello',1,'CLUSTER',NULL,NULL),
	(11,5,'hi',2,'UNIQUE',NULL,NULL),
	(12,4,'hello',1,'CLUSTER',5,NULL),
	(27,3,'MODIFIED',0,'CLUSTER',NULL,NULL),
	(28,3,'To impart need based education EXTRA',1,'CLUSTER',2,NULL),
	(29,6,'ELIMINATE',0,'UNIQUE',NULL,NULL),
	(30,3,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',2,'CLUSTER',2,NULL),
	(31,6,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',1,'CLUSTER',2,NULL),
	(34,8,'To impart need based education to meet the requirements of the industry and society.',0,'UNIQUE',NULL,NULL),
	(35,8,'To equip students for emerging technologies with global standards and ethics that aid insocietal sustainability.',1,'UNIQUE',NULL,NULL),
	(36,8,'To build technologically competent individuals for industry and entrepreneurialventures by providing infrastructure and human resources.',2,'UNIQUE',NULL,NULL),
	(37,7,'MODIFIED',0,'CLUSTER',3,NULL),
	(38,5,'check 1',3,'UNIQUE',NULL,NULL);

/*!40000 ALTER TABLE `curriculum_mission` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum_peos
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum_peos`;

CREATE TABLE `curriculum_peos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `peo_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `source_department_id` int DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_peos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum_peos` WRITE;
/*!40000 ALTER TABLE `curriculum_peos` DISABLE KEYS */;

INSERT INTO `curriculum_peos` (`id`, `curriculum_id`, `peo_text`, `position`, `visibility`, `source_curriculum_id`, `source_department_id`)
VALUES
	(1,2,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'CLUSTER',NULL,NULL),
	(2,2,'Pursue continuous learning in emerging technologies such as AI, data science, and cybersecurity to remain adaptable professionals.',1,'UNIQUE',NULL,NULL),
	(3,2,'Demonstrate leadership, teamwork, and ethical responsibility in developing sustainable software and computing solutions.',2,'UNIQUE',NULL,NULL),
	(6,4,'hound cwvwvvsf',0,'UNIQUE',NULL,NULL),
	(25,3,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'CLUSTER',2,NULL),
	(26,6,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'CLUSTER',2,NULL),
	(27,7,'Attain a strong grounding in computing fundamentals, algorithms, and system design to solve complex real-world problems.',0,'UNIQUE',NULL,NULL),
	(28,7,'Pursue continuous learning in emerging technologies such as AI, data science, and cybersecurity to remain adaptable professionals.',1,'UNIQUE',NULL,NULL),
	(29,7,'Demonstrate leadership, teamwork, and ethical responsibility in developing sustainable software and computing solutions.',2,'UNIQUE',NULL,NULL),
	(30,8,'Graduates will apply computer science and engineering principles and practices to solvereal- world problems with their technical competence.',0,'UNIQUE',NULL,NULL),
	(31,8,'Graduates will have the domain knowledge to pursue higher education and apply cuttingedge research to develop solutions for socially relevant problems.',1,'UNIQUE',NULL,NULL),
	(32,8,'Graduates will communicate effectively and practice their profession with ethics,integrity, leadership, teamwork, and social responsibility, and pursue lifelong learning throughout their careers.',2,'UNIQUE',NULL,NULL),
	(33,5,'check 1',0,'UNIQUE',NULL,NULL);

/*!40000 ALTER TABLE `curriculum_peos` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum_pos
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum_pos`;

CREATE TABLE `curriculum_pos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `po_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `source_department_id` int DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_pos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum_pos` WRITE;
/*!40000 ALTER TABLE `curriculum_pos` DISABLE KEYS */;

INSERT INTO `curriculum_pos` (`id`, `curriculum_id`, `po_text`, `position`, `visibility`, `source_curriculum_id`, `source_department_id`)
VALUES
	(1,2,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'CLUSTER',NULL,NULL),
	(2,2,'Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering sciences.',1,'UNIQUE',NULL,NULL),
	(3,2,'Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.',2,'UNIQUE',NULL,NULL),
	(4,2,'Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.',3,'UNIQUE',NULL,NULL),
	(5,2,'Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.',4,'UNIQUE',NULL,NULL),
	(6,2,'The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice',5,'UNIQUE',NULL,NULL),
	(7,2,'Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.',6,'UNIQUE',NULL,NULL),
	(8,2,'Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.',7,'UNIQUE',NULL,NULL),
	(9,2,'Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings.',8,'UNIQUE',NULL,NULL),
	(10,2,'Communication: Communicate effectively on complex engineering activities with the engineering community and with society at large, such as, being able to comprehend and write effective reports and design documentation, make effective presentations, and give and receive clear instructions.',9,'UNIQUE',NULL,NULL),
	(11,2,'Project Management and Finance: Demonstrate knowledge and understanding of the engineering and management principles and apply these to one’s own work, as a member and leader in a team, to manage projects and in multidisciplinary environments.',10,'UNIQUE',NULL,NULL),
	(12,2,'Life-long Learning: Recognize the need for, and have the preparation and ability to engage in independent and life-long learning in the broadest context of technological change.',11,'UNIQUE',NULL,NULL),
	(18,4,'sdcvvsvd',0,'UNIQUE',NULL,NULL),
	(19,4,'huii',1,'UNIQUE',NULL,NULL),
	(21,6,'MODIFIED',0,'UNIQUE',NULL,NULL),
	(22,3,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'CLUSTER',2,NULL),
	(23,7,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'UNIQUE',NULL,NULL),
	(24,7,'Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering sciences.',1,'UNIQUE',NULL,NULL),
	(25,7,'Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.',2,'UNIQUE',NULL,NULL),
	(26,7,'Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.',3,'UNIQUE',NULL,NULL),
	(27,7,'Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.',4,'UNIQUE',NULL,NULL),
	(28,7,'The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice',5,'UNIQUE',NULL,NULL),
	(29,7,'Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.',6,'UNIQUE',NULL,NULL),
	(30,7,'Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.',7,'UNIQUE',NULL,NULL),
	(31,7,'Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings',8,'UNIQUE',NULL,NULL),
	(32,7,'Communication: Communicate effectively on complex engineering activities with the engineering community and with society at large, such as, being able to comprehend and write effective reports and design documentation, make effective presentations, and give and receive clear instructions.',9,'UNIQUE',NULL,NULL),
	(33,7,'Project Management and Finance: Demonstrate knowledge and understanding of the engineering and management principles and apply these to one’s own work, as a member and leader in a team, to manage projects and in multidisciplinary environments.',10,'UNIQUE',NULL,NULL),
	(34,7,'Life-long Learning: Recognize the need for, and have the preparation and ability to engage in independent and life-long learning in the broadest context of technological change.',11,'UNIQUE',NULL,NULL),
	(35,8,'Engineering Knowledge: Apply the knowledge of mathematics, science, engineering fundamentals, and an engineering specialization to the solution of complex engineering problems.',0,'UNIQUE',NULL,NULL),
	(36,8,'Problem Analysis: Identify, formulate, review research literature, and analyse complex engineering problems reaching substantiated conclusions using first principles of mathematics, natural sciences, and engineering science',1,'UNIQUE',NULL,NULL),
	(37,8,'Design/ Development of Solutions: Design solutions for complex engineering problems and design system components or processes that meet the specified needs with appropriate consideration for public health and safety, and the cultural, societal, and environmental considerations.',2,'UNIQUE',NULL,NULL),
	(38,8,'Conduct Investigations of Complex Problems: Use research-based knowledge and research methods including design of experiments, analysis and interpretation of data, and synthesis of the information to provide valid conclusions.',3,'UNIQUE',NULL,NULL),
	(39,8,'Modern Tool Usage: Create, select, and apply appropriate techniques, resources, and modern engineering and IT tools including prediction and modeling to complex engineering activities with an understanding of the limitations.',4,'UNIQUE',NULL,NULL),
	(40,8,'The Engineer and Society: Apply reasoning informed by the contextual knowledge to assess societal, health, safety, legal and cultural issues and the consequent responsibilities relevant to the professional engineering practice',5,'UNIQUE',NULL,NULL),
	(41,8,'Environment and Sustainability: Understand the impact of the professional engineering solutions in societal and environmental contexts, and demonstrate the knowledge of, and need for sustainable development.',6,'UNIQUE',NULL,NULL),
	(42,8,'Ethics: Apply ethical principles and commit to professional ethics and responsibilities and norms of the engineering practice.',7,'UNIQUE',NULL,NULL),
	(43,8,'Individual and Team Work: Function effectively as an individual, and as a member or leader in diverse teams, and in multidisciplinary settings.',8,'UNIQUE',NULL,NULL),
	(44,5,'check 1',0,'UNIQUE',NULL,NULL);

/*!40000 ALTER TABLE `curriculum_pos` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum_psos
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum_psos`;

CREATE TABLE `curriculum_psos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `pso_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int NOT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `source_department_id` int DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `department_id` (`curriculum_id`,`position`) USING BTREE,
  CONSTRAINT `curriculum_psos_ibfk_1` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum_vision` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum_psos` WRITE;
/*!40000 ALTER TABLE `curriculum_psos` DISABLE KEYS */;

INSERT INTO `curriculum_psos` (`id`, `curriculum_id`, `pso_text`, `position`, `visibility`, `source_curriculum_id`, `source_department_id`)
VALUES
	(1,2,'Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.',0,'UNIQUE',NULL,NULL),
	(2,2,'Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.',1,'UNIQUE',NULL,NULL),
	(6,4,'s. csdcscsC',0,'CLUSTER',5,NULL),
	(9,5,'s. csdcscsC',0,'UNIQUE',NULL,NULL),
	(10,5,'s. csdcscsC',1,'UNIQUE',NULL,NULL),
	(15,7,'Apply algorithmic and data-driven reasoning to design efficient computing systems and intelligent applications.',0,'UNIQUE',NULL,NULL),
	(16,7,'Develop scalable and secure software using modern programming paradigms, tools, and cloud architectures.',1,'UNIQUE',NULL,NULL),
	(17,8,'Apply suitable algorithmic thinking and data management practices to design develop, and evaluate effective solutions for real-life and research problems.',0,'UNIQUE',NULL,NULL),
	(18,8,'Design and develop cost-effective solutions based on cutting-edge hardware and software tools and techniques to meet global requirements.',1,'UNIQUE',NULL,NULL),
	(19,5,'check 1',2,'UNIQUE',NULL,NULL);

/*!40000 ALTER TABLE `curriculum_psos` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table curriculum_vision
# ------------------------------------------------------------

DROP TABLE IF EXISTS `curriculum_vision`;

CREATE TABLE `curriculum_vision` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `vision` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `curriculum_vision` WRITE;
/*!40000 ALTER TABLE `curriculum_vision` DISABLE KEYS */;

INSERT INTO `curriculum_vision` (`id`, `curriculum_id`, `vision`)
VALUES
	(2,4,'To excel in the field of Computer Science and Engineering'),
	(3,6,''),
	(4,8,''),
	(5,7,'ccdcacdccasdcsdcCHECK 1'),
	(6,9,''),
	(7,10,''),
	(8,11,'To excel in the field of Computer Science and Engineering, to meet the emerging needsof the\nindustry, society, and beyond.'),
	(9,12,'');

/*!40000 ALTER TABLE `curriculum_vision` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table honour_cards
# ------------------------------------------------------------

DROP TABLE IF EXISTS `honour_cards`;

CREATE TABLE `honour_cards` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_regulation` (`curriculum_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `honour_cards` WRITE;
/*!40000 ALTER TABLE `honour_cards` DISABLE KEYS */;

INSERT INTO `honour_cards` (`id`, `curriculum_id`, `title`, `created_at`, `visibility`, `source_curriculum_id`)
VALUES
	(6,4,'Honour Vertical *','2026-01-13 09:52:21','UNIQUE',NULL),
	(7,14,'Honour card','2026-01-13 09:56:19','UNIQUE',NULL);

/*!40000 ALTER TABLE `honour_cards` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table honour_vertical_courses
# ------------------------------------------------------------

DROP TABLE IF EXISTS `honour_vertical_courses`;

CREATE TABLE `honour_vertical_courses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `honour_vertical_id` int NOT NULL,
  `course_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_course_vertical` (`honour_vertical_id`,`course_id`) USING BTREE,
  KEY `course_id` (`course_id`) USING BTREE,
  KEY `idx_vertical` (`honour_vertical_id`) USING BTREE,
  CONSTRAINT `honour_vertical_courses_ibfk_1` FOREIGN KEY (`honour_vertical_id`) REFERENCES `honour_verticals` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `honour_vertical_courses_ibfk_2` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `honour_vertical_courses` WRITE;
/*!40000 ALTER TABLE `honour_vertical_courses` DISABLE KEYS */;

INSERT INTO `honour_vertical_courses` (`id`, `honour_vertical_id`, `course_id`, `created_at`)
VALUES
	(3,3,1,'2026-01-13 09:53:10');

/*!40000 ALTER TABLE `honour_vertical_courses` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table honour_verticals
# ------------------------------------------------------------

DROP TABLE IF EXISTS `honour_verticals`;

CREATE TABLE `honour_verticals` (
  `id` int NOT NULL AUTO_INCREMENT,
  `honour_card_id` int NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_honour_card` (`honour_card_id`) USING BTREE,
  CONSTRAINT `honour_verticals_ibfk_1` FOREIGN KEY (`honour_card_id`) REFERENCES `honour_cards` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `honour_verticals` WRITE;
/*!40000 ALTER TABLE `honour_verticals` DISABLE KEYS */;

INSERT INTO `honour_verticals` (`id`, `honour_card_id`, `name`, `created_at`)
VALUES
	(3,6,'data','2026-01-13 09:52:29');

/*!40000 ALTER TABLE `honour_verticals` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table normal_cards
# ------------------------------------------------------------

DROP TABLE IF EXISTS `normal_cards`;

CREATE TABLE `normal_cards` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `semester_number` int DEFAULT NULL,
  `visibility` enum('UNIQUE','CLUSTER') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'UNIQUE',
  `source_curriculum_id` int DEFAULT NULL,
  `card_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'semester',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_semester_regulation` (`curriculum_id`) USING BTREE,
  CONSTRAINT `fk_semester_regulation` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `normal_cards` WRITE;
/*!40000 ALTER TABLE `normal_cards` DISABLE KEYS */;

INSERT INTO `normal_cards` (`id`, `curriculum_id`, `semester_number`, `visibility`, `source_curriculum_id`, `card_type`)
VALUES
	(3,4,1,'UNIQUE',NULL,'semester'),
	(33,6,1,'UNIQUE',NULL,'semester'),
	(41,4,NULL,'UNIQUE',NULL,'elective'),
	(42,4,1,'UNIQUE',NULL,'vertical'),
	(43,4,2,'UNIQUE',NULL,'vertical'),
	(44,4,2,'UNIQUE',NULL,'semester'),
	(45,6,2,'UNIQUE',NULL,'vertical'),
	(47,4,3,'UNIQUE',NULL,'vertical'),
	(48,10,1,'UNIQUE',NULL,'semester'),
	(49,10,2,'UNIQUE',NULL,'semester'),
	(50,10,3,'UNIQUE',NULL,'semester'),
	(51,10,4,'UNIQUE',NULL,'semester'),
	(52,10,5,'UNIQUE',NULL,'semester'),
	(53,10,6,'UNIQUE',NULL,'semester'),
	(54,10,7,'UNIQUE',NULL,'semester'),
	(55,10,8,'UNIQUE',NULL,'semester'),
	(71,14,1,'UNIQUE',NULL,'semester');

/*!40000 ALTER TABLE `normal_cards` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table peo_po_mapping
# ------------------------------------------------------------

DROP TABLE IF EXISTS `peo_po_mapping`;

CREATE TABLE `peo_po_mapping` (
  `id` int NOT NULL AUTO_INCREMENT,
  `curriculum_id` int NOT NULL,
  `peo_index` int NOT NULL,
  `po_index` int NOT NULL,
  `mapping_value` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_peopo_reg` (`curriculum_id`) USING BTREE,
  CONSTRAINT `fk_peopo_reg` FOREIGN KEY (`curriculum_id`) REFERENCES `curriculum` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=738 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `peo_po_mapping` WRITE;
/*!40000 ALTER TABLE `peo_po_mapping` DISABLE KEYS */;

INSERT INTO `peo_po_mapping` (`id`, `curriculum_id`, `peo_index`, `po_index`, `mapping_value`)
VALUES
	(11,4,1,1,3),
	(12,4,1,2,3),
	(13,4,1,3,3),
	(14,4,1,4,3),
	(15,4,1,5,3),
	(16,4,1,6,3),
	(17,4,1,7,3),
	(18,4,1,11,3),
	(19,4,1,12,3),
	(20,4,2,1,3),
	(21,4,2,2,3),
	(22,4,2,3,3),
	(23,4,2,4,3),
	(24,4,2,5,3),
	(25,4,2,6,3),
	(26,4,2,7,3),
	(27,4,2,10,3),
	(28,4,3,8,3),
	(29,4,3,9,3),
	(30,4,3,10,3),
	(31,4,3,11,3),
	(32,4,3,12,3),
	(693,6,1,1,3),
	(694,6,1,2,3),
	(695,6,1,3,3),
	(696,6,1,4,3),
	(697,6,1,5,3),
	(698,6,1,6,3),
	(699,6,1,7,3),
	(700,6,1,11,3),
	(701,6,1,12,3),
	(702,6,2,1,3),
	(703,6,2,2,3),
	(704,6,2,3,3),
	(705,6,2,4,3),
	(706,6,2,5,3),
	(707,6,2,6,3),
	(708,6,2,7,3),
	(709,6,2,10,3),
	(710,6,3,8,3),
	(711,6,3,9,3),
	(712,6,3,10,3),
	(713,6,3,11,3),
	(714,6,3,12,3),
	(737,7,1,1,2);

/*!40000 ALTER TABLE `peo_po_mapping` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table regulation_clause_history
# ------------------------------------------------------------

DROP TABLE IF EXISTS `regulation_clause_history`;

CREATE TABLE `regulation_clause_history` (
  `id` int NOT NULL AUTO_INCREMENT,
  `clause_id` int NOT NULL,
  `old_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `new_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `changed_by` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `changed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `change_reason` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `clause_id` (`clause_id`) USING BTREE,
  CONSTRAINT `regulation_clause_history_ibfk_1` FOREIGN KEY (`clause_id`) REFERENCES `regulation_clauses` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;



# Dump of table regulation_clauses
# ------------------------------------------------------------

DROP TABLE IF EXISTS `regulation_clauses`;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;



# Dump of table regulation_sections
# ------------------------------------------------------------

DROP TABLE IF EXISTS `regulation_sections`;

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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `regulation_sections` WRITE;
/*!40000 ALTER TABLE `regulation_sections` DISABLE KEYS */;

INSERT INTO `regulation_sections` (`id`, `regulation_id`, `section_no`, `title`, `display_order`, `created_at`, `updated_at`)
VALUES
	(1,1,1,'ADMISSION',1,'2025-12-29 04:27:34','2025-12-29 04:27:34');

/*!40000 ALTER TABLE `regulation_sections` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table regulations
# ------------------------------------------------------------

DROP TABLE IF EXISTS `regulations`;

CREATE TABLE `regulations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` enum('DRAFT','PUBLISHED','LOCKED') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'DRAFT',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `regulations` WRITE;
/*!40000 ALTER TABLE `regulations` DISABLE KEYS */;

INSERT INTO `regulations` (`id`, `code`, `name`, `status`, `created_at`, `updated_at`)
VALUES
	(1,'R2022','Academic Regulation 2022','DRAFT','2025-12-27 10:20:35','2025-12-27 10:20:35');

/*!40000 ALTER TABLE `regulations` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table sharing_tracking
# ------------------------------------------------------------

DROP TABLE IF EXISTS `sharing_tracking`;

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
) ENGINE=InnoDB AUTO_INCREMENT=84 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `sharing_tracking` WRITE;
/*!40000 ALTER TABLE `sharing_tracking` DISABLE KEYS */;

INSERT INTO `sharing_tracking` (`id`, `source_curriculum_id`, `target_curriculum_id`, `item_type`, `source_item_id`, `copied_item_id`, `created_at`)
VALUES
	(31,2,3,'mission',2,28,'2025-12-25 17:27:27'),
	(33,2,3,'mission',4,30,'2025-12-25 18:09:15'),
	(34,2,6,'mission',4,31,'2025-12-25 18:09:16'),
	(64,2,3,'peos',1,25,'2025-12-26 06:42:35'),
	(65,2,6,'peos',1,26,'2025-12-26 06:42:35'),
	(76,2,3,'semester',3,33,'2025-12-26 09:19:28'),
	(77,2,6,'semester',3,34,'2025-12-26 09:19:34'),
	(78,2,3,'pos',1,22,'2025-12-26 09:35:55'),
	(79,2,3,'semester',42,33,'2026-01-05 06:26:00'),
	(80,2,3,'semester',43,45,'2026-01-05 06:34:07'),
	(81,2,6,'semester',43,46,'2026-01-05 06:34:12'),
	(82,2,3,'semester',44,45,'2026-01-05 08:49:05'),
	(83,2,6,'semester',44,46,'2026-01-05 08:49:06');

/*!40000 ALTER TABLE `sharing_tracking` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table syllabus
# ------------------------------------------------------------

DROP TABLE IF EXISTS `syllabus`;

CREATE TABLE `syllabus` (
  `id` int NOT NULL AUTO_INCREMENT,
  `model_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `position` int DEFAULT '0',
  `course_id` int NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `syllabus_models_fk_courses` (`course_id`) USING BTREE,
  CONSTRAINT `syllabus_models_fk_courses` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `syllabus` WRITE;
/*!40000 ALTER TABLE `syllabus` DISABLE KEYS */;

INSERT INTO `syllabus` (`id`, `model_name`, `name`, `position`, `course_id`)
VALUES
	(6,'Module 1','Module 1',0,17),
	(7,'Experiment 2','Experiment 2',1,17),
	(8,'Module 1','Module 1',0,18),
	(11,'Unit 1','Unit 1',0,83),
	(12,'Unit 2','Unit 2',1,83);

/*!40000 ALTER TABLE `syllabus` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table syllabus_titles
# ------------------------------------------------------------

DROP TABLE IF EXISTS `syllabus_titles`;

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
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `syllabus_titles` WRITE;
/*!40000 ALTER TABLE `syllabus_titles` DISABLE KEYS */;

INSERT INTO `syllabus_titles` (`id`, `model_id`, `title_name`, `hours`, `title`, `position`)
VALUES
	(6,6,'Experiment 1',5,'Experiment 1',0),
	(7,8,'lineAR ',5,'lineAR ',0),
	(8,11,'MATHEMATICS MODELING OF LINEAR FUNCTIONS',8,'MATHEMATICS MODELING OF LINEAR FUNCTIONS',0);

/*!40000 ALTER TABLE `syllabus_titles` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table syllabus_topics
# ------------------------------------------------------------

DROP TABLE IF EXISTS `syllabus_topics`;

CREATE TABLE `syllabus_topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title_id` int NOT NULL,
  `topic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `position` int DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `title_id` (`title_id`) USING BTREE,
  CONSTRAINT `syllabus_topics_ibfk_1` FOREIGN KEY (`title_id`) REFERENCES `syllabus_titles` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `syllabus_topics` WRITE;
/*!40000 ALTER TABLE `syllabus_topics` DISABLE KEYS */;

INSERT INTO `syllabus_topics` (`id`, `title_id`, `topic`, `content`, `position`)
VALUES
	(14,6,'Rank of a Matrix','Rank of a Matrix',0),
	(15,7,'HSFFS','HSFFS',0);

/*!40000 ALTER TABLE `syllabus_topics` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `users` (`id`, `username`, `password_hash`, `full_name`, `email`, `role`, `is_active`, `created_at`, `updated_at`, `last_login`)
VALUES
	(1,'admin','$2a$10$H4BOz6nXYrnGQVYwC0eAQul.YF3LyhWTpb7xUf1HAKPT8y18DYDaq','System Administrator','admin@example.com','admin',1,'2026-01-07 05:52:45','2026-01-19 09:15:25','2026-01-19 09:15:25');

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
