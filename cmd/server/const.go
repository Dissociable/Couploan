package main

import (
	"github.com/Dissociable/Couploan/pkg/services"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

var (
	c *services.Container
)
