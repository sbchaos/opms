package parse_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sbchaos/consume/par"
	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/opms/cmd/optimus/internal/parse"
)

func TestParseDDL2(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		//  OPTIONS(description="Complete name of the partner")
		a := "CREATE EXTERNAL TABLE `p-gopay-gl-data-mart.gopay_open_loop.external_cross_pollination_reference_date`" + `
(
  brand_id STRING,
  partner_name STRING
)
OPTIONS(
  sheet_range="A:B",
  skip_leading_rows=1,
  format="GOOGLE_SHEETS",
  uris=["https://docs.google.com/spreadsheets/d/1TLYaxt5SpGXoVpy_1zrblXGkZI-oX6QU4Gx2aO-nZGM/edit#gid=1276072195"]
);`
		t3 := time.Now()
		par.Debug = true
		parser, err := parse.NewDDLParser(t, nil)
		assert.NoError(t, err)
		fmt.Println(time.Since(t3).String())

		t1 := time.Now()
		d, err := parser.ParseExternalTable(a)
		fmt.Printf("%+v\n", d)
		t2 := time.Since(t1)
		fmt.Printf("Time taken: %s", t2.String())
		assert.NoError(t, err)
		assert.Equal(t, "p_gopay_id_mart.gopay_open_loop.external_cross_pollination_reference_date", d.FullName)
		assert.Equal(t, "GOOGLE_SHEETS", d.Source.SourceType)
	})

	t.Run("test2", func(t *testing.T) {
		a := "CREATE EXTERNAL TABLE `p-gopay-gl-data-mart.external_aml.dim_external_aml_internal_suspicious_list`" + `	(
  aml_suspicious_user_id STRING OPTIONS(description="gopay user id"),
  aml_suspicious_wallet_id STRING OPTIONS(description="GoPay wallet_id"),
  aml_suspicious_source STRING OPTIONS(description="AML Internal Code, restricted description"),
  aml_source_category STRING,
  aml_source_party STRING,
  year_report STRING
)
OPTIONS(
  sheet_range="aml_internal_suspicious_list!A1:F",
  skip_leading_rows=1,
  format="GOOGLE_SHEETS",
  uris=["https://docs.google.com/spreadsheets/d/1c-nAQTBPSGlDl0Kb7pExa6t-lK6ut6XT0sVYFqJ0gdw/"]
);`

		// base.Debug = true
		parser, err := parse.NewDDLParser(t, nil)
		assert.NoError(t, err)

		bqet, err := parser.ParseExternalTable(a)
		fmt.Printf("%+v\n", bqet)
		assert.NoError(t, err)
		assert.Equal(t, bqet.Name, "p_gopay_id_mart.external_aml.dim_external_aml_internal_suspicious_list")
		assert.Equal(t, len(bqet.Schema), 6)
		assert.Equal(t, bqet.Schema[0].Name, "aml_suspicious_user_id")
		assert.Equal(t, bqet.Schema[0].Type, "STRING")
		assert.Equal(t, bqet.Schema[0].Description, "gopay user id")
	})
	t.Run("test3", func(t *testing.T) {
		input := "CREATE EXTERNAL TABLE `p-gopay-gl-data-mart.gopay_open_loop.external_gopay_strategic_merchant_alfamart_indomaret_hist_voucher_exclusion`" + `
		(
		  voucher_batch_id STRING
		)
		OPTIONS(
		  sheet_range="",
		  skip_leading_rows=1,
		  format="GOOGLE_SHEETS",
		  uris=["https://docs.google.com/spreadsheets/d/1D8IQPkW7TgZ3enrWdpb6MEzhHbmkef6UotrJGL7WV-I/edit#gid=1929690328"]
		);`

		// base.Debug = true
		parser, err := parse.NewDDLParser(t, nil)
		assert.NoError(t, err)

		bqet, err := parser.ParseExternalTable(input)
		fmt.Printf("%+v\n", bqet)
		assert.NoError(t, err)
		assert.Equal(t, bqet.Name, "p_gopay_id_mart.gopay_open_loop.external_gopay_strategic_merchant_alfamart_indomaret_hist_voucher_exclusion")
		assert.Equal(t, len(bqet.Schema), 1)
		assert.Equal(t, bqet.Schema[0].Name, "voucher_batch_id")
		assert.Equal(t, bqet.Schema[0].Type, "STRING")
	})
	t.Run("test4", func(t *testing.T) {
		input := "CREATE EXTERNAL TABLE `data-gojek-id-raw-internal.enterprise_esg_manual_collection.waste_associated_with_offices`" + `
(
  index_number STRING,
  disclosure_number STRING,
  disclosure_title STRING,
  question STRING,
  response STRING,
  supporting_link STRING,
  supporting_file STRING,
  data_owner_email STRING,
  peer_approval STRING,
  country_code STRING,
  entity STRING
)
OPTIONS(
  friendly_name="waste_associated_with_offices",
  description="This table contains information of esg manual data collection",
  labels=[("app-id", "c0ac1a7e-eeab-46ed-a30f-eed93d6b2f77"), ("environment", "production"), ("instance-id", "optimus_job"), ("product-group-id", "a24c9be0-b841-4a83-8727-37dac280df1e"), ("team-id", "a24c9be0-b841-4a83-8727-37dac280df1e")],
  sheet_range="All!A:k",
  format="GOOGLE_SHEETS",
  uris=["https://docs.google.com/spreadsheets/d/1US7gNNvtRnO2Wd_7oL7TTGOccvKQkZAUOHxSbXOmAyU/edit#gid=366027489"]
);`

		// base.Debug = true
		parser, err := parse.NewDDLParser(t, nil)
		assert.NoError(t, err)

		bqet, err := parser.ParseExternalTable(input)
		fmt.Printf("%+v\n", bqet)
		assert.NoError(t, err)
		assert.Equal(t, bqet.Name, "p_gojek_id_raw.enterprise_esg_manual_collection.waste_associated_with_offices")
		assert.Equal(t, len(bqet.Schema), 11)
		assert.Equal(t, bqet.Schema[0].Name, "index_number")
		assert.Equal(t, bqet.Schema[0].Type, "STRING")
	})

	t.Run("test5", func(t *testing.T) {
		input := "CREATE EXTERNAL TABLE `data-gojek-id-raw-internal.enterprise_esg_manual_collection.waste_associated_with_offices`" + `
(
)
OPTIONS(
  friendly_name="waste_associated_with_offices",
  sheet_range="All!A:k",
  format="GOOGLE_SHEETS",
  uris=["https://docs.google.com/spreadsheets/d/1US7gNNvtRnO2Wd_7oL7TTGOccvKQkZAUOHxSbXOmAyU/edit#gid=366027489"]
);`

		// base.Debug = true
		parser, err := parse.NewDDLParser(t, nil)
		assert.NoError(t, err)

		bqet, err := parser.ParseExternalTable(input)
		fmt.Printf("%+v\n", bqet)
		assert.NoError(t, err)
		assert.Equal(t, bqet.Name, "p_gojek_id_raw.enterprise_esg_manual_collection.waste_associated_with_offices")
		assert.Equal(t, len(bqet.Schema), 0)
	})
	t.Run("test6", func(t *testing.T) {
		input := "CREATE EXTERNAL TABLE `data-gojek-id-raw-internal.enterprise_esg_manual_collection.waste_associated_with_offices`" + `
(
  unique_id STRING OPTIONS(description="Identifier that represent unique id populated by backend team for specific event"),
  country_name STRING OPTIONS(description="Country name of GO-JEK operating countries."),
  employee_gender BYTES NOT NULL OPTIONS(description="gender"),
  worker_type_name STRING OPTIONS(description="worker_type_name"),
  entity_name STRING OPTIONS(description="entity name from salesforce."),
  jakarta_hire_date BYTES NOT NULL OPTIONS(description="hire date")
)
OPTIONS(
  description="This table contains information of diversity and inclusion data",
  labels=[("app-id", "c0ac1a7e-eeab-46ed-a30f-eed93d6b2f77"), ("environment", "production"), ("instance-id", "optimus_job"), ("product-group-id", "a24c9be0-b841-4a83-8727-37dac280df1e"), ("team-id", "a24c9be0-b841-4a83-8727-37dac280df1e")],
  sheet_range="SR Template!A:N",
  skip_leading_rows=1,
  format="GOOGLE_SHEETS",
  uris=["https://docs.google.com/spreadsheets/d/1ovTsYTfSkAoifW5Rh1svu72grfCchufEl8My8T0uoDs/edit#gid=0"]
);`

		par.Debug = true
		parser, err := parse.NewDDLParser(t, nil)
		assert.NoError(t, err)

		bqet, err := parser.ParseExternalTable(input)
		fmt.Printf("%+v\n", bqet)
		assert.NoError(t, err)
		assert.Equal(t, bqet.Name, "p_gojek_id_raw.enterprise_esg_manual_collection.waste_associated_with_offices")
		assert.Equal(t, len(bqet.Schema), 6)
	})

	t.Run("test7", func(t *testing.T) {
		input := "CREATE EXTERNAL TABLE `data-gojek-id-raw-internal.enterprise_esg_manual_collection.waste_associated_with_offices`" + `
OPTIONS(
  sheet_range="",
  skip_leading_rows=1,
  format="GOOGLE_SHEETS",
  uris=["https://docs.google.com/spreadsheets/d/1Y1LBpTLVqPo8kfOJo7ykR0RrFO1Cp4levfWltt1rl8M/edit#gid=0"]
);`

		par.Debug = true
		parser, err := parse.NewDDLParser(t, nil)
		assert.NoError(t, err)

		bqet, err := parser.ParseExternalTable(input)
		fmt.Printf("%+v\n", bqet)
		assert.NoError(t, err)
		assert.Equal(t, bqet.Name, "p_gojek_id_raw.enterprise_esg_manual_collection.waste_associated_with_offices")
		assert.Equal(t, len(bqet.Schema), 0)
	})

}
