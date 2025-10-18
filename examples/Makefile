
# Validator command (installed via go install)
VALIDATOR := n4l-validate
N4L := ../src/N4L

# List of all N4L files to process
N4L_FILES := doors.n4l Mary.n4l branches.n4l doubleslit.n4l \
             wardleymap.n4l kubernetes.n4l \
             SSTsearchexamples.n4l integral.n4l reasoning.n4l \
             moon.n4l astronomy.n4l music-catalogue.n4l reminders.n4l

# Helper function to validate then upload
# Usage: $(call validate_and_upload,filename,upload_flags)
define validate_and_upload
	@echo "🔍 Validating $(1)..."
	@if $(VALIDATOR) $(1) > /dev/null 2>&1; then \
		echo "✅ Validation passed"; \
		echo "📤 Uploading $(1)..."; \
		if $(N4L) -u $(2) $(1) 2>&1 | tail -5; then \
			echo "✓ $(1) uploaded successfully"; \
		else \
			echo "❌ $(1) upload FAILED"; \
			exit 1; \
		fi; \
	else \
		echo "❌ $(1) validation FAILED - aborting"; \
		echo "   Run: $(VALIDATOR) $(1)"; \
		echo "   to see detailed errors"; \
		exit 1; \
	fi
	@echo ""
endef

all: check-validator ../src/N4L validate-all
	@echo "========================================="
	@echo "  All files validated - starting upload"
	@echo "  Database will be wiped first!"
	@echo "========================================="
	@echo ""
	@echo "🗑️  Wiping database and uploading first file..."
	$(call validate_and_upload,doors.n4l,-wipe)
	@echo "📤 Uploading remaining files..."
	$(call validate_and_upload,Mary.n4l)
	$(call validate_and_upload,branches.n4l)
	$(call validate_and_upload,doubleslit.n4l)
	$(call validate_and_upload,wardleymap.n4l)
	$(call validate_and_upload,kubernetes.n4l)
# 	$(call validate_and_upload,SSTorytime.n4l)
	$(call validate_and_upload,SSTsearchexamples.n4l)
	$(call validate_and_upload,integral.n4l)
	$(call validate_and_upload,reasoning.n4l)
	$(call validate_and_upload,moon.n4l)
	$(call validate_and_upload,astronomy.n4l)
	$(call validate_and_upload,music-catalogue.n4l)
	$(call validate_and_upload,reminders.n4l)
	@echo "========================================="
	@echo "  ✅ All files uploaded successfully!"
	@echo "========================================="

# Validate all files before any upload (fail fast)
validate-all:
	@echo "========================================="
	@echo "  Pre-flight: Validating All Files"
	@echo "========================================="
	@echo ""
	@FAILED=0; \
	for file in $(N4L_FILES); do \
		echo "🔍 Validating $$file..."; \
		if $(VALIDATOR) $$file > /dev/null 2>&1; then \
			echo "✅ $$file passed"; \
		else \
			echo "❌ $$file FAILED"; \
			$(VALIDATOR) $$file; \
			FAILED=1; \
		fi; \
		echo ""; \
	done; \
	if [ $$FAILED -eq 1 ]; then \
		echo "=========================================" >&2; \
		echo "  ❌ Validation failed - fix errors first" >&2; \
		echo "=========================================" >&2; \
		exit 1; \
	fi
	@echo "========================================="
	@echo "  ✅ All files passed validation"
	@echo "========================================="
	@echo ""

# Check if validator is installed
check-validator:
	@command -v $(VALIDATOR) >/dev/null 2>&1 || { \
		echo "❌ n4l-validate not found!"; \
		echo ""; \
		echo "Install with:"; \
		echo "  cd /home/alex/SSTorytime"; \
		echo "  go install ./cmd/n4l-validate"; \
		echo ""; \
		exit 1; \
	}
	@echo "✓ Validator found: $$(which $(VALIDATOR))"
	@echo ""

# Validate all files without uploading
validate:
	@echo "========================================="
	@echo "  Validating All N4L Files"
	@echo "========================================="
	@echo ""
	@for file in $(N4L_FILES); do \
		echo "🔍 Validating $$file..."; \
		if $(VALIDATOR) $$file > /dev/null 2>&1; then \
			echo "✅ $$file passed"; \
		else \
			echo "❌ $$file FAILED"; \
			$(VALIDATOR) $$file; \
		fi; \
		echo ""; \
	done
	@echo "========================================="
	@echo "  Validation Complete"
	@echo "========================================="

# Upload a single file (validate first, use -force to overwrite)
# Usage: make upload FILE=myfile.n4l [FORCE=1]
upload:
	@if [ -z "$(FILE)" ]; then \
		echo "❌ Usage: make upload FILE=myfile.n4l [FORCE=1]"; \
		exit 1; \
	fi; \
	if [ "$(FORCE)" = "1" ]; then \
		$(call validate_and_upload,$(FILE),-force); \
	else \
		$(call validate_and_upload,$(FILE)); \
	fi

# Watch a file for changes (validation only)
# Usage: make watch FILE=myfile.n4l
watch:
	@if [ -z "$(FILE)" ]; then \
		echo "❌ Usage: make watch FILE=myfile.n4l"; \
		exit 1; \
	fi
	@echo "👁️  Starting watch mode for $(FILE)"
	@echo "Press Ctrl+C to stop"
	@echo ""
	@$(VALIDATOR) -w $(FILE)

# Show help
help:
	@echo "N4L Examples Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  make                 - Validate all, wipe DB, upload all (default)"
	@echo "  make validate        - Validate all files without uploading"
	@echo "  make upload FILE=<file> [FORCE=1] - Validate and upload single file"
	@echo "  make watch FILE=<file>            - Watch file for changes"
	@echo "  make clean           - Remove temporary files"
	@echo "  make help            - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make                               # Full rebuild (wipe + upload all)"
	@echo "  make validate                      # Just check syntax"
	@echo "  make upload FILE=astronomy.n4l     # Upload one file"
	@echo "  make upload FILE=moon.n4l FORCE=1  # Force overwrite existing"
	@echo "  make watch FILE=moon.n4l           # Live validation"
	@echo ""
	@echo "The validator (n4l-validate) must be installed first:"
	@echo "  cd .. && go install ./cmd/n4l-validate"

clean:
	rm -f *_test_log *~

.PHONY: all validate-all check-validator validate upload watch help clean
