package main

import (
	"fmt"
	"github.com/jmatosp/obclient/aspsp"
	"github.com/jmatosp/obclient/authorization"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const cliBanner = "Open Banking CLI v0.0.1"

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	storageFolder := viper.GetString("storageFolder")

	rootCmd := &cobra.Command{Use: "obcli"}

	clientRegister := &cobra.Command{
		Use:   "register",
		Short: "Dynamic register a new software client",
		Run: func(cmd *cobra.Command, args []string) {
			clientRegister(storageFolder)
		},
	}

	authorize := &cobra.Command{
		Use:   "auth",
		Short: "Authorize flow to use ASPSP services",
		Run: func(cmd *cobra.Command, args []string) {
			authorize(storageFolder)
		},
	}

	accountsCmd := &cobra.Command{
		Use:   "accounts",
		Short: "List accounts",
		Run: func(cmd *cobra.Command, args []string) {
			accountsList(storageFolder)
		},
	}

	rootCmd.AddCommand(clientRegister)
	rootCmd.AddCommand(authorize)
	rootCmd.AddCommand(accountsCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func accountsList(storageFolder string) {
	fmt.Println(cliBanner)
	fmt.Println("Accounts")
	tokenStorer := aspsp.NewFileTokenStorer(storageFolder)
	token, err := tokenStorer.Get()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	accountLister := makeAccountLister(token)
	_, err = accountLister.List()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	accounts, err := accountLister.List()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	aspsp.NewAccountsPrinter().Print(accounts)
}

func authorize(storageFolder string) {
	fmt.Println(cliBanner)
	fmt.Println("Authorize")
	storer := aspsp.NewClientStorer(storageFolder)
	client, err := storer.Get()
	if err == aspsp.ErrNotFound {
		fmt.Println("This software client is not registered yet, register first.")
		os.Exit(1)
	} else if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	authenticator, err := makeAuthenticator(client)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	token, err := authenticator.Authenticate()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	tokenStorer := aspsp.NewFileTokenStorer(storageFolder)
	err = tokenStorer.Store(token)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Got valid token")
}

func clientRegister(storageFolder string) {
	fmt.Println(cliBanner)
	storer := aspsp.NewClientStorer(storageFolder)
	_, err := storer.Get()
	if err == aspsp.ErrNotFound {
		register, err := makeClientRegister()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		client, err := register.Register()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		storer.Store(client)
		fmt.Printf("Client registered id: %s\n", client.Id)
		os.Exit(0)

	} else if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Client already registered, delete first to recreate")
	os.Exit(1)
}

func makeAccountLister(token authorization.Token) aspsp.AccountLister {
	return aspsp.NewAccountLister(
		makeSecuredTransport(),
		viper.GetString("endpoints"),
		viper.GetString("fapiFinancialId"),
		token,
	)
}

func makeSecuredTransport() authorization.Transport {
	return authorization.NewSecureTransport(
		viper.GetString("cerFile"),
		viper.GetString("keyFile"),
		viper.GetStringSlice("rootCAs"),
	)
}

func makeClientRegister() (authorization.ClientRegister, error) {
	return authorization.NewClientRegisterBuilder().
		WithWellKnown(viper.GetString("openidConfiguration")).
		WithSigPublicKeyFile(viper.GetString("sigPublicKeyFile")).
		WithSigPrivateKeyFile(viper.GetString("sigPrivateKeyFile")).
		WithCertFile(viper.GetString("cerFile")).
		WithKeyFile(viper.GetString("keyFile")).
		WithRootCAs(viper.GetStringSlice("rootCAs")).
		WithRedirectUrl(viper.GetString("redirectUrl")).
		WithSoftwareStatementID(viper.GetString("softwareStatementID")).
		WithSoftwareStatementName(viper.GetString("softwareStatementName")).
		Build()
}

func makeAuthenticator(client authorization.Client) (authorization.Authenticator, error) {
	return authorization.NewAuthenticatorBuilder().
		WithWellKnown(viper.GetString("openidConfiguration")).
		WithClient(client).
		WithFapiFinancialId(viper.GetString("fapiFinancialId")).
		WithAccessConsentEndpoint(viper.GetString("endpoints")).
		WithCertFile(viper.GetString("cerFile")).
		WithKeyFile(viper.GetString("keyFile")).
		WithRootCAs(viper.GetStringSlice("rootCAs")).
		WithRedirectUrl(viper.GetString("redirectUrl")).
		Build()
}
