package model

import "time"

// Lineup 用于记录玩家阵容信息
// 同一玩家的同一阵容（LineupKey）只保留最后一次更新的数据。
type Lineup struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	BattleUID  string    `json:"battle_uid" gorm:"uniqueIndex;size:128"`
	PlayerName string    `json:"player_name" gorm:"index:idx_player_lineup,unique;size:64"`
	PlayerID   string    `json:"player_id" gorm:"size:64"`
	PlayerRole string    `json:"player_role" gorm:"size:16"`
	UnionName  string    `json:"union_name" gorm:"size:64"`
	LineupKey  string    `json:"lineup_key" gorm:"index:idx_player_lineup,unique;size:191"`
	Hero1ID    int64     `json:"hero1_id"`
	Hero1Name  string    `json:"hero1_name" gorm:"size:64"`
	Hero1Level int64     `json:"hero1_level"`
	Hero1Star  int64     `json:"hero1_star"`
	Hero2ID    int64     `json:"hero2_id"`
	Hero2Name  string    `json:"hero2_name" gorm:"size:64"`
	Hero2Level int64     `json:"hero2_level"`
	Hero2Star  int64     `json:"hero2_star"`
	Hero3ID    int64     `json:"hero3_id"`
	Hero3Name  string    `json:"hero3_name" gorm:"size:64"`
	Hero3Level int64     `json:"hero3_level"`
	Hero3Star  int64     `json:"hero3_star"`
	BattleID   int64     `json:"battle_id"`
	RecordTime int64     `json:"record_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 覆盖默认表名
func (*Lineup) TableName() string {
	return "lineup"
}
