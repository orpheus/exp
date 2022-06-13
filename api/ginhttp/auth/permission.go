package auth

import (
	"fmt"
	"github.com/orpheus/exp/system/sysauth"
	"regexp"
	"strings"
)

// MakePermissionGuardian defines the necessary route-to-permission mappings for the whole system.
func MakePermissionGuardian() *Guardian {
	freedoms := []Freedom{
		open("signup", "POST"),
		open("login", "POST"),
		open("permission", "GET"),
		open("health", "GET"),
	}

	// These are ordered, put more specific endpoints higher than less specific ones
	restrictions := []Restriction{
		restrict("skill/addTxp", map[string]string{"POST": "exp.skill.addTxp"}),
		restrictAllSpecial("skills?", "exp.skill", map[string]string{"DELETE": "admin.skill.delete"}),
		restrictAll("skillConfigs?", "exp.skillConfig"),
		restrictAll("users?", "exp.user"),
		restrictAll("role?", "exp.role"),
	}

	return &Guardian{Freedoms: freedoms, Restrictions: restrictions}
}

type PermissionGuardian interface {
	sysauth.PermissionGetter
	HasOpenPermission(requestUri string, requestMethod string) bool
	GetRequiredPermission(requestUri string, requestMethod string) string
}

type Guardian struct {
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

func restrict(endpoint string, methodMap map[string]string) Restriction {
	return Restriction{
		UriRegex:  endpointRegex(endpoint),
		MethodMap: methodMap,
	}
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

func restrictAllSpecial(endpoint string, permissionPath string, specialRules map[string]string) Restriction {
	restriction := restrictAll(endpoint, permissionPath)
	for req, permission := range specialRules {
		restriction.MethodMap[req] = permission
	}
	return restriction
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

func (p *Guardian) GetPermissions() []string {
	var permissions []string

	for _, restriction := range p.Restrictions {
		for _, perm := range restriction.MethodMap {
			permissions = append(permissions, perm)
		}
	}

	return permissions
}

func (p *Guardian) HasOpenPermission(requestUri string, requestMethod string) bool {
	for _, freedom := range p.Freedoms {
		if freedom.UriRegex.MatchString(requestUri) {
			return true
		}
	}
	return false
}

func (p *Guardian) GetRequiredPermission(requestUri string, requestMethod string) string {
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

func HasPermission(requiredPermission string, userPermission string) bool {
	allStar := "*"

	if userPermission == allStar {
		return true
	}
	if requiredPermission == userPermission {
		return true
	}

	reqSplit := strings.Split(requiredPermission, ".")
	userSplit := strings.Split(userPermission, ".")

	reqLenIdx := len(reqSplit) - 1

	for idx, part := range userSplit {
		// i.e. exp.skill.read vs exp.skill.level.update
		if idx > reqLenIdx {
			return false
		}
		reqPart := reqSplit[idx]
		if reqPart != part {
			// i.e. exp.skill.read vs exp.skill.*
			if part == allStar {
				return true
			}
			// i.e. exp.skill.read vs exp.user.read
			return false
		}
		// i.e. exp.skill.read vs exp.skill.read
		if reqPart == part {
			continue
		}

	}

	panic("HasPermission fn failed")
}
