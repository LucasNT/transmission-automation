package TorrentEntryReader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/LucasNT/transmission-automation/interfaces"
)

type CsvTorrentEntryReader struct {
	csvReader *csv.Reader
}

func NewCsvTorrentEntryReader(csvFile *os.File) CsvTorrentEntryReader {
	var ret CsvTorrentEntryReader = CsvTorrentEntryReader{}
	ret.csvReader = csv.NewReader(csvFile)
	return ret
}

func (t CsvTorrentEntryReader) ReadTorrentEntry() (string, string, error) {
	reader, err := t.csvReader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "", "", interfaces.ErrNoTorrentEntry
		}
		return "", "", fmt.Errorf("Failed to read entry from csv file %w", err)
	}
	return reader[0], reader[1], nil
}
