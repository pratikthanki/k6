/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2020 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package modules

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
)

const extPrefix string = "k6/x/"

//nolint:gochecknoglobals
var (
	modules = make(map[string]interface{})
	mx      sync.RWMutex
)

// Register the given mod as an external JavaScript module that can be imported
// by name. The name must be unique across all registered modules and must be
// prefixed with "k6/x/", otherwise this function will panic.
func Register(name string, mod interface{}) {
	if !strings.HasPrefix(name, extPrefix) {
		panic(fmt.Errorf("external module names must be prefixed with '%s', tried to register: %s", extPrefix, name))
	}

	mx.Lock()
	defer mx.Unlock()

	if _, ok := modules[name]; ok {
		panic(fmt.Sprintf("module already registered: %s", name))
	}
	modules[name] = mod
}

// Module is the interface js modules should implement in order to get access to the VU
type Module interface {
	// NewModuleInstance will get modules.VU that should provide the module with a way to interact with the VU
	// This method will be called for *each* require/import and should return an unique instance for each call
	NewModuleInstance(VU) Instance
}

// GetJSModules returns a map of all registered js modules
func GetJSModules() map[string]interface{} {
	mx.Lock()
	defer mx.Unlock()
	result := make(map[string]interface{}, len(modules))

	for name, module := range modules {
		result[name] = module
	}

	return result
}

// Instance is what a module needs to return
type Instance interface {
	Exports() Exports
}

func getInterfaceMethods() []string {
	var t Instance
	T := reflect.TypeOf(&t).Elem()
	result := make([]string, T.NumMethod())

	for i := range result {
		result[i] = T.Method(i).Name
	}

	return result
}

// VU gives access to the currently executing VU to a module Instance
type VU interface {
	// Context return the context.Context about the current VU
	Context() context.Context

	// InitEnv returns common.InitEnvironment instance if present
	InitEnv() *common.InitEnvironment

	// State returns lib.State if any is present
	State() *lib.State

	// Runtime returns the goja.Runtime for the current VU
	Runtime() *goja.Runtime
}

// Exports is representation of ESM exports of a module
type Exports struct {
	// Default is what will be the `default` export of a module
	Default interface{}
	// Named is the named exports of a module
	Named map[string]interface{}
}
