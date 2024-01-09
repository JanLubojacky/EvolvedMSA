package pkg

import (
	"fmt"
	"math"
	// "math/rand"
)

// MSA represents the parsed data from the FASTA file
type MSA struct {
	TargetSequence Sequence
	OtherSequences []Sequence
	MsaLength      int
	Score          float64
}

// initialize the MSA by padding all the sequences with gaps
// inserting them in random positions
func (msa *MSA) Init() {
	msa.Score = math.Inf(-1)
	le := LengthEqualizer{P_mut: 0.2}
	le.Mutate(msa)

	// fmt.Println("Initialising")
}

func (msa *MSA) Print() {
	fmt.Println("Score: ", msa.Score)
	// fmt.Println("Target sequence:")
	// msa.TargetSequence.Print()
	fmt.Println("Sequences:")
	for _, sequence := range msa.OtherSequences {
		sequence.Print()
	}
}

// creates a copy of the MSA
// and returns a pointer to it
func (msa MSA) Copy() MSA {
	// Create a new MSA with the same TargetSequence and MsaLength
	newMSA := MSA{
		TargetSequence: msa.TargetSequence,
		MsaLength:      msa.MsaLength,
		Score:          msa.Score,
		OtherSequences: make([]Sequence, len(msa.OtherSequences)),
	}

	// Deep copy each sequence in OtherSequences
	for i, seq := range msa.OtherSequences {
		newMSA.OtherSequences[i] = seq.Copy()
	}

	return newMSA
}

func (msa *MSA) ResetIterPointers() {
	for i := 0; i < len(msa.OtherSequences); i++ {
		msa.OtherSequences[i].ResetIterPointer()
	}
}

func (msa *MSA) ComputeConsensus() Sequence {
	// iterate over all sequences in the MSA
	// and compute the consensus sequence
	// by taking the most frequent base in each column
	consensus := Sequence{}
	msa_len := msa.OtherSequences[0].Length

	msa.ResetIterPointers()

	for i := 0; i < msa_len-1; i++ {
		base_counts := make([]int, 4)

		for j := 0; j < len(msa.OtherSequences); j++ {
			seq_item := msa.OtherSequences[j].Yield()

			if seq_item == nil {
				fmt.Println("WARNING: Seq_item is nil but it shouldnt be!!")
				return consensus
			}

			switch seq_item.Value {
			case 'A':
				base_counts[0] += 1
			case 'C':
				base_counts[1] += 1
			case 'G':
				base_counts[2] += 1
			case 'T':
				base_counts[3] += 1
			}
		}
		// get the most frequent base in the column
		max_count := base_counts[0]
		max_index := 0
		for k := 1; k < len(base_counts); k++ {
			if base_counts[k] > max_count {
				max_count = base_counts[k]
				max_index = k
			}
		}

		switch max_index {
		case 0:
			consensus.InsertTail('A')
		case 1:
			consensus.InsertTail('C')
		case 2:
			consensus.InsertTail('G')
		case 3:
			consensus.InsertTail('T')
		}
	}

	return consensus
}

// iterate over all sequences in the MSA
// if a column contains only gaps, remove it
func (msa *MSA) DeleteSpaces() {
	gaps_count := 0

	// reset iter pointers to the beginning of the sequences
	msa.ResetIterPointers()

	i := 0
	null_count := 0

	for null_count < len(msa.OtherSequences) {
		// reset null_count
		null_count = 0

		// iterate over all sequences simultaneously
		for j := 0; j < len(msa.OtherSequences); j++ {
			// in yield iter pointer returns current node
			// and moves to the next one
			seq_item := msa.OtherSequences[j].Yield()

			if seq_item == nil {
				null_count++
				continue
			}

			if seq_item.Value == '-' {
				gaps_count++
			}
		}
		if gaps_count == len(msa.OtherSequences) {
			i += 1
			for j := 0; j < len(msa.OtherSequences); j++ {
				// if all sequences have a gap in the current column
				// delete the node
				msa.OtherSequences[j].Delete(msa.OtherSequences[j].IterPointer.Prev)
			}
		}

		// reset gaps count
		gaps_count = 0
		i += 1
	}
}
