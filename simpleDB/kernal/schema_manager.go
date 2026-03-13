package kernal

import (
	"fmt"
	"strings"

	json "github.com/json-iterator/go"
)

// ─── AlterTable Plan ─────────────────────────────────────────────────────────

// AlterTablePlan 描述一次 ALTER TABLE 操作的完整计划。
// 所有字段均可按需组合，同一次调用可包含多个操作。
type AlterTablePlan struct {
	// AddColumns 向表添加新列。新列可携带 Default 值，对已有行自动回填。
	// 若列名已存在则返回 ErrColumnAlreadyExists。
	AddColumns []Column `json:"addColumns,omitempty"`

	// DropColumns 按列名删除列，已有行的对应字段也会被清除。
	// 不允许删除主键列（返回 ErrCannotDropPrimaryKey）。
	// 若列名不存在则返回 ErrColumnNotFound。
	DropColumns []string `json:"dropColumns,omitempty"`

	// AddIndexes 为已有列添加普通索引（将列的 Indexed 置为 true）。
	AddIndexes []string `json:"addIndexes,omitempty"`

	// DropIndexes 移除已有列的普通索引（将列的 Indexed 置为 false）。
	// 若该列同时是唯一索引，唯一索引不受影响。
	DropIndexes []string `json:"dropIndexes,omitempty"`

	// AddUniques 为已有列添加唯一索引（同时隐式将 Indexed 置为 true）。
	AddUniques []string `json:"addUniques,omitempty"`

	// DropUniques 移除已有列的唯一索引（将列的 Unique 置为 false）。
	// 普通索引标志 Indexed 不受影响。
	DropUniques []string `json:"dropUniques,omitempty"`

	// AddForeignKeys 新增外键链路定义。
	// 会自动将外键字段标记为 Indexed（由 normalizeSchema 处理）。
	AddForeignKeys []ForeignKey `json:"addForeignKeys,omitempty"`

	// DropForeignKeys 删除外键链路，支持通过 Name / Field / Alias 指定。
	DropForeignKeys []string `json:"dropForeignKeys,omitempty"`
}

// isEmpty 判断计划是否为空（没有任何操作）。
func (p AlterTablePlan) isEmpty() bool {
	return len(p.AddColumns) == 0 &&
		len(p.DropColumns) == 0 &&
		len(p.AddIndexes) == 0 &&
		len(p.DropIndexes) == 0 &&
		len(p.AddUniques) == 0 &&
		len(p.DropUniques) == 0 &&
		len(p.AddForeignKeys) == 0 &&
		len(p.DropForeignKeys) == 0
}

// ─── DDL: HasSchema ───────────────────────────────────────────────────────────

// HasSchema 返回当前 SimpleDB 实例是否已配置 TableSchema。
// 此方法是线程安全的。
func (db *SimpleDB) HasSchema() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.schema != nil
}

// ─── DDL: CreateTable ─────────────────────────────────────────────────────────

// CreateTable 以严格模式创建表结构。
// 与 Configure 的区别：若表已存在 Schema，即使完全相同也返回 ErrSchemaAlreadyExists。
// 适用于明确的 DDL 创建语义（CREATE TABLE IF NOT EXISTS 请先调用 HasSchema 判断）。
func (db *SimpleDB) CreateTable(schema TableSchema) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	if db.schema != nil {
		return ErrSchemaAlreadyExists
	}

	schema = db.applySchemaDefaultsLocked(schema)
	normalized, err := normalizeSchema(schema)
	if err != nil {
		return err
	}

	return db.createTableLocked(normalized)
}

// ─── DDL: DropTable ───────────────────────────────────────────────────────────

