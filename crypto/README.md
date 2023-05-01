# Crypto package

This package provides a way to implement `RSA`, `ECDSA` and `ED25519` of which, `ED25519` has been implemented to
provide secure transfer of data out of box.

## Types

### KeyType
This is an enum for the supported key types which includes RSA, Ed25519, and ECDSA.

### Ed25519PrivateKey
This is an implementation of an ed25519 private key.

### Ed25519PublicKey
This is an implementation of an ed25519 public key.

### PubKeyUnmarshaller
This is a func that creates a PubKey from a given slice of bytes.

### PrivKeyUnmarshaller
This is a func that creates a PrivKey from a given slice of bytes.

### Key
This represents a crypto key that can be compared to another key. It contains three methods:

- Equals: checks whether two PubKeys are the same.
- Raw: returns the raw bytes of the key (not wrapped in the libp2p-crypto protobuf).
- Type: returns the protobuf key type.

### PrivKey
This represents a private key that can be used to generate a public key and sign data. It contains three methods:

- Key: includes the methods from Key.
- Sign: cryptographically signs the given bytes.
- GetPublic: returns a public key paired with this private key.

### PubKey
This is a public key that can be used to verify data signed with the corresponding private key. It also includes the methods from Key.

### GenSharedKey
This is a function that generates the shared key from a given private key.

## Functions

### UnmarshalEd25519PublicKey
This function returns a public key from input bytes.

### UnmarshalEd25519PrivateKey
This function returns a private key from input bytes.

### GenerateKeyPair
This function generates a private and public key.

## Usage
This package can be used to implement various cryptographic functionalities. The `GenerateKeyPair` function can be used to generate a private and public key, which can then be used to sign and verify data. The `UnmarshalEd25519PublicKey` and `UnmarshalEd25519PrivateKey` functions can be used to deserialize a public and private key, respectively, from byte slices.