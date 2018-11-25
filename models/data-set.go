package models

// DataSet gives a title to the data we want to export.
type DataSet struct {
	Title string
	Data  UserMap
}

// NewDataSet returns a new DataSet.
func NewDataSet(title string, data UserMap) DataSet {
	return DataSet{
		Title: title,
		Data:  data,
	}
}

type UserSet struct {
	Title string
	Users []User
}

func NewUserSet(title string, users []User) *UserSet {
	return &UserSet{
		Title: title,
		Users: users,
	}
}

type UserSetChan chan *UserSet
