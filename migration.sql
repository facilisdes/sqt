CREATE USER 'sqt_client'@'%' IDENTIFIED BY 'P@ssw0rd';
CREATE DATABASE sqt_test;
GRANT INSERT, SELECT ON sqt_test.* TO 'sqt_client'@'%';
USE sqt_test;
CREATE TABLE `test` (
                           `ID` int(11) unsigned NOT NULL,
                           `val` varchar(22) NOT NULL DEFAULT '',
                           PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
INSERT INTO `test` (`ID`, `val`) VALUES (1, "test");
***************************************
CREATE USER 'sqt_admin_1234'@'%' IDENTIFIED BY 'P@ssw0rd-12POss@*';
CREATE DATABASE sqt;
GRANT INSERT, SELECT ON sqt.* TO 'sqt_admin_1234'@'%';
GRANT INSERT, SELECT ON sqt_test.* TO 'sqt_admin_1234'@'%';
USE sqt;
CREATE TABLE `clients` (
                           `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                           `address` varchar(22) NOT NULL DEFAULT '',
                           `note` varchar(200) DEFAULT NULL,
                           `password` varchar(32) DEFAULT NULL,
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `address` (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `events` (
                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                          `IsExecuted` tinyint(1) DEFAULT NULL,
                          `Status` int(11) DEFAULT NULL,
                          `StatusText` varchar(200) DEFAULT NULL,
                          `Data` varchar(200) DEFAULT NULL,
                          `TimeElapsed` int(11) DEFAULT NULL,
                          `TimeQueuedMin` int(11) DEFAULT NULL,
                          `LocalData` varchar(11) DEFAULT NULL,
                          `TimeElapsedTotal` int(11) DEFAULT NULL,
                          `QueueSize` int(11) DEFAULT NULL,
                          `Command` int(11) DEFAULT NULL,
                          `RequestedKey` varchar(200) DEFAULT NULL,
                          `ValueIsValidated` tinyint(1) DEFAULT NULL,
                          `Client` varchar(22) DEFAULT NULL,
                          `TimeStart` int(11) DEFAULT NULL,
                          `TimeEnd` int(11) DEFAULT NULL,
                          PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;