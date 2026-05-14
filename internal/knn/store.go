package knn

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	VDIM = 14
	K    = 5
)

type Partition struct {
	Vectors []float32 // N*VDIM, layout contíguo cache-friendly
	Labels  []uint8   // 1=fraud, 0=legit
	N       int
}

type Store struct {
	parts [8]Partition
	total int
}

func (s *Store) Total() int          { return s.total }
func (s *Store) Parts() *[8]Partition { return &s.parts }

// Load lê o binário pré-processado (58 bytes por registro)
func Load(path string) (*Store, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	st := &Store{}
	br := bufio.NewReaderSize(f, 8*1024*1024)

	var key uint8
	var vec [VDIM]float32
	var label uint8

	for {
		if err := binary.Read(br, binary.LittleEndian, &key); err != nil {
			if err == io.EOF { break }
			return nil, err
		}
		binary.Read(br, binary.LittleEndian, &vec)
		binary.Read(br, binary.LittleEndian, &label)

		p := &st.parts[key&0x7]
		p.Vectors = append(p.Vectors, vec[:]...)
		p.Labels  = append(p.Labels, label)
		p.N++
		st.total++
	}

	return st, nil
}

// KNN retorna fraud_score = fraudes/5 usando distância L²
// (sem sqrt — monotônica, mesma ordenação)
// Zero alocações: tudo na stack.
func (p *Partition) KNN(q [VDIM]float32) float32 {
	if p.N == 0 { return 0 }

	var topDist  [K]float32
	var topLabel [K]uint8
	for i := range topDist { topDist[i] = 1e38 }
	worst    := float32(1e38)
	worstIdx := 0

	vecs   := p.Vectors
	labels := p.Labels

	for i := 0; i < p.N; i++ {
		b := i * VDIM
		var d float32
		// Unrolled 14 dimensões — o compilador vetoriza com SIMD
		d0 := q[0]-vecs[b+0];   d += d0*d0
		d1 := q[1]-vecs[b+1];   d += d1*d1
		d2 := q[2]-vecs[b+2];   d += d2*d2
		d3 := q[3]-vecs[b+3];   d += d3*d3
		d4 := q[4]-vecs[b+4];   d += d4*d4
		d5 := q[5]-vecs[b+5];   d += d5*d5
		d6 := q[6]-vecs[b+6];   d += d6*d6
		d7 := q[7]-vecs[b+7];   d += d7*d7
		d8 := q[8]-vecs[b+8];   d += d8*d8
		d9 := q[9]-vecs[b+9];   d += d9*d9
		d10 := q[10]-vecs[b+10]; d += d10*d10
		d11 := q[11]-vecs[b+11]; d += d11*d11
		d12 := q[12]-vecs[b+12]; d += d12*d12
		d13 := q[13]-vecs[b+13]; d += d13*d13

		if d < worst {
			topDist[worstIdx]  = d
			topLabel[worstIdx] = labels[i]
			worst = 0
			for j := 0; j < K; j++ {
				if topDist[j] > worst {
					worst    = topDist[j]
					worstIdx = j
				}
			}
		}
	}

	var frauds float32
	for i := 0; i < K; i++ {
		if topLabel[i] == 1 { frauds++ }
	}
	return frauds / K
}