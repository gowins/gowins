
---

## ✅ 主键与基础字段命名规范

| 字段 | 命名规范 | 示例 | 说明 |
|------|-----------|------|------|
| 主键 | `<entity>_id` | `user_id`, `order_id` | 避免多表 join 时冲突 |
| 外键 | `<related_entity>_id` | `user_id`, `product_id` | 明确表示关联关系 |
| UUID | `<entity>_uuid` | `user_uuid` | 如果不用自增主键，用于分布式唯一标识 |
| 自增序号 | `<entity>_seq` | `order_seq` | 用于展示用的可读编号 |

---

## ✅ 通用字段（创建、修改、删除等）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `created_at` | `timestamp` / `datetime` | 创建时间（GORM 会自动维护） |
| `updated_at` | `timestamp` / `datetime` | 最后更新时间（GORM 会自动维护） |
| `deleted_at` | `timestamp` / `datetime (nullable)` | 软删除标志（配合 GORM 的 `gorm.Model`） |
| `created_by` | `bigint` 或 `varchar` | 创建者（user_id 或用户名） |
| `updated_by` | `bigint` 或 `varchar` | 更新者 |

---

## ✅ 状态类字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `status` | `tinyint` / `varchar` | 通用状态（如 0-正常，1-禁用） |
| `is_active` | `bool` / `tinyint(1)` | 是否启用 |
| `is_deleted` | `bool` / `tinyint(1)` | 是否被删除（可配合逻辑删除） |
| `state` | `varchar` | 状态机（如订单的 `pending`, `paid`, `shipped`） |

---

## ✅ 时间区间 & 逻辑时间字段

| 字段名 | 示例 | 说明 |
|--------|------|------|
| `start_time` | 活动、计划开始时间 |
| `end_time` | 活动、计划结束时间 |
| `expire_at` | 过期时间，如优惠券 |
| `last_login_at` | 用户最后登录时间 |
| `paid_at` | 付款时间，订单场景 |

---

## ✅ 排序 & 分页 & 计数类字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `sort` / `order_num` | `int` | 排序字段 |
| `page_view` / `view_count` | `int` | 浏览量 |
| `click_count` | `int` | 点击数 |
| `version` | `int` | 乐观锁版本号（用于并发控制） |

---

## ✅ 用户行为/操作记录字段

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `last_login_ip` | `varchar` | 最后登录 IP |
| `login_count` | `int` | 登录次数 |
| `last_failed_login_at` | `datetime` | 最后登录失败时间 |
| `remark` / `notes` | `varchar/text` | 备注信息 |

---

## ✅ 常见实体字段命名规范（举例）

### 用户表（users）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `user_id` | bigint | 主键 |
| `username` | varchar | 用户名 |
| `email` | varchar | 邮箱地址 |
| `phone` | varchar | 手机号 |
| `password_hash` | varchar | 密码哈希 |
| `status` | tinyint | 状态（0 正常，1 禁用） |

### 订单表（orders）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| `order_id` | bigint | 主键 |
| `order_no` | varchar | 订单编号（业务展示用） |
| `user_id` | bigint | 下单用户 |
| `total_amount` | decimal | 总金额 |
| `paid_at` | datetime | 支付时间 |
| `status` | varchar | 订单状态 |

---

## 🔄 命名风格约定总结

| 维度 | 推荐 |
|------|------|
| 命名风格 | 全部使用小写、下划线分隔（snake_case） |
| 表名 | 复数形式（如 `orders`, `users`） |
| 字段名 | 清晰表达语义，带上前缀（如 `user_id`） |
| 布尔字段 | 以 `is_`、`has_` 开头，如 `is_deleted`、`has_bound` |
| 时间字段 | 以 `_at` 结尾，如 `created_at`, `expire_at` |

---

## 🔚 最佳实践推荐

- 实体主键都加实体名前缀：`user_id`, `order_id`
- 聚合字段中，字段名显式表达所属实体：如 `order.user_id`, `product.seller_id`
- 避免 `id`, `name`, `code` 等通用字段名无前缀使用
- 对外暴露 JSON 字段（例如 REST API）时也同步保留命名风格一致性

---
