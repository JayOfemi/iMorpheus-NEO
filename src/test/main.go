package main

import (
	"blockChainMorp/accounts"
	"fmt"
	"flag"
	"os"
	"log"
	"github.com/visualfc/goqt/ui"
	"blockChainMorp/gui"
	"blockChainMorp/client"
)

var noGUIFlag = true

// Run parses command line arguments and processes commands
func main() {
	if !checkFileIsExist(accounts.AccountFile) {
		client.CreateAccount()
		client.CreateBitcoinAccount()
	} else {
		fmt.Println("Accounts file already exists, you can create other accounts by cmd")
		fmt.Println("The exist account address list:")
		client.ListAddresses()
		fmt.Println("The exist account Bitcoin address list:")
		client.ListBTCAddresses()
	}

	if validateArgs() {
		getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
		createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
		createAccountCmd := flag.NewFlagSet("createaccount", flag.ExitOnError)
		listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
		printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
		reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
		sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
		startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
		showGUIFlag := flag.NewFlagSet("gui", flag.ExitOnError)

		getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
		createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
		sendFrom := sendCmd.String("from", "", "Source wallet address")
		sendTo := sendCmd.String("to", "", "Destination wallet address")
		sendAmount := sendCmd.Int("amount", 0, "Amount to send")
		sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
		startNodeMiner := startNodeCmd.String("miner", "", "Enable mining mode and send reward to ADDRESS")
		useBTCFlagCmd := createAccountCmd.Bool("usebtc", false, "Create account of Bitcoin")

		switch os.Args[1] {
		case "getbalance":
			err := getBalanceCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createblockchain":
			err := createBlockchainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "createaccount":
			err := createAccountCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "listaddresses":
			err := listAddressesCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "reindexutxo":
			err := reindexUTXOCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "send":
			err := sendCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "startnode":
			err := startNodeCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		case "gui":
			err := showGUIFlag.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}

		default:
			fmt.Println("Error reference!")
			os.Exit(1)
		}

		if getBalanceCmd.Parsed() {
			if *getBalanceAddress == "" {
				getBalanceCmd.Usage()
				os.Exit(1)
			}
			client.GetBalance(*getBalanceAddress)
			noGUIFlag = true
		}

		if createBlockchainCmd.Parsed() {
			if *createBlockchainAddress == "" {
				createBlockchainCmd.Usage()
				os.Exit(1)
			}
			if !accounts.ValidateAddress(*createBlockchainAddress) {
				log.Panic("ERROR: Address is not valid")
			}
			client.CreateBlockchain(*createBlockchainAddress)
			noGUIFlag = true
		}

		if createAccountCmd.Parsed() {
			if *useBTCFlagCmd {
				client.CreateBitcoinAccount()
			} else {
				client.CreateAccount()
			}
			noGUIFlag = true
		}

		if listAddressesCmd.Parsed() {
			client.ListAddresses()
			noGUIFlag = true
		}

		if printChainCmd.Parsed() {
			client.PrintChain()
			noGUIFlag = true
		}

		if reindexUTXOCmd.Parsed() {
			noGUIFlag = true
			client.ReindexUTXO()
		}

		if sendCmd.Parsed() {
			if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
				sendCmd.Usage()
				os.Exit(1)
			}

			client.Send(*sendFrom, *sendTo, *sendAmount, *sendMine, false)
			noGUIFlag = true
		}

		if startNodeCmd.Parsed() {
			client.StartNode(*startNodeMiner, false)
			noGUIFlag = true
		}

		if showGUIFlag.Parsed() {
			noGUIFlag = false
		}
	}

	if !noGUIFlag {
		ui.RunEx(os.Args, func() {
			w, err := gui.NewMainWindow()
			if err != nil {
				log.Fatalln(err)
			}
			w.Show()
		})
	}
}

func validateArgs() bool {
	if len(os.Args) < 2 {
		return false
	} else {
		return true
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}