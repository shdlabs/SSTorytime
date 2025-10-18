# N4L Guidelines for Structured Data

## A Comprehensive Guide to Converting Structured Data to N4L Format

This guide provides best practices for converting structured data formats (JSON, XML, HTML, YAML, etc.) into N4L (Notes for Learning) format using stable semantic arrows from SSTconfig.

---

## Table of Contents

1. [Core Principles](#core-principles)
2. [N4L Syntax Fundamentals](#n4l-syntax-fundamentals)
3. [Stable Arrow Types (SSTconfig)](#stable-arrow-types-sstconfig)
4. [Structured Data Patterns](#structured-data-patterns)
5. [Format-Specific Guidelines](#format-specific-guidelines)
6. [Best Practices](#best-practices)
7. [Common Patterns](#common-patterns)
8. [Anti-Patterns to Avoid](#anti-patterns-to-avoid)

---

## Core Principles

### The N4L Philosophy

1. **Simplicity First**: N4L is intentionally simple to encourage revisiting and organizing notes
2. **Graph-Based**: Everything is a node connected by semantic relationships (arrows)
3. **Context-Aware**: Use context tags `:: context ::` to group related concepts
4. **Sequential When Needed**: Use `_sequence_` mode for ordered relationships
5. **Human-Readable**: The format should be easy to read and edit in any text editor

### Key Design Goals

- **Avoid rigid schemas**: Don't force data into pre-approved models
- **Preserve semantics**: Relationships should be meaningful, not just structural
- **Enable reasoning**: The graph should support logical inference
- **Support iteration**: Easy to refine and reorganize notes

---

## N4L Syntax Fundamentals

### Basic Elements

```n4l
# Comments start with # or //
# Single-line only

- Section or Chapter Title
# Declares a section/chapter as the subject

:: context, tags, here ::
# Persistent context that applies to following statements
# Any number of :: marks is acceptable

+:: more, context ::
# Extend existing context

-:: remove, context ::
# Remove context tags

Statement text
# A simple node - any text not containing "("

Statement (arrow) Target
# A relationship between two nodes

Statement (arrow1) Target (arrow2) Another
# Chained relationships

" (arrow) Continuation
# Ditto mark " references previous statement

$1 (arrow) Reference
# Reference to first item in previous statement

@alias_name
# Create an alias for easy reference

$alias_name.1
# Reference the aliased line
```

### Reserved Symbols

- `()` - Arrows/relationships
- `+`, `-` - Add/remove operations
- `@` - Alias definition
- `$` - Reference (alias or position)
- `#`, `//` - Comments
- `"` - Ditto mark (references previous statement)
- `::` - Context boundaries

### Quoting Rules

```n4l
"Text with spaces needs quotes"
'Can use single quotes for "nested quotes"'
"Quote if text contains: spaces, !, [, ], or special chars"
```

---

## Stable Arrow Types (SSTconfig)

### CN-2: Contains/Containment Relationships

**Primary arrows for structured data:**

```n4l
(contain) / (belong)
# Object contains key, key belongs to object
# Array contains element, element belongs to array

(has-pt) / (pt-of)
# Entity has part, part is part of entity
# Use for compositional relationships

(consists) / (mkpt)
# Consists of / makes up part of
# For components and composition

(setof) / (in-set)
# Collection contains member
# Use for collections and sets

(grp-of) / (in-grp)
# Group contains element
# Use for groupings

(has-memb) / (is-memb)
# Has member / is member of
# For membership relationships
```

**Spatial containment:**

```n4l
(region) / (within)
# Region contains sub-region

(encl) / (is-encl)
# Encloses / is enclosed by

(house) / (is-housed)
# Houses / is housed by
```

### EP-3: Expression/Property Relationships

**For attributes and properties:**

```n4l
(hasX) / (isXof)
# Has value / is value of
# Generic property relationship

(propt) / (propt-of)
# Has property / is property of

(attr) / (attr-of)
# Has attribute / is attribute of

(name) / (nom-de)
# Has name / is name of

(ident) / (indent-for)
# Has identifier / is identifier for

(title) / (title of)
# Has title / is title of

(description) / (descrip-for)
# Has description / is description for
```

**For definitions and types:**

```n4l
(def) / (def-of)
# Defined as / is definition of

(means) / (meansb)
# Means / is meant by

(about) / (theme-of)
# Is about topic / is topic of
```

### LT-1: Logic/Temporal Relationships

**For sequences and ordering:**

```n4l
(then) / (before)
# Sequential relationship
# Reserved for _sequence_ mode

(sh-precede) / (sh-follow)
# Should precede / should follow

(leads to) / (comes from)
# Leads to / comes from
# Causal relationships
```

**For dependencies:**

```n4l
(dep-on) / (sust)
# Depends on / sustains

(det) / (det-by)
# Determines / is determined by

(constr) / (constr-by)
# Constrains / is constrained by
```

### NR-0: Narrative/Near Relationships

**For associations:**

```n4l
(near) / (near)
# Near in semantic space
# Symmetric relationship

(with) / (with)
# Associated with
# Symmetric relationship
```

---

## Structured Data Patterns

### Pattern 1: Object/Dictionary → N4L

**JSON Input:**

```json
{
  "user": {
    "name": "Alice",
    "age": 30,
    "active": true
  }
}
```

**N4L Output:**

```n4l
- User Data

:: user, profile ::

user
  " (contain) name
name
  " (contain) Alice
  " (belong) user

user
  " (contain) age
age
  " (contain) 30
  " (belong) user

user
  " (contain) active
active
  " (contain) true
  " (belong) user
```

**Alternative (more concise):**

```n4l
- User Data

user
  " (contain) name
    " (contain) Alice
  " (contain) age
    " (contain) 30
  " (contain) active
    " (contain) true
```

### Pattern 2: Array/List → N4L

**JSON Input:**

```json
{
  "tags": ["python", "golang", "rust"]
}
```

**N4L Output:**

```n4l
- Tags Collection

:: tags, languages ::

tags (setof) programming languages
  " (contain) "tags[0]"
"tags[0]"
  " (contain) python
  " (in-set) tags

  " (contain) "tags[1]"
"tags[1]"
  " (contain) golang
  " (in-set) tags

  " (contain) "tags[2]"
"tags[2]"
  " (contain) rust
  " (in-set) tags
```

**Sequential alternative:**

```n4l
- Tags Collection

:: _sequence_, tags, languages ::

tags (setof) programming languages

python (in-set) tags
golang (in-set) tags
rust (in-set) tags

-:: _sequence_ ::
```

### Pattern 3: Nested Objects → N4L

**JSON Input:**

```json
{
  "company": {
    "name": "TechCorp",
    "location": {
      "city": "Oslo",
      "country": "Norway"
    }
  }
}
```

**N4L Output:**

```n4l
- Company Information

:: company, organization ::

company
  " (contain) name
    " (contain) TechCorp
  " (contain) location

location (pt-of) company
  " (contain) city
    " (contain) Oslo
  " (contain) country
    " (contain) Norway

TechCorp (name) company
  " (based) Oslo
  " (based) Norway
```

### Pattern 4: Hierarchical Data → N4L

**XML/HTML Input:**

```xml
<book>
  <title>The Origin of Species</title>
  <author>Charles Darwin</author>
  <chapters>
    <chapter number="1">Variation Under Domestication</chapter>
    <chapter number="2">Variation Under Nature</chapter>
  </chapters>
</book>
```

**N4L Output:**

```n4l
- Book Structure

:: book, darwin, evolution ::

@book The Origin of Species
  " (author) Charles Darwin
  " (consists) chapters

chapters (pt-of) $book.1
  " (contain) "chapter 1"
  " (contain) "chapter 2"

"chapter 1"
  " (hasX) 1
  " (title) Variation Under Domestication
  " (pt-of) chapters

"chapter 2"
  " (hasX) 2
  " (title) Variation Under Nature
  " (pt-of) chapters
```

---

## Format-Specific Guidelines

### JSON → N4L

**Key Principles:**

1. **Objects** → Use `(contain)` for keys
2. **Arrays** → Use indexed names `array[0]` and `(in-set)` or `(contain)`
3. **Key-Value** → Key node with `(contain)` to value
4. **Nesting** → Preserve hierarchy with `(pt-of)` and `(belong)`
5. **Types** → Preserve primitive types (string, number, boolean, null)

**Example:**

```json
{
  "api": {
    "version": "1.0",
    "endpoints": ["/users", "/posts"]
  }
}
```

```n4l
- API Specification

api
  " (contain) version
    " (contain) 1.0
  " (contain) endpoints

endpoints (pt-of) api
  " (contain) "endpoints[0]"
    " (contain) /users
  " (contain) "endpoints[1]"
    " (contain) /posts
```

### XML/HTML → N4L

**Key Principles:**

1. **Elements** → Nodes with semantic names
2. **Attributes** → Use `(attr)` or `(hasX)` relationships
3. **Hierarchy** → Use `(contain)`, `(pt-of)`, `(consists)`
4. **Text Content** → Direct relationship to element
5. **Namespaces** → Include in context tags

**Example:**

```xml
<person id="123" status="active">
  <name>Alice</name>
  <email>alice@example.com</email>
</person>
```

```n4l
- Person Record

:: person, contact ::

@person_123 person
  " (ident) 123
  " (attr) status
    " (hasX) active
  " (contain) name
    " (contain) Alice
  " (contain) email
    " (contain) alice@example.com

Alice (name) $person_123.1
alice@example.com (email) $person_123.1
```

### YAML → N4L

**Key Principles:**

1. **Mappings** → Like JSON objects, use `(contain)`
2. **Sequences** → Like JSON arrays, use indexed approach
3. **Anchors/Aliases** → Use N4L `@alias` and `$alias` syntax
4. **Multiline Strings** → Quote properly or use separate nodes
5. **Type Tags** → Include in context if semantically important

**Example:**

```yaml
database:
  host: localhost
  port: 5432
  users:
    - admin
    - readonly
```

```n4l
- Database Configuration

:: database, config ::

database
  " (contain) host
    " (contain) localhost
  " (contain) port
    " (contain) 5432
  " (contain) users

users (pt-of) database
  " (grp-of) administrators and readers
  " (contain) "users[0]"
    " (contain) admin
    " (is-memb) users
  " (contain) "users[1]"
    " (contain) readonly
    " (is-memb) users
```

### CSV/Tabular → N4L

**Key Principles:**

1. **Rows** → Individual entities with `@alias`
2. **Columns** → Properties with `(hasX)` or `(attr)`
3. **Relationships** → Explicit with appropriate arrows
4. **Headers** → Include in context tags
5. **Table Name** → Use as section title

**Example:**

```csv
id,name,department
1,Alice,Engineering
2,Bob,Sales
```

```n4l
- Employee Records

:: employees, staff, company ::

@emp1 Employee 1
  " (ident) 1
  " (name) Alice
  " (works in) Engineering

@emp2 Employee 2
  " (ident) 2
  " (name) Bob
  " (works in) Sales

Engineering (has-memb) Alice
Sales (has-memb) Bob
```

---

## Best Practices

### 1. Choose Appropriate Arrows

**DO:**

```n4l
# Use semantic arrows from SSTconfig
user (contain) profile
profile (pt-of) user
name (attr) user
```

**DON'T:**

```n4l
# Avoid made-up arrows
user (has-a) profile
user (owns) name
```

### 2. Use Context Effectively

**DO:**

```n4l
:: user, authentication, security ::

user
  " (contain) password
  " (contain) token

-:: security ::  # Remove when done with security concepts
```

**DON'T:**

```n4l
# Don't leave context undefined
user
  " (contain) password
```

### 3. Leverage Sequence Mode Appropriately

**DO:**

```n4l
:: _sequence_, tutorial, steps ::

@step1 Install dependencies
@step2 Configure settings
@step3 Run application

-:: _sequence_ ::
```

**DON'T:**

```n4l
# Don't use sequence for unordered relationships
:: _sequence_ ::

user (has-attr) name
user (has-attr) email
```

### 4. Use Aliases for Clarity

**DO:**

```n4l
@user_alice Alice Johnson
  " (ident) alice@example.com
  " (role) administrator

$user_alice.1 (manages) Bob Smith
```

**DON'T:**

```n4l
# Don't repeat long references
Alice Johnson (manages) Bob Smith
Alice Johnson (created) Document
Alice Johnson (approved) Request
```

### 5. Quote Appropriately

**DO:**

```n4l
"array[0]" (contain) "value with spaces"
"special!chars" (note) "needs quotes"
```

**DON'T:**

```n4l
array[0] (contain) value with spaces  # Will break parser
```

### 6. Structure for Readability

**DO:**

```n4l
- Clear Section Title

:: relevant, context, tags ::

# Group related concepts
@main_concept Primary Concept
  " (contain) sub_concept_1
  " (contain) sub_concept_2

# Add spacing between groups

@related_concept Related Idea
  " (related to) $main_concept.1
```

**DON'T:**

```n4l
- section
user
" (has) name
" (has) email
" (has) id
otheruser
" (has) name
```

### 7. Preserve Semantic Meaning

**DO:**

```n4l
# Preserve what the data means
user (name) Alice
  " (role) administrator
  " (status) active
```

**DON'T:**

```n4l
# Don't just mirror structure
user (has) field1
field1 (has) Alice
user (has) field2
field2 (has) administrator
```

---

## Common Patterns

### Pattern: Configuration Files

```n4l
- Application Configuration

:: config, application, settings ::

@app_config application settings
  " (contain) database
  " (contain) logging
  " (contain) security

database (pt-of) $app_config.1
  " (attr) host
    " (hasX) localhost
  " (attr) port
    " (hasX) 5432

logging (pt-of) $app_config.1
  " (attr) level
    " (hasX) INFO
  " (attr) format
    " (hasX) json
```

### Pattern: API Responses

```n4l
- API Response Data

:: api, response, data ::

@response API response
  " (attr) status
    " (hasX) 200
  " (attr) success
    " (hasX) true
  " (contain) data

data (pt-of) $response.1
  " (contain) users

users (setof) user records
  " (contain) "users[0]"
"users[0]"
  " (ident) 1
  " (name) Alice
  " (in-set) users
```

### Pattern: Hierarchical Documents

```n4l
- Document Structure

:: document, hierarchy, structure ::

@doc Main Document
  " (title) System Architecture
  " (author) Engineering Team
  " (consists) sections

sections (pt-of) $doc.1
  " (contain) @sec1
  " (contain) @sec2

@sec1 Introduction
  " (sh-precede) $sec2.1
  " (about) system overview

@sec2 Implementation
  " (sh-follow) $sec1.1
  " (about) technical details
```

### Pattern: Entity Relationships

```n4l
- Entity Relationships

:: entities, relationships, data_model ::

@user_entity User
  " (has-pt) user_id
  " (has-pt) username
  " (has-pt) email

@post_entity Post
  " (has-pt) post_id
  " (has-pt) title
  " (has-pt) content
  " (has-pt) author_id

author_id (pt-of) $post_entity.1
  " (ref) user_id

user_id (pt-of) $user_entity.1
  " (ident-for) User
```

---

## Anti-Patterns to Avoid

### ❌ Don't Create Unstable Custom Arrows

```n4l
# BAD: Custom arrows not in SSTconfig
user !has-property! name
data !contains-key! value
```

**Instead use:**

```n4l
# GOOD: Stable arrows from SSTconfig
user (attr) name
data (contain) value
```

### ❌ Don't Lose Semantic Meaning

```n4l
# BAD: Pure structural mapping
object (has) key1
key1 (is) value1
```

**Instead use:**

```n4l
# GOOD: Semantic relationships
person (name) Alice
  " (role) administrator
```

### ❌ Don't Over-Nest

```n4l
# BAD: Too deep, hard to follow
root
  " (contain) level1
    " (contain) level2
      " (contain) level3
        " (contain) level4
          " (contain) value
```

**Instead use:**

```n4l
# GOOD: Flatten with explicit relationships
root (consist) major_components
level1 (pt-of) root
level2 (pt-of) level1
value (belong) level2
```

### ❌ Don't Ignore Context

```n4l
# BAD: No context
user
password
token
```

**Instead use:**

```n4l
# GOOD: Clear context
:: user, security, authentication ::

user
  " (contain) password
  " (contain) token
```

### ❌ Don't Abuse Sequence Mode

```n4l
# BAD: Using sequence for unordered attributes
:: _sequence_ ::

user
name
email
age
```

**Instead use:**

```n4l
# GOOD: Use proper containment
user
  " (contain) name
  " (contain) email
  " (contain) age
```

---

## Complete Example: Complex JSON to N4L

### Input: API Response JSON

```json
{
  "api": {
    "version": "2.0",
    "timestamp": "2025-10-17T12:00:00Z"
  },
  "data": {
    "users": [
      {
        "id": 1,
        "name": "Alice",
        "email": "alice@example.com",
        "roles": ["admin", "user"]
      },
      {
        "id": 2,
        "name": "Bob",
        "email": "bob@example.com",
        "roles": ["user"]
      }
    ],
    "metadata": {
      "total": 2,
      "page": 1
    }
  }
}
```

### Output: Comprehensive N4L

```n4l
- API User Response

:: api, users, response, v2 ::

# API Metadata
@api_info API Information
  " (hasX) version
    " (contain) 2.0
  " (hasX) timestamp
    " (contain) 2025-10-17T12:00:00Z

# Response Data Structure
@response_data Response Data
  " (contain) users
  " (contain) metadata

# Users Collection
users (pt-of) $response_data.1
  " (setof) user records
  " (contain) @user_1
  " (contain) @user_2

# User 1: Alice
@user_1 "users[0]"
  " (ident) 1
  " (name) Alice
  " (attr) email
    " (hasX) alice@example.com
  " (contain) roles

roles (pt-of) $user_1.1
  " (contain) "roles[0]"
    " (contain) admin
  " (contain) "roles[1]"
    " (contain) user

Alice (is-memb) users
  " (role) admin
  " (role) user

# User 2: Bob
@user_2 "users[1]"
  " (ident) 2
  " (name) Bob
  " (attr) email
    " (hasX) bob@example.com
  " (contain) roles

roles (pt-of) $user_2.1
  " (contain) "roles[0]"
    " (contain) user

Bob (is-memb) users
  " (role) user

# Metadata
metadata (pt-of) $response_data.1
  " (attr) total
    " (hasX) 2
  " (attr) page
    " (hasX) 1

# Semantic Relationships
admin (role) Alice
user (role) Alice
user (role) Bob
```

---

## Summary: Decision Tree for Arrows

```
Is it a container/collection relationship?
├─ Yes → Use CN-2 arrows
│   ├─ Object contains key → (contain) / (belong)
│   ├─ Whole has part → (has-pt) / (pt-of)
│   ├─ Set has member → (setof) / (in-set)
│   └─ Group has element → (grp-of) / (in-grp)
│
└─ No → Is it a property/attribute?
    ├─ Yes → Use EP-3 arrows
    │   ├─ Has property → (propt) / (propt-of)
    │   ├─ Has attribute → (attr) / (attr-of)
    │   ├─ Has value → (hasX) / (isXof)
    │   └─ Has name/identifier → (name) / (ident)
    │
    └─ No → Is it sequential/temporal?
        ├─ Yes → Use LT-1 arrows or _sequence_ mode
        │   ├─ Sequential order → (then) / (before)
        │   ├─ Dependency → (dep-on) / (sust)
        │   └─ Causal → (leads to) / (comes from)
        │
        └─ No → Is it an association?
            └─ Yes → Use NR-0 arrows
                ├─ Semantic proximity → (near) / (near)
                └─ Association → (with) / (with)
```

---

## Tools Reference

- **json2n4l**: Automated JSON to N4L parser (recommended)

  ```bash
  json2n4l -pretty input.json
  ```

- **N4L**: Compiler/validator and database uploader

  ```bash
  N4L -v file.n4l          # Validate
  N4L -u file.n4l          # Upload to database
  ```

- **text2N4L**: Extract significant sentences from plain text
  ```bash
  text2N4L -% 50 document.txt
  ```

---

## Additional Resources

- **SSTconfig/arrows-CN-2.sst**: Contains/Containment arrows
- **SSTconfig/arrows-EP-3.sst**: Expression/Property arrows
- **SSTconfig/arrows-LT-1.sst**: Logic/Temporal arrows
- **SSTconfig/arrows-NR-0.sst**: Narrative/Near arrows
- **docs/N4L.md**: Complete N4L language specification
- **examples/**: Sample N4L files for various domains

---

## Conclusion

N4L provides a flexible yet structured approach to representing knowledge from structured data sources. By following these guidelines and using stable arrows from SSTconfig, you can create meaningful, queryable knowledge graphs that preserve both structure and semantics.

Remember:

- **Simplicity** over complexity
- **Semantics** over structure
- **Stable arrows** over custom relations
- **Context** for grouping
- **Iteration** for improvement

The goal is not perfect modeling on the first try, but creating a foundation that can be refined and enhanced over time.
