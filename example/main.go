package main

import (
	"fmt"
	"github.com/jason-cn-dev/xuper-sdk-go/account"
	"github.com/jason-cn-dev/xuper-sdk-go/balance"
	"github.com/jason-cn-dev/xuper-sdk-go/contract"
	"github.com/jason-cn-dev/xuper-sdk-go/contract_account"
	"github.com/jason-cn-dev/xuper-sdk-go/crypto"
	"github.com/jason-cn-dev/xuper-sdk-go/network"
	"github.com/jason-cn-dev/xuper-sdk-go/transfer"
	hdapi "github.com/xuperchain/crypto/core/hdwallet/api"
)

var (
	node                  = "127.0.0.1:37101"
	bcname                = "xuper"
	mnemonic              = "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南"
	language              = 1
	strength        uint8 = 1
	password              = "123"
	keysPath              = "./keys/"
	to                    = "ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt"
	amount                = "10"
	fee                   = "0"
	desc                  = ""
	contractAccount       = "XC1234567890123456@xuper"
	contractName          = "counter1"
	runtime               = "c"
	codePath              = "example/contract_code/counter.wasm"
	txid                  = "b75825ca9b8847126022dceb5e8960d936511db4d952fe821c0dbd861df49eee"
	acc             *account.Account
)

func init() {
	var err error
	acc, err = account.RetrieveAccount(mnemonic, language)
	if err != nil {
		panic(err)
	}
}

func CreateAccount() {
	// create an account for the user,
	// strength 1 means that the number of mnemonics is 12
	// language 1 means that mnemonics is Chinese
	acc, err := account.CreateAccount(strength, language)
	if err != nil {
		panic(err)
	}

	fmt.Println("CreateAccount Address:", acc.Address)
	fmt.Println("CreateAccount Mnemonic:", acc.Mnemonic)

	// retrieve the account by mnemonics
	acc, err = account.RetrieveAccount(mnemonic, language)
	if err != nil {
		panic(err)
	}
	fmt.Println("RetrieveAccount Address:", acc.Address)

	// create an account, then encrypt using password and save it to a file
	acc, err = account.CreateAndSaveAccountToFile(keysPath, password, strength, language)
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateAndSaveAccountToFile Address:", acc.Address)

	// get the account from file, using password decrypt
	acc, err = account.GetAccountFromFile(keysPath, password)
	if err != nil {
		panic(err)
	}
	fmt.Println("GetAccountFromFile Address:", acc.Address)
}

func CreateContractAccount() {
	// define the name of the conrtact account to be created
	// conrtact account is (XC + 16 length digit + @xuper), like "XC1234567890123456@xuper"

	// initialize a client to operate the contract account
	ca := contractaccount.InitContractAccount(acc, node, bcname)

	// create contract account
	txid, err := ca.CreateContractAccount(contractAccount)
	if err != nil {
		panic(err)
	}

	/*
		// the 2nd way to create contract account
		preSelectUTXOResponse, err := ca.PreCreateContractAccount(contractAccount)
		if err != nil {
			panic(err)
		}
		txid, err := ca.PostCreateContractAccount(preSelectUTXOResponse)
		if err != nil {
			panic(err)
		}
	*/
	fmt.Println("CreateContractAccount txid:", txid)
}

func GetBalance() {
	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	b, err := trans.GetBalance()
	if err != nil {
		panic(err)
	}
	fmt.Println("GetBalance:", b)
}

func Transfer() {
	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// transfer destination address, amount, fee and description

	// transfer
	txid, gas, err := trans.Transfer(to, amount, fee, desc)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transfer txid:", txid, "gas:", gas)
}

func TransferByPlatform() {
	// retrieve the platform account by mnemonics
	accPlatform, err := account.RetrieveAccount(mnemonic, language)
	if err != nil {
		panic(err)
	}
	fmt.Println("platform account:", accPlatform.Address)

	// initialize a client to operate the transfer transaction
	trans := transfer.InitTransByPlatform(acc, accPlatform, node, bcname)

	// transfer destination address, amount, fee and description

	// transfer
	txid, gas, err := trans.Transfer(to, amount, fee, desc)
	if err != nil {
		panic(err)
	}
	fmt.Println("Transfer txid:", txid, "gas:", gas)
}

