#!/bin/bash

# Search for go.mod files with a replace directive pointing to internal/transport/domain
echo "Searching for go.mod files with conflicting replace directive..."

find . -name "go.mod" -type f -exec grep -l "replace.*github.com/theshop/ai/internal/domain.*internal/transport/domain" {} \;

echo "Done. Check the files listed above and remove or correct the replace directive."
