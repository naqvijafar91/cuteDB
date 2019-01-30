package cutedb

//DB - Handle exported by the package
type DB struct {
	storage *btree
}

//Open - Opens a new db connection at the file path
func Open(filePath string) (*DB, error) {
	storage, err := initializeBtree(filePath)
	if err != nil {
		return nil, err
	}
	return &DB{storage}, nil
}

//Put - Insert a key value pair in the database
func (db *DB) Put(key string, value string) error {
	pair := newPair(key, value)
	return db.storage.insert(pair)
}

//Get - Get the stored value from the database for the respective key
func (db *DB) Get(key string) (string, bool, error) {
	return db.Get(key)
}
