package auth

import (
	"fmt"
	"regexp"
	"strings"
)

type PermissionGuardian struct {
	Freedoms     []Freedom
	Restrictions []Restriction
}

type Freedom struct {
	UriRegex  *regexp.Regexp
	UriMethod string
}

type Restriction struct {
	UriRegex  *regexp.Regexp
	MethodMap map[string]string
}

func MakePermissionGuardian() PermissionGuardian {
	freedoms := []Freedom{
		open("signup", "POST"),
		open("login", "POST"),
	}

	restrictions := []Restriction{
		restrictAll("skills?", "exp.skill"),
		restrictAll("skillConfigs?", "exp.skillConfig"),
		restrictAll("users?", "exp.user"),
	}

	return PermissionGuardian{Freedoms: freedoms, Restrictions: restrictions}
}

func restrictAll(endpoint string, permissionPath string) Restriction {
	return Restriction{
		UriRegex: endpointRegex(endpoint),
		MethodMap: map[string]string{
			"GET":    fmt.Sprintf("%s.read", permissionPath),
			"PUT":    fmt.Sprintf("%s.update", permissionPath),
			"POST":   fmt.Sprintf("%s.create", permissionPath),
			"DELETE": fmt.Sprintf("%s.delete", permissionPath),
		}}
}

func open(uri string, requestMethod string) Freedom {
	return Freedom{
		UriRegex:  endpointRegex(uri),
		UriMethod: strings.ToUpper(requestMethod),
	}
}

func endpointRegex(endpoint string) *regexp.Regexp {
	// ?  - matches between 0 or 1 of the proceeding token (with or without endpoint slash)
	// .* - matches everything after
	// TODO(Handle multiple api versions)
	return regexp.MustCompile(fmt.Sprintf("/api/%s/?.*", endpoint))
}

func (p *PermissionGuardian) GetAllPermissions() []string {
	var permissions []string

	for _, restriction := range p.Restrictions {
		for _, perm := range restriction.MethodMap {
			permissions = append(permissions, perm)
		}
	}

	return permissions
}

func (p *PermissionGuardian) HasOpenPermission(requestUri string, requestMethod string) bool {
	for _, freedom := range p.Freedoms {
		if freedom.UriRegex.MatchString(requestUri) {
			return true
		}
	}
	return false
}

func (p *PermissionGuardian) GetRequiredPermission(requestUri string, requestMethod string) string {
	var methodMap map[string]string

	for _, restriction := range p.Restrictions {
		if restriction.UriRegex.MatchString(requestUri) {
			methodMap = restriction.MethodMap
			break
		}
	}

	if methodMap == nil {
		return ""
	}

	return methodMap[strings.ToUpper(requestMethod)]
}
