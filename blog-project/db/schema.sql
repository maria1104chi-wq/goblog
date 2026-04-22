-- 个人博客系统数据库结构
-- 字符集：utf8mb4，支持表情符号和中文
-- 存储引擎：InnoDB，支持事务

-- 创建数据库
CREATE DATABASE IF NOT EXISTS blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE blog_db;

-- 用户表（管理员）
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
    `phone` VARCHAR(20) NOT NULL UNIQUE COMMENT '手机号',
    `email` VARCHAR(100) COMMENT '邮箱',
    `avatar` VARCHAR(255) COMMENT '头像URL',
    `role` TINYINT DEFAULT 1 COMMENT '角色：1-超级管理员，2-普通管理员',
    `status` TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `last_login` DATETIME COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) COMMENT '最后登录IP',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_username` (`username`),
    INDEX `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 短信验证码表
CREATE TABLE `sms_codes` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
    `code` VARCHAR(10) NOT NULL COMMENT '验证码',
    `purpose` VARCHAR(50) NOT NULL COMMENT '用途：login-登录，register-注册，reset-重置',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',
    `used` TINYINT DEFAULT 0 COMMENT '是否已使用：0-未使用，1-已使用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `ip` VARCHAR(45) COMMENT '请求IP',
    INDEX `idx_phone_code` (`phone`, `code`),
    INDEX `idx_expires` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短信验证码表';

-- 文章分类表
CREATE TABLE `categories` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `slug` VARCHAR(50) NOT NULL UNIQUE COMMENT '分类别名',
    `description` VARCHAR(255) COMMENT '分类描述',
    `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父分类ID',
    `sort_order` INT DEFAULT 0 COMMENT '排序',
    `article_count` INT DEFAULT 0 COMMENT '文章数量',
    `status` TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_slug` (`slug`),
    INDEX `idx_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章分类表';

-- 文章标签表
CREATE TABLE `tags` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '标签ID',
    `name` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签名称',
    `slug` VARCHAR(50) NOT NULL UNIQUE COMMENT '标签别名',
    `article_count` INT DEFAULT 0 COMMENT '文章数量',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签表';

-- 文章表
CREATE TABLE `articles` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '文章ID',
    `title` VARCHAR(200) NOT NULL COMMENT '文章标题',
    `slug` VARCHAR(200) NOT NULL UNIQUE COMMENT '文章别名（伪静态URL）',
    `summary` VARCHAR(500) COMMENT '文章摘要',
    `content` LONGTEXT NOT NULL COMMENT '文章内容（Markdown格式）',
    `content_html` LONGTEXT COMMENT '文章内容（HTML格式）',
    `cover_image` VARCHAR(255) COMMENT '封面图片URL',
    `author_id` BIGINT UNSIGNED NOT NULL COMMENT '作者ID',
    `category_id` BIGINT UNSIGNED COMMENT '分类ID',
    `view_count` BIGINT DEFAULT 0 COMMENT '浏览次数',
    `like_count` BIGINT DEFAULT 0 COMMENT '点赞数',
    `share_count` BIGINT DEFAULT 0 COMMENT '分享数',
    `comment_count` BIGINT DEFAULT 0 COMMENT '评论数',
    `status` TINYINT DEFAULT 0 COMMENT '状态：0-草稿，1-发布，2-下架',
    `is_top` TINYINT DEFAULT 0 COMMENT '是否置顶：0-否，1-是',
    `seo_title` VARCHAR(200) COMMENT 'SEO标题',
    `seo_keywords` VARCHAR(500) COMMENT 'SEO关键词',
    `seo_description` VARCHAR(500) COMMENT 'SEO描述',
    `published_at` DATETIME COMMENT '发布时间',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME COMMENT '删除时间（软删除）',
    INDEX `idx_slug` (`slug`),
    INDEX `idx_author` (`author_id`),
    INDEX `idx_category` (`category_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_published` (`published_at`),
    INDEX `idx_view_count` (`view_count`),
    INDEX `idx_is_top` (`is_top`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';

-- 文章标签关联表
CREATE TABLE `article_tags` (
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `tag_id` BIGINT UNSIGNED NOT NULL COMMENT '标签ID',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`article_id`, `tag_id`),
    INDEX `idx_tag` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签关联表';

-- 文章多媒体附件表
CREATE TABLE `article_attachments` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '附件ID',
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `file_type` VARCHAR(20) NOT NULL COMMENT '文件类型：image, pdf, video',
    `file_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
    `file_path` VARCHAR(500) NOT NULL COMMENT '文件存储路径',
    `file_url` VARCHAR(500) NOT NULL COMMENT '文件访问URL',
    `file_size` BIGINT COMMENT '文件大小（字节）',
    `mime_type` VARCHAR(100) COMMENT 'MIME类型',
    `width` INT COMMENT '图片宽度',
    `height` INT COMMENT '图片高度',
    `duration` INT COMMENT '视频时长（秒）',
    `sort_order` INT DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_article` (`article_id`),
    INDEX `idx_file_type` (`file_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章附件表';

