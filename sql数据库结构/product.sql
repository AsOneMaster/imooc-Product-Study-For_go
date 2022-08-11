/*
 Navicat Premium Data Transfer

 Source Server         : 本机数据库
 Source Server Type    : MySQL
 Source Server Version : 50643
 Source Host           : 127.0.0.1:3306
 Source Schema         : imooc

 Target Server Type    : MySQL
 Target Server Version : 50643
 File Encoding         : 65001

 Date: 11/07/2019 18:43:43
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for product
-- ----------------------------
DROP TABLE IF EXISTS `product`;
CREATE TABLE `product` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `productName` varchar(255) DEFAULT NULL,
  `productNum` int(11) DEFAULT NULL,
  `productImage` varchar(255) DEFAULT NULL,
  `productUrl` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
