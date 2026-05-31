package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"stzbHelper/model"
	"time"
)

func processSocialPacket(cmdId int, decoded []byte) {
	if model.Conn == nil || len(decoded) == 0 {
		return
	}

	switch cmdId {
	case 2100:
		saveChatMessage(decoded)
	case 2200:
		saveWorldAnnouncement(decoded)
	}
}

func saveChatMessage(decoded []byte) {
	var raw []any
	if err := json.Unmarshal(trimDecodedSocial(decoded), &raw); err != nil {
		log.Printf("parse chat packet failed: %v", err)
		return
	}
	if len(raw) < 10 {
		return
	}

	msg := model.ChatMessage{
		MessageID: asInt64Social(raw[0]),
		UserID:    asInt64Social(raw[3]),
		Player:    asStringSocial(raw[4]),
		Content:   asStringSocial(raw[5]),
		Time:      asInt64Social(raw[6]),
		Channel:   asInt64Social(raw[7]),
		UnionName: asStringSocial(raw[8]),
		UnionID:   asInt64Social(raw[9]),
		Raw:       compactRawSocial(raw),
	}
	if len(raw) > 10 {
		msg.Voice = compactRawSocial(raw[10])
	}
	if len(raw) > 13 {
		msg.ServerID = asInt64Social(raw[13])
	}
	if len(raw) > 45 {
		msg.RoleTag = asStringSocial(raw[45])
	}
	if msg.MessageID == 0 || msg.Content == "" {
		return
	}

	var existing model.ChatMessage
	if err := model.Conn.Where("message_id = ?", msg.MessageID).First(&existing).Error; err == nil {
		msg.ID = existing.ID
	}
	if err := model.Conn.Save(&msg).Error; err != nil {
		log.Printf("save chat message failed: %v", err)
	}
}

func saveWorldAnnouncement(decoded []byte) {
	var raw []any
	if err := json.Unmarshal(trimDecodedSocial(decoded), &raw); err != nil {
		log.Printf("parse world announcement failed: %v", err)
		return
	}
	if len(raw) < 3 {
		return
	}

	eventType := asInt64Social(raw[0])
	landLevel := asInt64Social(raw[2])
	if eventType != 1 {
		return
	}

	item := model.WorldAnnouncement{
		Time:      time.Now().Unix(),
		EventType: eventType,
		Player:    asStringSocial(raw[1]),
		LandLevel: landLevel,
		Raw:       compactRawSocial(raw),
	}
	if item.Player == "" || item.LandLevel <= 0 {
		return
	}

	if err := model.Conn.Create(&item).Error; err != nil {
		log.Printf("save world announcement failed: %v", err)
	}
}

func trimDecodedSocial(decoded []byte) []byte {
	return []byte(strings.TrimRight(string(decoded), "\x00\ufeff\ufffd \r\n\t"))
}

func compactRawSocial(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}

func asInt64Social(v any) int64 {
	switch n := v.(type) {
	case float64:
		return int64(n)
	case float32:
		return int64(n)
	case int:
		return int64(n)
	case int64:
		return n
	case json.Number:
		i, _ := n.Int64()
		return i
	case string:
		i, _ := strconv.ParseInt(n, 10, 64)
		return i
	default:
		return 0
	}
}

func asStringSocial(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
