package emr_serverless //nolint

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/emrserverless"
	"github.com/aws/aws-sdk-go-v2/service/emrserverless/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bruin-data/bruin/pkg/git"
	"github.com/bruin-data/bruin/pkg/pipeline"
	"github.com/google/uuid"
)

type jobRunError struct {
	RunID   string
	Details string
	State   types.JobRunState
}

func (e jobRunError) Error() string {
	switch e.State { //nolint:exhaustive
	case types.JobRunStateFailed:
		return fmt.Sprintf("job run %s failed", e.RunID)
	case types.JobRunStateCancelled:
		return fmt.Sprintf("job run %s was cancelled", e.RunID)
	default:
		return fmt.Sprintf("job run %s is in an unknown state: %s", e.RunID, e.State)
	}
}

type JobRunParams struct {
	ApplicationID string
	ExecutionRole string
	Entrypoint    string
	Args          []string
	Config        string
	Timeout       time.Duration
	Region        string
	Logs          string
	Workspace     string
}

func parseParams(cfg *Client, params map[string]string) *JobRunParams {
	jobParams := JobRunParams{
		ApplicationID: cfg.ApplicationID,
		ExecutionRole: cfg.ExecutionRole,
		Region:        cfg.Region,
		Entrypoint:    params["entrypoint"],
		Config:        params["config"],
		Logs:          params["logs"],
		Workspace:     cfg.Workspace,
	}

	if params["timeout"] != "" {
		t, err := time.ParseDuration(params["timeout"])
		if err == nil {
			jobParams.Timeout = t
		}
	}
	if params["args"] != "" {
		arglist := strings.Split(strings.TrimSpace(params["args"]), " ")
		for _, arg := range arglist {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				jobParams.Args = append(jobParams.Args, arg)
			}
		}
	}
	return &jobParams
}

type Job struct {
	logger    *log.Logger
	emrClient *emrserverless.Client
	s3Client  *s3.Client
	asset     *pipeline.Asset
	pipeline  *pipeline.Pipeline
	params    *JobRunParams
	poll      *PollTimer
	env       map[string]string
}

func (job Job) buildJobRunConfig() *emrserverless.StartJobRunInput {
	driver := &types.JobDriverMemberSparkSubmit{
		Value: types.SparkSubmit{
			EntryPoint:          &job.params.Entrypoint,
			EntryPointArguments: job.params.Args,
		},
	}

	submitParams := job.params.Config
	for key, val := range job.env {
		submitParams += fmt.Sprintf(` --conf spark.executorEnv.%s=%s`, key, val)
		submitParams += fmt.Sprintf(` --conf spark.emr-serverless.driverEnv.%s=%s`, key, val)
	}

	if submitParams != "" {
		driver.Value.SparkSubmitParameters = &submitParams
	}

	cfg := &emrserverless.StartJobRunInput{
		ApplicationId:           &job.params.ApplicationID,
		Name:                    &job.asset.Name,
		ExecutionRoleArn:        &job.params.ExecutionRole,
		ExecutionTimeoutMinutes: aws.Int64(int64(job.params.Timeout.Minutes())),
		JobDriver:               driver,
	}

	if job.params.Logs != "" {
		cfg.ConfigurationOverrides = &types.ConfigurationOverrides{
			MonitoringConfiguration: &types.MonitoringConfiguration{
				S3MonitoringConfiguration: &types.S3MonitoringConfiguration{
					LogUri: aws.String(job.params.Logs),
				},
			},
		}
	}

	return cfg
}

type workspace struct {
	Root       *url.URL
	Entrypoint string
	Files      string
	Logs       string
}

// prepareWorkspace sets up an s3 bucket for a pyspark job run.
func (job Job) prepareWorkspace(ctx context.Context) (*workspace, error) {
	workspaceURI, err := url.Parse(job.params.Workspace)
	if err != nil {
		return nil, fmt.Errorf("error parsing workspace URL: %w", err)
	}
	jobID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("error generating job ID: %w", err)
	}
	jobURI := workspaceURI.JoinPath(job.pipeline.Name, jobID.String())

	scriptPath := job.asset.ExecutableFile.Path
	fd, err := os.Open(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %q: %w", scriptPath, err)
	}

	defer fd.Close()

	scriptURI := jobURI.JoinPath(job.asset.ExecutableFile.Name)
	_, err = job.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &scriptURI.Host,
		Key:    aws.String(strings.TrimPrefix(scriptURI.Path, "/")),
		Body:   fd,
	})
	if err != nil {
		return nil, fmt.Errorf("error uploading entrypoint %q: %w", scriptURI, err)
	}

	fd, err = os.CreateTemp("", "bruin-spark-context-*.zip")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file %w", err)
	}
	defer os.Remove(fd.Name())
	defer fd.Close()

	zipper := zip.NewWriter(fd)
	defer zipper.Close()

	repo, err := git.FindRepoFromPath(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("error finding project root: %w", err)
	}

	err = packageContextWithPrefix(
		zipper,
		os.DirFS(repo.Path),
	)
	if err != nil {
		return nil, fmt.Errorf("error packaging files: %w", err)
	}
	err = zipper.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing zip writer: %w", err)
	}
	_, err = fd.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error rewinding file %q: %w", fd.Name(), err)
	}

	contextURI := jobURI.JoinPath("context.zip")
	_, err = job.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &contextURI.Host,
		Key:    aws.String(strings.TrimPrefix(contextURI.Path, "/")),
		Body:   fd,
	})
	if err != nil {
		return nil, fmt.Errorf("error uploading context %q: %w", contextURI, err)
	}

	return &workspace{
		Root:       jobURI,
		Entrypoint: scriptURI.String(),
		Files:      contextURI.String(),
		Logs:       workspaceURI.JoinPath("logs").String(),
	}, nil
}

