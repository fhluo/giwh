# GIWH

《原神》祈愿记录导出工具

[English](readme.md) | 简体中文

GIWH 是一个帮你导出《原神》祈愿记录的命令行工具。

注：

- 由于查询API的限制，仅能获取最近六个月的祈愿记录
- 新记录的获取可能有1小时左右的延迟。
- 祈愿记录保存在 `%LocalAppData%\giwh\wish_history.json`。

## 使用

1. 首先打开游戏中的祈愿历史记录页面。
2. 使用 `giwh update` 命令更新祈愿记录。
3. 使用 `giwh` 命令查看统计信息或使用 `giwh export` 导出祈愿记录。

### giwh

显示统计信息。

### giwh update

更新祈愿记录。

### giwh import

导入祈愿记录。

```
giwh import <文件名>...
```

### giwh export

导出祈愿记录。

```
giwh export <文件名>
```

指定要导出的UID

```
giwh export <文件名> -u <UID>
```

指定要导出的祈愿类型

```
giwh export <文件名> -w 301,400
```

| ID | 祈愿类型 |
| :--: | :------------- |
| 100 | 新手祈愿 |
| 200 | 常驻祈愿 |
| 301 | 角色活动祈愿 |
| 302 | 武器活动祈愿 |
| 400 | 角色活动祈愿-2 |

### giwh merge

合并祈愿记录。

```
giwh merge <文件名>... -o <文件名>
```

### giwh stat

显示祈愿记录的统计信息。

```
giwh stat <文件名>
```

### giwh version

显示版本号。
