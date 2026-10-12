package errors

type Kind int

const (
	KindUnknown Kind = iota
	KindNotFound
	KindInvalidInput
	KindTimeout
	KindRateLimited
	KindConfig
	KindExternal
	KindCancel
)

type Error struct {
	Service string
	Op      string
	Kind    Kind
	Err     error
}

func (k Kind) String() string {
	switch k {
	case KindNotFound:
		return "not_found"
	case KindInvalidInput:
		return "invalid_input"
	case KindTimeout:
		return "timeout"
	case KindRateLimited:
		return "rate_limited"
	case KindConfig:
		return "config"
	case KindExternal:
		return "external"
	case KindCancel:
		return "cancel"
	default:
		return "unknown"
	}
}
