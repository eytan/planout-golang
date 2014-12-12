/*
 * Copyright 2014 URX
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package goplanout

import (
	"fmt"
)

type PlanOutCode interface{}

type Interpreter struct {
	ExperimentSalt             string
	Inputs, Outputs, Overrides map[string]interface{}
	Evaluated                  bool
}

func (interpreter *Interpreter) get(name string) (interface{}, bool) {
	value, ok := interpreter.Overrides[name]
	if ok {
		return value, true
	}

	value, ok = interpreter.Inputs[name]
	if ok {
		return value, true
	}

	value, ok = interpreter.Outputs[name]
	if ok {
		return value, true
	}
	return nil, false
}

func (interpreter *Interpreter) set(name string, value interface{}) {
	interpreter.Outputs[name] = value
}

func (interpreter *Interpreter) getOverrides() map[string]interface{} {
	return interpreter.Overrides
}

func (interpreter *Interpreter) hasOverrides(name string) bool {
	_, exists := interpreter.Overrides[name]
	return exists
}

func (interpreter *Interpreter) evaluate(code interface{}) interface{} {

	js, ok := code.(map[string]interface{})
	if ok {
		opptr, exists := isOperator(js)
		if exists {
			return opptr.execute(js, interpreter)
		}
	}

	arr, ok := code.([]interface{})
	if ok {
		v := make([]interface{}, len(arr))
		for i := range arr {
			v[i] = interpreter.evaluate(arr[i])
		}
		return v
	}

	return code
}

func (interpreter *Interpreter) Run(code interface{}) (map[string]interface{}, bool) {

	if interpreter.Evaluated {
		return interpreter.Outputs, true
	}

	defer func() (map[string]interface{}, bool) {
		if r := recover(); r != nil {
			fmt.Println("Recovered ", r)
			return nil, false
		}
		interpreter.Evaluated = true
		return interpreter.Outputs, true
	}()

	interpreter.evaluate(code)

	return interpreter.Outputs, true
}
