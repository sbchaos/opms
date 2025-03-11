package oss

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/cmdutil"
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
	output    string
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
	cmd.Flags().StringVarP(&read.output, "output", "o", "", "Write CSV before printing")
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

	proj := ""
	if r.tableName != "" {
		if name != "" {
			fmt.Println("ignoring the table name, name provided as well")
		} else {
			parts := strings.Split(r.tableName, ".")
			name = fmt.Sprintf("external-table/%s/%s/%s/file.csv", parts[0], parts[1], parts[2])
			proj = parts[0]
		}
	}

	if r.bucket == "" {
		fromVar, err := cmdutil.GetArgFromVar[string](r.cfg, "csv", proj, "bucket")
		if err != nil {
			return errors.New("bucket name is required, not found in vars")
		}
		r.bucket = fromVar
	}

	err = r.readFileFromBucket(ctx, client, name, printer)
	if err != nil {
		return err
	}

	printer.Render()
	return nil
}

func (r *readCommand) readFileFromBucket(ctx context.Context, client *oss.Client, objectKey string, printer table.Printer) error {
	result, err := client.GetObject(ctx, &oss.GetObjectRequest{
		Bucket: oss.Ptr(r.bucket),
		Key:    oss.Ptr(objectKey),
	})
	if err != nil {
		return err
	}

	defer result.Body.Close()

	content, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}

	if r.output != "" {
		cmdutil.WriteFile(r.output, content)
	}

	csvStr, err := csv.NewReader(bytes.NewReader(content)).ReadAll()
	if err != nil {
		return err
	}

	printer.AddHeader(append([]string{"row num"}, csvStr[0]...))

	for i, row := range csvStr[1:] {
		printer.AddField(strconv.Itoa(i + 1))
		for _, col := range row {
			printer.AddField(col)
		}
		printer.EndRow()
	}
	return nil
}
