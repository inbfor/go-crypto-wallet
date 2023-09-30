package manipulations

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func WalletCreation() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	prvKEy := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])

	return hexutil.Encode(hash.Sum(nil)[12:]), prvKEy, nil
}

func ShowBalance(ethcl *ethclient.Client, addr string) (string, error) {

	account := common.HexToAddress(addr)

	balance, err := ethcl.BalanceAt(context.Background(), account, nil)

	if err != nil {
		return "111", err
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	amount := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(amount)

	fmtString := fmt.Sprintf("%%.%df", 18)
	balanceEth := fmt.Sprintf(strings.TrimRight(strings.TrimRight(fmt.Sprintf(fmtString, amount), "0"), "."))

	fmt.Println(balanceEth)

	return balanceEth, nil
}

func TransferEthereum(ethcl *ethclient.Client, from string, to string, privateKey string) error {
	privateKeyCr, _ := crypto.HexToECDSA(privateKey)
	fromAddr := common.HexToAddress(from)
	toAddr := common.HexToAddress(to)

	nonce, err := ethcl.PendingNonceAt(context.Background(), fromAddr)

	if err != nil {
		return err
	}

	gasPrice, err := ethcl.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(100000000000000000)

	gasLimit := uint64(21000)

	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, nil)

	chainID, err := ethcl.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyCr)
	if err != nil {
		return err
	}

	err = ethcl.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	return nil
}
