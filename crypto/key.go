// Package crypto implements various cryptographic utilities used by libp2p.
// This includes a Public and Private key interface and key implementations
// for supported key algorithms.
package crypto

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
)

type KeyType int32

const (
	// RSA is an enum for the supported RSA key type
	RSA = iota
	// Ed25519 is an enum for the supported Ed25519 key type
	Ed25519
	// ECDSA is an enum for the supported ECDSA key type
	ECDSA
)

var (
	// ErrBadKeyType is returned when a key is not supported
	ErrBadKeyType = errors.New("invalid or unsupported key type")
)

// PubKeyUnmarshaller is a func that creates a PubKey from a given slice of bytes
type PubKeyUnmarshaller func(data []byte) (PubKey, error)

// PrivKeyUnmarshaller is a func that creates a PrivKey from a given slice of bytes
type PrivKeyUnmarshaller func(data []byte) (PrivKey, error)

// Ed25519PrivateKey is an ed25519 private key.
type Ed25519PrivateKey struct {
	k ed25519.PrivateKey
}

// Ed25519PublicKey is an ed25519 public key.
type Ed25519PublicKey struct {
	k ed25519.PublicKey
}

// PubKeyUnmarshallers is a map of unmarshallers by key type
var PubKeyUnmarshallers = map[int]PubKeyUnmarshaller{
	RSA:     UnmarshalRsaPublicKey,
	Ed25519: UnmarshalEd25519PublicKey,
	ECDSA:   UnmarshalECDSAPublicKey,
}

// UnmarshalEd25519PublicKey returns a public key from input bytes.
func UnmarshalEd25519PublicKey(data []byte) (PubKey, error) {
	if len(data) != 32 {
		return nil, errors.New("expect ed25519 public key data size to be 32")
	}

	return &Ed25519PublicKey{
		k: ed25519.PublicKey(data),
	}, nil
}
func UnmarshalECDSAPublicKey(_ []byte) (PubKey, error) {
	// TODO implement me
	panic("implement me")
}
func UnmarshalRsaPublicKey(_ []byte) (PubKey, error) {
	// TODO implement me
	panic("implement me")
}

// PrivKeyUnmarshallers is a map of unmarshallers by key type
var _ = map[int]PrivKeyUnmarshaller{
	RSA:     UnmarshalRsaPrivateKey,
	Ed25519: UnmarshalEd25519PrivateKey,
	ECDSA:   UnmarshalECDSAPrivateKey,
}

func UnmarshalEd25519PrivateKey(data []byte) (PrivKey, error) {
	switch len(data) {
	case ed25519.PrivateKeySize + ed25519.PublicKeySize:
		// Remove the redundant public key. See issue #36.
		redundantPk := data[ed25519.PrivateKeySize:]
		pk := data[ed25519.PrivateKeySize-ed25519.PublicKeySize : ed25519.PrivateKeySize]
		if subtle.ConstantTimeCompare(pk, redundantPk) == 0 {
			return nil, errors.New("expected redundant ed25519 public key to be redundant")
		}

		// No point in storing the extra data.
		newKey := make([]byte, ed25519.PrivateKeySize)
		copy(newKey, data[:ed25519.PrivateKeySize])
		data = newKey
	case ed25519.PrivateKeySize:
	default:
		return nil, fmt.Errorf(
			"expected ed25519 data size to be %d or %d, got %d",
			ed25519.PrivateKeySize,
			ed25519.PrivateKeySize+ed25519.PublicKeySize,
			len(data),
		)
	}

	return &Ed25519PrivateKey{
		k: ed25519.PrivateKey(data),
	}, nil
}
func UnmarshalECDSAPrivateKey(_ []byte) (PrivKey, error) {
	// TODO implement me
	panic("implement me")
}
func UnmarshalRsaPrivateKey(_ []byte) (PrivKey, error) {
	// TODO implement me
	panic("implement me")
}

// Key represents a crypto key that can be compared to another key
type Key interface {
	// Equals checks whether two PubKeys are the same
	Equals(Key) bool

	// Raw returns the raw bytes of the key (not wrapped in the
	// libp2p-crypto protobuf).
	//
	// This function is the inverse of {Priv,Pub}KeyUnmarshaler.
	Raw() ([]byte, error)

	// Type returns the protobuf key type.
	Type() KeyType
}

// PrivKey represents a private key that can be used to generate a public key and sign data
type PrivKey interface {
	Key

	// Sign Cryptographically sign the given bytes
	Sign([]byte) ([]byte, error)

	// GetPublic Return a public key paired with this private key
	GetPublic() PubKey
}

// PubKey is a public key that can be used to verify data signed with the corresponding private key
type PubKey interface {
	Key

	// Verify that 'sig' is the signed hash of 'data'
	Verify(data []byte, sig []byte) (bool, error)
}

// GenSharedKey generates the shared key from a given private key
type GenSharedKey func([]byte) ([]byte, error)

// GenerateKeyPair generates a private and public key
func GenerateKeyPair(typ, bits int) (PrivKey, PubKey, error) {
	return GenerateKeyPairWithReader(typ, bits, rand.Reader)
}

// GenerateKeyPairWithReader returns a keypair of the given type and bitsize
func GenerateKeyPairWithReader(typ, bits int, src io.Reader) (PrivKey, PubKey, error) {
	switch typ {
	case RSA:
		return GenerateRSAKeyPair(bits, src)
	case Ed25519:
		return GenerateEd25519KeyPair(src)
	case ECDSA:
		return GenerateECDSAKeyPair(src)
	default:
		return nil, nil, ErrBadKeyType
	}
}

