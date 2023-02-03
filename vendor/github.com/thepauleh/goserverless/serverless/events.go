package serverless

import (
	"encoding/json"
	"fmt"
)

// Events definition
type Events struct {
	HTTPEvent            *HTTPEvent            `json:"http,omitempty"`
	S3Event              *S3Event              `json:"s3,omitempty"`
	ScheduleEvent        *ScheduleEvent        `json:"schedule,omitempty"`
	SNSEvent             *SNSEvent             `json:"sns,omitempty"`
	SQSEvent             *SQSEvent             `json:"sqs,omitempty"`
	StreamEvent          *StreamEvent          `json:"stream,omitempty"`
	AlexaSkillEvent      *AlexaEvent           `json:"alexaSkill,omitempty"`
	AlexaSmartHomeEvent  *AlexaEvent           `json:"alexaSmartHome,omitempty"`
	IOTEvent             *IOTEvent             `json:"iot,omitempty"`
	CloudwatchEvent      *CloudwatchEvent      `json:"cloudwatchEvent,omitempty"`
	CloudwatchLogEvent   *CloudwatchLogEvent   `json:"cloudwatchLog,omitempty"`
	CognitoUserPoolEvent *CognitoUserPoolEvent `json:"cognitoUserPool,omitempty"`
}

// HTTPEvent definition
type HTTPEvent struct {
	Path       string                 `json:"path,omitempty"`
	Method     string                 `json:"method,omitempty"`
	Cors       bool                   `json:"cors,omitempty"`
	Private    bool                   `json:"private,omitempty"`
	Authorizer map[string]interface{} `json:"authorizer,omitempty"`
}

func (h *HTTPEvent) UnmarshalJSON(bytes []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(bytes, &event); err != nil {
		return err
	}
	if event["path"] != nil {
		h.Path = event["path"].(string)
	}
	if event["method"] != nil {
		h.Method = event["method"].(string)
	}
	if event["cors"] != nil {
		h.Cors = event["cors"].(bool)
	}
	if event["private"] != nil {
		h.Cors = event["private"].(bool)
	}
	if event["authorizer"] != nil {
		switch v:= event["authorizer"].(type) {
		case string:
			h.Authorizer = map[string]interface{}{"name": v}
		case map[string]interface{}:
			h.Authorizer = v
		default:
			return fmt.Errorf("http event authorizer set, but unparseable")
		}
	}
	return nil
}

// S3Event definition
type S3Event struct {
	Bucket string      `json:"bucket,omitempty"`
	Event  string      `json:"event,omitempty"`
	Rules  interface{} `json:"rules,omitempty"`
}

func (s *S3Event) UnmarshalJSON(bytes []byte) error {
	var event interface{}
	if err := json.Unmarshal(bytes, &event); err != nil {
		return err
	}
	switch e := event.(type) {
	case string:
		s.Bucket = e
	case map[string]interface{}:
		s.Bucket = e["bucket"].(string)
		s.Event = e["event"].(string)
		s.Rules = e["rules"]
	default:
		return fmt.Errorf("could not unmarshal type serverless.S3Event")
	}
	return nil
}

// ScheduleEvent definition
type ScheduleEvent struct {
	Rate    string                 `json:"rate,omitempty"`
	Enabled bool                   `json:"enabled,omitempty"`
	Input   map[string]interface{} `json:"input,omitempty"`
}

// SNSEvent definition
type SNSEvent struct {
	TopicName   string `json:"topicName,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// SQSEvent definition
type SQSEvent struct {
	ARN       string `json:"arn,omitempty"`
	BatchSize int    `json:"batchSize,omitempty"`
}

// StreamEvent definition
type StreamEvent struct {
	ARN              string `json:"arn,omitempty"`
	BatchSize        int    `json:"batchSize,omitempty"`
	StartingPosition string `json:"startingPosition,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
}

// AlexaEvent defines a AlexaSkillEvent or AlexaSmartHomeEvent
type AlexaEvent struct {
	AppID   string `json:"appId,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// IOTEvent definition
type IOTEvent struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	SQL         string `json:"sql,omitempty"`
	SQLVersion  string `json:"sqlVersion,omitempty"`
}

// CloudwatchEvent definition
type CloudwatchEvent struct {
	Event     interface{}            `json:"event,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
	InputPath string                 `json:"inputPath,omitempty"`
}

// CloudwatchLogEvent definition
type CloudwatchLogEvent struct {
	LogGroup string `json:"logGroup,omitempty"`
	Filter   string `json:"filter,omitempty"`
}

// CognitoUserPoolEvent definition
type CognitoUserPoolEvent struct {
	Pool    string `json:"pool,omitempty"`
	Trigger string `json:"trigger,omitempty"`
}
