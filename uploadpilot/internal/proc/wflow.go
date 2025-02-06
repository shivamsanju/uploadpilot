package proc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	flow "github.com/Azure/go-workflow"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/db/types"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
	commontasks "github.com/uploadpilot/uploadpilot/internal/proc/tasks/common"
)

type ProcWorkflowRunner struct {
	procRepo   *db.ProcessorRepo
	uploadRepo *db.UploadRepo
	wf         *flow.Workflow
}

func NewProcWorkflowRunner() *ProcWorkflowRunner {
	return &ProcWorkflowRunner{
		procRepo:   db.NewProcessorRepo(),
		uploadRepo: db.NewUploadRepo(),
		wf:         new(flow.Workflow),
	}
}

func (r *ProcWorkflowRunner) Run(ctx context.Context) error {
	if err := r.wf.Do(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProcWorkflowRunner) Build(ctx context.Context, workspaceID, processorID, uploadID string) error {
	r.wf.DontPanic = true
	infra.Log.Infof("creating workflow for upload: %s, processor: %s, workspace: %s", uploadID, processorID, workspaceID)
	processor, err := r.procRepo.Get(ctx, processorID)
	if err != nil {
		return err
	}

	upload, err := r.uploadRepo.Get(ctx, uploadID)
	if err != nil {
		infra.Log.Errorf("failed to get upload: %s", err.Error())
		return err
	}

	taskfnMap, procTaskMap, edges, err := r.getNodesAndEdgesMap(processor)
	if err != nil {
		return err
	}

	tmpDir := os.TempDir() + "/" + processorID + "/" + uploadID
	os.MkdirAll(tmpDir, os.ModePerm)
	firstTask := tasks.NewBaseTask(workspaceID, processorID, uploadID, tmpDir, upload.ID)

	steps := make([]flow.Builder, 0, len(taskfnMap))
	tsks := make([]flow.Steper, 0, len(taskfnMap))

	for id, task := range taskfnMap {
		depKeys := edges[id]
		var prevTasks []tasks.Task
		var PrevtasksKeys []string
		for _, depKey := range depKeys {
			prevTasks = append(prevTasks, taskfnMap[depKey])
			PrevtasksKeys = append(PrevtasksKeys, depKey)
		}

		pt, ok := procTaskMap[id]
		if !ok {
			return fmt.Errorf("failed to get task %s", id)
		}
		infra.Log.Infof("task: %s, prevTasks: %s", pt["key"], strings.Join(PrevtasksKeys, ","))

		build, err := r.buildStepWithDependencies(pt, task, prevTasks, firstTask)
		if err != nil {
			return err
		}
		steps = append(steps, build)
		tsks = append(tsks, task)
	}

	steps = r.addCommonSteps(steps, tsks, firstTask)
	r.wf.Add(steps...)

	return nil
}

func (r *ProcWorkflowRunner) buildStepWithDependencies(pt types.JSONB, task tasks.Task, prevTasks []tasks.Task, firstTask *tasks.BaseTask) (*flow.AddStep[tasks.Task], error) {
	taskData := pt["data"].(map[string]interface{})

	infra.Log.Infof("taskData: %+v", taskData)
	params := r.GetTaskParams(taskData)

	step := flow.Step(task).
		AfterStep(func(ctx context.Context, _ flow.Steper, err error) error {
			// throw errror only if continue on error is false
			prevTask := task.GetTask()
			if err != nil {
				infra.Log.Errorf(msg.ProcTaskFailed, pt["key"], prevTask.WorkspaceID, prevTask.ProcessorID, prevTask.UploadID, err.Error())
				if params.ContinueOnError {
					return nil
				}
				return err
			}
			return nil
		})

	if len(prevTasks) == 0 {
		step = step.Input(func(ctx context.Context, t tasks.Task) error {
			firstTask.TaskID = pt["id"].(string)
			firstTask.TaskData = taskData
			t.MakeTask(firstTask)
			return nil
		})
	} else {
		// There will be a single parent task only for now
		for _, prevTask := range prevTasks {
			step = step.DependsOn(prevTask).
				Input(func(ctx context.Context, t tasks.Task) error {
					inp := prevTask.GetTask()
					inp.TaskID = pt["id"].(string)
					inp.TaskData = taskData
					inp.Input = inp.Output
					t.MakeTask(inp)
					return nil
				})
		}
	}

	if params.Retry > 0 {
		step = step.Retry(func(ro *flow.RetryOption) {
			ro.Attempts = params.Retry
		})
	}

	if params.TimeoutMilSec > 0 {
		step = step.Timeout(time.Duration(params.TimeoutMilSec) * time.Second)
	}

	return &step, nil
}

func (r *ProcWorkflowRunner) addCommonSteps(steps []flow.Builder, tsks []flow.Steper, firstTask *tasks.BaseTask) []flow.Builder {
	success := flow.Step(commontasks.NewSuccessTask()).
		DependsOn(tsks...).
		When(flow.AllSucceededOrSkipped).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.MakeTask(firstTask)
			return nil
		})

	failure := flow.Step(commontasks.NewFailedTask()).
		DependsOn(tsks...).
		When(flow.AnyFailed).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.MakeTask(firstTask)
			return nil
		})

	cancelled := flow.Step(commontasks.NewCancelledTask()).
		DependsOn(tsks...).
		When(flow.BeCanceled).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.MakeTask(firstTask)
			return nil
		})

	cleanup := flow.Step(commontasks.NewCleanupTask()).
		DependsOn(tsks...).
		When(flow.Always).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.MakeTask(firstTask)
			return nil
		})

	return append(steps, success, failure, cancelled, cleanup)
}

