package storedvalue

type Contract struct {
	Body            []byte          `json:"body"`
	NamedKeys       []NamedKey      `json:"named_key"`
	ProtocolVersion ProtocolVersion `json:"protocol_version"`
}

func NewContract(body []byte, namedKeys []NamedKey, protocolVersion ProtocolVersion) Contract {
	return Contract{
		Body:            body,
		NamedKeys:       namedKeys,
		ProtocolVersion: protocolVersion,
	}
}

// TODO Need to check..
func (c Contract) FromBytes(src []byte) (contract Contract, err error, pos int) {
	return Contract{}, nil, pos
}

type ProtocolVersion struct {
	Major uint32 `json:"major,omitempty"`
	Minor uint32 `json:"minor,omitempty"`
	Patch uint32 `json:"patch,omitempty"`
}

func NewProtocolVersion(major uint32, minor uint32, patch uint32) ProtocolVersion {
	return ProtocolVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}
