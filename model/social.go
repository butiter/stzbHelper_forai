package model

type ChatMessage struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	MessageID int64  `json:"message_id" gorm:"column:message_id;uniqueIndex"`
	UserID    int64  `json:"user_id" gorm:"column:user_id;index"`
	Player    string `json:"player" gorm:"column:player;index"`
	Content   string `json:"content" gorm:"column:content"`
	Time      int64  `json:"time" gorm:"column:time;index"`
	Channel   int64  `json:"channel" gorm:"column:channel;index"`
	UnionName string `json:"union_name" gorm:"column:union_name;index"`
	UnionID   int64  `json:"union_id" gorm:"column:union_id;index"`
	ServerID  int64  `json:"server_id" gorm:"column:server_id"`
	Voice     string `json:"voice" gorm:"column:voice"`
	RoleTag   string `json:"role_tag" gorm:"column:role_tag"`
	Raw       string `json:"raw" gorm:"column:raw"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}

type WorldAnnouncement struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Time      int64  `json:"time" gorm:"column:time;index"`
	EventType int64  `json:"event_type" gorm:"column:event_type;index"`
	Player    string `json:"player" gorm:"column:player;index"`
	LandLevel int64  `json:"land_level" gorm:"column:land_level;index"`
	Raw       string `json:"raw" gorm:"column:raw"`
}

func (WorldAnnouncement) TableName() string {
	return "world_announcement"
}