// DropTable 删除表的所有行数据以及 Schema 元信息，但不关闭数据库文件。
// 调用后该 SimpleDB 实例进入"无 Schema"状态，可再次调用 CreateTable 重建。
func (db *SimpleDB) DropTable() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	if db.schema == nil {
		return ErrSchemaNotConfigured
	}

	// 收集需要删除的 key：行数据 + schema + sequence
	toDelete := make([]string, 0)
	for key, current := range db.index {
		if current.Deleted {
			continue
		}
		if strings.HasPrefix(key, rowKeyPrefix) ||
			key == metaSchemaKey ||
			key == metaSequenceKey {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		if err := db.deleteRawLocked(key); err != nil {
			return err
		}
	}

	// 重置内存状态
	db.schema = nil
	db.autoSeq = 0
	db.uniqueIdx = make(map[string]map[string]string)
	db.indexIdx = make(map[string]map[string]map[string]struct{})

	return nil
}

// ─── DDL: TruncateTable ───────────────────────────────────────────────────────

// TruncateTable 清空表中所有行数据，保留 Schema 定义及索引结构。
// 若表配置了自增主键，自增序列计数器也会重置为 0。
func (db *SimpleDB) TruncateTable() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	if db.schema == nil {
		return ErrSchemaNotConfigured
	}

	// 只删除行数据，不删除 schema/sequence 元数据
	toDelete := make([]string, 0)
	for key, current := range db.index {
		if current.Deleted {
			continue
		}
		if strings.HasPrefix(key, rowKeyPrefix) {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		if err := db.deleteRawLocked(key); err != nil {
			return err
		}
	}

	// 重置序列
	db.autoSeq = 0
	if db.schema.AutoIncrement && autoIncrementUsesSequence(*db.schema) {
		if err := db.persistSequenceLocked(0); err != nil {
			return err
		}
	}

	// 重建空的二级索引结构（保留索引定义，清空数据）
	db.resetSecondaryIndexesLocked()

	return nil
}

// ─── DDL: AlterTable ──────────────────────────────────────────────────────────

// AlterTable 根据 AlterTablePlan 修改表结构，操作按以下顺序执行：
//  1. DropColumns → 2. AddColumns → 3. DropUniques → 4. DropIndexes →
//  5. AddIndexes → 6. AddUniques
//
// 所有操作均在同一个锁保护下原子完成。已有行会按计划迁移（删列/补默认值）。
// 操作完成后自动重建二级索引。
func (db *SimpleDB) AlterTable(plan AlterTablePlan) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	if db.schema == nil {
		return ErrSchemaNotConfigured
	}

	if plan.isEmpty() {
		return fmt.Errorf("%w: 计划为空，没有任何操作", ErrAlterTableInvalid)
	}

	// ── 参数校验：提前发现错误 ───────────────────────────────────────────────

	// 校验 DropColumns
	for _, colName := range plan.DropColumns {
		if !db.hasColumn(colName) {
			return fmt.Errorf("%w: %s", ErrColumnNotFound, colName)
		}
		if colName == db.schema.PrimaryKey {
			return fmt.Errorf("%w: %s", ErrCannotDropPrimaryKey, colName)
		}
	}

	// 构建"操作后的列名集合"用于 AddColumns / index 操作的引用校验
	remainingCols := make(map[string]struct{}, len(db.schema.Columns))
	for _, col := range db.schema.Columns {
		remainingCols[col.Name] = struct{}{}
	}
	for _, colName := range plan.DropColumns {
		delete(remainingCols, colName)
	}
	for _, col := range plan.AddColumns {
		if _, exists := remainingCols[col.Name]; exists {
			return fmt.Errorf("%w: %s", ErrColumnAlreadyExists, col.Name)
		}
		remainingCols[col.Name] = struct{}{}
	}
	for _, field := range plan.AddIndexes {
		if _, exists := remainingCols[field]; !exists {
			return fmt.Errorf("%w: AddIndex 引用不存在的列 %s", ErrAlterTableInvalid, field)
		}
	}
	for _, field := range plan.DropIndexes {
		if _, exists := remainingCols[field]; !exists {
			return fmt.Errorf("%w: DropIndex 引用不存在的列 %s", ErrAlterTableInvalid, field)
		}
	}
	for _, field := range plan.AddUniques {
		if _, exists := remainingCols[field]; !exists {
			return fmt.Errorf("%w: AddUnique 引用不存在的列 %s", ErrAlterTableInvalid, field)
		}
	}
	for _, field := range plan.DropUniques {
		if _, exists := remainingCols[field]; !exists {
			return fmt.Errorf("%w: DropUnique 引用不存在的列 %s", ErrAlterTableInvalid, field)
		}
	}

	for _, foreignKey := range plan.AddForeignKeys {
		if _, exists := remainingCols[strings.TrimSpace(foreignKey.Field)]; !exists {
			return fmt.Errorf("%w: AddForeignKey 引用不存在的列 %s", ErrAlterTableInvalid, foreignKey.Field)
		}
	}

	if _, err := resolveForeignKeyDropIndices(db.schema.ForeignKeys, plan.DropForeignKeys); err != nil {
		return err
	}

	return db.alterTableLocked(plan)
}

