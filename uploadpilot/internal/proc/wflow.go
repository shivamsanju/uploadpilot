package proc

import (
	"context"
	"fmt"
	"time"

	flow "github.com/Azure/go-workflow"
	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/msg"
	"github.com/uploadpilot/uploadpilot/internal/proc/tasks"
	commontasks "github.com/uploadpilot/uploadpilot/internal/proc/tasks/common"
)

type ProcWorkflowRunner struct {
	procRepo *db.ProcessorRepo
	wf       *flow.Workflow
}

func NewProcWorkflowRunner() *ProcWorkflowRunner {
	return &ProcWorkflowRunner{
		procRepo: db.NewProcessorRepo(),
		wf:       new(flow.Workflow),
	}
}

func (r *ProcWorkflowRunner) Build(ctx context.Context, initialData *tasks.TaskData) error {
	infra.Log.Infof("hghghghg, %+v, %+v", r, initialData)
	processor, err := r.procRepo.Get(ctx, initialData.WorkspaceID, initialData.ProcessorID)
	if err != nil {
		return err
	}
	infra.Log.Errorf("some sss %+v", processor)

	nodes, edges, err := r.getNodesAndEdgesMap(processor)
	if err != nil {
		return err
	}

	steps := make([]flow.Builder, 0, len(nodes))
	tsks := make([]flow.Steper, 0, len(nodes))

	for id, task := range nodes {
		depKeys := edges[id]
		var prevTasks []tasks.Task
		for _, depKey := range depKeys {
			prevTasks = append(prevTasks, nodes[depKey])
		}

		build := r.buildStepWithDependencies(&processor.Tasks.Nodes[0], task, prevTasks, initialData)
		steps = append(steps, build)
		tsks = append(tsks, task)
	}

	steps = r.addCommonSteps(steps, tsks, initialData)
	r.wf.Add(steps...)

	return nil
}

func (r *ProcWorkflowRunner) Run(ctx context.Context) error {
	if err := r.wf.Do(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProcWorkflowRunner) buildStepWithDependencies(pt *models.ProcTask, task tasks.Task, prevTasks []tasks.Task, initialData *tasks.TaskData) *flow.AddStep[tasks.Task] {
	step := flow.Step(task).
		AfterStep(func(ctx context.Context, _ flow.Steper, err error) error {
			to := task.GetOutput()
			if err != nil {
				infra.Log.Errorf(msg.ProcTaskFailed, pt.Key, to.WorkspaceID, to.ProcessorID, to.UploadID, err.Error())
				if to.ContinueOnError {
					return nil
				}
				return err
			}
			return nil
		})

	if len(prevTasks) == 0 {
		step = step.Input(func(ctx context.Context, t tasks.Task) error {
			infra.Log.Infof("initial data %+v, task %+v", initialData, t)
			t.AddInput(initialData)
			return nil
		})
	} else {
		for _, prevTask := range prevTasks {
			step = step.DependsOn(prevTask).
				Input(func(ctx context.Context, t tasks.Task) error {
					pto := prevTask.GetOutput()
					pto.ContinueOnError = pt.ContinueOnError
					t.AddInput(pto)
					return nil
				})
		}
	}

	if pt.Retry > 0 {
		step = step.Retry(func(ro *flow.RetryOption) {
			ro.Attempts = pt.Retry
		})
	}

	if pt.TimeoutMilSec > 0 {
		step = step.Timeout(time.Duration(pt.TimeoutMilSec) * time.Second)
	}

	return &step
}

func (r *ProcWorkflowRunner) addCommonSteps(steps []flow.Builder, tsks []flow.Steper, initialData *tasks.TaskData) []flow.Builder {
	success := flow.Step(commontasks.NewSuccessTask()).
		DependsOn(tsks...).
		When(flow.AllSucceededOrSkipped).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.AddInput(initialData)
			return nil
		})

	failure := flow.Step(commontasks.NewFailedTask()).
		DependsOn(tsks...).
		When(flow.AnyFailed).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.AddInput(initialData)
			return nil
		})

	cancelled := flow.Step(commontasks.NewCancelledTask()).
		DependsOn(tsks...).
		When(flow.BeCanceled).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.AddInput(initialData)
			return nil
		})

	cleanup := flow.Step(commontasks.NewCleanupTask()).
		DependsOn(tsks...).
		When(flow.Always).
		Input(func(ctx context.Context, t tasks.Task) error {
			t.AddInput(initialData)
			return nil
		})

	return append(steps, success, failure, cancelled, cleanup)
}

func (r *ProcWorkflowRunner) getNodesAndEdgesMap(processor *models.Processor) (map[string]tasks.Task, map[string][]string, error) {
	nodes := make(map[string]tasks.Task)
	edges := make(map[string][]string)

	for _, t := range processor.Tasks.Nodes {
		task, ok := GetTaskFromRegistry(t.Key)
		if !ok {
			return nil, nil, fmt.Errorf("unknown task key: %s", t.Key)
		}

		nodes[t.ID] = task

		for _, edge := range processor.Tasks.Edges {
			if edge.Source == t.ID {
				edges[t.ID] = append(edges[t.ID], edge.Target)
			}
		}
	}

	return nodes, edges, nil
}
