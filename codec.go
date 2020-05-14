package main

import "fmt"

func (t GlobalMessageType) MarshalJSON() ([]byte, error) {
	name, ok := GlobalMessageType_Names[t]
	if !ok {
		return nil, ErrorTypeNotDefined
	}
	return []byte(fmt.Sprintf("\"%s\"", name)), nil
}
