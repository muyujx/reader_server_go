-- 用户
CREATE TABLE user
(
    id          INT PRIMARY KEY AUTO_INCREMENT,
    account     CHAR(30) NOT NULL,
    password    CHAR(64) NOT NULL,

    create_time BIGINT   NOT NULL,
    update_time BIGINT   NOT NULL,
    last_login  BIGINT   NOT NULL,


    UNIQUE INDEX uk_account(account)
)