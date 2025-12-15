package model

// HeroNameMap 供阵容统计使用的武将ID到名称映射。
// 如需补全可在此处追加映射，不影响使用（未命中时会回退到ID字符串）。
var HeroNameMap = map[int64]string{}
