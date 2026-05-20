-- 阅读历史记录
CREATE TABLE read_history
(
    id           INT PRIMARY KEY AUTO_INCREMENT,
    user_id      INT      NOT NULL,
    book_id      INT      NOT NULL,
    day_str      CHAR(10) NOT NULL,
    reading_cost INT      NOT NULL,
    start_page   INT      NOT NULL,
    end_page     INT      NOT NULL,
    create_time  BIGINT   NOT NULL,
    update_time  BIGINT   NOT NULL,
    UNIQUE uk_user_id_book_id_day_str (user_id, book_id, day_str)
);
