# GIWH

原神祈愿记录导出工具

[English](readme.md) | 简体中文

## 使用

使用前请确保打开过游戏中的祈愿历史记录页面

### 获取祈愿记录

```
>>> giwh 
```

- 该命令会将祈愿记录导出至 {uid}.json，其中 {uid} 是你在游戏中的 uid
- 如果 {uid}.json 已存在，则祈愿记录将被合并

```
>>> giwh xxx.json
```

- 指定文件名，祈愿记录将被保存至 .json
- 如果 xxx.json 已存在，则祈愿记录将被合并

```
>>> giwh input.json output.json
```

- 同时指定输入输出文件名，input.json 中的祈愿记录将与获取的新祈愿记录合并后输出至 output.json

### 合并祈愿记录

```
>>> giwh merge 1.json 2.json output.json
```

1.json 和 2.json 中的祈愿记录将被合并至 output.json