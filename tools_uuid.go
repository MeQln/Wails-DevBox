package main

import (
	"strings"

	"github.com/google/uuid"
)

const uuidMaxCount = 1000

// GenerateUuids 对齐 Tauri uuid::generate_uuids。
// version: 7 → v7（时间排序），其余 → v4。count 钳制到 [1, 1000]。
// hyphen=false 去掉连字符；uppercase=true 转大写。
func (a *App) GenerateUuids(version int, count int, uppercase bool, hyphen bool) []string {
	n := count
	if n < 1 {
		n = 1
	}
	if n > uuidMaxCount {
		n = uuidMaxCount
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		var id uuid.UUID
		if version == 7 {
			var err error
			id, err = uuid.NewV7()
			if err != nil {
				id = uuid.New()
			}
		} else {
			id = uuid.New()
		}
		s := id.String()
		if !hyphen {
			s = strings.ReplaceAll(s, "-", "")
		}
		if uppercase {
			s = strings.ToUpper(s)
		}
		out = append(out, s)
	}
	return out
}
