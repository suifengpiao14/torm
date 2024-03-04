package sourceprovider

import (
	"context"

	"github.com/pkg/errors"
)

type MemoryDB struct {
	InOutMap map[string]string
}

func (m *MemoryDB) TypeName() string {
	return "memory_db"
}

func (m *MemoryDB) ExecOrQueryContext(ctx context.Context, sql string) (out string, err error) {
	out, ok := m.InOutMap[sql]
	if !ok {
		err = errors.Errorf("not found by sql:%s", sql)
		return "", err
	}
	return out, nil
}
