package models

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "flag"
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
