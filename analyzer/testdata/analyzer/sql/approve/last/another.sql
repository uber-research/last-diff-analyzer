CREATE TABLE `testtable`
(
    `a`  INT(11)     NOT NULL,
    -- test
    `b`  VARCHAR(64) NOT NULL,
    `c`  VARCHAR(2)  NOT NULL,
    /* test */
    `d`  VARCHAR(32),
    PRIMARY KEY (`id`)
);
-- comment
SELECT * FROM testtable;
