package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/MaxHalford/eaopt"
)

// The name of the file to which the GA progress is appended to
const progressFileName = "progress.json"

// A Vector contains float64s.
type Vector []float64

// Evaluate a Vector.
func (X Vector) Evaluate() (y float64, err error) {
	y = math.Pow(X[0], 2) + math.Pow(X[1], 2)
	if math.Abs(X[0])+math.Abs(X[1]) < 4 {
		y += 10000
	}
	return
}

// Mutate a Vector by applying by resampling each element from a normal
// distribution with probability 0.8.
func (X Vector) Mutate(rng *rand.Rand) {
	eaopt.MutNormalFloat64(X, 0.8, rng)
}

// Crossover a Vector with another Vector by applying uniform crossover.
func (X Vector) Crossover(Y eaopt.Genome, rng *rand.Rand) {
	eaopt.CrossUniformFloat64(X, Y.(Vector), rng)
}

// Clone a Vector.
func (X Vector) Clone() eaopt.Genome {
	var XX = make(Vector, len(X))
	copy(XX, X)
	return XX
}

// MakeVector returns a random vector by generating 5 values uniformally
// distributed between -10 and 10.
func MakeVector(rng *rand.Rand) eaopt.Genome {
	return Vector(eaopt.InitUnifFloat64(2, -20, 20, rng))
}

func main() {
	var ga, err = eaopt.NewDefaultGAConfig().NewGA()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Empty the progress file
	f, _ := os.Create(progressFileName)
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprint(w, "")
	w.Flush()

	// Open the progress file in append mode
	f, _ = os.OpenFile(progressFileName, os.O_APPEND|os.O_WRONLY, 0666)
	defer f.Close()

	// Append the initial GA status to the progress file
	var bytes, _ = json.Marshal(ga)
	f.WriteString(string(bytes) + "\n")

	// Add a custom print function to track progress
	ga.Callback = func(ga *eaopt.GA) {
		fmt.Printf("Best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)
		// Append the current GA status to the progress file
		var bytes, _ = json.Marshal(ga)
		f.WriteString(string(bytes) + "\n")
	}

	// Run the GA
	ga.Minimize(MakeVector)
}