// ─── DDL: SchemaDiff ─────────────────────────────────────────────────────────

// SchemaDiff 计算从当前 Schema 迁移到目标 Schema 所需的最小 AlterTablePlan。
// 若当前无 Schema，返回 (nil, false, nil)，表示需要 CreateTable 而非 Alter。
// 若当前 Schema 与目标完全一致，返回 (nil, true, nil)（第二个值表示 Schema 已存在）。
// 否则返回 (plan, true, nil)，plan 包含所有需要执行的变更。
//
// 注意：SchemaDiff 检测以下变更：
//   - 新增列（target 有、current 无）
//   - 删除列（current 有、target 无，主键列除外）
//   - 索引变更（Indexed / Unique 标志的开启与关闭）
//   - 外键变更（ForeignKeys 增删）
//
// SchemaDiff 不检测列类型变更、列约束变更、主键变更。
// 如需这些变更，请直接调用 AlterTable 或重建表（DropTable + CreateTable）。
func (db *SimpleDB) SchemaDiff(target TableSchema) (*AlterTablePlan, bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if err := db.ensureOpen(); err != nil {
		return nil, false, err
	}

	if db.schema == nil {
		return nil, false, nil
	}

	normalized, err := normalizeSchema(target)
	if err != nil {
		return nil, true, err
	}

	plan := db.computeDiffLocked(*db.schema, normalized, true)
	if plan.isEmpty() {
		return nil, true, nil
	}
	return &plan, true, nil
}

