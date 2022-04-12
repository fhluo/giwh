# GIWH

Genshin Impact Wish History Exporter

English | [简体中文](readme_zh-Hans.md)

## Usage

Please make sure you have opened the wish history page in the game before using it.

Wish history is saved in `%LocalAppData%\giwh\wish_history.json`.

### giwh

Show stats.

### giwh update

Update wish history.

### giwh import

Import wish history.

```
giwh import <filename>...
```

### giwh export

Export wish history.

```
giwh export <filename>
```

Specify UID.

```
giwh export <filename> -u <UID>
```

Specify wish types.

```
giwh export <filename> -w 301,400
```

| ID | Wish Type |
| :--: | :--------------------- |
| 100 | Beginners' Wish |
| 200 | Standard Wish |
| 301 | Character Event Wish |
| 302 | Weapon Event Wish |
| 400 | Character Event Wish-2 |

### giwh merge

Merge wish histories.

```
giwh merge <filename>... -o <filename>
```

### giwh stat

Show stats.

```
giwh stat <filename>
```