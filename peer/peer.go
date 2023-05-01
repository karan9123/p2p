/**
 * @file peer.go
 * @brief This file contains the implementation of the Peer ID structure and related functions used in peer-to-peer (p2p) communication.
 */

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

// GenerateIDFromPubKey generates a peer ID from the given public key using multihash.
// @param pubKey The public key used to generate the ID.
// @return The generated peer ID and nil if successful, otherwise an error is returned.

func GenerateIDFromPubKey(pubKey cr.PubKey) (ID, error) {
	b, err := cr.MarshalPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	var alg uint64 = mh.SHA2_256

	if len(b) <= maxInlineKeyLength {
		alg = mh.IDENTITY
	}

	// Generate the multihash of the public key.
	hash, _ := mh.Sum(b, alg, -1)
	return ID(hash), nil
}

// CheckPublicKey verifies whether the peer ID was derived from the given public key.
// @param id The peer ID to check.
// @param pubKey The public key to compare against the ID.
// @return True if the ID was derived from the given public key, otherwise false.
func (id ID) CheckPublicKey(pubKey cr.PubKey) bool {
	// Generate the peer ID using the provided public key.
	oid, err := GenerateIDFromPubKey(pubKey)
	if err != nil {
		// Could not generate a public key from the given ID.
		fmt.Println("Could not generate Public Key from: ", id)
		return false
	}
	return oid == id
}

// ExtractPublicKey attempts to extract the public key from a given peer ID.
// @param id The peer ID to extract the public key from.
// @return The extracted public key and nil if successful, otherwise an error is returned.
func (id ID) ExtractPublicKey() (cr.PubKey, error) {
	decoded, err := mh.Decode([]byte(id))
	if err != nil {
		return nil, err
	}

	// Verify that the given ID contains a public key.
	if decoded.Code != mh.IDENTITY {
		return nil, errors.New(string("No Public Key for Peer:" + id))
	}

	// Extract the public key from the given ID.
	pk, err := cr.UnmarshalPublicKey(decoded.Digest)
	if err != nil {
		return nil, err
	}
	return pk, nil
}
