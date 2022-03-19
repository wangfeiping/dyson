package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yalp/jsonpath"

	"github.com/wangfeiping/dyson/config"
	"github.com/wangfeiping/log"
)

var starter = func() (cancel context.CancelFunc, err error) {
	log.Info("Start...")

	doJob()

	return
}

func doJob() {
	log.Debug("Exec: ")

	// ./gaiad q gov proposals --output json --count-total --limit 10 --status voting_period
	jsonStr := `{"proposals":[
				{"proposal_id":"62",
					"content":{
						"@type":"/cosmos.gov.v1beta1.TextProposal",
						"title":"Signal Proposal: Migration of Gravity DEX to a Separate Cosmos Chain",
						"description":"..."},
					"status":"PROPOSAL_STATUS_VOTING_PERIOD",
					"final_tally_result":{"yes":"0","abstain":"0","no":"0","no_with_veto":"0"},
					"submit_time":"2022-03-14T09:43:54.403555411Z",
					"deposit_end_time":"2022-03-28T09:43:54.403555411Z",
					"total_deposit":[{"denom":"uatom","amount":"64000000"}],
					"voting_start_time":"2022-03-14T09:43:54.403555411Z",
					"voting_end_time":"2022-03-28T09:43:54.403555411Z"}],
				"pagination":{"next_key":null,"total":"1"}}`
	// jsonStr := "{}"
	// filter, err := jsonpath.Prepare("$..proposal_id")
	// filter, err := jsonpath.Prepare("$.proposals[*].voting_start_time")
	filter, err := jsonpath.Prepare("$.proposals[0].voting_start_time")
	if err != nil {
		log.Error(err)
		return
	}
	var data interface{}
	if err = json.Unmarshal([]byte(jsonStr), &data); err != nil {
		log.Error(err)
		return
	}
	out, err := filter(data)
	if err != nil {
		log.Error("filter error: ", err)
		return
	}
	fmt.Println("voting_start_time: ", out)

	filter, err = jsonpath.Prepare("$.proposals[0].proposal_id")
	if err != nil {
		log.Error(err)
		return
	}
	out, err = filter(data)
	if err != nil {
		log.Error("filter error: ", err)
		return
	}
	fmt.Println("proposal_id", out)

	// {"balances":[{"denom":"uatom","amount":"1012728"}],"pagination":{"next_key":null,"total":"0"}}

	// config.Save()
	config.Load()
	execs := config.GetAll()

	for _, exec := range execs {
		log.Info("exec command: ", exec.Command)
	}

}

// NewStartCommand 创建 start/服务启动 命令
func NewStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   config.CmdStart,
		Short: "Start",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(starter, true)
		},
	}

	cmd.Flags().Int64P(config.FlegDuration, "d", 30, "The cycle time of the watch task")
	cmd.Flags().StringP(config.FlagListen, "l", ":9900", "The listening address(ip:port) of exporter")
	return cmd
}
