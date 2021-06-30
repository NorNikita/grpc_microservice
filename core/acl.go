package core

import (
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
)

func ParseAcl(rawAcl string) (map[string][]string, error) {
	var result map[string][]string

	err := json.Unmarshal([]byte(rawAcl), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *InnerStruct) checkPermissionForMethod(userName, methodName string) error {
	paths, ok := i.Acl[userName]
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	for _, name := range paths {
		components := strings.Split(name, "/")
		if methodName == name || (strings.Contains(methodName, components[1]) && components[2] == "*") {
			return nil
		}
	}
	return grpc.Errorf(codes.Unauthenticated, "unauthenticated")
}
