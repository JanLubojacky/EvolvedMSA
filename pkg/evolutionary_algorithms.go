package pkg

import (
	// "fmt"
	"math/rand"
	"sort"
)

type Crossover interface {
	Cross(p1 MSA, p2 MSA) MSA
}

// given an alignment, take sequences
// from two parents and create a child from them
type RandomShuffleCrossover struct {
	ShuffleProb float64 // probability of shuffling the alignment
}

func (self *RandomShuffleCrossover) Cross(p1 MSA, p2 MSA) MSA {

	// fmt.Println("CROSSING")

	child := MSA{}
	child.OtherSequences = make([]Sequence, len(p1.OtherSequences))

	for i := 0; i < len(p1.OtherSequences); i++ {
		if rand.Float64() < self.ShuffleProb {
			child.OtherSequences[i] = p1.OtherSequences[i].Copy()
		} else {
			child.OtherSequences[i] = p2.OtherSequences[i].Copy()
		}
	}

	return child
}

type Selector interface {
	Select(population []MSA) []MSA
}

type TournamentSelection struct {
	TournamentSize int // size of tournaments
	PopulationSize int // how many individuals to select
}

func (self *TournamentSelection) Select(population []MSA) []MSA {

	tournament_population := make([]MSA, self.PopulationSize)
	tournament := make([]MSA, self.TournamentSize)

	for i := 0; i < self.PopulationSize; i++ {
		// select tournament participants
		for j := 0; j < self.TournamentSize; j++ {
			tournament[j] = population[rand.Intn(len(population))]
		}

		// fmt.Println("tournament", i)
		// for _, child := range tournament {
		//   child.Print()
		// }

		// fmt.Println("tournament")
		// for _, child := range tournament {
		//   pr

		// order the participants by score
		sort.Slice(tournament, func(x, y int) bool {
			return tournament[x].Score > tournament[y].Score
		})

		// add the winner to the new population
		tournament_population[i] = tournament[0].Copy()
	}

	// fmt.Println("tournament population")
	// for _, child := range tournament_population {
	//   child.Print()
	// }

	return tournament_population
}

type GeneticAlgorithm struct {
	ProblemName    string       // name used to log the results
	PopulationSize int          // size of the population
	ChildrenSize   int          // size of the children population
	Population     []MSA        // population of MSAs
	BestCandidate  MSA          // keeps track of the best candidate
	EvalFunc       MSAEvaluator // function used to evaluate the MSAs
	Mutators       []Mutator    // function used to mutate the MSAs
	CrossFunc      Crossover    // function used to crossover the MSAs
	SelectionFunc  Selector     // function used to select the parents
	MaxIterations  int
	Log            Logger
}

func (self *GeneticAlgorithm) Init() {
	// ga.Evaluator.Init(&ga.BestCandidate)
	self.Population = make([]MSA, self.PopulationSize)
	for i := 0; i < self.PopulationSize; i++ {
		// copy the best candidate into the population
		self.Population[i] = self.BestCandidate.Copy()
		// initialize the individual
		self.Population[i].Init()
		self.EvalFunc.Evaluate(&self.Population[i])
	}

	self.EvalFunc.Init(&self.BestCandidate)
	self.EvalFunc.Evaluate(&self.BestCandidate)
}

func (self *GeneticAlgorithm) Optimise() {

	// initial eval

	// main loop
	for i := 0; i < self.MaxIterations; i++ {
		// mutations
		// for j := 0; j < self.PopulationSize; j++ {
		// 	// mutate the individual using all the mutators

		// create children
		children := make([]MSA, self.ChildrenSize)
		for j := 0; j < self.ChildrenSize; j++ {
			children[j] = MSA{}
		}

		// children[0].Print()
		// fmt.Println("population[0]")
		// self.Population[0].Print()

		for j := 0; j < self.ChildrenSize; j++ {

			// select parents
			p1 := rand.Intn(self.PopulationSize)
			p2 := rand.Intn(self.PopulationSize)

			// fmt.Println("p1", p1, "p2", p2, "j", j)

			// fmt.Println("children")
			// children[j].Print()
			// fmt.Println("population[p1]")
			// self.Population[p1].Print()
			// fmt.Println("population[p2]")
			// self.Population[p2].Print()

			// make a child
			children[j] = self.CrossFunc.Cross(self.Population[p1], self.Population[p2])

			// mutate the child
			// this has to happen last to repair the solution
			for _, mutator := range self.Mutators {
				mutator.Mutate(&children[j])
			}

			// evaluate the child
			self.EvalFunc.Evaluate(&children[j])
		}

		// fmt.Println("children after eval")
		// for _, child := range children {
		//   child.Print()
		// }

		// add childrens to the population
		new_population := make([]MSA, self.PopulationSize+self.ChildrenSize)
		copy(new_population, self.Population)
		copy(new_population[self.PopulationSize:], children)

		// fmt.Println("new population")
		// for _, child := range new_population {
		//   child.Print()
		// }

		// order the children by score
		sort.Slice(new_population, func(i, j int) bool {
			return new_population[i].Score > new_population[j].Score
		})

		// fmt.Println("new population sorted")
		// for _, child := range new_population {
		//   child.Print()
		// }

		if new_population[0].Score > self.BestCandidate.Score {
			self.BestCandidate = new_population[0]
		}

		// select the best individuals
		self.Population = self.SelectionFunc.Select(new_population)

		// fmt.Println("population after selection")
		// for _, child := range self.Population {
		//   child.Print()
		// }
		if i%10 == 0 {
			self.Log.WriteLog(self.BestCandidate.Score)
		}

	}

	sort.Slice(self.Population, func(i, j int) bool {
		return self.Population[i].Score > self.Population[j].Score
	})
}
