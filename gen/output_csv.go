package gen

import (
	"encoding/csv"
	"os"
)

func GenOutputCsv[T any](path string, headers []string, input []T, fn func(T) []string) error {
	var data [][]string
	data = append(data, headers)
	for _, v := range input {
		data = append(data, fn(v))
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	return nil
}
