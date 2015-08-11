package dmarc

import (
	"archive/zip"
	"compress/gzip"
	"encoding/xml"
	"io"
)

func Parse(r io.Reader) Feedback {
	data := Feedback{}
	xml.NewDecoder(r).Decode(&data)
	return data
}

func GzipParse(r io.Reader) Feedback {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}

	return Parse(gzr)
}

/* zip.Reader doesn't implement io.Reader :( */
func ZipParse(path string) Feedback {
	zr, err := zip.OpenReader(path)
	if err != nil {
		panic(err)
	}
	defer zr.Close()

	// Assuming each zip archive contains only one file
	fh, err := zr.File[0].Open()
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	return Parse(fh)
}