func DeployWasmContract() {
	// set contract account, contract will be installed in the contract account
	// initialize a client to operate the contract
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// set init args and contract file
	args := map[string]string{
		"creator": "xchain",
	}

	// deploy wasm contract
	txid, err := wasmContract.DeployWasmContract(args, codePath, runtime)
	if err != nil {
		panic(err)
	}

	/*
		// the 2nd way to deploy wasm contract, preDeploy and Post
		preSelectUTXOResponse, err := wasmContract.PreDeployWasmContract(args, codePath, "c")
		if err != nil {
			panic(err)
		}
		txid, err := wasmContract.PostWasmContract(preSelectUTXOResponse)
		if err != nil {
			panic(err)
		}
	*/
	fmt.Println("DeployWasmContract txid:", txid)
}

func InvokeWasmContract() {
	// initialize a client to operate the contract
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// set invoke function method and args
	args := map[string]string{
		"key": "counter",
	}
	methodName := "increase"

	// invoke contract
	txid, err := wasmContract.InvokeWasmContract(methodName, args)
	if err != nil {
		panic(err)
	}

	/*
		// the 2nd way to invoke wasm contract, preInvoke and Post
		preSelectUTXOResponse, err := wasmContract.PreInvokeWasmContract(methodName, args)
		if err != nil {
			log.Printf("InvokeWasmContract GetPreMethodWasmContractRes failed, err: %v", err)
			os.Exit(-1)
		}
		txid, err := wasmContract.PostWasmContract(preSelectUTXOResponse)
		if err != nil {
			log.Printf("InvokeWasmContract PostWasmContract failed, err: %v", err)
			os.Exit(-1)
		}
		log.Printf("txid: %v", txid)
	*/
	fmt.Println("InvokeWasmContract txid:", txid)
}

func QueryWasmContract() {
	// initialize a client to operate the contract
	contractAccount := ""
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// set query function method and args
	args := map[string]string{
		"key": "counter",
	}
	methodName := "get"

	// query contract
	preExeRPCRes, err := wasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		panic(err)
	}
	//gas := preExeRPCRes.GetResponse().GetGasUsed()
	//fmt.Println("gas:", gas)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Println("contract response:", string(res))
	}
}

func QueryTx() {
	// initialize a client to operate the transaction
	trans := transfer.InitTrans(nil, node, bcname)

	// query tx by txid
	tx, err := trans.QueryTx(txid)
	if err != nil {
		panic(err)
	}
	fmt.Println("QueryTx tx:", tx)
}

