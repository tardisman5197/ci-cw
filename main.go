package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

const executeTime = 3
const numberOfGoods = 20

// Artificial Immnune System
const aispopSize = 4
const cloneSizeFactor = 2
const bestFitness = 5250
const replacementSize = 2

// PSO
const psoPopSize = 2000
const cognativeCoeff = 0.8
const socialCoeff = 0.8
const intertia = 0.9

func main() {

	algorithmPtr := flag.String("alg", "PSO", "Sets algorithm to run: 'PSO', 'AIS'")
	fileExtentionPtr := flag.String("ext", "final", "Sets the file extention to the csv file")

	collectPtr := flag.Bool("c", false, "Collects data over a selection of seeds for a certain algorithm")

	seedPtr := flag.Int("seed", 0, "Set the default seed")

	flag.Parse()

	if *collectPtr {
		collectData(*algorithmPtr, *fileExtentionPtr)
	} else {
		demo(int64(*seedPtr))
	}
}

// collectData runs an algorithm with multiple seeds for the same
// amount of time and writes the result to a csv file.
func collectData(algorithm string, id string) {
	seeds := []int64{0, 50, 100, 150, 200, 250, 300, 350, 400,
		450, 500, 550, 600, 650, 700, 750, 800, 850, 900, 950,
		1000, 1050, 1100, 1150, 1200, 1250, 1300, 1350, 1400,
		1450, 1500, 1550, 1600}

	var data [][]string
	var averageRevenue float64

	data = append(data, []string{"executeTime", strconv.FormatInt(executeTime, 10)})
	data = append(data, []string{"numberOfGoods", strconv.FormatInt(numberOfGoods, 10)})

	switch algorithm {
	case "AIS":
		data = append(data, []string{"popSize", strconv.FormatInt(aispopSize, 10)})
		data = append(data, []string{"cloneSizeFactor", strconv.FormatInt(cloneSizeFactor, 10)})
		data = append(data, []string{"bestFitness", strconv.FormatInt(bestFitness, 10)})
		data = append(data, []string{"replacementSize", strconv.FormatInt(replacementSize, 10)})
	case "PSO":
		data = append(data, []string{"psoPopSize", strconv.FormatInt(psoPopSize, 10)})
		data = append(data, []string{"cognativeCoeff", strconv.FormatFloat(cognativeCoeff, 'f', -1, 64)})
		data = append(data, []string{"socialCoeff", strconv.FormatFloat(socialCoeff, 'f', -1, 64)})
		data = append(data, []string{"intertia", strconv.FormatFloat(intertia, 'f', -1, 64)})
	default:
		fmt.Printf("Algrothm not selected correctly: %v", algorithm)
		return
	}

	data = append(data, []string{})
	data = append(data, []string{"Seed", "Revenue"})

	for i := 0; i < len(seeds); i++ {
		fmt.Printf("Seed: %v\n", seeds[i])

		// Init Pricing Problem
		var p PricingProblem
		p.PricingProblem(numberOfGoods, seeds[i])

		var currentRevenue float64

		switch algorithm {
		case "AIS":
			// Run AIS on current seed
			fmt.Printf("\tStarting Artificial Immune System\n")
			_, currentRevenue = artificialImmuneSystem(numberOfGoods, p, seeds[i])
		case "PSO":
			// Run AIS on current seed
			fmt.Printf("\tStarting PSO\n")
			_, currentRevenue = PSO(numberOfGoods, p, seeds[i])
		}

		averageRevenue += currentRevenue

		// Append to data: seed, revenue
		data = append(data,
			[]string{
				strconv.FormatInt(seeds[i], 10),
				strconv.FormatFloat(currentRevenue, 'f', -1, 64),
			})

		fmt.Printf("\tRevenue %v\n\n", currentRevenue)
	}

	// Append to data: "Average", averageRevnue
	averageRevenue /= float64(len(seeds))
	data = append(data,
		[]string{
			"Average",
			strconv.FormatFloat(averageRevenue, 'f', -1, 64),
		})
	fmt.Printf("Average: %v\n", averageRevenue)

	writeCSV(algorithm+"_results_"+id+".csv", data)
}

// demo runs each of the algorithms on the default
// seed for the same amount of time.
func demo(defaultSeed int64) {
	var p PricingProblem
	p.PricingProblem(numberOfGoods, defaultSeed)

	// Run given random search
	fitnessTester(p, defaultSeed)
	fmt.Printf("\n\n")

	// Run Artificial Immune
	fmt.Printf("Starting Artificial Immune System\n")
	prices, revenue := artificialImmuneSystem(numberOfGoods, p, defaultSeed)
	fmt.Printf("Prices: %v\nRevenue %v\n\n", prices, revenue)

	// Run PSO algorithm
	fmt.Printf("Starting PSO\n")
	prices, revenue = PSO(numberOfGoods, p, defaultSeed)
	fmt.Printf("Prices: %v\nRevenue %v\n\n", prices, revenue)
}

// randomPrices generates a list of random prices.
func randomPrices(noOfGoods int) (prices []float64) {
	for i := 0; i < noOfGoods; i++ {
		prices = append(prices, rand.Float64()*10)
	}
	return
}

// generateRandomPopulation creates a population of prices.
func generateRandomPopulation(noOfGoods int, popSize int, p PricingProblem) (population []Prices) {
	for i := 0; i < popSize; i++ {
		currentPrices := Prices{}
		currentPrices.prices = randomPrices(noOfGoods)
		currentPrices.revenue = p.evaluate(currentPrices.prices)
		population = append(population, currentPrices)
	}
	return
}

// writeCSV writes data to a csv file with a given filename.
func writeCSV(fileName string, data [][]string) {
	fmt.Printf("Writing File: results/%v\n", fileName)
	file, err := os.Create("results/" + fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("Finished writing\n")
}
