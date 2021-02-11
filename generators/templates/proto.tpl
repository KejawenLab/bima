syntax = "proto3";

package grpcs;

import "google/api/annotations.proto";
import "pagination.proto";

option go_package = ".;grpcs";

message {{.Module}} {
    string id = 1;
{{range .Columns}}
    {{.ProtobufType}} {{.NameUnderScore}} = {{.Index}};
{{end}}
}

message {{.Module}}Response {
    int32 code = 1;
    {{.Module}} data = 2;
    string message = 3;
}

message {{.Module}}PaginatedResponse {
    int32 code = 1;
    repeated {{.Module}} data = 2;
    PaginationMetadata meta = 3;
}

service {{.Module}}s {
    rpc GetPaginated (Pagination) returns ({{.Module}}PaginatedResponse) {
        option (google.api.http) = {
            get: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}"
        };
    }

    rpc Create ({{.Module}}) returns ({{.Module}}Response) {
        option (google.api.http) = {
            post: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}"
            body: "*"
        };
    }

    rpc Update ({{.Module}}) returns ({{.Module}}Response) {
        option (google.api.http) = {
            put: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
            body: "*"

            additional_bindings {
                patch: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
                body: "*"
            }
        };
    }

    rpc Get ({{.Module}}) returns ({{.Module}}Response) {
        option (google.api.http) = {
            get: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
        };
    }

    rpc Delete ({{.Module}}) returns ({{.Module}}Response) {
        option (google.api.http) = {
            delete: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
        };
    }
}