// GenerateEd25519KeyPair generates a new ed25519 private and public key pair.
func GenerateEd25519KeyPair(src io.Reader) (PrivKey, PubKey, error) {
	pub, priv, err := ed25519.GenerateKey(src)
	if err != nil {
		return nil, nil, err
	}

	return &Ed25519PrivateKey{
			k: priv,
		},
		&Ed25519PublicKey{
			k: pub,
		},
		nil
}
func GenerateECDSAKeyPair(_ io.Reader) (PrivKey, PubKey, error) {
	// TODO implement me
	panic("implement me")
}
func GenerateRSAKeyPair(_ int, _ io.Reader) (PrivKey, PubKey, error) {
	// TODO implement me
	panic("implement me")
}

/*// UnmarshalPublicKey converts a protobuf serialized public key into its
// representative object
func UnmarshalPublicKey(data []byte) (PubKey, error) {
	// Create a new instance of the protocol buffer message type
	pmes := new(crypto.PublicKey)

	// Decode the data using a protobuf decoder
	decoder := NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(pmes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public key: %s", err)
	}

	// Convert the protocol buffer message into a PubKey object
	key, err := PublicKeyFromProto(pmes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert public key from proto: %s", err)
	}

	return key, nil
}

// PublicKeyFromProto converts an unserialized protobuf PublicKey message
// into its representative object.
func PublicKeyFromProto(pmes *crypto.PublicKey, proto int) (PubKey, error) {
	um, ok := MarshalPublicKey(proto)
	if !ok {
		return nil, ErrBadKeyType
	}

	data := pmes

	pk, err := um(data)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

// MarshalPublicKey converts a public key object into a protobuf serialized
// public key
func MarshalPublicKey(k PubKey) ([]byte, error) {
	pbmes, err := PublicKeyToProto(k)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(pbmes)
}

// PublicKeyToProto converts a public key object into an unserialized
// protobuf PublicKey message.
func PublicKeyToProto(k PubKey) (*crypto.PublicKey, error) {
	data, err := k.Raw()
	if err != nil {
		return nil, err
	}
	return &pb.PublicKey{
		Type: k.Type().Enum(),
		Data: data,
	}, nil
}*/

// UnmarshalPrivateKey converts a protobuf serialized private key into its
// representative object

// MarshalPrivateKey converts a key object into its protobuf serialized form.

func basicEquals(k1, k2 Key) bool {
	if k1.Type() != k2.Type() {
		return false
	}

	a, err := k1.Raw()
	if err != nil {
		return false
	}
	b, err := k2.Raw()
	if err != nil {
		return false
	}

	// Only the name is constant time,
	//the actual complexity depends on the length
	// , but it doesn't depend on the type of the input
	//because, array of bytes supplied.
	return subtle.ConstantTimeCompare(a, b) == 1
}

// Type of the private key (Ed25519).
func (k *Ed25519PrivateKey) Type() KeyType {
	return Ed25519
}

// Raw private key bytes.
func (k *Ed25519PrivateKey) Raw() ([]byte, error) {
	// The Ed25519 private key contains two 32-bytes curve points, the private
	// key and the public key.
	// It makes it more efficient to get the public key without re-computing an
	// elliptic curve multiplication.
	buf := make([]byte, len(k.k))
	copy(buf, k.k)

	return buf, nil
}

func (k *Ed25519PrivateKey) pubKeyBytes() []byte {
	return k.k[ed25519.PrivateKeySize-ed25519.PublicKeySize:]
}

// Equals compares two ed25519 private keys.
func (k *Ed25519PrivateKey) Equals(o Key) bool {
	edk, ok := o.(*Ed25519PrivateKey)
	if !ok {
		return basicEquals(k, o)
	}

	return subtle.ConstantTimeCompare(k.k, edk.k) == 1
}

// GetPublic returns an ed25519 public key from a private key.
func (k *Ed25519PrivateKey) GetPublic() PubKey {
	return &Ed25519PublicKey{k: k.pubKeyBytes()}
}

// Sign returns a signature from an input message.
func (k *Ed25519PrivateKey) Sign(msg []byte) (res []byte, err error) {
	//defer func() { catch.HandlePanic(recover(), &err, "ed15519 signing") }()

	return ed25519.Sign(k.k, msg), nil
}

// Type of the public key (Ed25519).
func (k *Ed25519PublicKey) Type() KeyType {
	return Ed25519
}

// Raw public key bytes.
func (k *Ed25519PublicKey) Raw() ([]byte, error) {
	return k.k, nil
}

// Equals compares two ed25519 public keys.
func (k *Ed25519PublicKey) Equals(o Key) bool {
	edk, ok := o.(*Ed25519PublicKey)
	if !ok {
		return basicEquals(k, o)
	}

	return bytes.Equal(k.k, edk.k)
}

// Verify checks a signature agains the input data.
func (k *Ed25519PublicKey) Verify(data []byte, sig []byte) (success bool, err error) {
	/*defer func() {
		catch.HandlePanic(recover(), &err, "ed15519 signature verification")

		// To be safe.
		if err != nil {
			success = false
		}
	}()*/
	return ed25519.Verify(k.k, data, sig), nil
}

func MarshalPublicKey(pubKey PubKey) ([]byte, error) {
	return pubKey.Raw()
}

func UnmarshalPublicKey(blst []byte) (PubKey, error) {
	return PubKeyUnmarshallers[1](blst)
}
