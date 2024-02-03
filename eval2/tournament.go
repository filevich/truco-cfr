package eval2

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval/dataset"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Tournament struct {
	NumPlayers     int
	NumDoubleGames int
	Agents         []bot.Agent
	Games          Table
}

func (t *Tournament) Start(ds dataset.Dataset, verbose bool) {

	t.NumDoubleGames = len(ds)

	if verbose {
		log.Printf("\nTournament %dp:\n", t.NumPlayers)
		for ix, agent := range t.Agents {
			log.Printf("\t%2d. %s\n", ix+1, agent.UID())
		}

		log.Printf("\nDone: ")
	}

	for i := 0; i < len(t.Agents)-1; i++ {
		agent1 := t.Agents[i]
		agent1.Initialize()

		for j := i + 1; j < len(t.Agents); j++ {
			agent2 := t.Agents[j]
			agent2.Initialize()

			// evaluar_bin.go
			res_partidas := PlayDoubleGames(ds, agent1, agent2, t.NumPlayers)
			t.Games.Add(
				agent1.UID(),
				agent2.UID(),
				res_partidas,
			)

			// termino de jugar contra agent2 -> ya no lo necesito
			agent2.Free()
			runtime.GC()
		}

		if verbose {
			log.Printf("%s", agent1.UID())
			if i == len(t.Agents)-2 {
				log.Printf("\n\n")
			} else {
				log.Printf(", ")
			}
		}

		// Just finished `agent1` evaluation
		// No longer needed
		// Free it
		agent1.Free()
		runtime.GC()
	}
}

func (t *Tournament) Inscriptos(subtitle string) []interface{} {
	// "0.857   8.437" ~> tiene len=20
	res := make([]interface{}, 0)
	for _, agent := range t.Agents {
		res = append(res, fmt.Sprintf("%s\n  %s ", agent.UID(), subtitle))
	}
	return res
}

func (torneo *Tournament) PrintWrTable(tabla Table) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		append(
			table.Row{"A\\B"},
			append(
				torneo.Inscriptos("WP    ADP"),
				"B/A",
			)...,
		),
	)

	for i := 0; i < len(torneo.Agents); i++ {

		agent1 := torneo.Agents[i]
		row := []interface{}{
			agent1.UID(),
		}

		for j := 0; j < len(torneo.Agents); j++ {
			if i == j {
				row = append(row, "             ")
				continue
			}
			agent2 := torneo.Agents[j]
			wp, adp := tabla.Metrics(agent1.UID(), agent2.UID())
			row = append(row, fmt.Sprintf("%.3f  %.3f", wp, adp))
		}

		row = append(row, agent1.UID())

		t.AppendRow(row)
		t.AppendSeparator()
	}

	// t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.Render()
}

func (torneo *Tournament) PrintWaldTable(tabla Table) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		append(
			table.Row{"A\\B"},
			append(
				torneo.Inscriptos("L    U"),
				"B/A",
			)...,
		),
	)

	for i := 0; i < len(torneo.Agents); i++ {

		agent1 := torneo.Agents[i]
		row := []interface{}{
			agent1.UID(),
		}

		for j := 0; j < len(torneo.Agents); j++ {
			if i == j {
				row = append(row, "             ")
				continue
			}
			agent2 := torneo.Agents[j]
			// wp, adp := tabla.Metrics(agent1.UID(), agent2.UID())
			// row = append(row, fmt.Sprintf("%.3f  %.3f", wp, adp))
			u, l := tabla.WaldInterval(agent1.UID(), agent2.UID())
			row = append(row, fmt.Sprintf("%.3f %.3f", l, u))
		}

		row = append(row, agent1.UID())

		t.AppendRow(row)
		t.AppendSeparator()
	}

	// t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.Render()
}

func (torneo *Tournament) Report() {
	if torneo.NumDoubleGames > 0 {
		log.Printf("%d Double games:\n", torneo.NumDoubleGames)

		log.Println()
		log.Printf("\nTABLE: WR (win ratio) & ADP (Avg. Diff. Points) for A vs B\n\n")
		torneo.PrintWrTable(torneo.Games)

		log.Println()
		log.Printf("\nTABLE: Adjusted Wald Intervals (90%%) for A vs B\n\n")
		torneo.PrintWaldTable(torneo.Games)
	}
}
