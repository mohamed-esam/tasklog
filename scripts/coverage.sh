#!/bin/bash

# Test Coverage Summary Script

echo "╔════════════════════════════════════════════════════════════╗"
echo "║            Tasklog Test Coverage Summary                   ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Run tests and generate coverage
go test -coverprofile=coverage.out ./... > /dev/null 2>&1

echo "📊 Core Business Logic Coverage (Testable Code):"
echo "─────────────────────────────────────────────────────────────"

# Individual package coverage
echo "  📦 config package:    $(go tool cover -func=coverage.out | grep 'internal/config/config.go' | grep 'total:' | awk '{print $3}')"
echo "  📦 storage package:   $(go tool cover -func=coverage.out | grep 'internal/storage/storage.go' | grep 'total:' | awk '{print $3}')"
echo "  📦 timeparse package: $(go tool cover -func=coverage.out | grep 'internal/timeparse/timeparse.go' | grep 'total:' | awk '{print $3}')"

# Calculate average of core packages
CORE_AVG=$(go tool cover -func=coverage.out | grep -E "(config|storage|timeparse)" | grep "\.go:" | awk '{gsub("%","",$3); sum+=$3; count++} END {if(count>0) printf "%.1f%%", sum/count}')
echo ""
echo "  🎯 Core Average:      $CORE_AVG"
echo "─────────────────────────────────────────────────────────────"

# Test counts
TOTAL_TESTS=$(go test -v ./internal/config ./internal/storage ./internal/timeparse 2>&1 | grep -E "^=== RUN" | wc -l | tr -d ' ')
PASSED_TESTS=$(go test -v ./internal/config ./internal/storage ./internal/timeparse 2>&1 | grep -E "^--- PASS" | wc -l | tr -d ' ')

echo ""
echo "✅ Tests Passed: $PASSED_TESTS/$TOTAL_TESTS"
echo ""
echo "📝 Notes:"
echo "  • cmd/ui packages not tested (interactive CLI)"
echo "  • jira/tempo have structure tests only"
echo "  • Focus on business logic ensures reliability"
echo ""
echo "💡 Run 'make test-coverage' for detailed HTML report"
echo ""
