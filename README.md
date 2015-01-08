[![Build Status](https://travis-ci.org/URXtech/planout-golang.svg?branch=master)](https://travis-ci.org/URXtech/planout-golang)

(Multi Variate Testing) Interpreter for PlanOut code written in Golang

Suppose we have a PlanOut experiment that randomly assigns users to a ranking
function:
```go
ranking = uniformChoice(choices=["relevance", "most_recent", "popularity"], unit=userid);
```
After [compiling the above PlanOut script into JSON](http://facebook.github.io/planout/docs/getting-started-with-the-interpreter.html),
and saving the file into `test/simple_ops.json`, we could execute the serialized
PlanOut code using the following Go code:

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/URXtech/planout-golang"
)

func main() {
	// Read PlanOut code from file on disk.
	data, _ := ioutil.ReadFile("test/simple_ops.json")

	// The PlanOut code is expected to use JSON, see the PlanOut compiler, see
	// http://facebook.github.io/planout/docs/planout-language.html for details
	var serialized_code map[string]interface{}
	json.Unmarshal(data, &serialized_code)

	// Initialize the expected input for the script, 'userid'
	// note that experiment_salt doesn't belong in the inputs.
	input := map[string]interface{} {"userid": 42, "experiment_salt" : "first_experiment"}

	// Construct an instance of the Interpreter object, which
	// includes the experiment-level salt and inputs.
	expt := &goplanout.Interpreter{
		Salt: "global_salt",
		Evaluated:      false,
		Inputs:         input,
		Outputs:        map[string]interface{}{},
		Overrides:      map[string]interface{}{},
	}
	
	// Call the Run(...) method on the Interpreter instance.
	// The output of the run will contain the dictionary 
	// of variables and associated values that were evaluated
	// as part of the experiment.
	output, ok := expt.Run(serialized_code)
	if !ok {
		fmt.Println("Failed to run the experiment")
	} else {
		fmt.Printf("Params: %v\n", params)
	}
}
```

Each execution of the above experiment will result in setting the variable `ranking`. The output to stdout will look like:

```go
Params: map[experiment_salt:expt userid:noocavzddw salt:id ranking:"most_recent"]
```
