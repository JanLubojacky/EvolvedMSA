package pkg

import (
  "math/rand"
  "time"
  // "fmt"
)

// accept a MSA and return a new mutated MSA
type Mutator interface {
	Mutate(msa *MSA)
}

type RandomInsertGap struct {
	P_mut float64
}

func (mutator *RandomInsertGap) Mutate(msa *MSA) {

	for i := 0; i < len(msa.OtherSequences); i++ {
		for node := msa.OtherSequences[i].Head; node != nil; node = node.Next {

			rand_roll := rand.Float64()

			if rand_roll < mutator.P_mut {
				// for head 50 % chance to insert a gap at the head
				if node == msa.OtherSequences[i].Head && rand_roll < mutator.P_mut/2 {
					msa.OtherSequences[i].InsertHead('-')
				} else {
					msa.OtherSequences[i].InsertAfter(node, '-')
				}
			}
		}
	}
}

type RandomDeleteGap struct {
	P_mut float64
}

func (mutator *RandomDeleteGap) Mutate(msa *MSA) {
	for i := 0; i < len(msa.OtherSequences); i++ {
		for node := msa.OtherSequences[i].Head; node != nil; node = node.Next {
			// if the random number is less than p_mut
			// delete the gap
			rand_roll := rand.Float64()
			if rand_roll < mutator.P_mut {
				if node.Value == '-' {
					msa.OtherSequences[i].Delete(node)
				}
			}
		}
	}
}

type LengthEqualizer struct {
	P_mut float64
}

func padSequence(sequence *Sequence, target_length int, P_mut float64) {

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for sequence.Length < target_length {
		for n := sequence.Head; n != nil; n = n.Next {
			rand_roll := r.Float64()

			if rand_roll < P_mut {
				// for head 50 % chance to insert a gap at the head
				if n == sequence.Head && rand_roll < P_mut/2 {
					sequence.InsertHead('-')
				} else {
					sequence.InsertAfter(n, '-')
				}

				if sequence.Length == target_length {
					return
				}
			}
		}
	}
}

// fills shorter sequences with gaps until they match the length of the longest sequence
func (mutator *LengthEqualizer) Mutate(msa *MSA) {
	target_length := 0
	for _, sequence := range msa.OtherSequences {
		if sequence.Length > target_length {
			target_length = sequence.Length
		}
	}

	// for all the other sequences
	// insert gaps in random positions until they are all of the same length
	for i := 0; i < len(msa.OtherSequences); i++ {
		padSequence(&msa.OtherSequences[i], target_length, mutator.P_mut)
	}
}
