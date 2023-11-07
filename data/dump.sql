DROP TABLE IF EXISTS `stock_quote`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `stock_quote` (
  `id` int NOT NULL AUTO_INCREMENT,
  `symbol` varchar(4) NOT NULL,
  `price` double DEFAULT NULL,
  `datepoint` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `stock_quote`
--

LOCK TABLES `stock_quote` WRITE;
/*!40000 ALTER TABLE `stock_quote` DISABLE KEYS */;
INSERT INTO `stock_quote` VALUES (2,'UBER',47.75,'2023-11-03 00:00:00'),(3,'UBER',46.48,'2023-11-02 00:00:00'),(4,'UBER',43.83,'2023-11-01 00:00:00'),(5,'UBER',43.28,'2023-10-31 00:00:00'),(6,'UBER',42.73,'2023-10-30 00:00:00'),(7,'UBER',41.23,'2023-10-27 00:00:00'),(8,'UBER',40.62,'2023-10-26 00:00:00'),(9,'UBER',42.35,'2023-10-25 00:00:00'),(10,'UBER',44.19,'2023-10-24 00:00:00'),(11,'UBER',43.04,'2023-10-21 00:00:00'),(12,'UBER',42.96,'2023-10-20 00:00:00'),(13,'UBER',42.72,'2023-10-19 00:00:00'),(14,'UBER',43,'2023-10-18 00:00:00'),(15,'UBER',44.38,'2023-10-17 00:00:00'),(16,'UBER',44.71,'2023-10-16 00:00:00'),(17,'UBER',43.48,'2023-10-13 00:00:00'),(18,'UBER',45.95,'2023-10-12 00:00:00'),(19,'UBER',46.64,'2023-10-11 00:00:00'),(20,'UBER',46.63,'2023-10-10 00:00:00'),(21,'UBER',45.45,'2023-10-09 00:00:00'),(22,'UBER',45.78,'2023-10-06 00:00:00'),(23,'UBER',44.61,'2023-10-05 00:00:00'),(24,'UBER',44.94,'2023-10-04 00:00:00');
/*!40000 ALTER TABLE `stock_quote` ENABLE KEYS */;
UNLOCK TABLES;