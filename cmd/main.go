package main

import (
	"EvolvedMSA/pkg"
	"fmt"
)

func main() {
	filename := "data/dna_158_al.fasta" // Replace with your FASTA file path
	msa, err := pkg.ParseFASTA(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ga_log := pkg.Logger{}
	ga_log.Init("logs/ga_dna_158_al.csv")
	defer ga_log.Close()
  // sa_log := pkg.Logger{}
  // sa_log.Init("logs/sa_dna_158.csv")
  // defer sa_log.Close()

	for run := 0; run < 10; run++ {

		// // SIMMULATE ANNEALING
		// localSearch := pkg.SimmulatedAnnealing{
		// 	BestCandidate: msa,
		// 	Mutators: []pkg.Mutator{
		// 		&pkg.RandomInsertGap{P_mut: 0.05},
		// 		&pkg.RandomDeleteGap{P_mut: 0.1},
		// 		&pkg.LengthEqualizer{P_mut: 0.1},
		// 	},
		// 	// Evaluator:     &pkg.UniformEvaluator{},
		// 	Evaluator:     &pkg.MeanColumnEntropy{},
		// 	MaxIterations: 20000,
		// 	Temp:          60,
		// 	CoolingRate:   0.9999,
		// 	MinimumTemp:   1e-8,
		// 	Verbose:       true,
		//     Log: sa_log,
		// 	// Verbose:       false,
		// }
		// localSearch.Init()
		// localSearch.BestCandidate.Print()
		//
		// localSearch.Optimize()
		//
		// fmt.Println("Final alignment:")
		// localSearch.AllTimeBestCandidate.Print()
		//
		// fmt.Println()
		// fmt.Println("Target sequence:")
		// localSearch.BestCandidate.TargetSequence.Print()
		//
		// fmt.Println("Consensus sequence:")
		//
		// consensus := localSearch.AllTimeBestCandidate.ComputeConsensus()
		// consensus.Print()
		//

		ps := 10
		// EA
		ea := pkg.GeneticAlgorithm{
			PopulationSize: ps,
			ChildrenSize:   10,
			BestCandidate:  msa,
			Mutators: []pkg.Mutator{
				&pkg.RandomInsertGap{P_mut: 0.1},
				&pkg.RandomDeleteGap{P_mut: 0.5},
				&pkg.LengthEqualizer{P_mut: 0.01},
			},
			EvalFunc: &pkg.MeanColumnEntropy{},
			SelectionFunc: &pkg.TournamentSelection{
				PopulationSize: ps,
				TournamentSize: 3,
			},
			CrossFunc: &pkg.RandomShuffleCrossover{},
			// Verbose:   true,
			MaxIterations: 1000,
      Log: ga_log,
			// Log: log,
		}

		ea.Init()

		// fmt.Println("Initial population: ")
		// for i := 0; i < ea.PopulationSize; i++ {
		// 	ea.Population[i].Print()
		// }
		ea.Optimise()

		// ea.BestCandidate.Print()

		// fmt.Println("=== \n Final result: \n ===")
		// for i := 0; i < ea.PopulationSize; i++ {
		// 	ea.Population[i].Print()
		// }
    ga_log.EndRow()
    // sa_log.EndRow()

	}
}
