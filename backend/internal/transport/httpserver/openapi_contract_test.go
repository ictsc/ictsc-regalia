package httpserver

import (
	"strings"
	"testing"

	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

func TestOpenAPIContractIsValidAndDeclaresTransportErrors(t *testing.T) {
	spec, err := api.GetSwagger()
	if err != nil {
		t.Fatalf("GetSwagger() error = %v", err)
	}
	if spec.OpenAPI != "3.1.0" {
		t.Fatalf("OpenAPI version = %q, want 3.1.0", spec.OpenAPI)
	}
	if err := spec.Validate(t.Context()); err != nil {
		t.Fatalf("OpenAPI validation failed: %v", err)
	}

	operationIDs := make(map[string]string)
	operationCount := 0
	for path, item := range spec.Paths.Map() {
		for method, operation := range item.Operations() {
			operationCount++
			location := strings.ToUpper(method) + " " + path
			if operation.OperationID == "" {
				t.Errorf("%s has no operationId", location)
			} else if previous, exists := operationIDs[operation.OperationID]; exists {
				t.Errorf("operationId %q is duplicated by %s and %s", operation.OperationID, previous, location)
			} else {
				operationIDs[operation.OperationID] = location
			}

			responses := operation.Responses.Map()
			if _, ok := responses["500"]; !ok {
				t.Errorf("%s does not declare RFC 9457 status 500", location)
			}
			hasInput := operation.RequestBody != nil || len(operation.Parameters) > 0 || len(item.Parameters) > 0
			if hasInput {
				if _, ok := responses["422"]; !ok {
					t.Errorf("%s accepts input but does not declare status 422", location)
				}
			}
			switch strings.ToUpper(method) {
			case "POST", "PUT", "PATCH", "DELETE":
				if _, ok := responses["403"]; !ok {
					t.Errorf("%s is a mutation but does not declare status 403", location)
				}
			}
		}
	}
	if operationCount != 61 {
		t.Errorf("operation count = %d, want 61", operationCount)
	}
}

func TestOAuthNextPatternRejectsCrossOriginRepresentations(t *testing.T) {
	spec, err := api.GetSwagger()
	if err != nil {
		t.Fatal(err)
	}
	parameter := spec.Paths.Value("/api/v1/auth/discord").Get.Parameters.GetByInAndName("query", "next")
	if parameter == nil || parameter.Schema == nil || parameter.Schema.Value == nil {
		t.Fatal("contestant OAuth next parameter schema is missing")
	}
	for _, invalid := range []string{"//evil.example", "/\\\\evil.example", "/path?query", "/path#fragment", "/path\x00control"} {
		if err := parameter.Schema.Value.VisitJSON(invalid); err == nil {
			t.Errorf("next pattern accepted %q", invalid)
		}
	}
}
