package peer

import (
	"errors"
	"fmt"
	mh "github.com/multiformats/go-multihash"
	cr "p2p/crypto"
)

// ID is a p2p peer identity.
// Peer IDs are derived by hashing a peer's public key and encoding the hash output as multihash.
type ID string

const maxInlineKeyLength = 42

func GenerateIDFromPubKey(pubKey cr.PubKey) (ID, error) {
	b, err := cr.MarshalPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	var alg uint64 = mh.SHA2_256

	if len(b) <= maxInlineKeyLength {
		alg = mh.IDENTITY
	}
	hash, _ := mh.Sum(b, alg, -1)
	return ID(hash), nil
}

// CheckPublicKey checks whether this ID was derived from pubKey.
func (id ID) CheckPublicKey(pubKey cr.PubKey) bool {
	oid, err := GenerateIDFromPubKey(pubKey)
	if err != nil {
		fmt.Println("Could not generate Public Key from: ", id)
		return false
	}
	return oid == id
}

// ExtractPublicKey attempts to extract the public key from an ID.
func (id ID) ExtractPublicKey() (cr.PubKey, error) {
	decoded, err := mh.Decode([]byte(id))
	if err != nil {
		return nil, err
	}
	if decoded.Code != mh.IDENTITY {
		return nil, errors.New(string("No Public Key for Peer:" + id))
	}
	pk, err := cr.UnmarshalPublicKey(decoded.Digest)
	if err != nil {
		return nil, err
	}
	return pk, nil
}
