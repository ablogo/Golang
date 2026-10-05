package services

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	repository "watcher/db"
	"watcher/models"
)

func CheckK8sStatus(repo *repository.Sqlyte, log *slog.Logger, command string, params ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, params...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	var description strings.Builder
	isError := true

	err := cmd.Run()
	if err != nil {
		description.WriteString("cmd run: ")
		s := command + " " + strings.Join(params, " ")
		description.WriteString(s)
		description.WriteString(err.Error())
	}

	stdoutStr := stdoutBuf.String()
	stderrStr := stderrBuf.String()

	if strings.Contains(stdoutStr, "Ready") {
		isError = false
	} else {

		if ctx.Err() == context.DeadlineExceeded {
			description.WriteString(ctx.Err().Error())
		} else {
			description.WriteString("stdout: ")
			description.WriteString(stdoutStr)

			description.WriteString("stderr: ")
			description.WriteString(stderrStr)
		}
	}

	repo.SaveRecord(models.Record{
		App:         "Kubernetes",
		Message:     "",
		Description: description.String(),
		IsError:     isError,
		CreatedAt:   time.Now(),
	})

}

func CheckWebSiteStatus(repo *repository.Sqlyte, log *slog.Logger, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	cmd := exec.CommandContext(ctx, "curl", "-sL", "-o", "/dev/null", "-w", "%{http_code}", url)
	stdout, err := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	var description strings.Builder
	isError := true

	if err := cmd.Start(); err != nil {
		description.WriteString("cmd start: ")
		description.WriteString(err.Error())
	}

	w := io.Writer(&bytes.Buffer{})
	scanner := bufio.NewScanner(stderr)

	io.Copy(w, stdout)
	output := w.(*bytes.Buffer).String()

	var outputErr strings.Builder
	for scanner.Scan() {
		outputErr.WriteString(scanner.Text())
	}

	if err = cmd.Wait(); err != nil {
		description.WriteString("cmd wait: ")
		description.WriteString(err.Error())
	}

	if strings.Contains(output, strconv.Itoa(http.StatusOK)) {
		isError = false
	} else {
		description.WriteString("cmd out: ")
		description.WriteString(outputErr.String())

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			description.WriteString(ctx.Err().Error())
		}
	}

	repo.SaveRecord(models.Record{
		App:         "Auth Service",
		Message:     "",
		Description: description.String(),
		IsError:     isError,
		CreatedAt:   time.Now(),
	})

}
