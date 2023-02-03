package serverless

import (
	"encoding/json"
	"sort"

	"github.com/awslabs/goformation/v4/cloudformation/utils"
)

// StateMachine_Policies is a helper struct that can hold either a String, String, IAMPolicyDocument, SAMPolicyTemplate, or IAMPolicyDocument value
type StateMachine_Policies struct {
	String *string

	StringArray *[]string

	IAMPolicyDocument *StateMachine_IAMPolicyDocument

	IAMPolicyDocumentArray *[]StateMachine_IAMPolicyDocument
	SAMPolicyTemplateArray *[]StateMachine_SAMPolicyTemplate
}

func (r StateMachine_Policies) value() interface{} {
	ret := []interface{}{}

	if r.String != nil {
		ret = append(ret, r.String)
	}

	if r.StringArray != nil {
		ret = append(ret, r.StringArray)
	}

	if r.IAMPolicyDocument != nil {
		ret = append(ret, *r.IAMPolicyDocument)
	}

	if r.IAMPolicyDocumentArray != nil {
		ret = append(ret, r.IAMPolicyDocumentArray)
	}

	if r.SAMPolicyTemplateArray != nil {
		ret = append(ret, r.SAMPolicyTemplateArray)
	}

	sort.Sort(utils.ByJSONLength(ret)) // Heuristic to select best attribute
	if len(ret) > 0 {
		return ret[0]
	}

	return nil
}

func (r StateMachine_Policies) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value())
}

// Hook into the marshaller
func (r *StateMachine_Policies) UnmarshalJSON(b []byte) error {

	// Unmarshal into interface{} to check it's type
	var typecheck interface{}
	if err := json.Unmarshal(b, &typecheck); err != nil {
		return err
	}

	switch val := typecheck.(type) {

	case string:
		r.String = &val

	case []string:
		r.StringArray = &val

	case map[string]interface{}:
		val = val // This ensures val is used to stop an error

		json.Unmarshal(b, &r.IAMPolicyDocument)

	case []interface{}:

		json.Unmarshal(b, &r.StringArray)

		json.Unmarshal(b, &r.IAMPolicyDocumentArray)

		json.Unmarshal(b, &r.SAMPolicyTemplateArray)

	}

	return nil
}
