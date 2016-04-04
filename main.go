package main

import (
	"fmt"
	"math/rand"
	"time"
)

type game struct {
	inning      int32
	side        string
	homeScore   int32
	awayScore   int32
	outs        int32
	strikes     int32
	balls       int32
	firstBase   bool
	secondBase  bool
	thirdBase   bool
	homePitches int32
	awayPitches int32
}

func newGame() game {
	return game{
		inning:      1,
		side:        "top",
		homeScore:   0,
		awayScore:   0,
		outs:        0,
		strikes:     0,
		balls:       0,
		firstBase:   false,
		secondBase:  false,
		thirdBase:   false,
		homePitches: 0,
		awayPitches: 0,
	}
}

func (g *game) switchSides() {
	if g.side == "top" {
		g.side = "bottom"
	} else {
		g.side = "top"
	}
	g.outs = 0
	g.firstBase = false
	g.secondBase = false
	g.thirdBase = false
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (g *game) throwPitch() {
	const pitchesThrown = 819
	const ballsThrown = 307
	const strikesThrown = 512

	probBall := float32(ballsThrown) / float32(pitchesThrown) // ~40%
	//probStrike := strikesThrown / pitchesThrown // ~60%

	num := rand.Float32()
	if num < probBall {
		g.balls++
	} else {
		g.strikes++
	}

	if g.side == "top" {
		g.homePitches++
	} else {
		g.awayPitches++
	}
}

func (g *game) strikeout() {
	g.balls = 0
	g.strikes = 0
	g.outs++
}

func (g *game) walk() {
	g.balls = 0
	g.strikes = 0
	g.advanceBatter()
}

// advanceBatter advances the batter to first base and handles any existing runners.
func (g *game) advanceBatter() {
	// is there a runner on first?
	if g.firstBase {
		// is there already a batter on second?
		if g.secondBase {
			// is there already a batter on third?
			if g.thirdBase {
				// current team scored a run!
				if g.side == "top" {
					g.awayScore++
				} else {
					g.homeScore++
				}
			}
			g.thirdBase = true

		}
		// move that batter to second
		g.secondBase = true
	}
	g.firstBase = true
}

func (g *game) batterUp() bool {
	return g.balls < 4 && g.strikes < 3
}

type batterResult int32

const (
	BATTER_UNDECIDED = iota
	BATTER_STRIKEOUT
	BATTER_WALK
)

func (g *game) batterResult() batterResult {
	if g.balls == 4 {
		return BATTER_WALK
	}

	if g.strikes == 3 {
		return BATTER_STRIKEOUT
	}

	return BATTER_UNDECIDED
}

func main() {
	fmt.Println("John Henry Baseball Simulator")

	const inningsToPlay = 1
	const simsToRun = 2

	for sim := 0; sim < simsToRun; sim++ {
		currentGame := newGame()
		printStateHeader()
		for ; currentGame.inning < inningsToPlay+1; currentGame.inning++ {
			for side := 0; side < 2; side++ {
				for currentGame.outs < 3 {
					printState(currentGame)

					for currentGame.batterUp() {
						currentGame.throwPitch()
						printState(currentGame)
					}

					switch currentGame.batterResult() {
					case BATTER_STRIKEOUT:
						currentGame.strikeout()
					case BATTER_WALK:
						currentGame.walk()
					}
				}
				currentGame.switchSides()
			}
		}
		printStateFooter()
		printStatTable(currentGame)
	}
}

func printStateHeader() {
	fmt.Printf(headerFormat,
		"inning", "side", "home", "away", "outs", "count", "1B", "2B", "3B")
}
func printState(g game) {

	first, second, third := "O", "O", "O"

	if g.firstBase {
		first = "X"
	}

	if g.secondBase {
		second = "X"
	}

	if g.thirdBase {
		third = "X"
	}

	fmt.Printf(rowFormat,
		g.inning, g.side, g.homeScore, g.awayScore, g.outs, g.balls, g.strikes,
		first, second, third)
}

func printStateFooter() {
	fmt.Printf(footerFormat)
}

const headerFormat = `
| %6s | %6s | %4s | %4s | %4s | %5s | %2s | %2s | %2s |
|--------|--------|------|------|------|-------|----|----|----|
`

const footerFormat = `|--------|--------|------|------|------|-------|----|----|----|
`

const rowFormat = `| %6d | %6s | %4d | %4d | %4d | %d - %d | %2s | %2s | %2s |
`

func printStatTable(g game) {
	fmt.Printf(statTable, g.homePitches, g.awayPitches)
}

const statTable = `
| Stat    | Home | Away |
|---------|------|------|
| Pitches | %4d | %4d |
|---------|------|------|
`
