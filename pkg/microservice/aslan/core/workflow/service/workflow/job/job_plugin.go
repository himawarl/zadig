/*
Copyright 2022 The KodeRover Authors.

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

package job

import (
	"strings"

	"github.com/koderover/zadig/pkg/microservice/aslan/config"
	commonmodels "github.com/koderover/zadig/pkg/microservice/aslan/core/common/repository/models"
	commonservice "github.com/koderover/zadig/pkg/microservice/aslan/core/common/service"
	"github.com/koderover/zadig/pkg/setting"
	"github.com/koderover/zadig/pkg/tool/log"
)

type PluginJob struct {
	job      *commonmodels.Job
	workflow *commonmodels.WorkflowV4
	spec     *commonmodels.PluginJobSpec
}

func (j *PluginJob) Instantiate() error {
	j.spec = &commonmodels.PluginJobSpec{}
	if err := commonmodels.IToiYaml(j.job.Spec, j.spec); err != nil {
		return err
	}
	j.job.Spec = j.spec
	return nil
}

func (j *PluginJob) SetPreset() error {
	j.spec = &commonmodels.PluginJobSpec{}
	if err := commonmodels.IToi(j.job.Spec, j.spec); err != nil {
		return err
	}
	j.job.Spec = j.spec
	return nil
}

func (j *PluginJob) MergeArgs(args *commonmodels.Job) error {
	if j.job.Name == args.Name && j.job.JobType == args.JobType {
		j.spec = &commonmodels.PluginJobSpec{}
		if err := commonmodels.IToi(j.job.Spec, j.spec); err != nil {
			return err
		}
		j.job.Spec = j.spec
		argsSpec := &commonmodels.PluginJobSpec{}
		if err := commonmodels.IToi(args.Spec, argsSpec); err != nil {
			return err
		}
		j.spec.Plugin.Inputs = argsSpec.Plugin.Inputs
		j.job.Spec = j.spec
	}
	return nil
}

func (j *PluginJob) ToJobs(taskID int64) ([]*commonmodels.JobTask, error) {
	logger := log.SugaredLogger()
	resp := []*commonmodels.JobTask{}
	j.spec = &commonmodels.PluginJobSpec{}
	if err := commonmodels.IToi(j.job.Spec, j.spec); err != nil {
		return resp, err
	}
	j.job.Spec = j.spec
	jobTask := &commonmodels.JobTask{
		Name:       j.job.Name,
		JobType:    string(config.JobPlugin),
		Properties: *j.spec.Properties,
		Plugin:     j.spec.Plugin,
		Outputs:    j.spec.Plugin.Outputs,
	}
	registries, err := commonservice.ListRegistryNamespaces("", true, logger)
	if err != nil {
		return resp, err
	}
	jobTask.Properties.Registries = registries

	renderedParams := []*commonmodels.Param{}
	for _, param := range j.spec.Plugin.Inputs {
		paramsKey := strings.Join([]string{"inputs", param.Name}, ".")
		renderedParams = append(renderedParams, &commonmodels.Param{Name: paramsKey, Value: param.Value, ParamsType: "string", IsCredential: false})
	}
	jobTask.Plugin = renderPlugin(jobTask.Plugin, renderedParams)

	jobTask.Outputs = j.spec.Plugin.Outputs
	return []*commonmodels.JobTask{jobTask}, nil
}

func renderPlugin(plugin *commonmodels.PluginTemplate, inputs []*commonmodels.Param) *commonmodels.PluginTemplate {
	for _, env := range plugin.Envs {
		env.Value = renderString(env.Value, setting.RenderPluginValueTemplate, inputs)
	}
	for i, arg := range plugin.Args {
		plugin.Args[i] = renderString(arg, setting.RenderPluginValueTemplate, inputs)
	}
	for i, cmd := range plugin.Cmds {
		plugin.Cmds[i] = renderString(cmd, setting.RenderPluginValueTemplate, inputs)
	}
	return plugin
}
