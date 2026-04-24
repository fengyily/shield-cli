package updater

import (
	"context"
	"sync"
	"time"
)

// Job tracks the state of an in-flight upgrade. Consumed by the Web UI via
// polling /api/version/upgrade/status.
type Job struct {
	mu sync.Mutex

	Running    bool      `json:"running"`
	Stage      string    `json:"stage"`   // download | verify | extract | install | restart | done | error
	Percent    int       `json:"percent"` // 0-100
	Error      string    `json:"error,omitempty"`
	TargetTag  string    `json:"target_tag,omitempty"`
	StartedAt  time.Time `json:"started_at,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}

func (j *Job) Snapshot() Job {
	j.mu.Lock()
	defer j.mu.Unlock()
	return Job{
		Running:    j.Running,
		Stage:      j.Stage,
		Percent:    j.Percent,
		Error:      j.Error,
		TargetTag:  j.TargetTag,
		StartedAt:  j.StartedAt,
		FinishedAt: j.FinishedAt,
	}
}

func (j *Job) set(stage string, pct int) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Stage = stage
	j.Percent = pct
}

func (j *Job) fail(err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Running = false
	j.Stage = "error"
	j.Error = err.Error()
	j.FinishedAt = time.Now()
}

func (j *Job) finish(stage string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Running = false
	j.Stage = stage
	j.Percent = 100
	j.FinishedAt = time.Now()
}

func (j *Job) start(tag string) bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.Running {
		return false
	}
	j.Running = true
	j.Stage = "starting"
	j.Percent = 0
	j.Error = ""
	j.TargetTag = tag
	j.StartedAt = time.Now()
	j.FinishedAt = time.Time{}
	return true
}

// Run performs Apply + Restart in the background. The returned job is also
// mutated in place so callers can poll Snapshot() for progress.
func (c *Checker) Run(ctx context.Context, job *Job, tag string, asService bool, port int) bool {
	if !job.start(tag) {
		return false
	}
	go func() {
		if err := Apply(ctx, tag, job.set); err != nil {
			job.fail(err)
			return
		}
		job.set("restart", 100)
		if err := Restart(asService, port); err != nil {
			job.fail(err)
			return
		}
		job.finish("restart")
	}()
	return true
}
