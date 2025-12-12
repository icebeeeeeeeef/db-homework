-- Demo seed data for WeBook DB homework
-- Idempotent inserts so you can rerun after docker restart

USE webook_db;

SET @now = UNIX_TIMESTAMP()*1000;

-- users (password: Passw0rd!)
INSERT INTO users (id, email, password, nickname, birthday, about_me, c_time, u_time)
VALUES
    (1, 'alice@example.com', '$2a$10$CVN0.Kikaor.4BHqEfRH.Odj/ZLcdzElL9BKl81sGMEwPb0oU17n6', 'Alice', '1995-05-05', 'DB homework demo user 1', @now, @now),
    (2, 'bob@example.com',   '$2a$10$CVN0.Kikaor.4BHqEfRH.Odj/ZLcdzElL9BKl81sGMEwPb0oU17n6', 'Bob',   '1992-08-20', 'DB homework demo user 2', @now, @now),
    (3, 'carol@example.com', '$2a$10$CVN0.Kikaor.4BHqEfRH.Odj/ZLcdzElL9BKl81sGMEwPb0oU17n6', 'Carol', '1990-12-12', 'DB homework demo user 3', @now, @now)
ON DUPLICATE KEY UPDATE u_time = VALUES(u_time);

-- articles (author_id references users.id)
-- status: 1=draft, 2=withdraw, 3=published
INSERT INTO articles (id, title, content, author_id, created_at, updated_at, status)
VALUES
    (101, 'Go 入门实践',       '<p>面向数据库作业的 Go 入门文章。</p>',            1, @now-86400000*3, @now-3600000, 3),
    (102, 'Redis 使用指南',    '<p>Redis 基础 + 互动计数示例。</p>',              1, @now-86400000*2, @now-7200000, 3),
    (103, '事务与隔离级别',    '<p>讨论事务、隔离级别与异常案例。</p>',           2, @now-86400000*5, @now-86400000, 3),
    (104, '草稿：缓存穿透',    '<p>草稿内容，未发布。</p>',                      2, @now-86400000*1, @now-86400000*1, 1),
    (105, '撤回：消息队列',    '<p>撤回的文章，用于状态演示。</p>',              3, @now-86400000*4, @now-86400000*2, 2),
    (106, '草稿：测试用例设计','<p>草稿文章，练习分页与过滤。</p>',              3, @now-86400000*2, @now-86400000*2, 1)
ON DUPLICATE KEY UPDATE title = VALUES(title), content = VALUES(content), updated_at = VALUES(updated_at), status = VALUES(status);

-- reader_articles mirrors已发布文章，用于公开详情和排行榜
INSERT INTO reader_articles (id, title, content, author_id, created_at, updated_at, status)
VALUES
    (101, 'Go 入门实践',    '<p>面向数据库作业的 Go 入门文章。</p>',   1, @now-86400000*3, @now-3600000, 3),
    (102, 'Redis 使用指南', '<p>Redis 基础 + 互动计数示例。</p>',     1, @now-86400000*2, @now-7200000, 3),
    (103, '事务与隔离级别', '<p>讨论事务、隔离级别与异常案例。</p>',  2, @now-86400000*5, @now-86400000, 3)
ON DUPLICATE KEY UPDATE title = VALUES(title), content = VALUES(content), updated_at = VALUES(updated_at), status = VALUES(status);

-- interactive aggregates
INSERT INTO interactives (biz_id, biz, created_at, updated_at, readcnt, likecnt, collectcnt)
VALUES
    (101, 'article', @now-3600000, @now-3600000, 120, 5, 2),
    (102, 'article', @now-7200000, @now-7200000, 90,  3, 1),
    (103, 'article', @now-86400000, @now-86400000, 45, 1, 0)
ON DUPLICATE KEY UPDATE readcnt = VALUES(readcnt), likecnt = VALUES(likecnt), collectcnt = VALUES(collectcnt), updated_at = VALUES(updated_at);

-- user like records
INSERT INTO user_like_somethings (biz_id, biz, uid, status, created_at, updated_at)
VALUES
    (101, 'article', 2, true, @now-1800000, @now-1800000),
    (101, 'article', 3, true, @now-1200000, @now-1200000),
    (102, 'article', 1, true, @now-3600000, @now-3600000)
ON DUPLICATE KEY UPDATE status = VALUES(status), updated_at = VALUES(updated_at);

-- user collect records (no folders, collect_id=0)
INSERT INTO user_collect_somethings (biz_id, biz, uid, collect_id, created_at, updated_at)
VALUES
    (101, 'article', 2, 0, @now-1000000, @now-1000000),
    (102, 'article', 3, 0, @now-2000000, @now-2000000)
ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at);
