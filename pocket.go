package wtsc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/hashicorp/go-cleanhttp"
	pocketGoProvider "github.com/pokt-foundation/pocket-go/provider"
	pocketGoSigner "github.com/pokt-foundation/pocket-go/signer"
	pocketGoUtils "github.com/pokt-foundation/pocket-go/utils"
	pocketCoreCodec "github.com/pokt-network/pocket-core/codec"
	pocketCoreCodecTypes "github.com/pokt-network/pocket-core/codec/types"
	pocketCoreCrypto "github.com/pokt-network/pocket-core/crypto"
	pocketCoreTypes "github.com/pokt-network/pocket-core/types"
	pocketCoreTypesModule "github.com/pokt-network/pocket-core/types/module"
	pocketCoreApps "github.com/pokt-network/pocket-core/x/apps"
	pocketCoreAuth "github.com/pokt-network/pocket-core/x/auth"
	pocketCoreAuthTypes "github.com/pokt-network/pocket-core/x/auth/types"
	pocketCoreGov "github.com/pokt-network/pocket-core/x/gov"
	pocketCoreNodes "github.com/pokt-network/pocket-core/x/nodes"
	pocketCoreNodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketCore "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/pokt-scan/wtsc/generated"
	"github.com/rs/zerolog/log"
	cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"math"
	"math/big"
	"strconv"
	"time"
)

func Codec() *pocketCoreCodec.Codec {
	if PocketCoreCodec == nil {
		MakeCodec()
	}
	return PocketCoreCodec
}

func MakeCodec() {
	// create a new codec
	PocketCoreCodec = pocketCoreCodec.NewCodec(pocketCoreCodecTypes.NewInterfaceRegistry())
	// register all the app module types
	pocketCoreTypesModule.NewBasicManager(
		pocketCoreApps.AppModuleBasic{},
		pocketCoreAuth.AppModuleBasic{},
		pocketCoreGov.AppModuleBasic{},
		pocketCoreNodes.AppModuleBasic{},
		pocketCore.AppModuleBasic{},
	).RegisterCodec(PocketCoreCodec)
	// register the sdk types
	pocketCoreTypes.RegisterCodec(PocketCoreCodec)
	// register the crypto types
	pocketCoreCrypto.RegisterAmino(PocketCoreCodec.AminoCodec().Amino)
	cryptoamino.RegisterAmino(PocketCoreCodec.AminoCodec().Amino)
	pocketCoreCodec.RegisterEvidences(PocketCoreCodec.AminoCodec(), PocketCoreCodec.ProtoCodec())
}

func IsValidChainPool(chains []string) bool {
	if len(chains) == 0 {
		return false
	}

	for _, chain := range chains {
		if e := pocketCoreNodesTypes.ValidateNetworkIdentifier(chain); e != nil {
			return false
		}
	}
	return true
}

func IsValidServicerList(servicerList []string) bool {
	if len(servicerList) == 0 {
		return false
	}

	for _, servicer := range servicerList {
		if !pocketGoUtils.ValidatePrivateKey(servicer) {
			return false
		}
	}

	return true
}

func IsValidMinServiceStake(minServiceStakeList MinServiceStake) bool {
	if len(minServiceStakeList) == 0 {
		return false
	}

	for _, minServiceStake := range minServiceStakeList {
		if err := pocketCoreNodesTypes.ValidateNetworkIdentifier(minServiceStake.Service); err != nil {
			return false
		}
	}

	return true
}

