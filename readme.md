# GIWH

Genshin Impact Wish History Exporter

English | [简体中文](readme_zh-Hans.md)

## Usage

Please make sure you have opened the wish history page in the game before using it.

Wish history is saved in `%LocalAppData%\giwh\wish_history.json`.

### `giwh`

Update wish history and show stats.

### `giwh import`

Import wish history.

```
giwh import <filename>...
```

### `giwh export`

Export wish history.

```
giwh export <filename>
```

### `giwh merge`

Merge wish histories.

```
giwh merge <filename>... -o <filename>
```

### `giwh stat`

Show stats for the given wish history.

```
giwh stat <filename>
```