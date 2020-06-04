package controller

import (
	"github.com/lab/labels-operator/pkg/controller/label"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, label.Add)
}
