//go:generate abigen --sol contract.sol --pkg main --out abi.go

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/rocksideio/rockside-sdk-go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"golang.org/x/crypto/acme/autocert"
)

var (
	portFlag            = flag.String("port", "9000", "Listening port")
	domainFlag          = flag.String("tls-domain", "", "Activates TLS on the given domain using acme (i.e. Let's Encrypt)")
	mainnetFlag         = flag.Bool("mainnet", false, "Use Ethereum mainnet instead of testnet")
	identityAddressFlag = flag.String("identity", os.Getenv("ROCKSIDE_IDENTITY"), "Rockside identity address")

	rocksideClient          *rockside.Client
	rocksideIdentityAddress common.Address
	contractAddress         = common.HexToAddress("0x5C2c0cBC0982c3545e16dd0d98F3F31f0eD2B22F") // default to testnet contract address
)

func main() {
	log.SetFlags(log.Ltime)
	flag.Parse()

	if key, found := os.LookupEnv("ROCKSIDE_API_KEY"); found {
		network := rockside.Testnet
		if *mainnetFlag {
			network = rockside.Mainnet
			contractAddress = common.HexToAddress("0x97cA295E85c997F7286F06E1d98B0939ff0D8aAA")
		}
		var err error
		rocksideClient, err = rockside.NewClientFromAPIKey(key, network)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("missing ROCKSIDE_API_KEY env variable")
	}

	identityAddr, err := hexutil.Decode(*identityAddressFlag)
	if err != nil {
		log.Fatalf("invalid identity address '%s': %s", *identityAddressFlag, err)
	}
	rocksideIdentityAddress = common.BytesToAddress(identityAddr)

	found, err := rocksideClient.Identities.Exists(rocksideIdentityAddress)
	if err != nil {
		log.Fatal(err)
	}
	if !found {
		log.Fatalf("cannot find %s in your Rockside identities listing", rocksideIdentityAddress.String())
	}

	if domain := *domainFlag; len(domain) > 0 {
		log.Printf("service starting on TLS port for domain %s", domain)
		log.Fatal(http.Serve(autocert.NewListener(domain), routes()))
	} else {
		log.Printf("service starting on %s port", *portFlag)
		log.Fatal(http.ListenAndServe(":"+*portFlag, routes()))
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.Handle("/upload", rateLimit(upload))
	mux.Handle("/lookup", rateLimit(lookup))
	return mux
}

func lookup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	fingerprint, err := getFingerprintFromQuery(r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contract, err := NewContractCaller(contractAddress, rocksideClient.RPCClient)
	if err != nil {
		log.Printf("cannot build contract caller: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timestamp, err := contract.Fingerprints(&bind.CallOpts{}, toByte32([]byte(fingerprint)))
	if err != nil {
		log.Printf("cannot build transactor: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v := struct {
		Fingerprint string
		Date        time.Time
	}{
		Fingerprint: fmt.Sprintf("%x", fingerprint),
		Date:        time.Unix(timestamp.Int64(), 0),
	}

	if err := json.NewEncoder(w).Encode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	fingerprint, err := digestSingleFileFromMultipart(r)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rocksideTransactor := rockside.NewTransactor(rocksideIdentityAddress, rocksideClient)
	contract, err := NewContractTransactor(
		contractAddress,
		rocksideTransactor,
	)
	if err != nil {
		log.Printf("cannot build transactor: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := contract.Register(rockside.TransactOptsWithoutReward(), fingerprint)
	if err != nil {
		log.Printf("cannot register: %s", err)
		http.Error(w, fmt.Sprintf("cannot register fingerprint on the blockchain: %s", err), http.StatusInternalServerError)
		return
	}

	txHash := rocksideTransactor.ReturnRocksideTransactionHash(tx.Hash())

	v := struct {
		Fingerprint     string
		TransactionHash string
		HashURL         string
	}{
		Fingerprint:     fmt.Sprintf("%x", fingerprint),
		TransactionHash: txHash,
		HashURL:         fmt.Sprintf("%s/tx/%s", rocksideClient.CurrentNetwork().ExplorerURL(), txHash),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")

	if err := enc.Encode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Printf("cannot execute template: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func getFingerprintFromQuery(r *http.Request) ([]byte, error) {
	hexFingerprint := r.FormValue("fingerprint")
	if len(hexFingerprint) != 64 {
		return nil, errors.New("fingerprint in query has invalid format")
	}

	fingerprint, err := hex.DecodeString(hexFingerprint)
	if err != nil {
		return nil, fmt.Errorf("cannot decode fingerprint: %s", err)
	}
	return fingerprint, nil
}

func digestSingleFileFromMultipart(r *http.Request) ([32]byte, error) {
	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		return [32]byte{}, err
	}

	var files []multipart.File
	for _, headers := range r.MultipartForm.File {
		for _, hh := range headers {
			f, err := hh.Open()
			if err != nil {
				return [32]byte{}, err
			}
			files = append(files, f)
		}
	}
	if l := len(files); l != 1 {
		return [32]byte{}, fmt.Errorf("expected single file in multipart POST, but got %d", l)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, files[0]); err != nil {
		return [32]byte{}, err
	}

	return toByte32(hasher.Sum(nil)), nil
}

func toByte32(b []byte) [32]byte {
	var out [32]byte
	copy(out[:], b)
	return out
}
