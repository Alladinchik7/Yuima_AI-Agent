package sqlite

type Store struct {
	db *Conn
}

func NewStore(db *Conn) *Store {
	return &Store{db: db}
}