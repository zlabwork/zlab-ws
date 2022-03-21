-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- 主机： 127.0.0.1
-- 生成日期： 2022-03-22 01:50:42
-- 服务器版本： 10.2.7-MariaDB
-- PHP 版本： 8.0.9

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `zlab_ws`
--

-- --------------------------------------------------------

--
-- 表的结构 `im_logs`
--

CREATE TABLE `im_logs` (
  `id` int(10) UNSIGNED NOT NULL,
  `uid` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `uuid` varchar(36) NOT NULL DEFAULT '',
  `os` varchar(32) NOT NULL DEFAULT '',
  `version` varchar(32) NOT NULL DEFAULT '',
  `token` varchar(255) NOT NULL DEFAULT '',
  `ctime` datetime NOT NULL DEFAULT '0001-01-01 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- 表的结构 `im_msg`
--

CREATE TABLE `im_msg` (
  `id` int(10) UNSIGNED NOT NULL,
  `mid` varchar(36) NOT NULL DEFAULT '',
  `type` tinyint(4) UNSIGNED NOT NULL DEFAULT 0,
  `sender` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `receiver` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `data` varchar(2048) NOT NULL DEFAULT '',
  `ctime` datetime NOT NULL DEFAULT '0001-01-01 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- --------------------------------------------------------

--
-- 表的结构 `im_todo`
--

CREATE TABLE `im_todo` (
  `id` int(10) UNSIGNED NOT NULL,
  `mid` varchar(36) NOT NULL DEFAULT '',
  `type` tinyint(4) UNSIGNED NOT NULL DEFAULT 0,
  `sender` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `receiver` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `data` varchar(2048) NOT NULL DEFAULT '',
  `ctime` datetime NOT NULL DEFAULT '0001-01-01 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- 转储表的索引
--

--
-- 表的索引 `im_logs`
--
ALTER TABLE `im_logs`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `im_msg`
--
ALTER TABLE `im_msg`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `im_todo`
--
ALTER TABLE `im_todo`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `im_logs`
--
ALTER TABLE `im_logs`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `im_msg`
--
ALTER TABLE `im_msg`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `im_todo`
--
ALTER TABLE `im_todo`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
