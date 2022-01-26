package meiwobuxing

import (
	"encoding/json"
	"errors"
	"github.com/dgraph-io/badger/v3"
)

var db *badger.DB

func InitDB(opts badger.Options) (err error) {
	db, err = badger.Open(opts)
	return err
}

func LevelDB() *badger.DB {
	if db == nil {
		panic("not init db")
	}
	return db
}

func SetJson(prefix string, key string, value interface{}) error {
	key = parseKey(prefix, key)
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if db != nil {
		return db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(key), v)
		})
	}
	return errors.New("not init db")
}

func GetJson(prefix string, key string, value interface{}) error {
	key = parseKey(prefix, key)
	var rawValue []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		rawValue, err = item.ValueCopy(rawValue)
		return err
	})

	if err != nil {
		return err
	}
	return json.Unmarshal(rawValue, value)
}

func Delete(prefix string, key string) error {
	key = parseKey(prefix, key)
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func Iter(prefix string) (items map[string]string, err error) {
	items = make(map[string]string)
	pLen := len(prefix)
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		p := []byte(prefix)

		defer it.Close()
		var err1 error
		for it.Seek(p); it.ValidForPrefix(p); it.Next() {
			item := it.Item()
			k := item.Key()
			var rawValue []byte
			rawValue, err1 = item.ValueCopy(rawValue)
			items[string(k[pLen:])] = string(rawValue)
		}
		return err1
	})
	return
}

func parseKey(prefix, key string) string {
	return prefix + ":" + key
}
