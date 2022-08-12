package wish

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/BurntSushi/toml"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
	"os"
)

type Importer interface {
	Import(filename string) (History, error)
}

var importers = map[string]Importer{
	".json": JSONImporter{},
	".toml": TOMLImporter{},
}

func RegisterImporter(ext string, i Importer) {
	importers[ext] = i
}

type Exporter interface {
	Export(h History, filename string) error
}

var exporters = map[string]Exporter{
	".json": JSONExporter{},
	".toml": TOMLExporter{},
	".csv":  CSVExporter{},
	".xlsx": XLSXExporter{},
}

func RegisterExporter(ext string, e Exporter) {
	exporters[ext] = e
}

type JSONImporter struct{}

func (i JSONImporter) Import(filename string) (History, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var items History
	return items, jsoniter.Unmarshal(data, &items)
}

type JSONExporter struct{}

func (e JSONExporter) Export(h History, filename string) error {
	data, err := jsoniter.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}

type TOMLImporter struct{}

func (i TOMLImporter) Import(filename string) (History, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result struct {
		List History `toml:"list"`
	}
	return result.List, toml.Unmarshal(data, &result)
}

type TOMLExporter struct{}

func (e TOMLExporter) Export(h History, filename string) error {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Indent = ""

	err := encoder.Encode(map[string]interface{}{"list": h})
	if err != nil {
		return err
	}

	return os.WriteFile(filename, buf.Bytes(), 0666)
}

type CSVExporter struct{}

func (e CSVExporter) Export(h History, filename string) error {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	if err := w.Write((&RawItem{}).ToCSVHeader()); err != nil {
		return err
	}

	if err := w.WriteAll(h.ToCSVRecords()); err != nil {
		return err
	}

	return os.WriteFile(filename, buf.Bytes(), 0666)
}

type XLSXExporter struct{}

func (e XLSXExporter) Export(h History, filename string) error {
	var err error

	f := excelize.NewFile()

	f.SetSheetName("Sheet1", SharedTypes[0].GetSharedWishName())

	for _, wish := range SharedTypes[1:] {
		f.NewSheet(wish.GetSharedWishName())
	}

	for _, wish := range SharedTypes {
		name := wish.GetSharedWishName()
		header := (&RawItem{}).ToCSVHeader()
		for i := range header {
			if err = f.SetCellValue(name, fmt.Sprintf("%c%d", 'A'+i, 1), header[i]); err != nil {
				return err
			}
		}

		var records [][]string
		if wish == CharacterEventWish {
			records = h.FilterByWishType(wish, CharacterEventWish2).ToCSVRecords()
		} else {
			records = h.FilterByWishType(wish).ToCSVRecords()
		}

		for i, record := range records {
			for j, value := range record {
				if err = f.SetCellValue(name, fmt.Sprintf("%c%d", 'A'+j, 2+i), value); err != nil {
					return err
				}
			}
		}
	}

	return f.SaveAs(filename)
}
