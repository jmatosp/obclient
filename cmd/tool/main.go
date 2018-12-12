package main

import (
	"fmt"
	"github.com/jmatosp/ob_security/aspsp"
	"github.com/jmatosp/ob_security/authorization"
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
		},
	}

	clientDetails := &cobra.Command{
		Use:   "client",
		Short: "View software client details",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cliBanner)
			storer := aspsp.NewClientStorer(storageFolder)
			client, err := storer.Get()
			if err == aspsp.ErrNotFound {
				fmt.Println("This software client is not registered yet, register first.")
				os.Exit(1)
			} else if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			fmt.Println("Client details")
			fmt.Printf("Client Id: %s\n", client.Id)
		},
	}
	//
	//tokenDetails := &cobra.Command{
	//	Use:   "token",
	//	Short: "View token details",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(cliBanner)
	//		storer := aspsp.NewFileTokenStorer(storageFolder)
	//		token, err := storer.Get()
	//		if err == aspsp.ErrNotFound {
	//			fmt.Println("Token does not exist, auth first to get a token")
	//			os.Exit(1)
	//		} else if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//		fmt.Println("Token details")
	//		fmt.Printf("Token: %s\n", token.Id)
	//		fmt.Printf("Expires: %d\n", token.ExpiresIn)
	//	},
	//}
	//
	//testTransport := &cobra.Command{
	//	Use:   "transport",
	//	Short: "Test certificates making a MTLS call to token",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(cliBanner)
	//		tester := aspsp.NewTransportTester(makeSecuredTransport(), configuration.TokenEndpoint)
	//		fmt.Println("Test certificates details")
	//		if err := tester.Test(); err != nil {
	//			fmt.Println(err)
	//			os.Exit(1)
	//		}
	//		fmt.Println("OK")
	//	},
	//}
	//
	//ssDetails := &cobra.Command{
	//	Use:   "statement",
	//	Short: "View software statement details",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(cliBanner)
	//		softwareStatement := makeSoftwareStatement()
	//		fmt.Println("Software Statement")
	//		fmt.Printf("Id: %s\n", softwareStatement.Id())
	//		fmt.Printf("Name: %s\n", softwareStatement.Name())
	//	},
	//}
	//
	//configurationCmd := &cobra.Command{
	//	Use:   "configuration",
	//	Short: "View ASPSP openid configuration",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(cliBanner)
	//		config := configuration
	//		fmt.Println("Ozone OpenId Configuration")
	//		fmt.Printf("TokenEndpoint: %s\n", config.TokenEndpoint)
	//		fmt.Printf("AuthorizationEndpoint: %s\n", config.AuthorizationEndpoint)
	//		fmt.Printf("RegistrationEndpoint: %s\n", config.RegistrationEndpoint)
	//		fmt.Printf("Issuer: %s\n", config.Issuer)
	//	},
	//}
	//
	//authorize := &cobra.Command{
	//	Use:   "auth",
	//	Short: "Authorize flow to use ASPSP services",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(cliBanner)
	//		fmt.Println("Authorize")
	//
	//		storer := aspsp.NewClientStorer(storageFolder)
	//		client, err := storer.Get()
	//		if err == aspsp.ErrNotFound {
	//			fmt.Println("This software client is not registered yet, register first.")
	//			os.Exit(1)
	//		} else if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//
	//		authenticator := aspsp.NewAuthenticator(
	//			makeCredentialsGranter(configuration, client),
	//			makeAccessConsenter(),
	//			makePSUAccessConsenter(client),
	//			makeTokenGenerator(configuration, client),
	//		)
	//		token, err := authenticator.Authenticate()
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//
	//		tokenStorer := aspsp.NewFileTokenStorer(storageFolder)
	//		err = tokenStorer.Store(token)
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//
	//		fmt.Println("Got valid token")
	//	},
	//}
	//
	//accountsCmd := &cobra.Command{
	//	Use:   "accounts",
	//	Short: "List accounts",
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println(cliBanner)
	//		fmt.Println("Accounts")
	//
	//		tokenStorer := aspsp.NewFileTokenStorer(storageFolder)
	//		token, err := tokenStorer.Get()
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//
	//		accountLister := makeAccountLister(token)
	//		_, err = accountLister.List()
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//
	//		accounts, err := accountLister.List()
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//
	//		aspsp.NewAccountsPrinter().Print(accounts)
	//	},
	//}

	rootCmd.AddCommand(clientRegister)
	rootCmd.AddCommand(clientDetails)
	//rootCmd.AddCommand(ssDetails)
	//rootCmd.AddCommand(configurationCmd)
	//rootCmd.AddCommand(testTransport)
	//rootCmd.AddCommand(authorize)
	//rootCmd.AddCommand(tokenDetails)
	//rootCmd.AddCommand(accountsCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//
//var transport aspsp.Transport
//
//
//var softwareStatement aspsp.SoftwareStatement
//
//
//func makeCredentialsGranter(configuration aspsp.Configuration, client aspsp.Client) aspsp.CredentialsGranter {
//	return aspsp.NewCredentialGrander(
//		makeSecuredTransport(),
//		configuration.TokenEndpoint,
//		client,
//	)
//}
//
//func makeAccessConsenter() aspsp.AccessConsenter {
//	return aspsp.NewAccessConsenter(
//		makeSecuredTransport(),
//		viper.GetString("endpoints"),
//		viper.GetString("fapiFinancialId"),
//	)
//}
//
//func makePSUAccessConsenter(client aspsp.Client) aspsp.PSUAccessConsenter {
//	return aspsp.NewPSUAccessConsenter(
//		makeSecuredTransport(),
//		viper.GetString("endpoints"),
//		viper.GetString("redirectUrl"),
//		client,
//	)
//}
//
//func makeTokenGenerator(configuration aspsp.Configuration, client aspsp.Client) aspsp.TokenGenerator {
//	return aspsp.NewTokenGenerator(
//		makeSecuredTransport(),
//		configuration.TokenEndpoint,
//		viper.GetString("redirectUrl"),
//		client,
//	)
//}
//
//func makeAccountLister(token aspsp.Token) aspsp.AccountLister {
//	return aspsp.NewAccountLister(
//		makeSecuredTransport(),
//		viper.GetString("endpoints"),
//		viper.GetString("fapiFinancialId"),
//		token,
//	)
//}

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