func EncryptedTransfer() {
	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// transfer destination address, amount, fee and description
	desc := "encrypted transfer tx"

	cryptoClient := crypto.GetCryptoClient()
	masterKey, err := cryptoClient.GenerateMasterKeyByMnemonic(mnemonic, language)
	if err != nil {
		panic(err)
	}

	privateKey, err := cryptoClient.GenerateChildKey(masterKey, hdapi.HardenedKeyStart+1)
	if err != nil {
		panic(err)
	}

	publicKey, err := cryptoClient.ConvertPrvKeyToPubKey(privateKey)
	if err != nil {
		panic(err)
	}

	// transfer
	txid, gas, err := trans.EncryptedTransfer(to, amount, fee, desc, publicKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("EncryptedTransfer txid:", txid, "gas:", gas)
}

func DecryptedTx() {
	// initialize a client to operate the transaction
	trans := transfer.InitTrans(nil, node, bcname)

	// query tx by txid
	TxStatus, err := trans.QueryTx(txid)
	if err != nil {
		panic(err)
	}
	encryptedTx := TxStatus.Tx

	xchainCryptoClient := crypto.GetXchainCryptoClient()
	masterKey, err := xchainCryptoClient.GenerateMasterKeyByMnemonic(mnemonic, language)
	if err != nil {
		panic(err)
	}

	decryptedDesc, err := trans.DecryptedTx(encryptedTx, masterKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("DecryptedTx desc:", decryptedDesc)
}

func BatchTransfer() {
	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// transfer destination address, amount, fee and description
	to1 := "alice"
	amount1 := "10"
	to2 := "bob"
	amount2 := "20"

	toAddressAndAmount := make(map[string]string)
	toAddressAndAmount[to1] = amount1
	toAddressAndAmount[to2] = amount2

	desc := "multi transfer test"

	// transfer
	txid, gas, err := trans.BatchTransfer(toAddressAndAmount, fee, desc)
	if err != nil {
		panic(err)
	}
	fmt.Println("BatchTransfer txid:", txid, "gas:", gas)
}

func BatchTransferByPlatform() {
	// retrieve the platform account by mnemonics
	accPlatform, err := account.RetrieveAccount(mnemonic, language)
	if err != nil {
		panic(err)
	}
	fmt.Println("platform account:", accPlatform.Address)

	// initialize a client to operate the transfer transaction
	trans := transfer.InitTransByPlatform(acc, accPlatform, node, bcname)

	// transfer destination address, amount, fee and description
	to1 := "alice"
	amount1 := "10"
	to2 := "bob"
	amount2 := "20"

	toAddressAndAmount := make(map[string]string)
	toAddressAndAmount[to1] = amount1
	toAddressAndAmount[to2] = amount2

	desc := "multi transfer test"

	// transfer
	txid, gas, err := trans.BatchTransfer(toAddressAndAmount, fee, desc)
	if err != nil {
		panic(err)
	}
	fmt.Println("BatchTransferByPlatform txid:", txid, "gas:", gas)
}

func CreateChain() {
	// initialize a client to operate the transfer transaction
	chain := network.InitChain(acc, node, bcname)

	// desc for creating a new blockchain

	// ./xchain-cli status -H 127.0.0.1:37801
	// ./xchain-cli account balance dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN -H 127.0.0.1:37801 --name TestChain

	// tdpos的desc demo
	//desc := `{
	//  "Module": "kernel",
	//  "Method": "CreateBlockChain",
	//  "Args": {
	//    "name": "HelloChain",
	//    "data": "{\"maxblocksize\": \"128\", \"award_decay\": {\"height_gap\": 31536000, \"ratio\": 1}, \"version\": \"1\", \"predistribution\": [{\"quota\": \"1000000000000000\", \"address\": \"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\"}], \"decimals\": \"8\", \"period\": \"3000\",\"award\": \"1000000\", \"genesis_consensus\": {\"config\": {\"init_proposer\": {\"1\": [\"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\", \"nYoKRf3jX7vhfSn4jUwHzUf5v5eVxdaNQ\", \"kGXLu6Kex54AJZcp5QPTQ5Hz4ebcUXLLB\"]}, \"timestamp\": \"1534928070000000000\", \"period\": \"500\", \"alternate_interval\": \"3000\", \"term_interval\": \"3000\", \"block_num\": \"10\", \"vote_unit_price\": \"1\", \"proposer_num\": \"3\"}, \"name\": \"tdpos\", \"type\":\"tdpos\"}}"
	//    }
	//}`

	// single的desc demo
	desc := `{
	   "Module": "kernel",
	   "Method": "CreateBlockChain",
	   "Args": {
	       "name": "TestChain",
	   	"data": "{\"version\": \"1\", \"consensus\": {\"miner\":\"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\", \"type\":\"single\"},\"predistribution\":[{\"address\": \"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\",\"quota\": \"1000000000000000\"}],\"maxblocksize\": \"128\",\"period\": \"3000\",\"award\": \"1000000\"}"
		    }
		}`

	// transfer
	txid, err := chain.CreateChain(desc)
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateChain txid:", txid)
}

func GetMultiChainBalance() {
	bcNames := []string{}
	bcNames = append(bcNames, "xuper")
	bcNames = append(bcNames, "HelloChain")

	// initialize a client to operate the transaction
	balanceUtil := balance.InitBalance(acc, node, bcNames)

	// get balance of the account
	balances, err := balanceUtil.GetBalanceDetails()
	if err != nil {
		panic(err)
	}
	fmt.Println("GetMultiChainBalance balances:", balances)
}

func main() {

	//v1.0
	//CreateAccount()
	//CreateContractAccount()
	//GetBalance()
	//Transfer()
	//DeployWasmContract()
	//InvokeWasmContract()
	//QueryWasmContract()
	//QueryTx()

	//v1.1
	//EncryptedTransfer()
	//DecryptedTx()
	//TransferByPlatform()
	//BatchTransfer()
	//BatchTransferByPlatform()
	//CreateChain()
	GetMultiChainBalance()
}
