package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var (
		out  []string
		fn   string = "data/CP936.txt"
		of   string = "data/CP936_TRANS.txt"
		file *os.File
		outF *os.File
		err  error
		br   *bufio.Reader
		line []byte
		ps   [][]byte
		k    string
		v    string
		i    int
	)

	if file, err = os.Open(fn); err != nil {
		log.Fatal(err)
	}

	if outF, err = os.OpenFile(of, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		log.Fatal(err)
	}

	defer func() {
		file.Close()
		outF.Close()
	}()

	br = bufio.NewReader(file)
	for {
		if line, _, err = br.ReadLine(); err != nil {
			if err != io.EOF {
				log.Fatal(err)
			} else {
				break
			}
		}

		if bytes.Index(line, []byte("#")) == 0 {
			continue
		}

		ps = bytes.Split(line, []byte("\t"))

		if len(ps) < 3 {
			continue
		}

		k = string(ps[1])
		if strings.Trim(k, " ") == "" {
			continue
		}

		v = string(ps[0]) + ":" + string(ps[1]) + ", //" + string(ps[2]) + "\n"
		out = append(out, v)
	}

	for i = 0; i < len(out); i++ {
		outF.Write([]byte(out[i]))
	}
}