func (job Job) deleteWorkspace(ws *workspace) {
	// todo(turtledev)
	//   * pagination
	//   * debug logs for errors
	//   * timeout for cleanup

	listArgs := &s3.ListObjectsV2Input{
		Bucket: &ws.Root.Host,
		Prefix: aws.String(strings.TrimPrefix(ws.Root.Path, "/")),
	}

	objs, err := job.s3Client.ListObjectsV2(context.Background(), listArgs)
	if err != nil {
		return
	}

	for _, obj := range objs.Contents {
		job.s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{ //nolint
			Bucket: &ws.Root.Host,
			Key:    obj.Key,
		})
	}
}

func (job Job) Run(ctx context.Context) (err error) {
	if job.asset.Type == pipeline.AssetTypeEMRServerlessPyspark {
		ws, err := job.prepareWorkspace(ctx)
		if err != nil {
			return fmt.Errorf("error preparing workspace: %w", err)
		}
		job.params.Entrypoint = ws.Entrypoint
		job.params.Config += " --conf spark.submit.pyFiles=" + ws.Files

		// only use workspace for logs if the
		// asset doesn't explicitly specify it
		if job.params.Logs == "" {
			job.params.Logs = ws.Logs
		}
		defer job.deleteWorkspace(ws) //nolint
	}

	run, err := job.emrClient.StartJobRun(ctx, job.buildJobRunConfig())
	if err != nil {
		return fmt.Errorf("error submitting job run: %w", err)
	}
	job.logger.Printf("created job run: %s", *run.JobRunId)
	defer func() { //nolint
		if err != nil && !errors.As(err, &jobRunError{}) {
			// todo(turtledev): timeout for cancellation
			job.logger.Printf("error detected. cancelling job run.")
			job.emrClient.CancelJobRun(context.Background(), &emrserverless.CancelJobRunInput{ //nolint
				ApplicationId: &job.params.ApplicationID,
				JobRunId:      run.JobRunId,
			})
		}
	}()

	var (
		previousState    = types.JobRunState("unknown")
		paginationToken  = ""
		maxAttemptsError = &retry.MaxAttemptsError{}
		jobLogs          = job.buildLogConsumer(ctx, run)
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(job.poll.Duration()):
			listJobArgs := &emrserverless.ListJobRunAttemptsInput{
				ApplicationId: &job.params.ApplicationID,
				JobRunId:      run.JobRunId,
			}
			if paginationToken != "" {
				listJobArgs.NextToken = &paginationToken
			}
			runs, err := job.emrClient.ListJobRunAttempts(ctx, listJobArgs)
			if errors.As(err, &maxAttemptsError) {
				job.poll.Increase()
				continue
			}
			job.poll.Reset()

			if err != nil {
				return fmt.Errorf("error checking job run status: %w", err)
			}
			if len(runs.JobRunAttempts) == 0 {
				return errors.New("job runs not found")
			}

			latestRun := runs.JobRunAttempts[len(runs.JobRunAttempts)-1]
			if previousState != latestRun.State {
				job.logger.Printf(
					"%s | %s | %s",
					*run.JobRunId,
					latestRun.State,
					*latestRun.StateDetails,
				)
				previousState = latestRun.State
			}
			for _, line := range jobLogs.Next() {
				job.logger.Printf("%s | %s | %s ", line.Source.Name, line.Source.Stream, line.Message)
			}

			switch latestRun.State { //nolint:exhaustive
			case types.JobRunStateFailed, types.JobRunStateCancelled:
				return jobRunError{
					RunID:   *run.JobRunId,
					State:   latestRun.State,
					Details: *latestRun.StateDetails,
				}
			case types.JobRunStateSuccess:
				return nil
			}

			if runs.NextToken != nil {
				paginationToken = *runs.NextToken
			}
		}
	}
}

type LogConsumer interface {
	Next() []LogLine
}

func (job Job) buildLogConsumer(ctx context.Context, run *emrserverless.StartJobRunOutput) LogConsumer {
	logURI := job.resolveLogURI(ctx, run)
	if logURI != "" {
		uri, err := url.Parse(logURI)
		if err == nil {
			return &S3LogConsumer{
				Ctx:   ctx,
				URI:   uri,
				S3cli: job.s3Client,
				RunID: *run.JobRunId,
				AppID: *run.ApplicationId,
			}
		}
	}

	return &NoOpLogConsumer{}
}

func (job Job) resolveLogURI(ctx context.Context, run *emrserverless.StartJobRunOutput) string {
	if job.params.Logs != "" {
		return job.params.Logs
	}

	app, err := job.emrClient.GetApplication(
		ctx,
		&emrserverless.GetApplicationInput{
			ApplicationId: run.ApplicationId,
		},
	)
	if err != nil {
		return ""
	}

	monitoringCfg := *app.Application.MonitoringConfiguration
	if monitoringCfg.S3MonitoringConfiguration != nil && monitoringCfg.S3MonitoringConfiguration.LogUri != nil {
		return *monitoringCfg.S3MonitoringConfiguration.LogUri
	}

	return ""
}