func (r *ProcWorkflowRunner) getNodesAndEdgesMap(processor *models.Processor) (map[string]tasks.Task, map[string]types.JSONB, map[string][]string, error) {
	taskfnMap := make(map[string]tasks.Task)
	taskMap := make(map[string]types.JSONB)
	prevTasks := make(map[string][]string)

	for _, node := range processor.Canvas["nodes"].([]interface{}) {
		node := node.(map[string]interface{})
		task, ok := GetTaskFromRegistry(models.TaskKey(node["key"].(string)))
		if !ok {
			return nil, nil, nil, fmt.Errorf("unknown task key: %s", node["key"].(string))
		}
		taskData, ok := processor.Data[node["id"].(string)]
		if ok {
			td := taskData.(map[string]interface{})
			nd := node["data"].(map[string]interface{})
			for key, value := range td {
				nd[key] = value
			}
			node["data"] = nd
		}

		taskfnMap[node["id"].(string)] = task
		taskMap[node["id"].(string)] = node
	}

	for _, t := range processor.Canvas["nodes"].([]interface{}) {
		for _, edge := range processor.Canvas["edges"].([]interface{}) {
			edge := edge.(map[string]interface{})
			t := t.(map[string]interface{})
			if edge["source"] == t["id"].(string) {
				prevTasks[edge["target"].(string)] = append(prevTasks[edge["target"].(string)], t["id"].(string))
			}
		}
	}

	return taskfnMap, taskMap, prevTasks, nil
}

func (r *ProcWorkflowRunner) GetTaskParams(taskData map[string]interface{}) *TaskParams {
	params := &TaskParams{}

	ce, ok := taskData["continueOnError"]
	if !ok {
		params.ContinueOnError = false
	} else {
		params.ContinueOnError = ce.(bool)
	}

	tms, ok := taskData["timeoutMilSec"]
	if !ok {
		params.TimeoutMilSec = 0
	} else {
		params.TimeoutMilSec = uint64(tms.(uint64))
	}

	retry, ok := taskData["retry"]
	if !ok {
		params.Retry = 0
	} else {
		params.Retry = uint64(retry.(uint64))
	}

	return params
}

type TaskParams struct {
	ContinueOnError bool
	TimeoutMilSec   uint64
	Retry           uint64
}
