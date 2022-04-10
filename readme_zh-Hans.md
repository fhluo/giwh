# GIWH

原神祈愿记录导出工具

[English](readme.md) | 简体中文

## 使用

使用前请确保打开过游戏中的祈愿历史记录页面。

祈愿记录保存在 `%LocalAppData%\giwh\wish_history.json`。

### `giwh`

更新祈愿记录并显示统计信息。

### `giwh import`

导入祈愿记录。

```
giwh import <文件名>...
```

### `giwh export`

导出祈愿记录。

```
giwh export <文件名>
```

### `giwh merge`

合并祈愿记录。

```
giwh merge <文件名>... -o <文件名>
```

### `giwh stat`

显示祈愿记录的统计信息。

```
giwh stat <文件名>
```