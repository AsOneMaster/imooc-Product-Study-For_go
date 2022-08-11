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

 Date: 11/07/2019 18:43:54
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `nickName` varchar(255) DEFAULT NULL,
  `userName` varchar(255) DEFAULT NULL,
  `passWord` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
