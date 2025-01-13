package osx

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Cli struct {
	IP         string
	Username   string
	Password   string
	Port       int
	client     *ssh.Client
	LastResult string
}

func NewSSH(ip string, username string, password string, port ...int) *Cli {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	switch {
	case len(port) <= 0:
		cli.Port = 22
	case len(port) > 0:
		cli.Port = port[0]
	}
	return cli
}

func (c *Cli) connect() error {
	config := ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}
func (c *Cli) RunSH(cmd string) (_err error) {
	if c.client == nil {
		if _err = c.connect(); _err != nil {
			return _err
		}
	}
	session, _err := c.client.NewSession()
	if _err != nil {
		return _err
	}
	defer session.Close()

	//执行远程命令
	combo, _err := session.CombinedOutput(cmd)
	if _err != nil {
		log.Fatal("远程执行cmd 失败", _err)
	}
	log.Println("命令输出:", string(combo))
	return
}

func (c *Cli) RunTerminal(stdout, stderr io.Writer) error {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := term.GetSize(fd)
	if err != nil {
		panic(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	session.Shell()
	session.Wait()
	return nil
}
