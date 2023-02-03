# Go Serverless
[![Build Status](https://travis-ci.org/thepauleh/goserverless.svg?branch=master)](https://travis-ci.org/thepauleh/goserverless)   

GoFormation for users wishing to transfer to serverless! 

`GoFormation` is a Go library for working with AWS CloudFormation / AWS Serverless Application Model (SAM) templates. 
- [AWS GoFormation](#aws-goformation)
    - [Main features](#main-features)
    - [Installation](#installation)
    - [Usage](#usage)
        - [Marshalling Serverless described with Go structs, into YAML](#marshalling-serverless-with-go-structs-into-yaml)

## Main features

 * Describe Serverless templates as Go objects (structs), and then turn it into YAML.

## Installation

As with other Go libraries, GoFormation can be installed with `go get`.

```
$ go get github.com/thepauleh/goserverless/serverless
```

## Usage

### Marshalling Serverless with Go structs, into YAML

Below is an example of building a CloudFormation template programmatically, then outputting the resulting JSON

```go
package main

import (
	"fmt"

	"github.com/thepauleh/goserverless/serverless"
)

func main() {

	// Create a new CloudFormation template
    template := serverless.NewTemplate("myService")
    
    template.Service = "myService"

    template.Provider = &serverless.Provider{
        Name: "aws",
        Runtime: "nodejs6.10",
        MemorySize: 512,
        Timeout: 10,
        VersionFunctions: false
    }

	// An example function
	template.Functions["users"] = &serverless.AWSServerlessFunction{
		Handler: "service.o",
        Name:   "${self:provider.stage}-users",
        Description: "Description of what the lambda function does",
        Runtime: "go1.x",
        MemorySize: 128,
        ReservedConcurrency: 5,
        Timeout: 30,
		Events: []serverless.Events{
			serverless.HttpEvent{
				Path: "users/create",
				Method: "post",
			},
		},
	}

	y, err := template.YAML()
	if err != nil {
		fmt.Printf("Failed to generate YAML: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(y))
	}
}
```

Would output the following YAML template:

```yaml
service: myService

provider:
  name: aws
  runtime: nodejs6.10
  memorySize: 512
  timeout: 10
  versionFunctions: false

functions:
  users:
    handler: service.o 
    name: ${self:provider.stage}-users
    description: Description of what the lambda function does
    runtime: go1.x
    memorySize: 128
    timeout: 30
    reservedConcurrency: 5
    events:
      - http:
          path: users/create
          method: post
```
