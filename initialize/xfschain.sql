/*
Navicat MySQL Data Transfer

Source Server         : eth
Source Server Version : 50737
Source Host           : localhost:3306
Source Database       : xfschain

Target Server Type    : MYSQL
Target Server Version : 50737
File Encoding         : 65001

Date: 2022-03-05 18:24:31
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for chain_address
-- ----------------------------
DROP TABLE IF EXISTS `chain_address`;
CREATE TABLE `chain_address` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `address` varchar(34) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `balance` decimal(65,0) unsigned NOT NULL DEFAULT '0',
  `nonce` bigint(20) unsigned NOT NULL,
  `extra` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `code` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `state_root` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `alias` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` tinyint(4) unsigned zerofill NOT NULL DEFAULT '0000',
  `display` tinyint(3) unsigned NOT NULL DEFAULT '1',
  `from_state_root` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from_block_height` bigint(20) unsigned NOT NULL DEFAULT '0',
  `from_block_hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `create_from_address` varchar(34) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_from_block_height` bigint(20) unsigned DEFAULT NULL,
  `create_from_block_hash` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_from_state_root` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_from_tx_hash` varchar(66) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `create_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `tx_count` int(10) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `only_address` (`address`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=325 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for chain_block_header
-- ----------------------------
DROP TABLE IF EXISTS `chain_block_header`;
CREATE TABLE `chain_block_header` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `height` bigint(20) unsigned NOT NULL,
  `hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `version` int(10) unsigned NOT NULL,
  `hash_prev_block` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `timestamp` bigint(20) unsigned NOT NULL DEFAULT '0',
  `coinbase` varchar(34) COLLATE utf8mb4_unicode_ci NOT NULL,
  `state_root` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `transactions_root` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `receipts_root` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `gas_limit` bigint(20) NOT NULL,
  `gas_used` bigint(20) NOT NULL,
  `bits` bigint(20) unsigned NOT NULL,
  `nonce` bigint(20) unsigned NOT NULL,
  `extra_nonce` decimal(65,0) unsigned NOT NULL,
  `tx_count` int(10) unsigned NOT NULL DEFAULT '0',
  `rewards` decimal(65,0) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `only_height_hash` (`height`,`hash`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=2279 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for chain_block_tx
-- ----------------------------
DROP TABLE IF EXISTS `chain_block_tx`;
CREATE TABLE `chain_block_tx` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `block_hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `block_height` bigint(20) unsigned NOT NULL,
  `block_time` bigint(20) unsigned NOT NULL DEFAULT '0',
  `version` int(10) unsigned NOT NULL DEFAULT '0',
  `tx_from` varchar(34) COLLATE utf8mb4_unicode_ci NOT NULL,
  `tx_to` varchar(34) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gas_price` decimal(65,0) unsigned NOT NULL DEFAULT '0',
  `gas_limit` decimal(65,0) unsigned NOT NULL DEFAULT '0',
  `gas_used` decimal(65,0) unsigned NOT NULL DEFAULT '0',
  `gas_fee` decimal(65,0) unsigned NOT NULL DEFAULT '0',
  `data` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nonce` bigint(20) unsigned NOT NULL DEFAULT '0',
  `value` decimal(65,0) unsigned NOT NULL,
  `signature` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1',
  `type` tinyint(4) unsigned zerofill NOT NULL DEFAULT '0001',
  `create_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `only_txhash` (`hash`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=52628 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for chain_token
-- ----------------------------
DROP TABLE IF EXISTS `chain_token`;
CREATE TABLE `chain_token` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `symbol` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `total_supply` decimal(65,0) unsigned NOT NULL,
  `decimals` int(10) unsigned NOT NULL,
  `address` varchar(34) COLLATE utf8mb4_unicode_ci NOT NULL,
  `creator` varchar(34) COLLATE utf8mb4_unicode_ci NOT NULL,
  `tx_count` bigint(20) unsigned NOT NULL,
  `from_tx_hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from_block_height` bigint(20) unsigned NOT NULL,
  `from_block_hash` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from_state_root` varchar(66) COLLATE utf8mb4_unicode_ci NOT NULL,
  `create_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;
