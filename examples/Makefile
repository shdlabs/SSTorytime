
# Validator command (installed via go install)
VALIDATOR := n4l-validate
N4L := ../src/N4L

# List of all N4L files to process
N4L_FILES := doors.n4l Mary.n4l branches.n4l doubleslit.n4l \
             wardleymap.n4l kubernetes.n4l SSTorytime.n4l \
             SSTsearchexamples.n4l integral.n4l reasoning.n4l \
             moon.n4l astronomy.n4l music-catalogue.n4l reminders.n4l

# Helper function to try upload and record failure
# Usage: $(call try_upload,filename,upload_flags)
# Returns: Appends to FAILED_UPLOADS if upload fails
define try_upload
	@echo "📤 Uploading $(1)..."
	@if $(N4L) -u $(2) $(1) > /tmp/n4l_upload_$(1).log 2>&1; then \
		echo "✓ $(1) uploaded successfully"; \
		tail -2 /tmp/n4l_upload_$(1).log | grep -v "^$$"; \
		rm -f /tmp/n4l_upload_$(1).log; \
		echo "$(1):SUCCESS" >> /tmp/n4l_results.txt; \
	else \
		echo "❌ $(1) upload FAILED - continuing with others..."; \
		echo "   Error log saved to /tmp/n4l_upload_$(1).log"; \
		echo "$(1):FAILED" >> /tmp/n4l_results.txt; \
	fi
	@echo ""
endef

all: check-validator ../src/N4L
	@echo "========================================="
	@echo "  Batch Upload (continues on failures)"
	@echo "  Database will be wiped first!"
	@echo "========================================="
	@echo ""
	@rm -f /tmp/n4l_results.txt /tmp/n4l_upload_*.log
	@echo "🗑️  Wiping database and uploading first file..."
	$(call try_upload,doors.n4l,-wipe)
	@echo "📤 Uploading remaining files..."
	$(call try_upload,Mary.n4l)
	$(call try_upload,branches.n4l)
	$(call try_upload,doubleslit.n4l)
	$(call try_upload,wardleymap.n4l)
	$(call try_upload,kubernetes.n4l)
	$(call try_upload,SSTorytime.n4l)
	$(call try_upload,SSTsearchexamples.n4l)
	$(call try_upload,integral.n4l)
	$(call try_upload,reasoning.n4l)
	$(call try_upload,moon.n4l)
	$(call try_upload,astronomy.n4l)
	$(call try_upload,music-catalogue.n4l)
	$(call try_upload,reminders.n4l)
	@echo "========================================="
	@echo "  Upload Summary"
	@echo "========================================="
	@echo ""
	@if [ -f /tmp/n4l_results.txt ]; then \
		SUCCESS=`grep -c ":SUCCESS" /tmp/n4l_results.txt 2>/dev/null | tr -d '\n' || echo 0`; \
		FAILED=`grep -c ":FAILED" /tmp/n4l_results.txt 2>/dev/null | tr -d '\n' || echo 0`; \
		echo "✅ Successfully uploaded: $$SUCCESS files"; \
		echo ""; \
		if [ "$$FAILED" != "0" ] && [ -n "$$FAILED" ]; then \
			echo "❌ Failed uploads ($$FAILED files):"; \
			echo ""; \
			grep ":FAILED" /tmp/n4l_results.txt | while read line; do \
				FILE=$${line%%:*}; \
				echo "  • $$FILE"; \
				if [ -f /tmp/n4l_upload_$$FILE.log ]; then \
					echo "    Error:"; \
					head -5 /tmp/n4l_upload_$$FILE.log | sed 's/^/      /'; \
					echo ""; \
					echo "    📋 To diagnose: $(VALIDATOR) $$FILE"; \
					echo "    📤 To retry: $(N4L) -u -force $$FILE"; \
					echo "    📄 Full log: /tmp/n4l_upload_$$FILE.log"; \
					echo ""; \
				fi; \
			done; \
			echo "=========================================" >&2; \
			echo "⚠️  Some files failed - see details above" >&2; \
			echo "=========================================" >&2; \
		else \
			echo "🎉 All files uploaded successfully!"; \
		fi; \
	fi
	@echo ""

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

# Upload a single file (with optional force flag)
# Usage: make upload FILE=myfile.n4l [FORCE=1]
upload:
	@if [ -z "$(FILE)" ]; then \
		echo "❌ Usage: make upload FILE=myfile.n4l [FORCE=1]"; \
		exit 1; \
	fi
ifeq ($(FORCE),1)
	$(call try_upload,$(FILE),-force)
else
	$(call try_upload,$(FILE))
endif

# Validate files that failed upload
# Usage: make validate-failed
validate-failed:
	@if [ ! -f /tmp/n4l_results.txt ]; then \
		echo "❌ No upload results found. Run 'make' first."; \
		exit 1; \
	fi
	@FAILED=$$(grep -c ":FAILED" /tmp/n4l_results.txt 2>/dev/null || echo 0); \
	if [ $$FAILED -eq 0 ]; then \
		echo "✅ No failed uploads to validate"; \
	else \
		echo "========================================"; \
		echo "  Validating Failed Uploads ($$FAILED files)"; \
		echo "========================================"; \
		echo ""; \
		grep ":FAILED" /tmp/n4l_results.txt | while read line; do \
			FILE=$${line%%:*}; \
			echo "🔍 Validating $$FILE..."; \
			if $(VALIDATOR) $$FILE; then \
				echo "✅ $$FILE passed validation (N4L semantic error)"; \
			else \
				echo "❌ $$FILE has syntax errors"; \
			fi; \
			echo ""; \
		done; \
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
	@echo "Resilient Upload Strategy:"
	@echo "  - Uploads all files, continues on failures"
	@echo "  - Shows summary at the end with failed files"
	@echo "  - Saves error logs for manual inspection"
	@echo ""
	@echo "Targets:"
	@echo "  make                       - Wipe DB, upload all (continues on failures)"
	@echo "  make validate              - Validate all files without uploading"
	@echo "  make validate-failed       - Validate only files that failed upload"
	@echo "  make upload FILE=<file> [FORCE=1] - Upload single file"
	@echo "  make watch FILE=<file>     - Watch file for changes (validation)"
	@echo "  make clean                 - Remove temporary files and logs"
	@echo "  make help                  - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make                                # Upload all (continues on failures)"
	@echo "  make validate-failed                # Diagnose failed uploads"
	@echo "  make upload FILE=astronomy.n4l      # Upload one file"
	@echo "  make upload FILE=moon.n4l FORCE=1   # Force overwrite"
	@echo "  make watch FILE=moon.n4l            # Live validation"
	@echo ""
	@echo "After 'make' completes:"
	@echo "  - Check summary for failed files"
	@echo "  - Use 'make validate-failed' to diagnose"
	@echo "  - Check /tmp/n4l_upload_<file>.log for details"
	@echo "  - Fix and retry: ../src/N4L -u -force <file>"
	@echo ""
	@echo "The validator (n4l-validate) must be installed first:"
	@echo "  cd .. && go install ./cmd/n4l-validate"

clean:
	rm -f *_test_log *~
	rm -f /tmp/n4l_results.txt /tmp/n4l_upload_*.log

.PHONY: all validate-all check-validator validate validate-failed upload watch help clean
