package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

const MnemonicPath = `/Users/ping/Dev/mnemonic`

func main() {
	cpuFile, err := os.Create("cpu.profile")
	if err != nil {
		fmt.Printf("create cpu profile error, error : %v\n", err)
		os.Exit(0)
	}

	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		return
	}
	defer pprof.StopCPUProfile()

	ch := make(chan map[string]string, 2)

	for {
		select {
		case expectData := <-ch:
			for k, v := range expectData {
				fmt.Printf("address %s  mnemonic %s\n", k, v)
				saveExpectMnemonic(k, v)
			}
			break
		case <-time.Tick(time.Second):
			for i := 0; i < 100; i++ {
				go func() {
					_, _, err := mnemonicFun(ch)
					if err != nil {

					}
				}()
			}
			fmt.Printf("new epoch\n")
		}
	}
}

func mnemonicFun(ch chan map[string]string) (mnemonic, addr string, err error) {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatal(err)
	}

	mnemonic, err = bip39.NewMnemonic(entropy)

	seed := bip39.NewSeed(mnemonic, "")
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Fatal(err)
	}
	//path := hdwallet.MustParseDerivationPath()
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		log.Fatal(err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	addr = account.Address.Hex()
	if isExpectAddr(addr[2:]) {
		ch <- map[string]string{
			addr: mnemonic,
		}
	}
	return
}

func NewEthAddress(ch chan map[string]string) (privkey, addr string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privkey = hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	addr = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	if isExpectAddr(addr[2:]) {
		ch <- map[string]string{
			addr: privkey,
		}
	}
	return
}

func saveExpectMnemonic(addr, mnemonic string) {
	f, err := os.OpenFile(MnemonicPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("mnemonic: %s \t address: %s\n", mnemonic, addr))
}

func isExpectAddr(address string) bool {
	if strings.HasPrefix(address, "00000000") ||
		strings.HasSuffix(address, "00000000") ||
		strings.HasPrefix(address, "88888888") ||
		strings.HasSuffix(address, "88888888") ||
		strings.HasPrefix(address, "99999999") ||
		strings.HasSuffix(address, "99999999") ||
		strings.HasPrefix(address, "66666666") ||
		strings.HasSuffix(address, "66666666") ||
		strings.HasPrefix(address, "8888") ||
		strings.HasSuffix(address, "8888") ||
		strings.HasPrefix(address, "9999") ||
		strings.HasSuffix(address, "9999") ||
		strings.HasPrefix(address, "0000") ||
		strings.HasSuffix(address, "0000") ||
		strings.HasPrefix(address, "666666") ||
		strings.HasSuffix(address, "666666") ||
		strings.HasSuffix(address, "6666") ||
		strings.HasPrefix(address, "6666") {
		return true
	}
	return false
}
