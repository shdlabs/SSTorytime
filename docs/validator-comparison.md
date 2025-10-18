# N4L Validator Comparison Report

Date: October 18, 2025
Branch: feature/n4l-validator

## Test Summary

### Files Tested

1. `astronomy.n4l` - ✅ PASSED both validators, uploaded successfully
2. `moon.n4l` - ✅ PASSED both validators, uploaded successfully
3. `music-catalogue.n4l` - ✅ PASSED both validators, uploaded successfully
4. `test_validation.n4l` - ⚠️ Different error detection

## Validator Comparison

### Our New Validator (`cmd/n4l-validate/`)

**Purpose:** Lexical and syntax validation without database upload

**Features:**

- ✅ Compiler-quality error messages with line numbers and context
- ✅ Shows exact source code location
- ✅ Provides helpful hints
- ✅ Detects unterminated strings
- ✅ Recognizes N4L special syntax (ditto `"`, aliases `@`, references `$`)
- ✅ Fast validation without database overhead

**Example Output:**

```
test_validation.n4l:7:25: error: unterminated string
     7 |  Missing closing quote: "this is broken
       |                         ^
  hint: add closing quote
```

**Current Limitations:**

- Only does lexical/syntactic checks (not semantic)
- Doesn't validate arrow definitions yet (parser not complete)
- Doesn't check node references
- Core arrows are hardcoded, not loaded from .sst files

### Mark's N4L Tool (`src/N4L`)

**Purpose:** Full parsing, validation, and database upload

**Features:**

- ✅ Complete semantic validation
- ✅ Loads arrow definitions from SSTconfig/\*.sst files
- ✅ Validates arrow usage against configuration
- ✅ Checks node references and relationships
- ✅ Handles database upload

**Example Output:**

```
11:N4L test_validation.n4l No such arrow has been declared in the configuration: (next) at line 11
```

**Issues (that motivated our validator):**

- ❌ No line/column context shown in error messages
- ❌ Doesn't show source code
- ❌ No helpful hints
- ❌ Must attempt upload to find errors (risky)
- ❌ Error format: `line:tool filename message at line X` (redundant, unclear)

## Key Findings

### 1. Core Arrows Discovery

We initially hardcoded these "core" arrows:

- `(contain)` / `(belong)` ✅ EXISTS in arrows-CN-2.sst
- `(hasX)` / `(isXof)` ✅ EXISTS in arrows-EP-3.sst
- `(next)` / `(prev)` ❌ DOESN'T EXIST (theoretical)

**Actual arrows from SSTconfig:**

- **CN-2 (Containment):** `(contain)`, `(belong)`, `(setof)`, `(in-set)`, `(has-pt)`, `(pt-of)`
- **EP-3 (Expression/Property):** `(propt)`, `(propt-of)`, `(hasX)`, `(isXof)`
- **LT-1 (Logic/Temporal):** `(fwd)`, `(bwd)`, `(=>)`, `(<=)`, `(cause)`, `(cause-by)`
- **NR-0 (Near/Relevance):** `(near)`, `(see-also)`, `(similar)`, `(contrast)`

### 2. Error Detection Differences

**Test case: Unterminated string**

```n4l
Missing closing quote: "this is broken
```

- **Our validator:** ✅ Caught at lexical stage, clear error with context
- **Mark's N4L:** ⚠️ Skipped (possibly treats quote as ditto marker)

**Test case: Undefined arrow**

```n4l
This node (next) that node
```

- **Our validator:** ⏭️ Skipped (parser not complete, no semantic analysis yet)
- **Mark's N4L:** ✅ Caught: "No such arrow has been declared"

### 3. Successful Uploads

All three real example files uploaded successfully:

- `astronomy.n4l` - "Astronomy by the numbers poem"
- `moon.n4l` - "Moon tidal effect"
- `music-catalogue.n4l` - Music collection data

Both validators passed these files, confirming they are syntactically correct.

## Recommendations

### Phase 1: Complete Parser (NEXT)

```
internal/n4l/parser.go
- Build AST from tokens
- Track indentation for hierarchy
- Parse arrow syntax: node (arrow) node
- Parse aliases and references
```

### Phase 2: Load Arrow Definitions

```
internal/n4l/arrows.go
- Read SSTconfig/*.sst files
- Parse arrow definitions
- Build arrow registry
- Merge with core arrows
```

### Phase 3: Semantic Validation

```
internal/n4l/validator.go
- Check arrow definitions exist
- Validate node references
- Check context usage
- Detect dangling arrows
```

### Phase 4: Integration

```
- Make N4L.go use our compiler for validation phase
- Keep database logic in N4L.go
- Share arrow definitions between both
- Unified error reporting
```

## Usage Examples

### Validate Before Upload (Recommended Workflow)

```bash
# Step 1: Validate syntax locally (fast, no DB risk)
cd /home/alex/SSTorytime
./cmd/n4l-validate/n4l-validate examples/myfile.n4l

# Step 2: If validation passes, upload to database
cd examples
../src/N4L -u myfile.n4l
```

### Current Arrows Available

Check `SSTconfig/*.sst` for complete list. Common arrows:

- Containment: `(contain)`, `(belong)`, `(setof)`, `(in-set)`
- Properties: `(hasX)`, `(isXof)`, `(propt)`, `(propt-of)`
- Causality: `(cause)`, `(cause-by)`, `(fwd)`, `(bwd)`, `(=>)`
- Association: `(near)`, `(see-also)`, `(similar)`, `(contrast)`

## Next Action Items

1. ✅ **DONE:** Build lexer with proper error reporting
2. ✅ **DONE:** Test on real N4L files
3. ✅ **DONE:** Compare against Mark's N4L
4. ⏳ **TODO:** Implement parser (build AST)
5. ⏳ **TODO:** Load arrow definitions from .sst files
6. ⏳ **TODO:** Implement semantic validator
7. ⏳ **TODO:** Fix godoc2n4l breadcrumb extraction
8. ⏳ **TODO:** Generate golang-stdlib.n4l documentation

## Conclusion

The new validator successfully:

- ✅ Provides clear, actionable error messages
- ✅ Validates syntax without database upload risk
- ✅ Handles N4L special features (ditto, aliases)
- ✅ Passes all real-world example files

Improvements needed:

- Complete parser implementation for arrow/node validation
- Load real arrow definitions from SSTconfig
- Add semantic validation phase
- Integration with Mark's N4L for full workflow