func StakeServicer(
	signer *pocketGoSigner.Signer,
	servicer *generated.GetWhatToStakeGetWhatToStakeWtsOptimizationResponseServicersWtsStakeNode,
) func() {
	return func() {
		// stake
		log.Debug().Str("address", signer.GetAddress()).Msg("reading node from rpc")
		node, err := PocketRpcProvider.GetNode(servicer.Address, &pocketGoProvider.GetNodeOptions{Height: 0})
		if err != nil {
			log.Error().Err(err).Str("address", servicer.Address).Msg("failed to get pocket node")
			return
		}
		nodeTokens, err := strconv.ParseInt(node.Tokens, 10, 64)
		if err != nil {
			log.Error().Err(err).Str("address", servicer.Address).Str("tokens", node.Tokens).Msg("failed to parse pocket node tokens")
			return
		}

		// --- @NOTE: this should be the required code using pocket-go package but fail due to some unnecessary imports
		//txBuilder := pocketGoTxBuilder.NewTransactionBuilder(pocketRpcProvider, signer)
		//
		//txMsg, err := pocketGoTxBuilder.NewStakeNode(
		//	signer.GetPublicKey(),
		//	node.ServiceURL,
		//	node.OutputAddress,
		//	servicer.Chains,
		//	nodeTokens,
		//)
		//txOptions := pocketGoTxBuilder.TransactionOptions{
		//	Memo: AppConfig.TxMemo,
		//	Fee:  txFee,
		//}
		//
		//result, err := txBuilder.SubmitWithCtx(ctx, AppConfig.NetworkID, txMsg, &txOptions)
		// --- @END

		// Instead I basically copy&paste just importing the codec things to allows handle the transaction properly.
		if err != nil {
			log.Error().Err(err).Msg("failed to create stake node tx message")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// value is already validated
		txFee, _ := AppConfig.TxFee.Int64()

		feeStruct := pocketCoreTypes.Coins{
			pocketCoreTypes.Coin{
				Amount: pocketCoreTypes.NewInt(txFee),
				Denom:  "upokt",
			},
		}

		entropy, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			log.Error().Err(err).Msg("failed to generate entropy")
			return
		}

		cryptoPublicKey, err := pocketCoreCrypto.NewPublicKey(signer.GetPublicKey())
		if err != nil {
			log.Error().Err(err).Msg("failed to create crypto")
			return
		}

		decodedAddress, err := hex.DecodeString(node.OutputAddress)
		if err != nil {
			log.Error().Err(err).Msg("failed to decode output address")
		}

		txMsg := &pocketCoreNodesTypes.MsgStake{
			PublicKey:  cryptoPublicKey,
			Chains:     servicer.Services, // aka chains on morse
			Value:      pocketCoreTypes.NewInt(nodeTokens),
			ServiceUrl: node.ServiceURL,
			Output:     decodedAddress,
		}

		signBytes, err := pocketCoreAuth.StdSignBytes(AppConfig.NetworkID, entropy.Int64(), feeStruct, txMsg, AppConfig.TxMemo)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}

		signature, err := signer.SignBytes(signBytes)
		if err != nil {
			log.Error().Err(err).Msg("failed to sign transaction")
			return
		}

		signatureStruct := pocketCoreAuthTypes.StdSignature{PublicKey: cryptoPublicKey, Signature: signature}
		tx := pocketCoreAuthTypes.NewTx(txMsg, feeStruct, signatureStruct, AppConfig.TxMemo, entropy.Int64())

		txBytes, err := pocketCoreAuth.DefaultTxEncoder(Codec())(tx, -1)

		signedTX := hex.EncodeToString(txBytes)

		sendTransactionInput := &pocketGoProvider.SendTransactionInput{
			Address:     signer.GetAddress(),
			RawHexBytes: signedTX,
		}

		txResult, txErr := PocketRpcProvider.SendTransactionWithCtx(ctx, sendTransactionInput)

		if txErr != nil {
			log.Error().Err(txErr).Msg("failed to submit stake node transaction")
			return
		}

		log.Info().
			Str("address", signer.GetAddress()).
			Strs("chains", servicer.Services).
			Str("height", txResult.Height).
			Str("hash", txResult.Txhash).
			Str("raw_log", txResult.RawLog).
			Msg("successfully submitted stake node transaction")
	}
}

func NewPocketRpcProvider(url string, maxRetries, maxTimeout uint) {
	// create a pocket rpc provider to reuse it
	PocketRpcProvider = pocketGoProvider.NewProvider(url)
	PocketRpcProvider.UpdateRequestConfig(pocketGoProvider.RequestConfigOpts{
		Retries:   int(maxRetries),
		Timeout:   time.Duration(maxTimeout) * time.Millisecond,
		Transport: cleanhttp.DefaultPooledTransport(),
	})
}
