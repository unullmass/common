/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package auth

import (
	types "intel/isecl/lib/common/types/aas"
	"strings"
)

func ValidatePermissionAndGetRoleContext(privileges []types.RoleInfo, reqRoles []types.RoleInfo,
	retNilCtxForEmptyCtx bool) (*map[string]*types.RoleInfo, bool) {

	ctx := make(map[string]*types.RoleInfo)
	foundMatchingRole := false
	for _, role := range privileges {
		for _, reqRole := range reqRoles {
			if role.Service == reqRole.Service && role.Name == reqRole.Name {
				if strings.TrimSpace(role.Context) == "" && retNilCtxForEmptyCtx == true {
					return nil, true
				}
				if strings.TrimSpace(role.Context) != "" {
					ctx[strings.TrimSpace(role.Context)] = &role
				}
				foundMatchingRole = true
			}
		}

	}
	return &ctx, foundMatchingRole
}
