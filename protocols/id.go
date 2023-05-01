package protocol

// ID is an identifier used to write protocol headers in streams.
type ID string

// These are reserved protocol.IDs.
const (
	TestingID ID = "/p2p/_testing"
)

// ConvertFromStrings used to testing
func ConvertFromStrings(ids []string) (res []ID) {
	res = make([]ID, 0, len(ids))
	for _, id := range ids {
		res = append(res, ID(id))
	}
	return res
}

// ConvertToStrings used for testing
func ConvertToStrings(ids []ID) (res []string) {
	res = make([]string, 0, len(ids))
	for _, id := range ids {
		res = append(res, string(id))
	}
	return res
}