-- 敏感词库表
CREATE TABLE `sensitive_words` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '敏感词ID',
    `word` VARCHAR(100) NOT NULL UNIQUE COMMENT '敏感词内容',
    `level` TINYINT DEFAULT 1 COMMENT '敏感级别：1-低，2-中，3-高',
    `category` VARCHAR(50) COMMENT '敏感词分类',
    `status` TINYINT DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_word` (`word`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='敏感词库表';

-- 评论表
CREATE TABLE `comments` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '评论ID',
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父评论ID（用于回复）',
    `user_id` BIGINT UNSIGNED COMMENT '用户ID（管理员评论）',
    `nickname` VARCHAR(50) COMMENT '访客昵称（匿名或自定义）',
    `avatar` VARCHAR(255) COMMENT '访客头像',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `content_html` TEXT COMMENT '评论内容（过滤后）',
    `ip_address` VARCHAR(45) COMMENT '评论者IP',
    `ip_location` VARCHAR(100) COMMENT 'IP归属地',
    `user_agent` VARCHAR(500) COMMENT '用户代理',
    `status` TINYINT DEFAULT 0 COMMENT '状态：0-待审核，1-通过，2-拒绝',
    `like_count` INT DEFAULT 0 COMMENT '评论点赞数',
    `is_admin` TINYINT DEFAULT 0 COMMENT '是否管理员评论：0-否，1-是',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME COMMENT '删除时间（软删除）',
    INDEX `idx_article` (`article_id`),
    INDEX `idx_parent` (`parent_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 评论点赞记录表（防止重复点赞）
CREATE TABLE `comment_likes` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `comment_id` BIGINT UNSIGNED NOT NULL COMMENT '评论ID',
    `ip_address` VARCHAR(45) NOT NULL COMMENT '点赞者IP',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uk_comment_ip` (`comment_id`, `ip_address`),
    INDEX `idx_comment` (`comment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论点赞记录表';

-- 文章点赞记录表（防止重复点赞）
CREATE TABLE `article_likes` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `ip_address` VARCHAR(45) NOT NULL COMMENT '点赞者IP',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uk_article_ip` (`article_id`, `ip_address`),
    INDEX `idx_article` (`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章点赞记录表';

-- 文章分享记录表
CREATE TABLE `article_shares` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `article_id` BIGINT UNSIGNED NOT NULL COMMENT '文章ID',
    `platform` VARCHAR(20) COMMENT '分享平台：wechat, weibo, qq, link',
    `ip_address` VARCHAR(45) COMMENT '分享者IP',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_article` (`article_id`),
    INDEX `idx_platform` (`platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章分享记录表';

-- 系统配置表
CREATE TABLE `system_configs` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
    `config_key` VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    `config_value` TEXT COMMENT '配置值',
    `config_type` VARCHAR(20) DEFAULT 'string' COMMENT '配置类型：string, number, boolean, json',
    `group_name` VARCHAR(50) COMMENT '配置分组',
    `description` VARCHAR(255) COMMENT '配置说明',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_group` (`group_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 操作日志表
CREATE TABLE `operation_logs` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
    `user_id` BIGINT UNSIGNED COMMENT '操作用户ID',
    `action` VARCHAR(100) NOT NULL COMMENT '操作行为',
    `module` VARCHAR(50) COMMENT '模块名称',
    `request_method` VARCHAR(10) COMMENT '请求方法',
    `request_url` VARCHAR(500) COMMENT '请求URL',
    `request_params` TEXT COMMENT '请求参数',
    `response_code` INT COMMENT '响应状态码',
    `ip_address` VARCHAR(45) COMMENT '操作IP',
    `user_agent` VARCHAR(500) COMMENT '用户代理',
    `duration` INT COMMENT '执行时长（毫秒）',
    `remark` VARCHAR(500) COMMENT '备注',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_user` (`user_id`),
    INDEX `idx_action` (`action`),
    INDEX `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 统计数据表（每日统计）
