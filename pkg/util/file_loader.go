package util

import (
	"bufio"
	"io"
	"os"
)

func LoadLineString(n string) ([]string, error) {
	f, err := os.Open(n)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	var s []string
	for {
		b, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		s = append(s, string(b))
	}
	return s, nil
}
