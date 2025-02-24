package oss

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	tableName string
}

func NewReadCSVCommand(cfg *config.Config) *cobra.Command {
	read := &readCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "csv",
		Short:   "Read content from OSS",
		Example: "opms oss csv",
		RunE:    read.RunE,
	}

	cmd.Flags().StringVarP(&read.name, "name", "n", "", "File name")
	cmd.Flags().StringVarP(&read.tableName, "table", "t", "", "Table name")
	cmd.Flags().StringVarP(&read.bucket, "bucket", "b", "", "Bucket name")
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

	if r.name == "" && r.tableName == "" {
		return errors.New("either name or table name is required")
	}

	name := ""
	if r.name != "" {
		name = r.name
	}

	if r.tableName != "" {
		if name != "" {
			fmt.Println("ignoring the table name, name provided as well")
		} else {
			parts := strings.Split(r.tableName, ".")
			name = fmt.Sprintf("external-table/%s/%s/%s/file.csv", parts[0], parts[1], parts[2])
		}
	}

	err = readFileFromBucket(ctx, client, r.bucket, name, printer)
	if err != nil {
		return err
	}

	err = printer.Render()
	return nil
}

func readFileFromBucket(ctx context.Context, client *oss.Client, bucketName, objectKey string, printer table.Printer) error {
	result, err := client.GetObject(ctx, &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectKey),
	})
	if err != nil {
		return err
	}

	defer result.Body.Close()

	content, err := csv.NewReader(result.Body).ReadAll()
	if err != nil {
		return err
	}

	printer.AddHeader(append([]string{"row num"}, content[0]...))

	for i, row := range content[1:] {
		printer.AddField(strconv.Itoa(i))
		for _, col := range row {
			printer.AddField(col)
		}
		printer.EndRow()
	}
	return nil
}
