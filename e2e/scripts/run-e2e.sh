#!/usr/bin/env bash

echo "Generating temporary cards.yml..."

TMPD=$(mktemp -d)
echo "Made $TMPD"

cat > cards.yml <<EOF
---
output: table
scorers:
  - name: Validate node_exporter_build_info exists
    description: "node_exporter_build_info should always exist for us to test against"
    criticality: 100
    scoring-method:
      type: family_name_scorer
      criteria:
        - node_exporter_build_info
EOF

cp ../bin/cards $TMPD
OUTPUT=$($TMPD/cards 2>&1)

FOUND=$(echo "$OUTPUT" | grep -c node_exporter_build_info)
EC=$?

echo "grep rc=$EC (found=$FOUND)"

rm -rf $TMPD cards.yml
exit $EC