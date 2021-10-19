/*
Copyright 2021 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/koderover/zadig/pkg/microservice/policy/core/service"
	internalhandler "github.com/koderover/zadig/pkg/shared/handler"
	e "github.com/koderover/zadig/pkg/tool/errors"
)

func CreateRole(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	args := &service.Role{}
	if err := c.ShouldBindJSON(args); err != nil {
		ctx.Err = err
		return
	}

	projectName := c.Query("projectName")
	if projectName == "" {
		ctx.Err = e.ErrInvalidParam.AddDesc("projectName is empty")
		return
	}

	ctx.Err = service.CreateRole(projectName, args, ctx.Logger)
}

func ListRole(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	projectName := c.Query("projectName")
	if projectName == "" {
		ctx.Err = e.ErrInvalidParam.AddDesc("args projectName can't be empty")
		return
	}
	ctx.Resp, ctx.Err = service.ListRole(projectName)
	return
}

func CreateGlobalRole(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	args := &service.Role{}
	if err := c.ShouldBindJSON(args); err != nil {
		ctx.Err = err
		return
	}

	ctx.Err = service.CreateRole("", args, ctx.Logger)
}

func ListGlobalRole(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()

	ctx.Resp, ctx.Err = service.ListGlobalRole(ctx.Logger)
	return
}

func DeleteRole(c *gin.Context) {
	ctx := internalhandler.NewContext(c)
	defer func() { internalhandler.JSONResponse(c, ctx) }()
	name := c.Param("name")
	if name == "" {
		ctx.Err = e.ErrInvalidParam.AddDesc("args name can't be empty")
		return
	}
	projectName := c.Query("projectName")
	if name == "" {
		ctx.Err = e.ErrInvalidParam.AddDesc("args projectName can't be empty")
		return
	}
	ctx.Err = service.DeleteRole(name, projectName, ctx.Logger)
	return
}
