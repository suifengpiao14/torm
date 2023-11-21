package tormstream

import (
	"context"
	"encoding/json"

	"github.com/suifengpiao14/stream"
)

// TormPackHandler 执行模板返回SQL
func TormPackHandler(torm TormI) (packHandler stream.PackHandler) {
	packHandler = stream.NewPackHandler(func(ctx context.Context, input []byte) (out []byte, err error) {
		volume := make(VolumeMap)
		err = json.Unmarshal(input, &volume)
		if err != nil {
			return nil, err
		}
		sqls, _, _, err := GetSQL(torm.Identity(), torm.TplName(), &volume)
		if err != nil {
			return nil, err
		}
		out = []byte(sqls)
		return out, nil
	}, nil)
	return packHandler
}
