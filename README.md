[![Build Status](https://travis-ci.org/URXtech/planout-golang.svg?branch=master)](https://travis-ci.org/URXtech/planout-golang)

(Multi Variate Testing) Interpreter for PlanOut code written in Golang

Here's an example program that consumes compiled PlanOut code and executes 
the associated experiment using the golang interpreter.


```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/URXtech/planout-golang"
)

// Helper function to generate random string.
func generateString() string {
	s := make([]byte, 10)
	for j := 0; j < 10; j++ {
		s[j] = 'a' + byte(rand.Int()%26)
	}
	return string(s)
}

func main() {
	// Read PlanOut code from file on disk.
	data, _ := ioutil.ReadFile("test/simple_ops.json")

	// The PlanOut code is expected to use json.
	// This format is the same as the output of
	// the PlanOut compiler webapp
	// http://facebook.github.io/planout/demo/planout-compiler.html
	var js map[string]interface{}
	json.Unmarshal(data, &js)

	// Set the necessary input parameters required to run
	// the experiments. For instance, simple_ops.json expects
	// the value for 'userid' to be set.
	params := make(map[string]interface{})
	params["experiment_salt"] = "expt"
	params["userid"] = generateString()

	// Construct an instance of the Interpreter object.
	// Initialize ExperimentSalt and set Inputs to params.
	expt := &goplanout.Interpreter{
		ExperimentSalt: "global_salt",
		Evaluated:      false,
		Inputs:         params,
		Outputs:        map[string]interface{}{},
		Overrides:      map[string]interface{}{},
	}
	
	// Call the Run(...) method on the Interpreter instance.
	// The output of the run will contain the dictionary 
	// of variables and associated values that were evaluated
	// as part of the experiment.
	output, ok := expt.Run(js)
	if !ok {
		fmt.Println("Failed to run the experiment")
	} else {
		fmt.Printf("Params: %v\n", params)
	}
}
```

Suppose we want to run the following experiment:
```go
id = uniformChoice(choices=[1, 2, 3, 4], unit=userid);
```

The PlanOut code generated by the compiler looks like:

```json
{
  "op": "seq",
  "seq": [
    {
      "op": "set",
      "var": "id",
      "value": {
        "choices": {
          "op": "array",
          "values": [
            1,
            2,
            3,
            4
          ]
        },
        "unit": {
          "op": "get",
          "var": "userid"
        },
        "op": "uniformChoice"
      }
    }
  ]
}
```

Each execution of the above experiment will result in setting the variable 'id'. The output to stdout will look like:

```go
Params: map[experiment_salt:expt userid:noocavzddw salt:id id:2]
Params: map[experiment_salt:expt userid:cuncjyqmmz salt:id id:1]
```
