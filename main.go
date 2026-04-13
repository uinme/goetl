package main

import (
	"etl/liquid"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	err := run()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}
}

func run() error {
	// Liquid template 解析
	parserd, err := liquid.Parse("C:Users/uinme/go_workspace/etl/t_goetl_test.yml.liquid")
	if err != nil {
		return fmt.Errorf("%w\n", err)
	}

	// Yaml 解析
	var config ConfigSource
	err = yaml.Unmarshal(parserd, &config)
	if err != nil {
		return fmt.Errorf("YAMLパースエラー: %w", err)
	}

	// Etl初期化
	etl, err := NewEtl(config)
	if err != nil {
		return fmt.Errorf("Etl 初期化エラー: %w", err)
	}

	// Etl実行
	etl.Run()

	return nil
}
