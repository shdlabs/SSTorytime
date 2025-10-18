# N4L Quick Reference Card for Structured Data

## Arrow Selection Guide

### CN-2: Contains/Containment

```
(contain) / (belong)         # Primary: Object contains key, array contains element
(has-pt) / (pt-of)          # Part-whole: Entity has part, part is part of entity
(consists) / (mkpt)         # Composition: Consists of, makes up part of
(setof) / (in-set)          # Collections: Set has member, member in set
(grp-of) / (in-grp)         # Groups: Group has element, element in group
(has-memb) / (is-memb)      # Membership: Has member, is member of
```

### EP-3: Expression/Property

```
(hasX) / (isXof)            # Generic property: Has value, is value of
(propt) / (propt-of)        # Property: Has property, is property of
(attr) / (attr-of)          # Attribute: Has attribute, is attribute of
(name) / (nom-de)           # Name: Has name, is name of
(ident) / (indent-for)      # Identifier: Has ID, is ID for
(def) / (def-of)            # Definition: Defined as, is definition of
```

### LT-1: Logic/Temporal

```
(then) / (before)           # Sequence: Comes after, comes before (use _sequence_)
(dep-on) / (sust)          # Dependency: Depends on, sustains
(leads to) / (comes from)   # Causal: Leads to, comes from
```

---

## Quick Patterns

### JSON Object

```json
{ "key": "value" }
```

```n4l
object
  " (contain) key
    " (contain) value
```

### JSON Array

```json
{ "list": [1, 2, 3] }
```

```n4l
list
  " (contain) "list[0]"
    " (contain) 1
  " (contain) "list[1]"
    " (contain) 2
  " (contain) "list[2]"
    " (contain) 3
```

### Nested Object

```json
{ "user": { "name": "Alice", "age": 30 } }
```

```n4l
user
  " (contain) name
    " (contain) Alice
  " (contain) age
    " (contain) 30
```

### XML Element

```xml
<person id="123">
  <name>Alice</name>
</person>
```

```n4l
@person_123 person
  " (ident) 123
  " (contain) name
    " (contain) Alice
```

---

## Syntax Cheat Sheet

```n4l
# Comment

- Section Title

:: context, tags ::                    # Set context
+:: more, tags ::                      # Add to context
-:: remove, tags ::                    # Remove from context

Statement                              # Simple node
Statement (arrow) Target               # Relationship
Statement (arrow1) Target (arrow2) More # Chain

" (arrow) Continuation                 # Ditto: reference previous
$1 (arrow) Reference                   # Reference position 1
$2 (arrow) Reference                   # Reference position 2

@alias Name                            # Create alias
$alias.1 (arrow) Target                # Use alias

"Quoted text with spaces"              # Quote if spaces/special chars
'Single quotes for "nested quotes"'    # Alternative quoting

:: _sequence_ ::                       # Start sequence mode
# Statements linked by (then)
-:: _sequence_ ::                      # End sequence mode
```

---

## Common Anti-Patterns

### ❌ DON'T

```n4l
# Custom arrows not in SSTconfig
user !has-property! name

# Missing context
user
password

# Over-nesting
a (contain) b (contain) c (contain) d (contain) e

# Sequence for unordered data
:: _sequence_ ::
user
name
email
```

### ✅ DO

```n4l
# Stable arrows from SSTconfig
user (attr) name

# Clear context
:: user, security ::
user
  " (contain) password

# Flatten with relationships
a (contain) b
b (pt-of) a
c (belong) b

# Proper containment for attributes
user
  " (contain) name
  " (contain) email
```

---

## Decision Tree

```
What relationship?
│
├─ Container? → CN-2
│  ├─ Object/Array → (contain)/(belong)
│  ├─ Part-Whole → (has-pt)/(pt-of)
│  └─ Set/Group → (setof)/(in-set)
│
├─ Property? → EP-3
│  ├─ Attribute → (attr)/(attr-of)
│  ├─ Value → (hasX)/(isXof)
│  └─ Name/ID → (name)/(ident)
│
├─ Sequential? → LT-1 or _sequence_
│  ├─ Ordered → use :: _sequence_ ::
│  └─ Causal → (leads to)/(comes from)
│
└─ Associated? → NR-0
   └─ Related → (near)/(with)
```

---

## Format Mapping

| Input Format  | Primary Arrow | Secondary  | Context       |
| ------------- | ------------- | ---------- | ------------- |
| JSON object   | (contain)     | (belong)   | json, api     |
| JSON array    | (contain)     | (in-set)   | collection    |
| XML element   | (contain)     | (pt-of)    | xml, document |
| XML attribute | (attr)        | (hasX)     | metadata      |
| YAML mapping  | (contain)     | (belong)   | yaml, config  |
| YAML sequence | (contain)     | (in-set)   | list          |
| HTML element  | (contain)     | (consists) | html, web     |
| CSV row       | (is-memb)     | (attr)     | data, table   |

---

## Tools

```bash
# Convert JSON to N4L
json2n4l -pretty input.json

# Validate N4L
N4L -v file.n4l

# Upload to database
N4L -u file.n4l

# Extract from text
text2N4L -% 50 document.txt
```

---

## Resources

- **Full Guide**: `docs/N4L_STRUCTURED_DATA_GUIDE.md`
- **Language Spec**: `docs/N4L.md`
- **Arrow Definitions**: `SSTconfig/arrows-*.sst`
- **Examples**: `examples/*.n4l`
