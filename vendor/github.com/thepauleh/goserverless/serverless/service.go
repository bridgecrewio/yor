package serverless

import (
	"encoding/json"
	"fmt"
)

// Service is the Serverless Service definition
type Service struct {
	Name         string
	AwsKmsKeyArn string
}

func (s *Service) UnmarshalJSON(data []byte) error {
	// From Unmarshaler doc, []byte("null") is a no-op.
	if string(data) == "null" {
		return nil
	}
	var service interface{}
	if err := json.Unmarshal(data, &service); err != nil {
		return err
	}
	switch service.(type){
	case string:
		s.Name = service.(string)
	case map[string]interface{}:
		serviceMap := service.(map[string]interface{})
		s.Name = serviceMap["name"].(string)
		if serviceMap["awsKmsKeyArn"] != nil {
			s.AwsKmsKeyArn = serviceMap["awsKmsKeyArn"].(string)
		}
	default:
		return fmt.Errorf("unable to parse %s as service data", string(data))
	}
	return nil
}

func (s *Service) MarshalJSON() ([]byte, error) {
	if len(s.AwsKmsKeyArn) == 0 {
		return []byte(fmt.Sprintf(`"%s"`, s.Name)), nil
	}
	return json.Marshal(map[string]string{"name": s.Name, "awsKmsKeyArn": s.AwsKmsKeyArn})
}

