package main

import "fmt"

type ConfigSource struct {
	In  InputSource  `yaml:"in"`
	Out OutputSource `yaml:"out"`
}

type InputSource struct {
	Type_           string `yaml:"type"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Schema          string `yaml:"schema"`
	DefaultTimezone string `yaml:"default_timezone"`
	Table           string `yaml:"table"`
	Query           string `yaml:"query"`
	BeforeSelect    string `yaml:"before_select"`
	AfterSelect     string `yaml:"after_select"`
}

type OutputSource struct {
	Type_           string   `yaml:"type"`
	Host            string   `yaml:"host"`
	Port            int      `yaml:"port"`
	User            string   `yaml:"user"`
	Password        string   `yaml:"password"`
	Database        string   `yaml:"databse"`
	Schema          string   `yaml:"schema"`
	DefaultTimezone string   `yaml:"default_timezone"`
	Table           string   `yaml:"table"`
	Mode            string   `yaml:"mode"`
	MergeMode       string   `yaml:"merge_mode"`
	MergeKeys       []string `yaml:"merge_keys"`
	BeforeLoad      string   `yaml:"before_load"`
	AfterLoad       string   `yaml:"after_load"`
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

func NewEtl(config ConfigSource) (Etl, error) {
	etl := Etl{}

	etl.config = config

	switch config.In.Type_ {
	case "postgresql":
		etl.inputPlugin = &PgInputPlugin{}
	default:
		return etl, fmt.Errorf("%sは実装されていません", config.In.Type_)
	}

	switch config.Out.Type_ {
	case "postgresql":
		etl.outputPlugin = &PgOutputPlugin{}
	default:
		return etl, fmt.Errorf("%sは実装されていません", config.Out.Type_)
	}

	return etl, nil
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
