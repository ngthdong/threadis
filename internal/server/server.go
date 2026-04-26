package server

import (
	"io"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/ngthdong/threadis/internal/config"
	"github.com/ngthdong/threadis/internal/constant"
	"github.com/ngthdong/threadis/internal/core"
	"github.com/ngthdong/threadis/internal/core/iomultiplexing"
)

func readCommand(fd int) (*core.Command, error) {
	buf := make([]byte, 512)

	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}

	return core.ParseCmd(buf[:n])
}

func RunServer() {
	log.Println("Starting server on", config.Port)

	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("listener is not a TCPListener")
	}

	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	ioMultiplexer, err := iomultiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMultiplexer.Close()

	if err = ioMultiplexer.Monitor(iomultiplexing.Event{
		Fd: serverFd,
		Op: iomultiplexing.OpRead,
	}); err != nil {
		log.Fatal(err)
	}

	lastActiveExpireExecTime := time.Now()

	for {
		if time.Since(lastActiveExpireExecTime) >= constant.ActiveExpireFrequency {
			core.ActiveDeleteExpiredKeys()
			lastActiveExpireExecTime = time.Now()
		}

		events, err := ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for i := 0; i < len(events); i++ {
			fd := events[i].Fd

			if fd == serverFd {
				connFd, _, err := syscall.Accept(serverFd)
				if err != nil {
					continue
				}

				if err = ioMultiplexer.Monitor(iomultiplexing.Event{
					Fd: connFd,
					Op: iomultiplexing.OpRead,
				}); err != nil {
					_ = syscall.Close(connFd)
				}
				continue
			}

			cmd, err := readCommand(fd)
			if err != nil {
				if err == io.EOF || err == syscall.ECONNRESET {
					_ = syscall.Close(fd)
				}
				continue
			}

			if cmd == nil {
				continue
			}

			if err = core.ExecuteAndResponse(cmd, fd); err != nil {
				_ = syscall.Close(fd)
			}
		}
	}
}