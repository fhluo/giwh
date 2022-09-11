package wish

import (
	jsoniter "github.com/json-iterator/go"
	"os"
)

type Importer interface {
	Import(filename string) (Items, error)
}

var importers = map[string]Importer{
	".json": JSONImporter{},
}

func RegisterImporter(ext string, i Importer) {
	importers[ext] = i
}

type Exporter interface {
	Export(items Items, filename string) error
}

var exporters = map[string]Exporter{
	".json": JSONExporter{},
}

func RegisterExporter(ext string, e Exporter) {
	exporters[ext] = e
}

type JSONImporter struct{}

func (i JSONImporter) Import(filename string) (Items, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var items Items
	return items, jsoniter.Unmarshal(data, &items)
}

type JSONExporter struct{}

func (e JSONExporter) Export(items Items, filename string) error {
	data, err := jsoniter.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0666)
}
