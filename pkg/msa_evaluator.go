package pkg

import (
	"fmt"
  "math"
)

// evaluator used to evaluate the quality of the MSA
type MSAEvaluator interface {
  // initialize the evaluator
  // depending on the evaluator this might need to happen
  // before the candidate is initialized
	Init(msa *MSA)
  // evaluate the MSA
	Evaluate(msa *MSA) float64
}

type MeanColumnEntropy struct{}

func computeEntropy (counts []int, colSize int) float64 {
  H := 0.0
  for i := 0; i < len(counts); i++ {
    if counts[i] != 0 {
      p := float64(counts[i])
      H += p * math.Log2(p)
    }
  }

  return H / float64(colSize)
}

func (me *MeanColumnEntropy) Init(msa *MSA) {}

func (me *MeanColumnEntropy) Evaluate(msa *MSA) float64 {
  meanEntropy := 0.0
	baseCounts := make([]int, 4)
	msaLen := msa.OtherSequences[0].Length - 1

  // RESET ITER POINTERS
  msa.ResetIterPointers()

	// iterate over all sequences simultaneously
	for i := 0; i < msaLen; i++ {
    // reset base_counts
    baseCounts = []int{0, 0, 0, 0}

		// for each sequence
		for j := 0; j < len(msa.OtherSequences); j++ {
			// get the base at the current position
			base := msa.OtherSequences[j].Yield()

			switch base.Value {
			case 'A':
				baseCounts[0]++
			case 'C':
				baseCounts[1]++
			case 'G':
				baseCounts[2]++
			case 'T':
				baseCounts[3]++
			}
		}

    // fmt.Println("i", i, "baseCounts", baseCounts)

    meanEntropyVal := computeEntropy(baseCounts, len(msa.OtherSequences))
    // fmt.Println("meanEntropyVal", meanEntropyVal)
    meanEntropy += meanEntropyVal
	}

  msa.Score = meanEntropy / float64(msaLen)

  return meanEntropy / float64(msaLen)
}

type GLOCSA struct {
}

func (e *GLOCSA) Init(msa *MSA) {
}

func meanColumnHomogenity(msa *MSA) float64 {

  // columnHomogenity := 0.0
	baseCounts := make([]int, 5)
	msaLen := msa.OtherSequences[0].Length

  // RESET ITER POINTERS
  msa.ResetIterPointers()

	// iterate over all sequences simultaneously
	for i := 0; i < msaLen; i++ {
    // reset base_counts
    baseCounts = []int{0, 0, 0, 0, 0}

		// for each sequence
		for j := 0; j < len(msa.OtherSequences)-1; j++ {
			// get the base at the current position
			base := msa.OtherSequences[j].Yield()
      if base == nil {
        fmt.Println("sequence ", j, "is nil at ", i)
        continue
      }

			switch base.Value {
			case 'A':
				baseCounts[0]++
			case 'C':
				baseCounts[1]++
			case 'G':
				baseCounts[2]++
			case 'T':
				baseCounts[3]++
			}
		}
	}
  return 0.0
}

func gapConcentration(msa *MSA) float64 {
  return 0.0
}

func columnIncrement(msa *MSA) float64 {
  return 0.0
}

// compute the entropy of each column in the MSA
// and sum them up
func (e *GLOCSA) Evaluate(msa *MSA) float64 {

  return 0.0
}

type UniformEvaluator struct {
	SumL float64
}

func (ue *UniformEvaluator) Init(msa *MSA) {
	ue.SumL = computeSumL(msa)	
}

// computes and assignes the score to the MSA.Score field
// also returns the value in case it is needed
func (ue *UniformEvaluator) Evaluate(msa *MSA) float64 {

	CS := float64(computeCoverage(msa))
	// fmt.Println("CS", CS)

	msa_len := float64(msa.OtherSequences[0].Length)
	// fmt.Println("msa_len", msa_len)
	// fmt.Println("ue.SumL", ue.SumL)
	LS := (ue.SumL - msa_len) / ue.SumL

	// fmt.Println("LS", LS)

	msa.Score = CS + LS

	return CS + LS
}

// computes the sum of the length of all the sequences
func computeSumL(msa *MSA) float64 {
	sumL := 0.0
	for i := 0; i < len(msa.OtherSequences); i++ {
		sumL += float64(msa.OtherSequences[i].Length)
	}
	return sumL
}

func findMax(numbers []int) int {
	max_count := numbers[0]

	for _, num := range numbers {
		if num > max_count {
			max_count = num
		}
	}

	return max_count
}

// compute the coverage of the MSA by the consensus sequence
func computeCoverage(msa *MSA) int {

	coverage := 0
	base_counts := make([]int, 4)
	msa_len := msa.OtherSequences[0].Length - 1

  // RESET ITER POINTERS
  msa.ResetIterPointers()

	// iterate over all sequences simultaneously
	for i := 0; i < msa_len; i++ {
    // reset base_counts
    base_counts = []int{0, 0, 0, 0}

		// for each sequence
		for j := 0; j < len(msa.OtherSequences); j++ {
			// get the base at the current position
			base := msa.OtherSequences[j].Yield()

      // if base == nil {
      //   fmt.Println("sequence ", j, "is nil at ", i)
      //   continue
      // }

			switch base.Value {
			case 'A':
				base_counts[0]++
			case 'C':
				base_counts[1]++
			case 'G':
				base_counts[2]++
			case 'T':
				base_counts[3]++
			}
		}

		coverage += findMax(base_counts)
	}

	return coverage
}
