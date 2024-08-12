package tasks

import "github.com/Dissociable/Couploan/pkg/services"

func StartTasksRunner(c *services.Container, async bool) {
	// Determine the database based on the environment
	ReconnectRedisPool(c, async)

	// Keep watching the redis client and pinging it, if it fails, try to reconnect
	PoolWatchdog(c)
}
