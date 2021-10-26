major-imp/descriptive-scorers:

**Problem**: https://github.com/ful09003/cards/issues/1

Thoughts:
  - Wrap scoring fns in a struct with extra metadata (e.g. name, purpose)
  - Refactor Evaluate() to 


Example scorers that make real-life sense:
- "Do any labels match a regex (UUID, IP?)"
- "Do any metric families have total label count > some number?"
- "Do any metric families have a name in [some list of known-bad-for-me]?"

Example config:

---
- name: uuids-are-discouraged
  description: "Prefer not emitting UUIDs; if not needed drop the label"
  criticality: 3 # In our fictional scenario, UUIDs are more critically 'bad'
  scoring-method:
    type: regex_scorer
    regex: "some-uuid-regex-here"
- name: racist-terms-are-discouraged
  description: "Prefer not using label name/values including the word 'slave' or 'master'"
  criticality: 10 # Treat this as a much higher-priority concern
  scoring-method:
    type: regex_scorer
    regex: "desired-regex-here"
- name: limit-label-lengths
  description: "Prefer a maximum of 5 labels (including 2 built-in labels)"
  criticality: 1 # In our fictional scenario, this isn't as critical as other concerns
  scoring-method:
    type: label_length_scorer
    max-length: 5
- name: confusing-metric
  description: "Prefer not collecting metrics that we know cause confusion"
  criticality: 10
  scoring-method:
    type: metric_exclusion
    exclude:
      - some_metric
      - some_other_metric