package worldlogo

import (
	"encoding/csv"
	"os"
)

func WriteToSCV(name string, items []WorldLogo) error {
	// write to scv
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	csvwriter := csv.NewWriter(f)

	for _, item := range items {
		if err = csvwriter.Write([]string{item.Name, item.Key, item.Src}); err != nil {
			return err
		}
	}
	csvwriter.Flush()
	return nil
}
