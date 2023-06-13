package address

type StringID struct {
	ID string
}

func (id StringID) String() string {
	return id.ID
}
