package cmd

import (
	"fmt"
	"os"

	eirinix "github.com/SUSE/eirinix"
	ingress "github.com/mudler/eirini-ingress/extensions/ingress"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var kubeconfig string
var namespace string
var rootCmd = &cobra.Command{
	Use:   "eirini-ingress",
	Short: "eirini-ingress creates ingress and services for apps pushed in Cloud Foundry",
	PreRun: func(cmd *cobra.Command, args []string) {

		viper.BindPFlag("kubeconfig", cmd.Flags().Lookup("kubeconfig"))
		viper.BindPFlag("namespace", cmd.Flags().Lookup("namespace"))

		viper.BindEnv("kubeconfig")
		viper.BindEnv("namespace", "NAMESPACE")
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		filter := false
		x := eirinix.NewManager(
			eirinix.ManagerOptions{
				Namespace:           viper.GetString("namespace"),
				KubeConfig:          viper.GetString("kubeconfig"),
				OperatorFingerprint: "eirini-ingress", // Not really used for now, but setting it up for future
				FilterEiriniApps:    &filter,
			})
		x.GetLogger().Info("Starting watcher in", x.GetManagerOptions().Namespace)
		x.GetLogger().Info(" Kubeconfig ", x.GetManagerOptions().KubeConfig)
		x.AddWatcher(ingress.NewPodWatcher())

		err = x.Watch()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "eirini", "Namespace to watch for Eirini apps")
	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to a kubeconfig, not required in-cluster")
}
