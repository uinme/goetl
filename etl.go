package main

import "fmt"

type ConfigSource struct {
	In  InputSource
	Out OutputSource
}

type InputSource struct {
	Type_           string
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	Schema          string
	DefaultTimezone string
	Table           string
	Query           string
}

type OutputSource struct {
	Type_           string
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	Schema          string
	DefaultTimezone string
	Table           string
	Mode            string
	MergeMode       string
}

type DbColumn struct {
	Name     string
	TypeName string
}

type Schema struct {
	Column []DbColumn
}

type CommitStatus int

const (
	Success CommitStatus = iota
	Failed
)

type CommitReport struct {
	Status CommitStatus
}

type Page struct {
}

type OutputWriter interface {
	Write([]any) error
	Commit() (CommitReport, error)
	Rollback() error
	Close() error
}

type InputPlugin interface {
	Transaction(ConfigSource, func(Schema) (CommitReport, error)) (CommitReport, error)
	Run(ConfigSource, Schema, OutputWriter) (CommitReport, error)
}

type OutputPlugin interface {
	Transaction(ConfigSource, Schema, func(Schema) (CommitReport, error)) (CommitReport, error)
	Open(ConfigSource, Schema) (OutputWriter, error)
}

type Etl struct {
	config       ConfigSource
	inputPlugin  InputPlugin
	outputPlugin OutputPlugin
}

func NewEtl(path string) (Etl, error) {

	return Etl{}, nil
}

func (e *Etl) Run() error {
	fmt.Println("=== ETL 開始 ===")

	e.config = ConfigSource{}
	e.inputPlugin = &PgInputPlugin{}
	e.outputPlugin = &PgOutputPlugin{}

	// 1. Input トランザクション開始
	inputReport, err := e.inputPlugin.Transaction(e.config, e.inputControl)
	if err != nil {
		return fmt.Errorf("Input transaction error")
	}

	fmt.Print(inputReport)

	fmt.Println("=== ELT 終了 ===")

	return nil
}

func (e *Etl) inputControl(schema Schema) (CommitReport, error) {
	fmt.Println("Input control 開始")

	// 2. Output トランザクション開始
	outputReport, err := e.outputPlugin.Transaction(e.config, schema, e.outputControl)
	if err != nil {
		return CommitReport{Status: Failed}, fmt.Errorf("Output transaction error")
	}

	fmt.Println(outputReport)

	fmt.Println("Input control 終了")

	return outputReport, nil
}

func (e *Etl) outputControl(schema Schema) (CommitReport, error) {
	fmt.Println("Output control 開始")

	txOutput, err := e.outputPlugin.Open(e.config, schema)
	if err != nil {
		return CommitReport{}, fmt.Errorf("transactional ouput error")
	}

	// 3. InputPlugin.Run 実行
	inputReport, err := e.inputPlugin.Run(e.config, schema, txOutput)
	if err != nil {
		return CommitReport{}, fmt.Errorf("failed to input.run")
	}

	fmt.Println("Output contorl 終了")

	return inputReport, nil
}
