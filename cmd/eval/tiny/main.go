package main

import (
	"fmt"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval"
)

func main() {
	const (
		tiny_eval   = 1_000
		num_players = 2
		b           = "/media/jp/6e5bdfb0-c84b-4144-8d6d-4688934f1afe/models/6p/48np-multi6/a1"
	)

	fmt.Printf("loading T1K22...")
	ds := eval.Load_dataset("/Users/jp/Workspace/go/truco-ai/truco-cfr/eval/t1k22.json")
	fmt.Println(" [done]")

	agents := []bot.Agent{
		&bot.BotRandom{},
		&bot.BotSimple{},

		// &bot.BotCFR{
		// 	N: "final_es-lmccfr_d25h0m_D48h0m_t24878_p0_a1_2208092259",
		// 	F: b + "/final_es-lmccfr_d25h0m_D48h0m_t24878_p0_a1_2208092259.model",
		// },
	}

	for i, agent := range agents {
		agent.Initialize()
		fmt.Printf("[%2d/%2d] tiny evaluating %s...", i+1, len(agents), agent.UID())
		wr_ale, wr_det, di_ale, di_det, wu_ale, wd_ale, wu_det, wd_det, delta := eval.Tiny_eval_float(agent, num_players, ds[:tiny_eval])
		fmt.Println(" -> " + eval.Format_Tiny_eval(wr_ale, wr_det, di_ale, di_det, wu_ale, wd_ale, wu_det, wd_det, delta))
		agent.Free()
	}
}
