package db

const (
	MaxLimit     = 1000
	DefaultLimit = 100
)

var DefaultPaginator = NewPaginator(DefaultLimit, 0)

type Paginator struct {
	Limit  uint64 `json:"limit,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}

func NewPaginator(limit, offset uint64) *Paginator {
	return &Paginator{
		Limit:  NormalizeLimit(limit),
		Offset: offset,
	}
}

func NormalizeLimit(limit uint64) uint64 {
	switch {
	case limit == 0:
		return DefaultLimit
	case limit > MaxLimit:
		return MaxLimit
	default:
		return limit
	}
}
