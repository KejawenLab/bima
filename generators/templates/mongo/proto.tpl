syntax = "proto3";

package grpcs;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "bima/pagination.proto";

option go_package = ".;grpcs";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
	info: {
		title: "{{.Module}} Service";
		version: "{{.ApiVersion}}";
	};
};

message {{.Module}} {
    string id = 1;
{{range .Columns}}
    {{.ProtobufType}} {{.NameUnderScore}} = {{.Index}};
{{end}}
}

message {{.Module}}PaginatedResponse {
    repeated {{.Module}} data = 1;
    PaginationMetadata meta = 2;
}

service {{.Module}}s {
    rpc GetPaginated (Pagination) returns ({{.Module}}PaginatedResponse) {
        option (google.api.http) = {
            get: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}"
        };
    }

    rpc Create ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            post: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}"
            body: "*"
        };
    }

    rpc Update ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            put: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
            body: "*"

            additional_bindings {
                patch: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
                body: "*"
            }
        };
    }

    rpc Get ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            get: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
        };
    }

    rpc Delete ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            delete: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
        };
    }
}
