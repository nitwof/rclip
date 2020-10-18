package cmd

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	copyData string
)

var copyCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"cp", "c"},
	Short:   "Copy content and sends it to RClip server",
	PreRun: func(cmd *cobra.Command, args []string) {
		registerViperKey(
			"client.target",
			"CLIENT_TARGET",
			cmd.Flags().Lookup("target"),
			ServerDefaultAddr,
		)
		registerViperKey(
			"client.copy.clipboard",
			"CLIENT_COPY_CLIPBOARD",
			cmd.Flags().Lookup("clipboard"),
			false,
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		data := []byte(copyData)

		if len(copyData) == 0 {
			var err error

			data, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to read stdin")
			}
		}

		targetAddr := viper.GetString("client.target")

		conn, err := grpc.Dial(targetAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatal().
				Err(err).
				Str("target", targetAddr).
				Msg("Failed to connect to the server")
		}

		defer conn.Close()

		client := api.NewClipboardAPIClient(conn)
		_, err = client.Push(context.Background(), &api.PushRequest{Value: data})
		if err != nil {
			log.Fatal().
				Err(err).
				Str("method", "Push").
				Msg("Failed to execute RPC method")
		}

		if viper.GetBool("client.copy.clipboard") {
			clipboard.Write(data)
		}
	},
}

func init() {
	copyCmd.Flags().StringP(
		"target", "t", ServerDefaultAddr,
		"Target server address",
	)
	copyCmd.Flags().BoolP(
		"clipboard", "c", false,
		"Also copy value to the system clipboard",
	)
	copyCmd.Flags().StringVarP(
		&copyData, "data", "d", "",
		"Use the given data instead of stdin",
	)
}
