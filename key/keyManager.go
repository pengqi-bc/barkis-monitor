package key

import (
	"fmt"
	"strings"

	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/auth"
	"github.com/cosmos/go-bip39"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	defaultBIP39Passphrase = ""
)

type KeyManager interface {
	Sign(auth.StdSignMsg) ([]byte, error)
	GetPrivKey() crypto.PrivKey
	GetAddr() sdk.AccAddress
}

func NewMnemonicKeyManager(mnemonic string) (KeyManager, error) {
	k := keyManager{}
	err := k.recoveryFromMnemonic(mnemonic, FullPath)
	return &k, err
}

type keyManager struct {
	privKey  crypto.PrivKey
	addr     sdk.AccAddress
	mnemonic string
}

func (m *keyManager) recoveryFromMnemonic(mnemonic, keyPath string) error {
	words := strings.Split(mnemonic, " ")
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("mnemonic length should either be 12 or 24")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, defaultBIP39Passphrase)
	if err != nil {
		return err
	}
	// create master key and derive first key:
	masterPriv, ch := ComputeMastersFromSeed(seed)
	derivedPriv, err := DerivePrivateKeyForPath(masterPriv, ch, keyPath)
	if err != nil {
		return err
	}
	priKey := secp256k1.PrivKeySecp256k1(derivedPriv)
	addr := sdk.AccAddress(priKey.PubKey().Address())
	if err != nil {
		return err
	}
	m.addr = addr
	m.privKey = priKey
	m.mnemonic = mnemonic
	return nil
}

func (m *keyManager) Sign(msg auth.StdSignMsg) ([]byte, error) {
	return m.privKey.Sign(msg.Bytes())
}

func (m *keyManager) GetPrivKey() crypto.PrivKey {
	return m.privKey
}

func (m *keyManager) GetAddr() sdk.AccAddress {
	return m.addr
}
