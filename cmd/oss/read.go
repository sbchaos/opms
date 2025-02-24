package oss

import (
	"context"
	"encoding/csv"
	"os"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

const timeout = 1 * time.Minute

type readCommand struct {
	cfg *config.Config

	name   string
	bucket string
}

func NewReadCommand(cfg *config.Config) *cobra.Command {
	read := &readCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "read",
		Short:   "Read content from OSS",
		Example: "opms oss read",
		RunE:    read.RunE,
	}

	cmd.Flags().StringVarP(&read.name, "name", "n", "", "File name")
	cmd.Flags().StringVarP(&read.name, "bucket", "b", "", "Bucket name")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("bucket")
	return cmd
}

func (r *readCommand) RunE(_ *cobra.Command, _ []string) error {
	client, err := mc.NewOSSClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)

	err = readFileFromBucket(ctx, client, r.bucket, r.name, printer)
	if err != nil {
		return err
	}

	err = printer.Render()
	return nil
}

func readFileFromBucket(ctx context.Context, client *oss.Client, bucketName, objectKey string, printer table.Printer) error {
	getObject, err := client.GetObject(ctx, &oss.GetObjectRequest{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		return err
	}

	defer getObject.Body.Close()

	content, err := csv.NewReader(getObject.Body).ReadAll()
	if err != nil {
		return err
	}

	printer.AddHeader(content[0])

	for _, row := range content[1:] {
		for _, col := range row {
			printer.AddField(col)
		}
		printer.EndRow()
	}
	return nil
}
