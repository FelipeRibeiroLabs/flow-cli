/*
 * Flow CLI
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package project

import (
	"fmt"
	"path"

	"github.com/onflow/flow-go-sdk"
)

// ImportReplacer implements file import replacements functionality for the project contracts with optionally included aliases.
type ImportReplacer struct {
	contracts []*Contract
	aliases   Aliases
}

func NewImportReplacer(contracts []*Contract, aliases Aliases) *ImportReplacer {
	return &ImportReplacer{
		contracts: contracts,
		aliases:   aliases,
	}
}

func (i *ImportReplacer) Replace(program *Program) (*Program, error) {
	imports := program.Imports()
	contractsLocations := i.getContractsLocations()

	for _, imp := range imports {
		importLocation := path.Clean(absolutePath(program.Location(), imp))
		target, found := contractsLocations[importLocation]
		if !found {
			return nil, fmt.Errorf("import %s could not be resolved from the configuration", imp)
		}
		program.ReplaceImport(imp, target)
	}

	return program, nil
}

// getContractsLocations return a map with contract locations as keys and addresses where they are deployed as values.
func (i *ImportReplacer) getContractsLocations() map[string]string {
	sourceTarget := make(map[string]string)
	for _, contract := range i.contracts {
		sourceTarget[path.Clean(contract.Location())] = contract.AccountAddress.String()
	}

	for source, target := range i.aliases {
		sourceTarget[path.Clean(source)] = flow.HexToAddress(target).String()
	}

	return sourceTarget
}

func absolutePath(basePath, relativePath string) string {
	return path.Join(path.Dir(basePath), relativePath)
}
