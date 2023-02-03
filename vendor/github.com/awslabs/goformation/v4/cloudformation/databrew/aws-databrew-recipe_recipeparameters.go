package databrew

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Recipe_RecipeParameters AWS CloudFormation Resource (AWS::DataBrew::Recipe.RecipeParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html
type Recipe_RecipeParameters struct {

	// AggregateFunction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-aggregatefunction
	AggregateFunction string `json:"AggregateFunction,omitempty"`

	// Base AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-base
	Base string `json:"Base,omitempty"`

	// CaseStatement AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-casestatement
	CaseStatement string `json:"CaseStatement,omitempty"`

	// CategoryMap AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-categorymap
	CategoryMap string `json:"CategoryMap,omitempty"`

	// CharsToRemove AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-charstoremove
	CharsToRemove string `json:"CharsToRemove,omitempty"`

	// CollapseConsecutiveWhitespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-collapseconsecutivewhitespace
	CollapseConsecutiveWhitespace string `json:"CollapseConsecutiveWhitespace,omitempty"`

	// ColumnDataType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-columndatatype
	ColumnDataType string `json:"ColumnDataType,omitempty"`

	// ColumnRange AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-columnrange
	ColumnRange string `json:"ColumnRange,omitempty"`

	// Count AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-count
	Count string `json:"Count,omitempty"`

	// CustomCharacters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-customcharacters
	CustomCharacters string `json:"CustomCharacters,omitempty"`

	// CustomStopWords AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-customstopwords
	CustomStopWords string `json:"CustomStopWords,omitempty"`

	// CustomValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-customvalue
	CustomValue string `json:"CustomValue,omitempty"`

	// DatasetsColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-datasetscolumns
	DatasetsColumns string `json:"DatasetsColumns,omitempty"`

	// DateAddValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-dateaddvalue
	DateAddValue string `json:"DateAddValue,omitempty"`

	// DateTimeFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-datetimeformat
	DateTimeFormat string `json:"DateTimeFormat,omitempty"`

	// DateTimeParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-datetimeparameters
	DateTimeParameters string `json:"DateTimeParameters,omitempty"`

	// DeleteOtherRows AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-deleteotherrows
	DeleteOtherRows string `json:"DeleteOtherRows,omitempty"`

	// Delimiter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-delimiter
	Delimiter string `json:"Delimiter,omitempty"`

	// EndPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-endpattern
	EndPattern string `json:"EndPattern,omitempty"`

	// EndPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-endposition
	EndPosition string `json:"EndPosition,omitempty"`

	// EndValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-endvalue
	EndValue string `json:"EndValue,omitempty"`

	// ExpandContractions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-expandcontractions
	ExpandContractions string `json:"ExpandContractions,omitempty"`

	// Exponent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-exponent
	Exponent string `json:"Exponent,omitempty"`

	// FalseString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-falsestring
	FalseString string `json:"FalseString,omitempty"`

	// GroupByAggFunctionOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-groupbyaggfunctionoptions
	GroupByAggFunctionOptions string `json:"GroupByAggFunctionOptions,omitempty"`

	// GroupByColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-groupbycolumns
	GroupByColumns string `json:"GroupByColumns,omitempty"`

	// HiddenColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-hiddencolumns
	HiddenColumns string `json:"HiddenColumns,omitempty"`

	// IgnoreCase AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-ignorecase
	IgnoreCase string `json:"IgnoreCase,omitempty"`

	// IncludeInSplit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-includeinsplit
	IncludeInSplit string `json:"IncludeInSplit,omitempty"`

	// Input AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-input
	Input interface{} `json:"Input,omitempty"`

	// Interval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-interval
	Interval string `json:"Interval,omitempty"`

	// IsText AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-istext
	IsText string `json:"IsText,omitempty"`

	// JoinKeys AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-joinkeys
	JoinKeys string `json:"JoinKeys,omitempty"`

	// JoinType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-jointype
	JoinType string `json:"JoinType,omitempty"`

	// LeftColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-leftcolumns
	LeftColumns string `json:"LeftColumns,omitempty"`

	// Limit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-limit
	Limit string `json:"Limit,omitempty"`

	// LowerBound AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-lowerbound
	LowerBound string `json:"LowerBound,omitempty"`

	// MapType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-maptype
	MapType string `json:"MapType,omitempty"`

	// ModeType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-modetype
	ModeType string `json:"ModeType,omitempty"`

	// MultiLine AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-multiline
	MultiLine bool `json:"MultiLine,omitempty"`

	// NumRows AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-numrows
	NumRows string `json:"NumRows,omitempty"`

	// NumRowsAfter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-numrowsafter
	NumRowsAfter string `json:"NumRowsAfter,omitempty"`

	// NumRowsBefore AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-numrowsbefore
	NumRowsBefore string `json:"NumRowsBefore,omitempty"`

	// OrderByColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-orderbycolumn
	OrderByColumn string `json:"OrderByColumn,omitempty"`

	// OrderByColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-orderbycolumns
	OrderByColumns string `json:"OrderByColumns,omitempty"`

	// Other AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-other
	Other string `json:"Other,omitempty"`

	// Pattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-pattern
	Pattern string `json:"Pattern,omitempty"`

	// PatternOption1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-patternoption1
	PatternOption1 string `json:"PatternOption1,omitempty"`

	// PatternOption2 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-patternoption2
	PatternOption2 string `json:"PatternOption2,omitempty"`

	// PatternOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-patternoptions
	PatternOptions string `json:"PatternOptions,omitempty"`

	// Period AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-period
	Period string `json:"Period,omitempty"`

	// Position AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-position
	Position string `json:"Position,omitempty"`

	// RemoveAllPunctuation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeallpunctuation
	RemoveAllPunctuation string `json:"RemoveAllPunctuation,omitempty"`

	// RemoveAllQuotes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeallquotes
	RemoveAllQuotes string `json:"RemoveAllQuotes,omitempty"`

	// RemoveAllWhitespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeallwhitespace
	RemoveAllWhitespace string `json:"RemoveAllWhitespace,omitempty"`

	// RemoveCustomCharacters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removecustomcharacters
	RemoveCustomCharacters string `json:"RemoveCustomCharacters,omitempty"`

	// RemoveCustomValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removecustomvalue
	RemoveCustomValue string `json:"RemoveCustomValue,omitempty"`

	// RemoveLeadingAndTrailingPunctuation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeleadingandtrailingpunctuation
	RemoveLeadingAndTrailingPunctuation string `json:"RemoveLeadingAndTrailingPunctuation,omitempty"`

	// RemoveLeadingAndTrailingQuotes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeleadingandtrailingquotes
	RemoveLeadingAndTrailingQuotes string `json:"RemoveLeadingAndTrailingQuotes,omitempty"`

	// RemoveLeadingAndTrailingWhitespace AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeleadingandtrailingwhitespace
	RemoveLeadingAndTrailingWhitespace string `json:"RemoveLeadingAndTrailingWhitespace,omitempty"`

	// RemoveLetters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removeletters
	RemoveLetters string `json:"RemoveLetters,omitempty"`

	// RemoveNumbers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removenumbers
	RemoveNumbers string `json:"RemoveNumbers,omitempty"`

	// RemoveSourceColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removesourcecolumn
	RemoveSourceColumn string `json:"RemoveSourceColumn,omitempty"`

	// RemoveSpecialCharacters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-removespecialcharacters
	RemoveSpecialCharacters string `json:"RemoveSpecialCharacters,omitempty"`

	// RightColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-rightcolumns
	RightColumns string `json:"RightColumns,omitempty"`

	// SampleSize AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-samplesize
	SampleSize string `json:"SampleSize,omitempty"`

	// SampleType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sampletype
	SampleType string `json:"SampleType,omitempty"`

	// SecondInput AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-secondinput
	SecondInput string `json:"SecondInput,omitempty"`

	// SecondaryInputs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-secondaryinputs
	SecondaryInputs []Recipe_SecondaryInput `json:"SecondaryInputs,omitempty"`

	// SheetIndexes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sheetindexes
	SheetIndexes []int `json:"SheetIndexes,omitempty"`

	// SheetNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sheetnames
	SheetNames []string `json:"SheetNames,omitempty"`

	// SourceColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumn
	SourceColumn string `json:"SourceColumn,omitempty"`

	// SourceColumn1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumn1
	SourceColumn1 string `json:"SourceColumn1,omitempty"`

	// SourceColumn2 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumn2
	SourceColumn2 string `json:"SourceColumn2,omitempty"`

	// SourceColumns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-sourcecolumns
	SourceColumns string `json:"SourceColumns,omitempty"`

	// StartColumnIndex AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startcolumnindex
	StartColumnIndex string `json:"StartColumnIndex,omitempty"`

	// StartPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startpattern
	StartPattern string `json:"StartPattern,omitempty"`

	// StartPosition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startposition
	StartPosition string `json:"StartPosition,omitempty"`

	// StartValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-startvalue
	StartValue string `json:"StartValue,omitempty"`

	// StemmingMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stemmingmode
	StemmingMode string `json:"StemmingMode,omitempty"`

	// StepCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stepcount
	StepCount string `json:"StepCount,omitempty"`

	// StepIndex AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stepindex
	StepIndex string `json:"StepIndex,omitempty"`

	// StopWordsMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-stopwordsmode
	StopWordsMode string `json:"StopWordsMode,omitempty"`

	// Strategy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-strategy
	Strategy string `json:"Strategy,omitempty"`

	// TargetColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetcolumn
	TargetColumn string `json:"TargetColumn,omitempty"`

	// TargetColumnNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetcolumnnames
	TargetColumnNames string `json:"TargetColumnNames,omitempty"`

	// TargetDateFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetdateformat
	TargetDateFormat string `json:"TargetDateFormat,omitempty"`

	// TargetIndex AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-targetindex
	TargetIndex string `json:"TargetIndex,omitempty"`

	// TimeZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-timezone
	TimeZone string `json:"TimeZone,omitempty"`

	// TokenizerPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-tokenizerpattern
	TokenizerPattern string `json:"TokenizerPattern,omitempty"`

	// TrueString AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-truestring
	TrueString string `json:"TrueString,omitempty"`

	// UdfLang AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-udflang
	UdfLang string `json:"UdfLang,omitempty"`

	// Units AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-units
	Units string `json:"Units,omitempty"`

	// UnpivotColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-unpivotcolumn
	UnpivotColumn string `json:"UnpivotColumn,omitempty"`

	// UpperBound AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-upperbound
	UpperBound string `json:"UpperBound,omitempty"`

	// UseNewDataFrame AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-usenewdataframe
	UseNewDataFrame string `json:"UseNewDataFrame,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-value
	Value string `json:"Value,omitempty"`

	// Value1 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-value1
	Value1 string `json:"Value1,omitempty"`

	// Value2 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-value2
	Value2 string `json:"Value2,omitempty"`

	// ValueColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-valuecolumn
	ValueColumn string `json:"ValueColumn,omitempty"`

	// ViewFrame AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-databrew-recipe-recipeparameters.html#cfn-databrew-recipe-recipeparameters-viewframe
	ViewFrame string `json:"ViewFrame,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Recipe_RecipeParameters) AWSCloudFormationType() string {
	return "AWS::DataBrew::Recipe.RecipeParameters"
}
