package cmd

import (
	"fmt"
	"github.com/rebuy-de/aws-nuke/pkg/config"
	"github.com/spf13/cobra"
	"sort"

	origin "github.com/rebuy-de/aws-nuke/cmd"
	"github.com/rebuy-de/aws-nuke/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/resources"
	log "github.com/sirupsen/logrus"
)

type IamUsernameFilters []config.Filter
type IamUserPolicyAttachmentFilters []config.Filter
type IamUserAccessKeyFilters []config.Filter
type IamUsernames []string
type AccessKey string

func (names IamUsernames) toIamUsernameFilters() IamUsernameFilters {
	filters := make(IamUsernameFilters, 0, len(names))
	for _, n := range names {
		filters = append(filters, config.Filter{
			Type:  config.FilterTypeGlob,
			Value: n,
		})
	}
	return filters
}

func (names IamUsernames) toIamUserPolicyAttachmentFilters() IamUserPolicyAttachmentFilters {
	filters := make(IamUserPolicyAttachmentFilters, 0, len(names))
	for _, n := range names {
		filters = append(filters, config.Filter{
			Type:  config.FilterTypeGlob,
			Value: fmt.Sprintf("%s -> AdministratorAccess", n),
		})
	}
	return filters
}

func (names IamUsernames) toIamUserAccessKeyFilters() IamUserAccessKeyFilters {
	filters := make(IamUserAccessKeyFilters, 0, len(names))
	for _, n := range names {
		filters = append(filters, config.Filter{
			Type:     config.FilterTypeGlob,
			Property: "UserName",
			Value:    n,
		})
	}
	return filters
}

func NewRootCommand() *cobra.Command {
	var (
		accountId    string
		iamUsernames IamUsernames
		params       origin.NukeParameters
		creds        awsutil.Credentials
		verbose      bool
	)

	command := &cobra.Command{
		Use:   "aws-nuke",
		Short: "aws-nuke removes every resource from AWS",
		Long:  `A tool which removes every resource from an AWS account.  Use it with caution, since it cannot distinguish between production and non-production.`,
	}

	command.PreRun = func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.InfoLevel)
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		var err error

		err = creds.Validate()
		if err != nil {
			return err
		}

		command.SilenceUsage = true

		account, err := awsutil.NewAccount(creds)
		if err != nil {
			return err
		}

		n := origin.NewNuke(params, *account)

		fmt.Println(iamUsernames.toIamUsernameFilters())
		fmt.Println(iamUsernames.toIamUserPolicyAttachmentFilters())
		fmt.Println(iamUsernames.toIamUserAccessKeyFilters())

		n.Config = &config.Nuke{
			Accounts: map[string]config.Account{
				accountId: {
					Filters: config.Filters{
						"IAMUser":                 iamUsernames.toIamUsernameFilters(),
						"IAMUserPolicyAttachment": iamUsernames.toIamUserPolicyAttachmentFilters(),
						"IAMUserAccessKey":        iamUsernames.toIamUserAccessKeyFilters(),
						"EC2VPC": {
							config.Filter{
								Property: "IsDefault",
								Value:    "true",
							},
						},
					},
				},
			},
			AccountBlacklist: []string{""},
			Regions: []string{
				"global",
				"us-east-2",
				"us-east-1",
				"us-west-1",
				"us-west-2",
				"af-south-1",
				"ap-east-1",
				"ap-south-2",
				"ap-southeast-3",
				"ap-southeast-4",
				"ap-south-1",
				"ap-northeast-3",
				"ap-northeast-2",
				"ap-southeast-1",
				"ap-southeast-2",
				"ap-northeast-1",
				"ca-central-1",
				"eu-central-1",
				"eu-west-1",
				"eu-west-2",
				"eu-south-1",
				"eu-west-3",
				"eu-south-2",
				"eu-north-1",
				"eu-central-2",
				"me-south-1",
				"me-central-1",
				"sa-east-1",
			},
		}

		return n.Run()
	}

	command.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false,
		"Enables debug output.")

	command.PersistentFlags().StringVarP(
		&params.ConfigPath, "config", "c", "",
		"(required) Path to the nuke config file.")

	command.PersistentFlags().StringVar(
		&creds.Profile, "profile", "",
		"Name of the AWS profile name for accessing the AWS API. "+
			"Cannot be used together with --access-key-id and --secret-access-key.")
	command.PersistentFlags().StringVar(
		&creds.AccessKeyID, "access-key-id", "",
		"AWS access key ID for accessing the AWS API. "+
			"Must be used together with --secret-access-key. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&creds.SecretAccessKey, "secret-access-key", "",
		"AWS secret access key for accessing the AWS API. "+
			"Must be used together with --access-key-id. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&creds.SessionToken, "session-token", "",
		"AWS session token for accessing the AWS API. "+
			"Must be used together with --access-key-id and --secret-access-key. "+
			"Cannot be used together with --profile.")
	command.PersistentFlags().StringVar(
		&accountId, "account-id", "",
		"AWS account id that you want to run nuke on")
	command.PersistentFlags().StringSliceVar(
		(*[]string)(&iamUsernames), "iam-username", []string{}, "")

	command.PersistentFlags().StringSliceVarP(
		&params.Targets, "target", "t", []string{},
		"Limit nuking to certain resource types (eg IAMServerCertificate). "+
			"This flag can be used multiple times.")
	command.PersistentFlags().StringSliceVarP(
		&params.Excludes, "exclude", "e", []string{},
		"Prevent nuking of certain resource types (eg IAMServerCertificate). "+
			"This flag can be used multiple times.")
	command.PersistentFlags().BoolVar(
		&params.NoDryRun, "no-dry-run", true,
		"If specified, it actually deletes found resources. "+
			"Otherwise it just lists all candidates.")
	command.PersistentFlags().BoolVar(
		&params.Force, "force", true,
		"Don't ask for confirmation before deleting resources. "+
			"Instead it waits 15s before continuing. Set --force-sleep to change the wait time.")
	command.PersistentFlags().IntVar(
		&params.ForceSleep, "force-sleep", 3,
		"If specified and --force is set, wait this many seconds before deleting resources. "+
			"Defaults to 15.")
	command.PersistentFlags().IntVar(
		&params.MaxWaitRetries, "max-wait-retries", 0,
		"If specified, the program will exit if resources are stuck in waiting for this many iterations. "+
			"0 (default) disables early exit.")

	command.AddCommand(origin.NewVersionCommand())
	command.AddCommand(NewResourceTypesCommand())

	return command
}

func NewResourceTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-types",
		Short: "lists all available resource types",
		Run: func(cmd *cobra.Command, args []string) {
			names := resources.GetListerNames()
			sort.Strings(names)

			for _, resourceType := range names {
				fmt.Println(resourceType)
			}
		},
	}

	return cmd
}
