package domain

import "time"

type RequestLog struct {
	Method     string    `json:"method" bson:"method"`
	Path       string    `json:"path" bson:"path"`
	StatusCode int       `json:"status_code" bson:"status_code"`
	Duration   int64     `json:"duration_ms" bson:"duration_ms"`
	IP         string    `json:"ip" bson:"ip"`
	UserAgent  string    `json:"user_agent" bson:"user_agent"`
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`
}

func (r *RequestLog) BeforeInsert() {
	r.Timestamp = time.Now()
}
