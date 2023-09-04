# GIWH

Genshin Impact Wish History Exporter

[简体中文](readme_zh-Hans.md)

## Features

- Fetch wish history through API and save it to `%LocalAppData%\giwh\wish_history.json`.
- Update wish history.
- Show the progress of getting the next 4-Star and 5-Star.
- Show the number of pulls used to acquire a 5-Star character or weapon.
- Export wish history to json, toml, csv or xlsx format.
- Import wish history in json or toml format.

### Note

- Due to the limitation of the query API, only the last six months of wish history can be fetched.
- There may be a delay of about 1 hour in fetching new records.

## Usage

1. First, open the wish history page in the game.
2. Use `giwh update`  to update wish history.
3. Use `giwh` to view the stats or use `giwh export` to export the wish history.

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

### giwh version

Show version number.
