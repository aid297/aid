package driver

func (tx *Tx) Get(key string) ([]byte, bool, error) {
	value, ok, err := tx.core.Get(key)
	return value, ok, wrapError(err)
}

func (tx *Tx) Put(key string, value []byte) error { return wrapError(tx.core.Put(key, value)) }

func (tx *Tx) Delete(key string) error { return wrapError(tx.core.Delete(key)) }

func (tx *Tx) Query(prefix string) (map[string][]byte, error) {
	rows, err := tx.core.Query(prefix)
	return rows, wrapError(err)
}

func (tx *Tx) Keys() ([]string, error) {
	keys, err := tx.core.Keys()
	return keys, wrapError(err)
}

func (tx *Tx) FindRow(primaryKey any) (Row, bool, error) {
	row, ok, err := tx.core.FindRow(primaryKey)
	return row, ok, wrapError(err)
}

func (tx *Tx) FindByConditions(conditions []QueryCondition) ([]Row, error) {
	rows, err := tx.core.FindByConditions(conditions)
	return rows, wrapError(err)
}

func (tx *Tx) Find(conditions ...QueryCondition) ([]Row, error) {
	rows, err := tx.core.Find(conditions...)
	return rows, wrapError(err)
}

func (tx *Tx) FindOne(conditions ...QueryCondition) (Row, bool, error) {
	row, ok, err := tx.core.FindOne(conditions...)
	return row, ok, wrapError(err)
}

func (tx *Tx) InsertRow(values Row) (Row, error) {
	row, err := tx.core.InsertRow(values)
	return row, wrapError(err)
}

func (tx *Tx) InsertRows(values []Row) ([]Row, error) {
	rows, err := tx.core.InsertRows(values)
	return rows, wrapError(err)
}

func (tx *Tx) UpdateRow(primaryKey any, updates Row) (Row, error) {
	row, err := tx.core.UpdateRow(primaryKey, updates)
	return row, wrapError(err)
}

func (tx *Tx) UpdateRows(updates []RowUpdate) ([]Row, error) {
	rows, err := tx.core.UpdateRows(updates)
	return rows, wrapError(err)
}

func (tx *Tx) DeleteRow(primaryKey any) error { return wrapError(tx.core.DeleteRow(primaryKey)) }

func (tx *Tx) DeleteRows(primaryKeys []any) error { return wrapError(tx.core.DeleteRows(primaryKeys)) }

func (tx *Tx) RemoveByCondition(conditions ...QueryCondition) (int, error) {
	count, err := tx.core.RemoveByCondition(conditions...)
	return count, wrapError(err)
}

func (tx *Tx) RemoveOneByCondition(conditions ...QueryCondition) (bool, error) {
	ok, err := tx.core.RemoveOneByCondition(conditions...)
	return ok, wrapError(err)
}

func (tx *Tx) Commit() error { return wrapError(tx.core.Commit()) }

func (tx *Tx) Rollback() error { return wrapError(tx.core.Rollback()) }
