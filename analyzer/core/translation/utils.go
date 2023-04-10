//  Copyright (c) 2023 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package translation

import (
	"fmt"

	"analyzer/core/mast"
	ts "analyzer/core/treesitter"
)

func childrenNumberError(node *ts.Node) error {
	return fmt.Errorf("node %q has unexpected number of children: %d", node.Type, len(node.Children))
}

func nilChildError(node *ts.Node) error {
	return fmt.Errorf("unexpected nil child node %q", node.Type)
}

func nodeTypeError(node mast.Node) error {
	return fmt.Errorf("node has unexpected type: %T", node)
}
