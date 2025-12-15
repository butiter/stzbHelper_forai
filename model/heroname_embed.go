package model

import (
	_ "embed"
	"encoding/json"
	"log"
	"strconv"
)

//go:embed herocfg.json
var heroCfgJSON []byte

func init() {
	if len(heroCfgJSON) == 0 {
		return
	}

	type heroCfg struct {
		Name string `json:"name"`
	}

	var raw map[string]heroCfg
	if err := json.Unmarshal(heroCfgJSON, &raw); err != nil {
		log.Printf("解析武将名称映射失败: %v", err)
		return
	}

	HeroNameMap = make(map[int64]string, len(raw))
	for idStr, cfg := range raw {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			HeroNameMap[id] = cfg.Name
		}
	}
}
