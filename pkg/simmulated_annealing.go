package pkg

import (
	// "fmt"
	"math"
	"math/rand"
)

type SimmulatedAnnealing struct {
	BestCandidate        MSA
	AllTimeBestCandidate MSA
	Mutators             []Mutator
	Evaluator            MSAEvaluator
	MaxIterations        int
	// sa params
	Temp        float64
	CoolingRate float64
	MinimumTemp float64
	// verbose
	Verbose bool
	// Logging
	Log Logger
}

func (sa *SimmulatedAnnealing) Init() {
	sa.Evaluator.Init(&sa.BestCandidate)
	sa.BestCandidate.Init()
	sa.Evaluator.Evaluate(&sa.BestCandidate)
}

func (sa *SimmulatedAnnealing) Optimize() {

	p_accept := 0.0
	sa.AllTimeBestCandidate = sa.BestCandidate.Copy()

	// fmt.Println("Initial alignment:")
	sa.BestCandidate.Print()
	sa.AllTimeBestCandidate.Print()

	for i := 0; i < sa.MaxIterations; i++ {
		// create a copy of the individual
		candidate_msa := sa.BestCandidate.Copy()

		// mutate the individual using all the mutators
		for _, mutator := range sa.Mutators {
			mutator.Mutate(&candidate_msa)
		}

		// delete columns containing only gaps
		candidate_msa.DeleteSpaces()

		// evaluate the candidate
		sa.Evaluator.Evaluate(&candidate_msa)

		// fmt.Println("Candidate score: ", candidate_msa.Score)
		// fmt.Println("All time best: ", sa.AllTimeBestCandidate.Score)

		delta := sa.BestCandidate.Score - candidate_msa.Score

		if candidate_msa.Score > sa.AllTimeBestCandidate.Score {
			sa.AllTimeBestCandidate = candidate_msa.Copy()
		}

		if delta <= 0 { // better or equal solution, WORKS ONLY FOR MAXIMIZATION
			sa.BestCandidate = candidate_msa.Copy()
		} else { // accept worse solution with probability e^(-delta/temp)
			p_accept = math.Exp(-(delta) / sa.Temp)

			if rand.Float64() < p_accept {
				sa.BestCandidate = candidate_msa.Copy()
			}
		}

		// if candidate_msa.Score > sa.BestCandidate.Score {
		// 	fmt.Println("Iteration: ", i)
		// 	fmt.Println("Candidate score: ", candidate_msa.Score)
		// 	fmt.Println("new best")
		// 	candidate_msa.Print()
		// 	sa.BestCandidate = candidate_msa.Copy()
		// 	fmt.Println("====================================")
		// }

		// cooling schedule
		if sa.Temp > sa.MinimumTemp {
			sa.Temp *= sa.CoolingRate
		}

		if i%200 == 0 {
			sa.Log.WriteLog(sa.AllTimeBestCandidate.Score)
		}
	}
}
