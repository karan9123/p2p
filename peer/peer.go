package peer

// ID is a p2p peer identity.
// Peer IDs are derived by hashing a peer's public key and encoding the hash output as multihash.
type ID string
