CREATE TABLE `testtable`
(
    `a`  INT(11)     NOT NULL,
    `b`  VARCHAR(64) NOT NULL,
    `c`  VARCHAR(2)  NOT NULL,
    `d`  VARCHAR(32),
    PRIMARY KEY (`id`)
);

SELECT * FROM testtable;
SELECT a FROM testtable;
