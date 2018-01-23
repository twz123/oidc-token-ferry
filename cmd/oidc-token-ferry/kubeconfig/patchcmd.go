package kubeconfig

import (
	"io/ioutil"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/twz123/oidc-token-ferry/pkg/api"
	"github.com/twz123/oidc-token-ferry/pkg/cli"
)

const useStdInOut = "-"

type patchCmd struct {
	UserName      func(string) `long:"user-name" description:"User name to use when generating client configuration. Either user-name or user-claim-name may be specified."`
	UserClaimName func(string) `long:"user-claim-name" description:"Claim that defines the user name to use when generating client configuration. Either user-name or user-claim-name may be specified."`

	Args struct {
		Kubeconfig   flags.Filename `positional-arg-name:"KUBECONFIG_FILE" description:"Path to the kubeconfig file to be patched. Uses the default discovery mechanism if omitted/empty. Special value '-' (hyphen) means read from STDIN."`
		Outputconfig flags.Filename `positional-arg-name:"OUTPUT_FILE" description:"Path to the patched kubeconfig file to be written. Overwrites kubeconfig if omitted/empty. Special value '-' (hyphen) means write to STDOUT."`
	} `positional-args:"yes"`

	cli cli.CLI

	internalError     error
	determineUserName func(*api.TokenFerry) (string, error)
}

func PatchCmd(cli cli.CLI) interface{} {
	cmd := &patchCmd{cli: cli}
	cmd.UserName = cmd.makeUserSelector(selectStaticUserName)
	cmd.UserClaimName = cmd.makeUserSelector(selectUserNameFromClaim)
	return cmd
}

func (cmd *patchCmd) Execute([]string) error {
	if cmd.internalError != nil {
		return cmd.internalError
	}

	if cmd.determineUserName == nil {
		return &flags.Error{
			Type:    flags.ErrRequired,
			Message: "either --user-name or --user-claim-name need to be specified",
		}
	}

	config, defaultOutputPath, err := loadClientConfig(string(cmd.Args.Kubeconfig))
	if err != nil {
		return errors.Wrap(err, "failed to patch kubeconfig")
	}

	ferry, err := cmd.cli.PerformChallenge()
	if err != nil {
		return err
	}

	userName, err := cmd.determineUserName(ferry)
	if err != nil {
		return errors.Wrap(err, "failed to determine user name")
	}

	patchClientConfig(config, userName, ferry)

	outputPath := string(cmd.Args.Outputconfig)
	if outputPath == "" {
		outputPath = defaultOutputPath
	}

	if outputPath == useStdInOut {
		content, err := clientcmd.Write(*config)
		if err != nil {
			return err
		}
		if _, err := os.Stdout.Write(content); err != nil {
			return errors.Wrap(err, "failed to write patched kubeconfig to STDOUT")
		}

		return nil
	}

	if err := clientcmd.WriteToFile(*config, outputPath); err != nil {
		return errors.Wrapf(err, "failed to write patched kubeconfig to %s", outputPath)
	}

	return nil
}

func loadClientConfig(kubeconfig string) (*clientcmdapi.Config, string, error) {
	switch kubeconfig {
	case "":
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		switch len(rules.Precedence) {
		case 1:
			kubeconfig = rules.Precedence[0]
		case 0:
			return nil, "", errors.New("there is no default client configuration")
		default:
			return nil, "", errors.Errorf("several default client configuration files available: %s", rules.Precedence)
		}

		config, err := rules.Load()
		if err != nil {
			return nil, "", err
		}

		return config, kubeconfig, nil

	case useStdInOut:
		kubeconfigBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, "", errors.Wrap(err, "failed to read kubeconfig from STDIN")
		}

		config, err := clientcmd.Load(kubeconfigBytes)
		if err != nil {
			return nil, "", err
		}

		return config, useStdInOut, nil

	default:
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		rules.ExplicitPath = kubeconfig
		config, err := rules.Load()
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to read kubeconfig from %s", kubeconfig)
		}

		return config, kubeconfig, nil
	}
}

func (cmd *patchCmd) makeUserSelector(selector func(*api.TokenFerry, string) (string, error)) func(string) {
	return func(value string) {
		if cmd.determineUserName != nil {
			cmd.internalError = errors.New("either user-name or user-claim-name may be specified")
			return
		}

		cmd.determineUserName = func(ferry *api.TokenFerry) (string, error) {
			return selector(ferry, value)
		}
	}
}

func selectStaticUserName(ferry *api.TokenFerry, userName string) (string, error) {
	return userName, nil
}

func selectUserNameFromClaim(ferry *api.TokenFerry, claimName string) (string, error) {
	if claimValue, ok := ferry.Creds.IDToken.Claims[claimName]; ok {
		if claimValueString, ok := claimValue.(string); ok {
			return claimValueString, nil
		}

		return "", errors.Errorf("claim %s is not a string: %v", claimName, claimValue)
	}

	return "", errors.Errorf("no such claim: %s", claimName)
}

func patchClientConfig(config *clientcmdapi.Config, userName string, ferry *api.TokenFerry) {
	if config.AuthInfos == nil {
		config.AuthInfos = make(map[string]*clientcmdapi.AuthInfo)
	}

	config.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		AuthProvider: &clientcmdapi.AuthProviderConfig{
			Name: "oidc",
			Config: map[string]string{
				"client-id":      ferry.ClientID,
				"client-secret":  ferry.ClientSecret,
				"id-token":       ferry.Creds.IDToken.Value,
				"idp-issuer-url": ferry.Creds.IDToken.Issuer,
				"refresh-token":  ferry.Creds.RefreshToken,
			},
		},
	}
}
