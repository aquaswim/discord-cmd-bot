package main

import (
	"discord-cmd-bot/config"
	"fmt"
	"os/exec"
)

type CommandRunner interface {
	RunCommand(cmd string) (string, error)
	HasCommand(cmd string) bool
}

type commandRunner struct {
	cmdMap map[string]*config.Command
}

func (c *commandRunner) HasCommand(cmd string) bool {
	_, found := c.cmdMap[cmd]
	return found
}

func (c *commandRunner) RunCommand(cmd string) (string, error) {
	command, err := c.getCommand(cmd)
	if err != nil {
		return "", err
	}
	_cmd := exec.Command(command.Syntax, command.Args...)
	output, err := _cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (c *commandRunner) getCommand(cmd string) (*config.Command, error) {
	command := c.cmdMap[cmd]
	if command != nil {
		return command, nil
	}
	return nil, fmt.Errorf("command %s not found", cmd)
}

func NewCommandRunner(cfg *config.Config) CommandRunner {
	cmdMap := map[string]*config.Command{}
	for i := 0; i < len(cfg.Commands); i++ {
		_cmd := &cfg.Commands[i]
		cmdMap[_cmd.Command] = _cmd
	}
	return &commandRunner{cmdMap: cmdMap}
}