// computeDiffLocked 计算从 current → target 的 AlterTablePlan。
// withDrops 控制是否计算 DropColumns/DropIndexes/DropUniques。
// 调用方必须持有至少读锁。
func (db *SimpleDB) computeDiffLocked(current, target TableSchema, withDrops bool) AlterTablePlan {
	var plan AlterTablePlan

	// 构建当前列的 map（name → Column）
	currentCols := make(map[string]Column, len(current.Columns))
	for _, col := range current.Columns {
		currentCols[col.Name] = col
	}

	// 构建目标列的 map（name → Column）
	targetCols := make(map[string]Column, len(target.Columns))
	for _, col := range target.Columns {
		targetCols[col.Name] = col
	}

	// 1. 新增列（target 有、current 无）
	for _, col := range target.Columns {
		if _, exists := currentCols[col.Name]; !exists {
			plan.AddColumns = append(plan.AddColumns, col)
		}
	}

	// 2. 删除列（current 有、target 无，且不是主键）
	if withDrops {
		for _, col := range current.Columns {
			if _, exists := targetCols[col.Name]; !exists {
				if col.Name != current.PrimaryKey {
					plan.DropColumns = append(plan.DropColumns, col.Name)
				}
			}
		}
	}

	// 3. 对共有列，比较 Indexed/Unique 变更
	for _, targetCol := range target.Columns {
		currentCol, exists := currentCols[targetCol.Name]
		if !exists {
			continue // 新增列已在上面处理
		}

		goingUnique := !currentCol.Unique && targetCol.Unique
		droppingUnique := currentCol.Unique && !targetCol.Unique

		// ── Unique 变更 ───────────────────────────────────────────────────────
		if goingUnique {
			plan.AddUniques = append(plan.AddUniques, targetCol.Name)
			// AddUniques 会自动将 Indexed=true，无需额外 AddIndex
		} else if droppingUnique && withDrops {
			plan.DropUniques = append(plan.DropUniques, targetCol.Name)
			// DropUniques 不改变 Indexed 标志；原来 email.Indexed=false，
			// 独立索引会随 unique 一起消失，无需额外 DropIndex。
			// 但若 target 想保留显式 Indexed，则需要 AddIndex。
			if targetCol.Indexed && !currentCol.Indexed {
				plan.AddIndexes = append(plan.AddIndexes, targetCol.Name)
			}
		}

		// ── Indexed 变更（仅在 Unique 状态不变时比较显式 Indexed 标志）────────
		if !goingUnique && !droppingUnique {
			if !currentCol.Indexed && targetCol.Indexed {
				plan.AddIndexes = append(plan.AddIndexes, targetCol.Name)
			} else if currentCol.Indexed && !targetCol.Indexed && withDrops {
				plan.DropIndexes = append(plan.DropIndexes, targetCol.Name)
			}
		}
	}

	// 4. 外键变更（使用规范签名对比）
	currentFKSet := make(map[string]ForeignKey, len(current.ForeignKeys))
	for _, foreignKey := range current.ForeignKeys {
		currentFKSet[foreignKeySignature(foreignKey)] = foreignKey
	}
	targetFKSet := make(map[string]ForeignKey, len(target.ForeignKeys))
	for _, foreignKey := range target.ForeignKeys {
		targetFKSet[foreignKeySignature(foreignKey)] = foreignKey
	}

	for signature, foreignKey := range targetFKSet {
		if _, exists := currentFKSet[signature]; !exists {
			plan.AddForeignKeys = append(plan.AddForeignKeys, foreignKey)
		}
	}

	if withDrops {
		for signature, foreignKey := range currentFKSet {
			if _, exists := targetFKSet[signature]; !exists {
				dropRef := strings.TrimSpace(foreignKey.Name)
				if dropRef == "" {
					dropRef = strings.TrimSpace(foreignKey.Field)
				}
				plan.DropForeignKeys = append(plan.DropForeignKeys, dropRef)
			}
		}
	}

	return plan
}

// ─── DDL: AutoMigrate ────────────────────────────────────────────────────────

// AutoMigrate 以**保守策略**将表结构迁移到目标 Schema：
//   - 若当前无 Schema → 等同于 CreateTable（严格建表）
//   - 若当前 Schema 与目标完全一致 → 幂等，无操作
//   - 若当前 Schema 与目标不同 → **只新增**列/索引，绝不删除任何列或索引
//
// 适合应用启动时安全地追加字段，不会破坏现有数据。
// 若需要删除列，请手动调用 AlterTable 或使用 SyncSchema。
func (db *SimpleDB) AutoMigrate(schema TableSchema) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	schema = db.applySchemaDefaultsLocked(schema)
	normalized, err := normalizeSchema(schema)
	if err != nil {
		return err
	}

	// 禁止修改引擎类型
	if db.schema != nil && db.schema.Engine != normalized.Engine {
		return fmt.Errorf("cannot change table engine from %s to %s", db.schema.Engine, normalized.Engine)
	}

	// 无 Schema → 建表
	if db.schema == nil {
		return db.createTableLocked(normalized)
	}

	// 计算保守 diff（withDrops=false，只增不减）
	plan := db.computeDiffLocked(*db.schema, normalized, false)
	if plan.isEmpty() {
		return nil // 已是最新，幂等
	}

	return db.alterTableLocked(plan)
}

// ─── DDL: SyncSchema ─────────────────────────────────────────────────────────

