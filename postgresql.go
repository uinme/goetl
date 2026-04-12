package main

import "fmt"

type PgInputPlugin struct {
	InputPlugin
}

func (p *PgInputPlugin) Transaction(config ConfigSource, run func(schema Schema) (CommitReport, error)) (CommitReport, error) {
	fmt.Println("Inputトランザクション開始")

	schema := Schema{}
	fmt.Println("[PgInputPlugin] Resolved Schema")

	report, err := run(schema)
	if err != nil {
		return CommitReport{}, fmt.Errorf("Failed to control run: %w", err)
	}

	fmt.Println("Inputコミット")

	fmt.Println("Inputトランザクション終了")

	return report, nil
}

func (p *PgInputPlugin) Run(config ConfigSource, schema Schema, writer OutputWriter) (CommitReport, error) {
	fmt.Println("Input処理開始")

	fmt.Println("Inputデータ取得")
	data := []any{}

	// 取得データをOutputへ書き込み
	writer.Write(data)

	fmt.Println("Input処理終了")
	return CommitReport{}, nil
}

type PgOutputWriter struct {
	OutputWriter
}

func (p *PgOutputWriter) Write([]any) error {
	fmt.Println("データ書き込み開始")
	fmt.Println("データ書き込み終了")
	return nil
}

func (p *PgOutputWriter) Commit() (CommitReport, error) {
	fmt.Println("TransactionOutput.Commit() start")
	return CommitReport{Status: Success}, nil
}

func (p *PgOutputWriter) Rollback() error {
	fmt.Println("transactionOutput.Rollback() start")
	return nil
}

func (p *PgOutputWriter) Close() error {
	fmt.Println("transactionOutput.Close() start")
	return nil
}

type PgOutputPlugin struct {
	OutputPlugin
}

func (p *PgOutputPlugin) Transaction(config ConfigSource, schema Schema, run func(schema Schema) (CommitReport, error)) (CommitReport, error) {
	fmt.Println("Output トランザクション開始")

	report, err := run(schema)
	if err != nil {
		return CommitReport{}, fmt.Errorf("Failed to output control: %w", err)
	}

	fmt.Println("Outputコミット")

	fmt.Println("Output トランザクション終了")
	return report, nil
}

func (p *PgOutputPlugin) Open(config ConfigSource, schema Schema) (OutputWriter, error) {
	fmt.Println("OuputWriter オープン")
	pg := PgOutputWriter{}
	return &pg, nil
}
