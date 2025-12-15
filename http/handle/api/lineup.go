package api

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"stzbHelper/http/common"
	"stzbHelper/model"
)

func LineupList(c *gin.Context) {
	var lineups []model.Lineup

	query := applyLineupFilter(c, model.Conn.Limit(30).Order("updated_at DESC"))
	if next := c.Query("nextid"); next != "" {
		if id, err := strconv.Atoi(next); err == nil && id > 0 {
			query = query.Where("id < ?", id)
		}
	}

	query.Find(&lineups)

	var total int64
	model.Conn.Model(&model.Lineup{}).Count(&total)

	common.Response{Data: gin.H{
		"list":  lineups,
		"total": total,
	}}.Success(c)
}

func ExportLineupCSV(c *gin.Context) {
	var lineups []model.Lineup
	applyLineupFilter(c, model.Conn.Order("updated_at DESC")).Find(&lineups)

	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	_ = writer.Write([]string{"玩家", "同盟", "角色", "阵容", "大营", "中军", "前锋", "战报ID", "记录时间"})

	for _, item := range lineups {
		_ = writer.Write([]string{
			item.PlayerName,
			item.UnionName,
			item.PlayerRole,
			item.LineupKey,
			formatHeroCell(item.Hero1Name, item.Hero1Star, item.Hero1Level),
			formatHeroCell(item.Hero2Name, item.Hero2Star, item.Hero2Level),
			formatHeroCell(item.Hero3Name, item.Hero3Star, item.Hero3Level),
			strconv.FormatInt(item.BattleID, 10),
			time.Unix(item.RecordTime, 0).Format("2006-01-02 15:04:05"),
		})
	}

	writer.Flush()

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=lineup.csv")
	c.String(http.StatusOK, buffer.String())
}

func applyLineupFilter(c *gin.Context, query *gorm.DB) *gorm.DB {
	name := c.Query("playername")
	union := c.Query("unionname")
	lineup := c.Query("lineup")
	role := c.Query("role")

	if name != "" {
		query = query.Where("player_name LIKE ?", "%"+name+"%")
	}
	if union != "" {
		query = query.Where("union_name LIKE ?", "%"+union+"%")
	}
	if lineup != "" {
		query = query.Where("lineup_key LIKE ?", "%"+lineup+"%")
	}
	if role != "" {
		query = query.Where("player_role = ?", role)
	}

	return query
}

func formatHeroCell(name string, star, level int64) string {
	display := name
	if display == "" {
		display = "-"
	}
	return display + " " + strconv.FormatInt(star, 10) + "红 " + strconv.FormatInt(level, 10) + "级"
}