// SyncSchema 以**完全同步策略**将表结构精确对齐到目标 Schema：
//   - 若当前无 Schema → 等同于 CreateTable
//   - 若当前 Schema 与目标完全一致 → 幂等，无操作
//   - 若当前 Schema 与目标不同 → 全量同步，包括**删除**目标中已移除的列/索引
//
// ⚠️ 删列是不可逆操作，会永久丢失该列的所有数据，请谨慎使用。
// 若只需安全追加字段，请使用 AutoMigrate。
func (db *SimpleDB) SyncSchema(schema TableSchema) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.ensureOpen(); err != nil {
		return err
	}

	schema = db.applySchemaDefaultsLocked(schema)
	normalized, err := normalizeSchema(schema)
	if err != nil {
		return err
	}

	// 无 Schema → 建表
	if db.schema == nil {
		return db.createTableLocked(normalized)
	}

	// 已一致 → 幂等
	if schemasEqual(*db.schema, normalized) {
		return nil
	}

	// 计算完整 diff（withDrops=true）
	plan := db.computeDiffLocked(*db.schema, normalized, true)
	if plan.isEmpty() {
		return nil
	}

	return db.alterTableLocked(plan)
}

// ─── 内部无锁实现（供 AutoMigrate/SyncSchema 共用，避免重复加锁）────────────

// createTableLocked 在已持有写锁的情况下建表，假设 schema 已规范化。
func (db *SimpleDB) createTableLocked(normalized TableSchema) error {
	payload, err := json.Marshal(normalized)
	if err != nil {
		return err
	}

	// 提前设置 schema，以便 putRawLocked 能够识别引擎类型
	db.schema = &normalized

	if normalized.Engine == EngineDisk || (normalized.Engine == EngineMem && normalized.Disk) {
		if err = db.ensureFileOpenLocked(); err != nil {
			return err
		}
	}

	if err = db.putRawLocked(metaSchemaKey, payload); err != nil {
		return err
	}

	if normalized.Engine == EngineMem && normalized.Disk {
		// 内存引擎开启落盘时，强制立即将 Schema 写入磁盘以防掉电丢失
		if err = db.persistMemToDiskLocked(false); err != nil {
			return err
		}
	}

	if normalized.AutoIncrement && autoIncrementUsesSequence(normalized) {
		if err = db.persistSequenceLocked(0); err != nil {
			return err
		}
	}
	db.autoSeq = 0
	db.resetSecondaryIndexesLocked()
	return nil
}

