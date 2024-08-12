package tasks

import (
	"fmt"
	"github.com/Dissociable/Couploan/config"
	"github.com/Dissociable/Couploan/pkg/services"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

func StartInspector(
	sharedRedisPool *SharedRedisPool, c *services.Container, inspector *asynq.Inspector,
) {
	// Determine the database based on the environment
	db := c.Config.Cache.Database
	if c.Config.App.Environment == config.EnvTest {
		db = c.Config.Cache.TestDatabase
	}
	if sharedRedisPool != nil {
		*inspector = *asynq.NewInspector(sharedRedisPool)
	} else {
		*inspector = *asynq.NewInspector(
			asynq.RedisClientOpt{
				Addr:     fmt.Sprintf("%s:%d", c.Config.Cache.Hostname, c.Config.Cache.Port),
				DB:       db,
				Username: c.Config.Cache.Username,
				Password: c.Config.Cache.Password,
			},
		)
	}
}

// CheckIncompleteTasksForTaskExistence Check incomplete tasks for existence of a task id in the specified queue
func CheckIncompleteTasksForTaskExistence(
	queue string, taskId string,
) (exists bool, err error) {
	funcs := []func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error){
		asynqInspector.ListPendingTasks,
		asynqInspector.ListRetryTasks,
		asynqInspector.ListArchivedTasks,
		asynqInspector.ListScheduledTasks,
		asynqInspector.ListActiveTasks,
	}
	for _, fn := range funcs {
		exists, err = CheckTaskQueueForTaskExistence(fn, queue, taskId)
		if err != nil {
			err = errors.Wrapf(err, "failed to check task queue for task id %s", taskId)
			return
		}
		if exists {
			return
		}
	}
	return false, nil
}

func CheckTaskQueueForTaskExistence(
	listFunc func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error),
	queue string, taskId string,
) (exists bool, err error) {
	tasks, err := listFunc(queue, asynq.PageSize(0))
	if err != nil {
		if errors.Is(err, asynq.ErrQueueNotFound) {
			return false, nil
		}
		err = errors.Wrap(err, "failed to list tasks")
		return
	}
	for _, task := range tasks {
		return task.ID == taskId, nil
	}
	return false, nil
}