CREATE TABLE `daily_statistics` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    `stat_date` DATE NOT NULL UNIQUE COMMENT '统计日期',
    `page_views` BIGINT DEFAULT 0 COMMENT '页面浏览量',
    `unique_visitors` BIGINT DEFAULT 0 COMMENT '独立访客数',
    `new_articles` INT DEFAULT 0 COMMENT '新增文章数',
    `new_comments` INT DEFAULT 0 COMMENT '新增评论数',
    `total_articles` INT DEFAULT 0 COMMENT '文章总数',
    `total_comments` INT DEFAULT 0 COMMENT '评论总数',
    `total_users` INT DEFAULT 0 COMMENT '用户总数',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_date` (`stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='每日统计表';

-- 备份记录表
CREATE TABLE `backup_records` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '备份ID',
    `backup_type` VARCHAR(20) NOT NULL COMMENT '备份类型：database, file, full',
    `file_path` VARCHAR(500) NOT NULL COMMENT '备份文件路径',
    `file_size` BIGINT COMMENT '文件大小',
    `status` TINYINT DEFAULT 0 COMMENT '状态：0-进行中，1-完成，2-失败',
    `remark` VARCHAR(500) COMMENT '备注',
    `started_at` DATETIME COMMENT '开始时间',
    `completed_at` DATETIME COMMENT '完成时间',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_type` (`backup_type`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='备份记录表';

-- 初始化系统配置数据
INSERT INTO `system_configs` (`config_key`, `config_value`, `config_type`, `group_name`, `description`) VALUES
('site_name', '我的个人博客', 'string', 'basic', '网站名称'),
('site_title', '个人博客 - 分享技术与生活', 'string', 'basic', '网站标题'),
('site_keywords', '博客,技术,生活,分享', 'string', 'basic', '网站关键词'),
('site_description', '这是一个个人博客网站，分享技术与生活点滴', 'string', 'basic', '网站描述'),
('site_logo', '/logo.png', 'string', 'basic', '网站Logo'),
('site_footer', '© 2024 我的个人博客. All rights reserved.', 'string', 'basic', '页脚信息'),
('icp_number', '', 'string', 'basic', 'ICP备案号'),
('aliyun_sms_access_key', '', 'string', 'sms', '阿里云短信AccessKey ID'),
('aliyun_sms_access_secret', '', 'string', 'sms', '阿里云短信AccessKey Secret'),
('aliyun_sms_sign_name', '我的博客', 'string', 'sms', '阿里云短信签名'),
('aliyun_sms_template_code', 'SMS_123456789', 'string', 'sms', '阿里云短信模板代码'),
('sms_code_expire_minutes', '5', 'number', 'sms', '短信验证码有效期（分钟）'),
('sms_code_length', '6', 'number', 'sms', '短信验证码长度'),
('article_page_size', '10', 'number', 'article', '文章列表每页数量'),
('hot_article_limit', '10', 'number', 'article', '热门文章显示数量'),
('comment_need_audit', '0', 'boolean', 'comment', '评论是否需要审核：0-否，1-是'),
('allow_anonymous_comment', '1', 'boolean', 'comment', '允许匿名评论：0-否，1-是'),
('max_upload_image_size', '5242880', 'number', 'upload', '图片最大上传大小（字节）'),
('max_upload_pdf_size', '10485760', 'number', 'upload', 'PDF最大上传大小（字节）'),
('max_upload_video_size', '52428800', 'number', 'upload', '视频最大上传大小（字节）'),
('allowed_image_types', 'jpg,jpeg,png,gif,webp', 'string', 'upload', '允许的图片类型'),
('enable_seo_auto', '1', 'boolean', 'seo', '启用自动SEO优化：0-否，1-是');

-- 插入默认管理员账户（密码需要后端加密后更新）
-- 初始密码：admin123（实际使用时请修改）
INSERT INTO `users` (`username`, `password_hash`, `phone`, `email`, `role`, `status`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '13800138000', 'admin@example.com', 1, 1);

-- 插入默认分类
INSERT INTO `categories` (`name`, `slug`, `description`, `sort_order`) VALUES
('技术分享', 'tech', '技术相关文章', 1),
('生活随笔', 'life', '生活感悟与随笔', 2),
('读书笔记', 'reading', '读书心得与笔记', 3),
('旅行见闻', 'travel', '旅行中的见闻', 4);

-- 插入默认标签
INSERT INTO `tags` (`name`, `slug`) VALUES
('Go语言', 'go'),
('Vue.js', 'vue'),
('MySQL', 'mysql'),
('Docker', 'docker'),
('Linux', 'linux'),
('前端开发', 'frontend'),
('后端开发', 'backend');
