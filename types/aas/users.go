/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package aas

type UserCreate struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type UserCreateResponse struct {
	ID   string `json:"user_id"`
	Name string `json:"username"`
}

type UserRoleCreate struct {
	ID      string `json:"user_id"`
	RoleIds RoleIDs
}

type UserCred struct {
        UserName string `json:"username"`
        Password string `json:"password"`
}
