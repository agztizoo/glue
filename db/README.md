# db

实现 gorm v2 相关封装.

## 主要能力

[x] 数据库配置支持
  [x] 主/从配置
[x] 数据库路由
  [x] 通过 Context 数据库路由
[x] 事务管理器实现
  [x] 事务闭包
  [x] 事务逃逸
  [x] OnCommitted 回调
[x] 扩展能力
  [x] 全局 Scope: 从 Context 注入检索字段
  [x] 初始化插件: 加/解密支持

## 初始化

支持自定义驱动适配, 驱动装饰. 通过实现 Dialector, 进行自定义数据库适配.

[x] WithInitializeHook 数据库初始化插桩
