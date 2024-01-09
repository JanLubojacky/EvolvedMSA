package pkg

import "fmt"

type Optimizer interface {
  Optimize()
}

type LocalSearch struct {
  BestCandidate MSA
  Mutators []Mutator
  Evaluator MSAEvaluator
  MaxIterations int
}

func (ls LocalSearch) Init() {
  // initialize the evaluator
  // depending on the evaluator this might need to happen
  // before the candidate is initialized
  ls.Evaluator.Init(&ls.BestCandidate)
  // initialize the MSA by padding all the sequences with gaps
  // inserting them in random positions
  ls.BestCandidate.Init()
  // // evaluate the initial MSA
  // ls.Evaluator.Evaluate(&ls.BestCandidate)
}

func (ls LocalSearch) Optimize() {
  for i := 0; i < ls.MaxIterations; i++ {
    // create a copy of the individual
    candidate_msa := ls.BestCandidate.Copy()

    // mutate the individual using all the mutators
    for _, mutator := range ls.Mutators {
       mutator.Mutate(&candidate_msa)
    }

    // delete columns containing only gaps
    candidate_msa.DeleteSpaces()

    // evaluate the candidate
    ls.Evaluator.Evaluate(&candidate_msa)

    if candidate_msa.Score > ls.BestCandidate.Score {
      ls.BestCandidate = candidate_msa
    }

    // if i % 10 == 0 {
      // fmt.Println("Iteration: ", i)
      // fmt.Println("Score: ", ls.BestCandidate.Score)
      // ls.BestCandidate.Print()
    // }
  }

  fmt.Println("Final result: ")
  ls.BestCandidate.Print()
}
