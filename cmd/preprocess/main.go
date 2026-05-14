package main

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

const VDIM = 14

type refEntry struct {
	Vector [VDIM]float32 `json:"vector"`
	Label  string        `json:"label"`
}

func main() {
	inPath  := flag.String("in",  "resources/references.json.gz", "")
	outPath := flag.String("out", "resources/references.bin",     "")
	flag.Parse()

	f, err := os.Open(*inPath)
	if err != nil { log.Fatal(err) }
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil { log.Fatal(err) }
	defer gr.Close()

	out, err := os.Create(*outPath)
	if err != nil { log.Fatal(err) }
	defer out.Close()
	bw := bufio.NewWriterSize(out, 8*1024*1024)

	dec := json.NewDecoder(gr)
	dec.Token() // consome '['

	count := 0
	var entry refEntry
	for dec.More() {
		if err := dec.Decode(&entry); err != nil {
			if err == io.EOF { break }
			log.Fatalf("entry %d: %v", count, err)
		}

		// Partition key: bit2=is_online(idx9), bit1=card_present(idx10), bit0=unknown_merchant(idx11)
		var key uint8
		if entry.Vector[9]  > 0.5 { key |= 0x4 }
		if entry.Vector[10] > 0.5 { key |= 0x2 }
		if entry.Vector[11] > 0.5 { key |= 0x1 }

		var label uint8
		if strings.EqualFold(entry.Label, "fraud") { label = 1 }

		// 1 byte key + 56 bytes float32×14 + 1 byte label = 58 bytes/registro
		binary.Write(bw, binary.LittleEndian, key)
		binary.Write(bw, binary.LittleEndian, entry.Vector)
		binary.Write(bw, binary.LittleEndian, label)
		count++
	}

	bw.Flush()
	log.Printf("Concluído: %d vetores → %s", count, *outPath)
}