// alterTableLocked 在已持有写锁的情况下执行 AlterTable，假设 plan 非空且已验证。
// 注意：此方法不做重复的参数校验，调用方需保证 plan 合法。
func (db *SimpleDB) alterTableLocked(plan AlterTablePlan) error {
	// 直接复用公开的 AlterTable 逻辑，但 AlterTable 会自己加锁，
	// 因此这里直接内联阶段性操作以避免死锁。
	// ── 构建 dropSet ──────────────────────────────────────────────────────────
	dropSet := toStringSet(plan.DropColumns)

	// ── 构建新 Schema ─────────────────────────────────────────────────────────
	newSchema := cloneSchema(*db.schema)

	// 删列
	if len(plan.DropColumns) > 0 {
		newCols := make([]Column, 0, len(newSchema.Columns))
		for _, col := range newSchema.Columns {
			if _, dropped := dropSet[col.Name]; !dropped {
				newCols = append(newCols, col)
			}
		}
		newSchema.Columns = newCols
	}

	// 加列
	newSchema.Columns = append(newSchema.Columns, plan.AddColumns...)

	// index/unique 修改
	addIdxSet := toStringSet(plan.AddIndexes)
	dropIdxSet := toStringSet(plan.DropIndexes)
	addUniqSet := toStringSet(plan.AddUniques)
	dropUniqSet := toStringSet(plan.DropUniques)

	for i := range newSchema.Columns {
		name := newSchema.Columns[i].Name
		if _, ok := dropUniqSet[name]; ok {
			newSchema.Columns[i].Unique = false
		}
		if _, ok := addUniqSet[name]; ok {
			newSchema.Columns[i].Unique = true
			newSchema.Columns[i].Indexed = true
		}
		if _, ok := dropIdxSet[name]; ok {
			if !newSchema.Columns[i].Unique {
				newSchema.Columns[i].Indexed = false
			}
		}
		if _, ok := addIdxSet[name]; ok {
			newSchema.Columns[i].Indexed = true
		}
	}

	// 外键链路修改
	if len(plan.DropForeignKeys) > 0 {
		dropIndices, err := resolveForeignKeyDropIndices(newSchema.ForeignKeys, plan.DropForeignKeys)
		if err != nil {
			return err
		}
		if len(dropIndices) > 0 {
			nextForeignKeys := make([]ForeignKey, 0, len(newSchema.ForeignKeys)-len(dropIndices))
			for index, foreignKey := range newSchema.ForeignKeys {
				if _, dropped := dropIndices[index]; dropped {
					continue
				}
				nextForeignKeys = append(nextForeignKeys, foreignKey)
			}
			newSchema.ForeignKeys = nextForeignKeys
		}
	}

	if len(plan.AddForeignKeys) > 0 {
		newSchema.ForeignKeys = append(newSchema.ForeignKeys, plan.AddForeignKeys...)
	}

	// ── 规范化新 Schema ───────────────────────────────────────────────────────
	normalized, err := normalizeSchema(newSchema)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrAlterTableInvalid, err)
	}

	// ── 迁移已有行 ────────────────────────────────────────────────────────────
	rowKeys := make([]string, 0)
	for key, current := range db.index {
		if !current.Deleted && strings.HasPrefix(key, rowKeyPrefix) {
			rowKeys = append(rowKeys, key)
		}
	}
	for _, key := range rowKeys {
		row, decErr := decodeRow(db.index[key].Value)
		if decErr != nil {
			return decErr
		}
		for colName := range dropSet {
			delete(row, colName)
		}
		for _, col := range plan.AddColumns {
			if _, exists := row[col.Name]; !exists && col.Default != nil {
				row[col.Name] = col.Default
			}
		}
		encodedRow, encErr := encodeRow(row)
		if encErr != nil {
			return encErr
		}
		if putErr := db.putRawLocked(key, encodedRow); putErr != nil {
			return putErr
		}
	}

	if normalized.Engine == EngineDisk {
		if err = db.ensureFileOpenLocked(); err != nil {
			return err
		}
	}

	// ── 持久化新 Schema ───────────────────────────────────────────────────────
	payload, err := json.Marshal(normalized)
	if err != nil {
		return err
	}
	if err = db.putRawLocked(metaSchemaKey, payload); err != nil {
		return err
	}

	return db.rebuildStructuredState()
}

// ─── 辅助函数 ─────────────────────────────────────────────────────────────────

// toStringSet 将字符串切片转为 set。
func toStringSet(items []string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

func foreignKeySignature(foreignKey ForeignKey) string {
	return strings.TrimSpace(foreignKey.Name) + "|" +
		strings.TrimSpace(foreignKey.Field) + "|" +
		strings.TrimSpace(foreignKey.RefTable) + "|" +
		strings.TrimSpace(foreignKey.RefField) + "|" +
		strings.TrimSpace(foreignKey.Alias)
}

func resolveForeignKeyDropIndices(foreignKeys []ForeignKey, dropRefs []string) (map[int]struct{}, error) {
	dropIndices := make(map[int]struct{})
	if len(dropRefs) == 0 {
		return dropIndices, nil
	}

	for _, dropRef := range dropRefs {
		requested := strings.TrimSpace(dropRef)
		if requested == "" {
			return nil, fmt.Errorf("%w: DropForeignKey 引用为空", ErrAlterTableInvalid)
		}

		matchedIndex := -1
		for index, foreignKey := range foreignKeys {
			if requested != strings.TrimSpace(foreignKey.Name) &&
				requested != strings.TrimSpace(foreignKey.Field) &&
				requested != strings.TrimSpace(foreignKey.Alias) {
				continue
			}
			if matchedIndex >= 0 {
				return nil, fmt.Errorf("%w: DropForeignKey %s 命中多个外键，请使用 Name 精确指定", ErrAlterTableInvalid, requested)
			}
			matchedIndex = index
		}

		if matchedIndex < 0 {
			return nil, fmt.Errorf("%w: DropForeignKey 未找到 %s", ErrAlterTableInvalid, requested)
		}

		dropIndices[matchedIndex] = struct{}{}
	}

	return dropIndices, nil
}
