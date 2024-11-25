package serve

import (
	"context"
	"fmt"

	"github.com/eddycharly/generic-auth-server/apis/v1alpha1"
	"github.com/eddycharly/generic-auth-server/pkg/auth"
	"github.com/eddycharly/generic-auth-server/pkg/policy"
	healthz "github.com/eddycharly/generic-auth-server/pkg/probes"
	"github.com/eddycharly/generic-auth-server/pkg/signals"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func Command() *cobra.Command {
	var healthAddress string
	var authAddress string
	var healthCertFile string
	var healthKeyFile string
	var authCertFile string
	var authKeyFile string
	var kubeConfigOverrides clientcmd.ConfigOverrides
	command := &cobra.Command{
		Use:   "serve",
		Short: "Start the authz server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetLogger(zap.New(zap.UseDevMode(true)))
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var httpErr, mgrErr error
				err := func(ctx context.Context) error {
					// create a rest config
					kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
						clientcmd.NewDefaultClientConfigLoadingRules(),
						&kubeConfigOverrides,
					)
					config, err := kubeConfig.ClientConfig()
					if err != nil {
						return err
					}
					// create a wait group
					var group wait.Group
					// wait all tasks in the group are over
					defer group.Wait()
					// create a controller manager
					scheme := runtime.NewScheme()
					if err := v1alpha1.Install(scheme); err != nil {
						return err
					}
					mgr, err := ctrl.NewManager(config, ctrl.Options{
						Scheme: scheme,
					})
					if err != nil {
						return fmt.Errorf("failed to construct manager: %w", err)
					}
					// create compiler
					compiler := policy.NewCompiler()
					// create provider
					provider, err := policy.NewKubeProvider(mgr, compiler)
					if err != nil {
						return err
					}
					// create a cancellable context
					ctx, cancel := context.WithCancel(ctx)
					// start manager
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						mgrErr = mgr.Start(ctx)
					})
					if !mgr.GetCache().WaitForCacheSync(ctx) {
						defer cancel()
						return fmt.Errorf("failed to wait for cache sync")
					}
					// create servers
					health := healthz.NewServer(healthAddress, healthCertFile, healthKeyFile)
					auth := auth.NewHttpServer(authAddress, authCertFile, authKeyFile, provider)
					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						httpErr = health.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						httpErr = auth.Run(ctx)
					})
					return nil
				}(ctx)
				return multierr.Combine(err, httpErr, mgrErr)
			})
		},
	}
	command.Flags().StringVar(&healthAddress, "health-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&healthCertFile, "health-cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&healthKeyFile, "health-key-file", "", "File containing tls private key")
	command.Flags().StringVar(&authAddress, "auth-address", ":9081", "Address to listen on for auth checks")
	command.Flags().StringVar(&authCertFile, "auth-cert-file", "", "File containing tls certificate")
	command.Flags().StringVar(&authKeyFile, "auth-key-file", "", "File containing tls private key")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return command
}
