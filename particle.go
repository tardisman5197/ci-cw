package main

import (
	"math/rand"
)

// Particle models a valid set of values for the antennae array problem.
type Particle struct {
	currentPostion  []float64
	currentVelocity []float64
	pBest           []float64
	pBestRevenue    float64
}

// update calculates the new velocity and position of a particle.
// The postion is calculated by taking the previous postion and adding the
// current velocity.
//		x(t + 1) = x(t) + v(t + 1)
//
// The velocity is calculated by taking into account inertia,
// congnative attraction (pBest) and social attraction (gBest)
//		v(t + 1) = n.v(t) + Phi1.r1.(p(t) - x(t)) + Phi2.r2.(g(t) - x(t))
//
// where n and Phi are  the intertia, conganative and social coefficents and
// r is uniform random vector.
func (p *Particle) update(gBest []float64) {
	// When updating the postion or velocity the last element in the postion is
	// not changed as one antennae has to remain fixed.

	// update the position
	for i := 0; i < len(p.currentPostion)-1; i++ {
		p.currentPostion[i] = p.currentPostion[i] + p.currentVelocity[i]
	}

	// update the velocity
	for i := 0; i < len(p.currentVelocity)-1; i++ {
		// r is a random number between 0 and 1
		r1 := rand.Float64()
		r2 := rand.Float64()
		p.currentVelocity[i] = (intertia * p.currentVelocity[i]) +
			(cognativeCoeff * r1 * (p.pBest[i] - p.currentPostion[i])) +
			(socialCoeff * r2 * (gBest[i] - p.currentPostion[i]))
	}
}

// evaluate calaculates the fitness of the particles current position and updates
// the personal best of the particle.
func (p *Particle) evalulate(pP PricingProblem) bool {
	currentRevenue := pP.evaluate(p.currentPostion)

	if currentRevenue > p.pBestRevenue {
		p.pBest = make([]float64, len(p.currentPostion))
		copy(p.pBest, p.currentPostion)
		p.pBestRevenue = currentRevenue
		return true
	}
	return false
}
